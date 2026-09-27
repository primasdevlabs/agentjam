package runtime

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/primasdevlabs/agentjam/pkg/core"
	"github.com/primasdevlabs/agentjam/pkg/dispatcher"
)

func setupWorkspace(t *testing.T) string {
	t.Helper()
	root := t.TempDir()

	write := func(rel, content string) {
		p := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	write("agents/tester/agent.yaml", `name: tester
type: agent
version: 1.0.0
description: Test agent
skills: []
tools: []
`)
	write("workflows/simple/workflow.yaml", `name: simple
type: workflow
version: 1.0.0
description: Simple test workflow
steps:
  - id: step-1
    name: First
    agent: tester
  - id: step-2
    name: Second
    agent: tester
`)
	return root
}

func TestExecutorRunPrepared(t *testing.T) {
	root := setupWorkspace(t)
	rt := NewRuntime(root, dispatcher.ToolDispatcherOptions{})
	exec := NewExecutor(rt, nil)

	run, err := exec.Run("simple")
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if run.Status != "success" {
		t.Fatalf("expected success, got %s", run.Status)
	}
	if len(run.Steps) != 2 {
		t.Fatalf("expected 2 steps, got %d", len(run.Steps))
	}
	for _, s := range run.Steps {
		if s.Status != "prepared" {
			t.Fatalf("step %s should be prepared, got %s", s.StepID, s.Status)
		}
		if s.Instruction == "" {
			t.Fatalf("step %s missing assembled instruction", s.StepID)
		}
	}
}

func TestExecutorRunWithRunner(t *testing.T) {
	root := setupWorkspace(t)
	rt := NewRuntime(root, dispatcher.ToolDispatcherOptions{})
	exec := NewExecutor(rt, func(step core.WorkflowStep, ctx StepContext) (interface{}, error) {
		return map[string]string{"echo": step.Name}, nil
	})

	run, err := exec.Run("simple")
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if run.Steps[0].Status != "success" {
		t.Fatalf("expected success, got %s", run.Steps[0].Status)
	}
}

func TestExecutorAbortsOnFailure(t *testing.T) {
	root := setupWorkspace(t)
	rt := NewRuntime(root, dispatcher.ToolDispatcherOptions{})
	exec := NewExecutor(rt, func(step core.WorkflowStep, ctx StepContext) (interface{}, error) {
		if step.ID == "step-1" {
			return nil, errors.New("boom")
		}
		return nil, nil
	})

	run, err := exec.Run("simple")
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if run.Status != "failed" {
		t.Fatalf("expected failed, got %s", run.Status)
	}
	if run.Steps[1].Status != "skipped" {
		t.Fatalf("step-2 should be skipped, got %s", run.Steps[1].Status)
	}
}

func TestExecutorContinueOnFailure(t *testing.T) {
	root := setupWorkspace(t)
	rt := NewRuntime(root, dispatcher.ToolDispatcherOptions{})

	// Patch the workflow so step-1 continues on failure.
	wfPath := filepath.Join(root, "workflows", "cont", "workflow.yaml")
	if err := os.MkdirAll(filepath.Dir(wfPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(wfPath, []byte(`name: cont
type: workflow
version: 1.0.0
description: Continue-on-failure workflow
steps:
  - id: step-1
    name: First
    agent: tester
    onFailure: continue
  - id: step-2
    name: Second
    agent: tester
`), 0o644); err != nil {
		t.Fatal(err)
	}

	exec := NewExecutor(rt, func(step core.WorkflowStep, ctx StepContext) (interface{}, error) {
		if step.ID == "step-1" {
			return nil, errors.New("boom")
		}
		return "ok", nil
	})

	run, err := exec.Run("cont")
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if run.Status != "partial" {
		t.Fatalf("expected partial, got %s", run.Status)
	}
	if run.Steps[1].Status != "success" {
		t.Fatalf("step-2 should have run, got %s", run.Steps[1].Status)
	}
}

func TestExecutorUnknownAgent(t *testing.T) {
	root := setupWorkspace(t)
	rt := NewRuntime(root, dispatcher.ToolDispatcherOptions{})
	exec := NewExecutor(rt, nil)

	wfPath := filepath.Join(root, "workflows", "bad", "workflow.yaml")
	if err := os.MkdirAll(filepath.Dir(wfPath), 0o755); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(wfPath, []byte(`name: bad
type: workflow
steps:
  - id: step-1
    name: First
    agent: ghost
`), 0o644)

	run, err := exec.Run("bad")
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if run.Steps[0].Status != "error" {
		t.Fatalf("expected error status, got %s", run.Steps[0].Status)
	}
}

func TestExecutorNotFound(t *testing.T) {
	root := setupWorkspace(t)
	rt := NewRuntime(root, dispatcher.ToolDispatcherOptions{})
	exec := NewExecutor(rt, nil)
	if _, err := exec.Run("nonexistent"); err == nil {
		t.Fatal("expected error for unknown workflow")
	}
}

func TestListWorkflows(t *testing.T) {
	root := setupWorkspace(t)
	names, err := ListWorkflows(filepath.Join(root, "workflows"))
	if err != nil {
		t.Fatal(err)
	}
	if len(names) != 1 || names[0] != "simple" {
		t.Fatalf("unexpected workflows: %v", names)
	}
}

func TestExecutorRecordsToMemory(t *testing.T) {
	root := setupWorkspace(t)
	rt := NewRuntime(root, dispatcher.ToolDispatcherOptions{})
	exec := NewExecutor(rt, nil)
	if _, err := exec.Run("simple"); err != nil {
		t.Fatal(err)
	}
	if got := rt.GetMemoryManager().Count(core.MemoryScopeEpisodic); got < 2 {
		t.Fatalf("expected >=2 episodic entries, got %d", got)
	}
}
