package dispatcher

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/agentjam/agentjam/pkg/core"
	"github.com/agentjam/agentjam/pkg/memory"
)

// ToolHandler defines tool execution handler in Go.
type ToolHandler func(args map[string]interface{}) (interface{}, error)

type registeredTool struct {
	manifest core.ToolManifest
	handler  ToolHandler
}

// ToolDispatcherOptions configures tool security parameters.
type ToolDispatcherOptions struct {
	AllowDestructive bool
	AllowAdmin       bool
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

// ListTools returns all registered tool manifests.
func (td *ToolDispatcher) ListTools() []core.ToolManifest {
	td.mu.RLock()
	defer td.mu.RUnlock()

	manifests := make([]core.ToolManifest, 0, len(td.tools))
	for _, t := range td.tools {
		manifests = append(manifests, t.manifest)
	}
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

func (td *ToolDispatcher) validatePathSafety(args map[string]interface{}) bool {
	pathKeys := []string{"filePath", "path", "directoryPath", "targetFile", "cwd"}
	rootClean := filepath.Clean(td.workspaceRoot)

	for _, key := range pathKeys {
		if val, ok := args[key].(string); ok && filepath.IsAbs(val) {
			valClean := filepath.Clean(val)
			if !strings.HasPrefix(valClean, rootClean) {
				return false
			}
		}
	}
	return true
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

	out, err := t.handler(call.Arguments)
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
