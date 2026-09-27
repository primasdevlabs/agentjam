package mcp

import (
	"fmt"
	"os"
	"path/filepath"
	goruntime "runtime"
	"sort"
	"strings"
	"time"

	"github.com/primasdevlabs/agentjam/pkg/core"
	"github.com/primasdevlabs/agentjam/pkg/dispatcher"
	"github.com/primasdevlabs/agentjam/pkg/parser"
	agentruntime "github.com/primasdevlabs/agentjam/pkg/runtime"
)

// RegisterBuiltinTools wires workspace tool manifests to concrete Go
// implementations on the runtime's dispatcher: filesystem read/write/list,
// git inspection, and a run_workflow meta-tool backed by the Executor.
func RegisterBuiltinTools(rt *agentruntime.AgentJamRuntime) {
	td := rt.GetToolDispatcher()
	root := rt.RootPath()

	td.RegisterTool(loadToolManifest(root, "filesystem", "Read, write, and list files in the workspace.",
		core.SafetySafeWrite), filesystemHandler(root))
	td.RegisterTool(loadToolManifest(root, "git", "Inspect git status, diffs, branches, and history.",
		core.SafetyReadOnly), gitHandler(rt))
	td.RegisterTool(loadToolManifest(root, "terminal", "Execute shell commands in the workspace.",
		core.SafetyDestructive), terminalHandler(rt))
	td.RegisterTool(loadToolManifest(root, "github", "Query GitHub repositories via the gh CLI.",
		core.SafetyReadOnly), githubHandler(rt))
	td.RegisterTool(core.ToolManifest{
		Name:        "run_workflow",
		Version:     "1.0.0",
		Type:        "tool",
		Description: "Run an AgentJam workflow by name and return its execution report.",
		Category:    "orchestration",
		SafetyLevel: core.SafetyReadOnly,
		Parameters: map[string]interface{}{
			"properties": map[string]interface{}{
				"workflow": map[string]interface{}{"type": "string", "description": "Workflow name to run."},
			},
			"required": []string{"workflow"},
		},
	}, workflowHandler(rt))
}

// loadToolManifest finds a tool manifest by name in the workspace tree,
// falling back to a minimal synthetic manifest when absent.
func loadToolManifest(root, name, fallbackDesc string, safety core.ToolSafetyLevel) core.ToolManifest {
	for _, res := range parser.DiscoverResources(root) {
		if res.Type != core.ResourceTypeTool || res.ID != name {
			continue
		}
		if m, err := parser.ParseTool(res.Path); err == nil {
			return m
		}
	}
	return core.ToolManifest{
		Name:        name,
		Version:     "1.0.0",
		Type:        "tool",
		Description: fallbackDesc,
		SafetyLevel: safety,
	}
}

func resolvePath(root, p string) (string, error) {
	if p == "" {
		return "", fmt.Errorf("path argument is required")
	}
	abs := p
	if !filepath.IsAbs(abs) {
		abs = filepath.Join(root, p)
	}
	abs = filepath.Clean(abs)
	rootClean := filepath.Clean(root)
	cmpAbs, cmpRoot := abs, rootClean
	if goruntime.GOOS == "windows" {
		cmpAbs, cmpRoot = strings.ToLower(cmpAbs), strings.ToLower(cmpRoot)
	}
	if cmpAbs != cmpRoot && !strings.HasPrefix(cmpAbs, cmpRoot+string(os.PathSeparator)) {
		return "", fmt.Errorf("path %q escapes workspace root", p)
	}
	return abs, nil
}

func filesystemHandler(root string) dispatcher.ToolHandler {
	return func(args map[string]interface{}) (interface{}, error) {
		op, _ := args["operation"].(string)
		if op == "" {
			op = "read"
		}
		path, err := resolvePath(root, strArg(args, "path"))
		if err != nil {
			return nil, err
		}

		switch op {
		case "read":
			data, err := os.ReadFile(path)
			if err != nil {
				return nil, err
			}
			return map[string]interface{}{"path": path, "content": string(data)}, nil

		case "write":
			content, _ := args["content"].(string)
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				return nil, err
			}
			if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
				return nil, err
			}
			return map[string]interface{}{"path": path, "bytes": len(content)}, nil

		case "list":
			entries, err := os.ReadDir(path)
			if err != nil {
				return nil, err
			}
			names := make([]string, 0, len(entries))
			for _, e := range entries {
				name := e.Name()
				if e.IsDir() {
					name += string(os.PathSeparator)
				}
				names = append(names, name)
			}
			sort.Strings(names)
			return map[string]interface{}{"path": path, "entries": names}, nil

		default:
			return nil, fmt.Errorf("unsupported filesystem operation %q (read|write|list)", op)
		}
	}
}

