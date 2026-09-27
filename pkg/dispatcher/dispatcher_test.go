package dispatcher_test

import (
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/primasdevlabs/agentjam/pkg/core"
	"github.com/primasdevlabs/agentjam/pkg/dispatcher"
)

func newDispatcher(t *testing.T, opts dispatcher.ToolDispatcherOptions) (*dispatcher.ToolDispatcher, string) {
	t.Helper()
	root := t.TempDir()
	return dispatcher.NewToolDispatcher(root, opts, nil), root
}

func manifest(name string, level core.ToolSafetyLevel) core.ToolManifest {
	return core.ToolManifest{
		Name: name, Version: "1.0.0", Type: "tool",
		Description: name + " tool", SafetyLevel: level,
	}
}

func TestDispatchSuccess(t *testing.T) {
	td, _ := newDispatcher(t, dispatcher.ToolDispatcherOptions{})
	td.RegisterTool(manifest("add", core.SafetyReadOnly),
		func(args map[string]interface{}) (interface{}, error) {
			return args["a"].(float64) + args["b"].(float64), nil
		})

	res := td.Dispatch(core.ToolCall{ID: "1", ToolName: "add",
		Arguments: map[string]interface{}{"a": 2.0, "b": 3.0}})
	if res.Status != "success" || res.Output != 5.0 {
		t.Errorf("Expected success with 5.0, got %+v", res)
	}
}

func TestDispatchUnknownTool(t *testing.T) {
	td, _ := newDispatcher(t, dispatcher.ToolDispatcherOptions{})
	res := td.Dispatch(core.ToolCall{ID: "x", ToolName: "nope"})
	if res.Status != "error" {
		t.Errorf("Expected error status for unregistered tool, got %s", res.Status)
	}
}

func TestSafetyLevelEnforcement(t *testing.T) {
	td, _ := newDispatcher(t, dispatcher.ToolDispatcherOptions{})
	td.RegisterTool(manifest("rm", core.SafetyDestructive),
		func(args map[string]interface{}) (interface{}, error) { return "deleted", nil })
	td.RegisterTool(manifest("admin", core.SafetyAdmin),
		func(args map[string]interface{}) (interface{}, error) { return "ok", nil })

	if res := td.Dispatch(core.ToolCall{ID: "1", ToolName: "rm"}); res.Status != "blocked" {
		t.Errorf("Destructive tool without permission should be blocked, got %s", res.Status)
	}
	if res := td.Dispatch(core.ToolCall{ID: "2", ToolName: "admin"}); res.Status != "blocked" {
		t.Errorf("Admin tool without permission should be blocked, got %s", res.Status)
	}

	td2, _ := newDispatcher(t, dispatcher.ToolDispatcherOptions{AllowDestructive: true})
	td2.RegisterTool(manifest("rm", core.SafetyDestructive),
		func(args map[string]interface{}) (interface{}, error) { return "deleted", nil })
	if res := td2.Dispatch(core.ToolCall{ID: "3", ToolName: "rm"}); res.Status != "success" {
		t.Errorf("Destructive tool with permission should succeed, got %+v", res)
	}
}

func TestPathBoundaryEnforcement(t *testing.T) {
	td, root := newDispatcher(t, dispatcher.ToolDispatcherOptions{})
	td.RegisterTool(manifest("read", core.SafetyReadOnly),
		func(args map[string]interface{}) (interface{}, error) { return "data", nil })

	inside := filepath.Join(root, "sub", "file.txt")
	if res := td.Dispatch(core.ToolCall{ID: "1", ToolName: "read",
		Arguments: map[string]interface{}{"filePath": inside}}); res.Status != "success" {
		t.Errorf("Path inside workspace should be allowed, got %+v", res)
	}

	if res := td.Dispatch(core.ToolCall{ID: "2", ToolName: "read",
		Arguments: map[string]interface{}{"filePath": filepath.Join(root, "..", "outside.txt")}}); res.Status == "success" {
		t.Error("Path escaping workspace should be blocked")
	}

	// Sibling directory sharing the root's name prefix must be rejected
	sibling := root + "-evil"
	if res := td.Dispatch(core.ToolCall{ID: "3", ToolName: "read",
		Arguments: map[string]interface{}{"filePath": filepath.Join(sibling, "f.txt")}}); res.Status == "success" {
		t.Error("Prefix-sibling path must be blocked")
	}
}

func TestExecutionTimeout(t *testing.T) {
	td, _ := newDispatcher(t, dispatcher.ToolDispatcherOptions{ExecutionTimeoutMs: 50})
	td.RegisterTool(manifest("slow", core.SafetyReadOnly),
		func(args map[string]interface{}) (interface{}, error) {
			time.Sleep(500 * time.Millisecond)
			return "done", nil
		})

	res := td.Dispatch(core.ToolCall{ID: "t", ToolName: "slow"})
	if res.Status != "error" {
		t.Errorf("Expected timeout error, got %+v", res)
	}
}

func TestHandlerErrorPropagates(t *testing.T) {
	td, _ := newDispatcher(t, dispatcher.ToolDispatcherOptions{})
	td.RegisterTool(manifest("fail", core.SafetyReadOnly),
		func(args map[string]interface{}) (interface{}, error) {
			return nil, errors.New("boom")
		})
	res := td.Dispatch(core.ToolCall{ID: "e", ToolName: "fail"})
	if res.Status != "error" || res.Error != "boom" {
		t.Errorf("Expected handler error propagation, got %+v", res)
	}
}

func TestUnregisterAndList(t *testing.T) {
	td, _ := newDispatcher(t, dispatcher.ToolDispatcherOptions{})
	td.RegisterTool(manifest("b-tool", core.SafetyReadOnly), nil)
	td.RegisterTool(manifest("a-tool", core.SafetyReadOnly), nil)

	tools := td.ListTools()
	if len(tools) != 2 || tools[0].Name != "a-tool" {
		t.Errorf("Expected sorted tool list, got %+v", tools)
	}
	if !td.UnregisterTool("a-tool") {
		t.Error("Unregister should return true for registered tool")
	}
	if td.UnregisterTool("a-tool") {
		t.Error("Unregister should return false for missing tool")
	}
}
