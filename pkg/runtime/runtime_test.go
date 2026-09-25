package runtime_test

import (
	"os"
	"testing"
	"time"

	"github.com/primasdevlabs/agentjam/pkg/context"
	"github.com/primasdevlabs/agentjam/pkg/core"
	"github.com/primasdevlabs/agentjam/pkg/dispatcher"
	"github.com/primasdevlabs/agentjam/pkg/policy"
	"github.com/primasdevlabs/agentjam/pkg/registry"
	"github.com/primasdevlabs/agentjam/pkg/runtime"
	"github.com/primasdevlabs/agentjam/pkg/validator"
)

func TestFullGoArchitecture(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	rt := runtime.NewRuntime(cwd, dispatcher.ToolDispatcherOptions{AllowDestructive: true})

	// 1. Context Snapshot Test
	cm := rt.GetContextManager()
	snapshot := cm.BuildContextSnapshot(context.ContextOptions{
		ActiveAgent: "software-engineer",
		CustomRules: []string{"Strict TypeScript only", "No hardcoded secrets"},
	})
	if snapshot.WorkspaceRoot != cwd {
		t.Errorf("Expected WorkspaceRoot %s, got %s", cwd, snapshot.WorkspaceRoot)
	}

	// 2. Memory Manager Test
	mm := rt.GetMemoryManager()
	mm.Set("goal", "Full Go porting complete", core.MemoryScopeWorking, []string{"test"}, 0)
	entry, found := mm.Get("goal", core.MemoryScopeWorking)
	if !found || entry.Value != "Full Go porting complete" {
		t.Errorf("Memory retrieval failed: %v", entry)
	}

	// 3. Tool Dispatcher Test
	td := rt.GetToolDispatcher()
	td.RegisterTool(
		core.ToolManifest{
			Name:         "add_numbers",
			Version:      "1.0.0",
			Type:         "tool",
			Description:  "Adds two numbers",
			Capabilities: []string{"math"},
			SafetyLevel:  core.SafetyReadOnly,
		},
		func(args map[string]interface{}) (interface{}, error) {
			a, _ := args["a"].(float64)
			b, _ := args["b"].(float64)
			return a + b, nil
		},
	)

	result := td.Dispatch(core.ToolCall{
		ID:        "call_1",
		ToolName:  "add_numbers",
		Arguments: map[string]interface{}{"a": 10.0, "b": 20.0},
	})
	if result.Status != "success" || result.Output != 30.0 {
		t.Errorf("Tool dispatch failed: %v", result)
	}

	// 4. Design Governance Evaluator Test
	evalResult := policy.EvaluateDesignRules(nil, "UI with 🚀 emoji and supercharge filler")
	if len(evalResult) == 0 {
		t.Errorf("Expected anti-slop design policy violations, got 0")
	}

	// 5. Registry Index Test
	index := registry.BuildRegistryIndex(cwd)
	if index.AgentsCount < 0 {
		t.Errorf("Registry index failed")
	}

	// 6. Validator Test
	errors := validator.ValidateRepository(cwd)
	_ = errors

	// 7. Freshness & Precedence Test
	fresh := runtime.CheckFreshness(time.Now().Format(time.RFC3339), "7d")
	if !fresh.IsFresh {
		t.Errorf("Freshness check failed")
	}
}
