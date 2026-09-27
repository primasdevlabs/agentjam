package dispatcher

import (
	"fmt"
	"os"
	"path/filepath"
	goruntime "runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/primasdevlabs/agentjam/pkg/core"
	"github.com/primasdevlabs/agentjam/pkg/memory"
)

// ToolHandler defines tool execution handler in Go.
type ToolHandler func(args map[string]interface{}) (interface{}, error)

type registeredTool struct {
	manifest core.ToolManifest
	handler  ToolHandler
}

// ToolDispatcherOptions configures tool security parameters.
type ToolDispatcherOptions struct {
	AllowDestructive   bool
	AllowAdmin         bool
	ExecutionTimeoutMs int64
}

// ToolDispatcher handles tool registration, schema validation, safety enforcement, and execution.
type ToolDispatcher struct {
	mu            sync.RWMutex
	workspaceRoot string
	tools         map[string]registeredTool
	memoryManager *memory.MemoryManager
	options       ToolDispatcherOptions
}

// NewToolDispatcher initializes a ToolDispatcher.
func NewToolDispatcher(workspaceRoot string, opts ToolDispatcherOptions, mm *memory.MemoryManager) *ToolDispatcher {
	td := &ToolDispatcher{
		workspaceRoot: workspaceRoot,
		tools:         make(map[string]registeredTool),
		memoryManager: mm,
		options:       opts,
	}
	return td
}

// RegisterTool registers a tool manifest and handler.
func (td *ToolDispatcher) RegisterTool(manifest core.ToolManifest, handler ToolHandler) {
	td.mu.Lock()
	defer td.mu.Unlock()
	td.tools[manifest.Name] = registeredTool{
		manifest: manifest,
		handler:  handler,
	}
}

// UnregisterTool removes a previously registered tool.
func (td *ToolDispatcher) UnregisterTool(name string) bool {
	td.mu.Lock()
	defer td.mu.Unlock()
	if _, ok := td.tools[name]; !ok {
		return false
	}
	delete(td.tools, name)
	return true
}

// ListTools returns all registered tool manifests.
func (td *ToolDispatcher) ListTools() []core.ToolManifest {
	td.mu.RLock()
	defer td.mu.RUnlock()

	manifests := make([]core.ToolManifest, 0, len(td.tools))
	for _, t := range td.tools {
		manifests = append(manifests, t.manifest)
	}
	sort.Slice(manifests, func(i, j int) bool { return manifests[i].Name < manifests[j].Name })
	return manifests
}

func (td *ToolDispatcher) checkSafety(manifest core.ToolManifest) (bool, string) {
	level := manifest.SafetyLevel
	if level == "" {
		level = core.SafetyReadOnly
	}

	if level == core.SafetyDestructive && !td.options.AllowDestructive {
		return false, fmt.Sprintf("Tool '%s' requires destructive safety permission.", manifest.Name)
	}
	if level == core.SafetyAdmin && !td.options.AllowAdmin {
		return false, fmt.Sprintf("Tool '%s' requires admin safety permission.", manifest.Name)
	}
	return true, ""
}

// withinRoot reports whether path p is the workspace root or nested inside it.
func withinRoot(root, p string) bool {
	root = filepath.Clean(root)
	p = filepath.Clean(p)
	if goruntime.GOOS == "windows" {
		root = strings.ToLower(root)
		p = strings.ToLower(p)
	}
	return p == root || strings.HasPrefix(p, root+string(os.PathSeparator))
}

func (td *ToolDispatcher) validatePathSafety(args map[string]interface{}) bool {
	pathKeys := []string{"filePath", "path", "directoryPath", "targetFile", "cwd"}
	for _, key := range pathKeys {
		if val, ok := args[key].(string); ok && filepath.IsAbs(val) && !withinRoot(td.workspaceRoot, val) {
			return false
		}
	}
	return true
}

// execute runs a tool handler, enforcing the configured execution timeout.
func (td *ToolDispatcher) execute(handler ToolHandler, args map[string]interface{}) (out interface{}, err error) {
	timeoutMs := td.options.ExecutionTimeoutMs
	if timeoutMs <= 0 {
		return handler(args)
	}

	type result struct {
		out interface{}
		err error
	}
	done := make(chan result, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				done <- result{err: fmt.Errorf("tool handler panicked: %v", r)}
			}
		}()
		o, e := handler(args)
		done <- result{out: o, err: e}
	}()

	select {
	case res := <-done:
		return res.out, res.err
	case <-time.After(time.Duration(timeoutMs) * time.Millisecond):
		return nil, fmt.Errorf("tool execution timed out after %dms", timeoutMs)
	}
}

// Dispatch executes a tool call safely.
func (td *ToolDispatcher) Dispatch(call core.ToolCall) core.ToolResult {
	startTime := time.Now()

	td.mu.RLock()
	t, exists := td.tools[call.ToolName]
	td.mu.RUnlock()

	if !exists {
		return core.ToolResult{
			ID:         call.ID,
			ToolName:   call.ToolName,
			Status:     "error",
			Error:      fmt.Sprintf("Tool '%s' is not registered.", call.ToolName),
			DurationMs: time.Since(startTime).Milliseconds(),
		}
	}

	if allowed, reason := td.checkSafety(t.manifest); !allowed {
		return core.ToolResult{
			ID:         call.ID,
			ToolName:   call.ToolName,
			Status:     "blocked",
			Error:      reason,
			DurationMs: time.Since(startTime).Milliseconds(),
		}
	}

	if !td.validatePathSafety(call.Arguments) {
		return core.ToolResult{
			ID:         call.ID,
			ToolName:   call.ToolName,
			Status:     "blocked",
			Error:      fmt.Sprintf("Path arguments must reside within workspace root '%s'", td.workspaceRoot),
			DurationMs: time.Since(startTime).Milliseconds(),
		}
	}

	out, err := td.execute(t.handler, call.Arguments)
	durationMs := time.Since(startTime).Milliseconds()

	if err != nil {
		result := core.ToolResult{
			ID:         call.ID,
			ToolName:   call.ToolName,
			Status:     "error",
			Error:      err.Error(),
			DurationMs: durationMs,
		}
		if td.memoryManager != nil {
			td.memoryManager.Set("tool_call:"+call.ID, result, core.MemoryScopeEpisodic, []string{"tool_call_error"}, 0)
		}
		return result
	}

	result := core.ToolResult{
		ID:         call.ID,
		ToolName:   call.ToolName,
		Status:     "success",
		Output:     out,
		DurationMs: durationMs,
	}

	if td.memoryManager != nil {
		td.memoryManager.Set("tool_call:"+call.ID, result, core.MemoryScopeEpisodic, []string{"tool_call"}, 0)
	}

	return result
}
