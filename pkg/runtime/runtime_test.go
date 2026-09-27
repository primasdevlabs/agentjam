package runtime_test

import (
	"os"
	"path/filepath"
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

func repoRoot(t *testing.T) string {
	t.Helper()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	return filepath.Clean(filepath.Join(cwd, "..", ".."))
}

func TestRuntimeWiresSubsystems(t *testing.T) {
	rt := runtime.NewRuntime(t.TempDir(), dispatcher.ToolDispatcherOptions{AllowDestructive: true})
	if rt.GetPolicyEngine() == nil || rt.GetContextManager() == nil ||
		rt.GetMemoryManager() == nil || rt.GetToolchainManager() == nil ||
		rt.GetToolDispatcher() == nil {
		t.Fatal("Runtime must wire all subsystems")
	}
}

func TestRuntimeLoadsPolicies(t *testing.T) {
	rt := runtime.NewRuntime(repoRoot(t), dispatcher.ToolDispatcherOptions{})
	if len(rt.GetPolicyEngine().GetPolicies()) == 0 {
		t.Error("Expected runtime to auto-load policies from policies/")
	}
}

func TestFullGoArchitecture(t *testing.T) {
	root := repoRoot(t)
	rt := runtime.NewRuntime(root, dispatcher.ToolDispatcherOptions{AllowDestructive: true})

	// 1. Context Snapshot
	snapshot := rt.GetContextManager().BuildContextSnapshot(context.ContextOptions{
		ActiveAgent: "software-engineer",
		CustomRules: []string{"Strict TypeScript only", "No hardcoded secrets"},
	})
	if snapshot.WorkspaceRoot != root {
		t.Errorf("Expected WorkspaceRoot %s, got %s", root, snapshot.WorkspaceRoot)
	}

	// 2. Memory
	mm := rt.GetMemoryManager()
	mm.Set("goal", "Full Go porting complete", core.MemoryScopeWorking, []string{"test"}, 0)
	entry, found := mm.Get("goal", core.MemoryScopeWorking)
	if !found || entry.Value != "Full Go porting complete" {
		t.Errorf("Memory retrieval failed: %v", entry)
	}

	// 3. Dispatcher
	td := rt.GetToolDispatcher()
	td.RegisterTool(
		core.ToolManifest{
			Name: "add_numbers", Version: "1.0.0", Type: "tool",
			Description: "Adds two numbers", Capabilities: []string{"math"},
			SafetyLevel: core.SafetyReadOnly,
		},
		func(args map[string]interface{}) (interface{}, error) {
			a, _ := args["a"].(float64)
			b, _ := args["b"].(float64)
			return a + b, nil
		},
	)
	result := td.Dispatch(core.ToolCall{
		ID: "call_1", ToolName: "add_numbers",
		Arguments: map[string]interface{}{"a": 10.0, "b": 20.0},
	})
	if result.Status != "success" || result.Output != 30.0 {
		t.Errorf("Tool dispatch failed: %v", result)
	}

	// 4. Design governance
	if evalResult := policy.EvaluateDesignRules(nil, "UI with \U0001F680 emoji and supercharge filler"); len(evalResult) == 0 {
		t.Error("Expected anti-slop design policy violations, got 0")
	}

	// 5. Registry
	index := registry.BuildRegistryIndex(root)
	if index.AgentsCount == 0 {
		t.Error("Expected non-zero agents in registry index")
	}

	// 6. Validator
	if validator.HasErrors(validator.ValidateRepository(root)) {
		t.Error("Repository validation produced errors")
	}

	// 7. Freshness & precedence
	fresh := runtime.CheckFreshness(time.Now().Format(time.RFC3339), "7d")
	if !fresh.IsFresh {
		t.Error("Freshness check failed")
	}
}

func TestCheckFreshnessDurations(t *testing.T) {
	now := time.Now()

	if r := runtime.CheckFreshness(now.Add(-2*time.Hour).Format(time.RFC3339), "1d"); !r.IsFresh {
		t.Error("2h-old doc should be fresh under 1d policy")
	}
	if r := runtime.CheckFreshness(now.Add(-48*time.Hour).Format(time.RFC3339), "1d"); r.IsFresh {
		t.Error("48h-old doc should be stale under 1d policy")
	}
	if r := runtime.CheckFreshness(now.Add(-20*24*time.Hour).Format(time.RFC3339), "30d"); !r.IsFresh {
		t.Error("20d-old doc should be fresh under 30d policy")
	}
	if r := runtime.CheckFreshness(now.Add(-10*24*time.Hour).Format(time.RFC3339), "1w"); r.IsFresh {
		t.Error("10d-old doc should be stale under 1w policy")
	}
	if r := runtime.CheckFreshness("not-a-date", "7d"); !r.IsFresh {
		t.Error("Unparseable timestamp should default to fresh")
	}
}

func TestResolveHighest(t *testing.T) {
	best := runtime.ResolveHighest([]runtime.PrecedenceCandidate{
		{Level: "Model Knowledge", SourceName: "llm", Value: "v1"},
		{Level: "Current Codebase State", SourceName: "code", Value: "v2"},
		{Level: "Pinned Rules", SourceName: "pins", Value: "v3"},
	})
	if best.SourceName != "code" {
		t.Errorf("Expected codebase state to win, got %+v", best)
	}
	if empty := runtime.ResolveHighest(nil); empty != (runtime.PrecedenceCandidate{}) {
		t.Error("Empty candidates should return zero value")
	}
}

func TestDetectEnvironments(t *testing.T) {
	dir := t.TempDir()
	for _, f := range []string{".cursorrules", "CLAUDE.md", ".windsurfrules"} {
		if err := os.WriteFile(filepath.Join(dir, f), []byte("x"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	envs := runtime.DetectEnvironments(dir)
	names := map[string]bool{}
	for _, e := range envs {
		names[e.Name] = true
		if len(e.DetectedFiles) == 0 {
			t.Errorf("Environment %s missing detected files", e.Name)
		}
	}
	for _, want := range []string{"Cursor IDE", "Claude Code CLI", "Windsurf IDE"} {
		if !names[want] {
			t.Errorf("Expected environment %q to be detected, got %v", want, names)
		}
	}

	if envs := runtime.DetectEnvironments(t.TempDir()); len(envs) != 0 {
		t.Error("Empty dir should detect no environments")
	}
}