func gitHandler(rt *agentruntime.AgentJamRuntime) dispatcher.ToolHandler {
	return func(args map[string]interface{}) (interface{}, error) {
		op, _ := args["operation"].(string)
		var cmd string
		switch op {
		case "status":
			cmd = "git status --short --branch"
		case "diff":
			cmd = "git --no-pager diff"
		case "log":
			cmd = "git --no-pager log --oneline -20"
		case "branch":
			cmd = "git branch --all"
		default:
			return nil, fmt.Errorf("unsupported git operation %q (status|diff|log|branch)", op)
		}
		ok, output, ms := rt.GetToolchainManager().RunCommand(cmd, 30*time.Second)
		if !ok {
			return nil, fmt.Errorf("git %s failed: %s", op, strings.TrimSpace(output))
		}
		return map[string]interface{}{"operation": op, "output": strings.TrimSpace(output), "durationMs": ms}, nil
	}
}

// dangerousTerminalPatterns is a minimal denylist applied on top of the
// dispatcher's destructive-safety gate for shell commands.
var dangerousTerminalPatterns = []string{
	"rm -rf /", "rm -rf ~", "mkfs", "dd if=", ":(){ :|:& };:",
	"shutdown", "reboot", "format ", "del /f /s /q c:\\",
}

func terminalHandler(rt *agentruntime.AgentJamRuntime) dispatcher.ToolHandler {
	return func(args map[string]interface{}) (interface{}, error) {
		cmd := strArg(args, "command")
		if cmd == "" {
			return nil, fmt.Errorf("command argument is required")
		}
		lower := strings.ToLower(cmd)
		for _, pat := range dangerousTerminalPatterns {
			if strings.Contains(lower, pat) {
				return nil, fmt.Errorf("command matched denylist pattern %q", pat)
			}
		}
		timeoutMs := time.Duration(intArg(args, "timeoutMs")) * time.Millisecond
		ok, output, ms := rt.GetToolchainManager().RunCommand(cmd, timeoutMs)
		return map[string]interface{}{
			"command":    cmd,
			"exitOk":     ok,
			"output":     strings.TrimSpace(output),
			"durationMs": ms,
		}, nil
	}
}

// githubHandler proxies whitelisted read operations to the gh CLI.
func githubHandler(rt *agentruntime.AgentJamRuntime) dispatcher.ToolHandler {
	return func(args map[string]interface{}) (interface{}, error) {
		if !rt.GetToolchainManager().HasBinary("gh") {
			return nil, fmt.Errorf("github tool requires the gh CLI in PATH")
		}
		op, _ := args["operation"].(string)
		var cmd string
		switch op {
		case "pr-list":
			cmd = "gh pr list --limit 20"
		case "issue-list":
			cmd = "gh issue list --limit 20"
		case "repo-view":
			cmd = "gh repo view --json name,description,url,defaultBranchRef"
		case "pr-view":
			num := strArg(args, "number")
			if num == "" {
				return nil, fmt.Errorf("number argument is required for pr-view")
			}
			cmd = "gh pr view " + num + " --json title,state,author,url"
		default:
			return nil, fmt.Errorf("unsupported github operation %q (pr-list|issue-list|repo-view|pr-view)", op)
		}
		ok, output, ms := rt.GetToolchainManager().RunCommand(cmd, 30*time.Second)
		if !ok {
			return nil, fmt.Errorf("gh %s failed: %s", op, strings.TrimSpace(output))
		}
		return map[string]interface{}{"operation": op, "output": strings.TrimSpace(output), "durationMs": ms}, nil
	}
}

func intArg(args map[string]interface{}, key string) int {
	switch v := args[key].(type) {
	case float64:
		return int(v)
	case int:
		return v
	}
	return 0
}

func workflowHandler(rt *agentruntime.AgentJamRuntime) dispatcher.ToolHandler {
	return func(args map[string]interface{}) (interface{}, error) {
		name := strArg(args, "workflow")
		if name == "" {
			return nil, fmt.Errorf("workflow argument is required")
		}
		exec := agentruntime.NewExecutor(rt, nil)
		return exec.Run(name)
	}
}

func strArg(args map[string]interface{}, key string) string {
	if v, ok := args[key].(string); ok {
		return v
	}
	return ""
}
