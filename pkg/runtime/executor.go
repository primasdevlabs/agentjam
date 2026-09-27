package runtime

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/primasdevlabs/agentjam/pkg/context"
	"github.com/primasdevlabs/agentjam/pkg/core"
	"github.com/primasdevlabs/agentjam/pkg/parser"
)

// StepResult captures the outcome of a single workflow step.
type StepResult struct {
	StepID      string      `json:"stepId"`
	Name        string      `json:"name"`
	Agent       string      `json:"agent,omitempty"`
	Skill       string      `json:"skill,omitempty"`
	Status      string      `json:"status"` // prepared | success | skipped | error
	Instruction string      `json:"instruction,omitempty"`
	Output      interface{} `json:"output,omitempty"`
	Error       string      `json:"error,omitempty"`
	DurationMs  int64       `json:"durationMs"`
}

// WorkflowRun is the execution report for one workflow invocation.
type WorkflowRun struct {
	Workflow   string       `json:"workflow"`
	Status     string       `json:"status"` // success | failed | partial
	Steps      []StepResult `json:"steps"`
	StartedAt  string       `json:"startedAt"`
	FinishedAt string       `json:"finishedAt"`
	DurationMs int64        `json:"durationMs"`
}

// StepContext carries the resolved per-step environment handed to a runner.
type StepContext struct {
	Runtime   *AgentJamRuntime
	Snapshot  core.ContextSnapshot
	StepIndex int
	Inputs    map[string]interface{}
}

// StepRunner executes a resolved step. When nil on the Executor, steps are
// "prepared" — their agent instruction context is assembled but not executed.
type StepRunner func(step core.WorkflowStep, ctx StepContext) (interface{}, error)

// Executor runs workflow manifests against the runtime.
type Executor struct {
	runtime *AgentJamRuntime
	runner  StepRunner
}

// NewExecutor creates a workflow executor. runner may be nil, in which case
// every step is resolved to its assembled agent instruction and marked
// "prepared" without execution.
func NewExecutor(rt *AgentJamRuntime, runner StepRunner) *Executor {
	return &Executor{runtime: rt, runner: runner}
}

// LoadWorkflow locates a workflow manifest by name or directory basename.
func LoadWorkflow(workflowDir, name string) (*core.WorkflowManifest, error) {
	var found *core.WorkflowManifest
	err := filepath.Walk(workflowDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || info.Name() != "workflow.yaml" {
			return err
		}
		bundle, perr := parser.ParseWorkflow(path)
		if perr != nil {
			return nil
		}
		m := bundle.Manifest
		if m.Name == name || filepath.Base(filepath.Dir(path)) == name {
			cp := m
			found = &cp
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if found == nil {
		return nil, fmt.Errorf("workflow %q not found under %s", name, workflowDir)
	}
	return found, nil
}

// ListWorkflows returns sorted workflow names discoverable in the workspace.
func ListWorkflows(workflowDir string) ([]string, error) {
	seen := map[string]bool{}
	err := filepath.Walk(workflowDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || info.Name() != "workflow.yaml" {
			return err
		}
		if bundle, perr := parser.ParseWorkflow(path); perr == nil {
			name := bundle.Manifest.Name
			if name == "" {
				name = filepath.Base(filepath.Dir(path))
			}
			seen[name] = true
		}
		return nil
	})
	out := make([]string, 0, len(seen))
	for n := range seen {
		out = append(out, n)
	}
	sort.Strings(out)
	return out, err
}

// Run executes the named workflow. Steps run in manifest order; each step
// resolves its agent, builds a context snapshot, and is dispatched to the
// configured StepRunner (or prepared when no runner is set). Step outcomes
// are recorded in episodic memory.
func (e *Executor) Run(workflowName string) (*WorkflowRun, error) {
	root := e.runtime.RootPath()
	workflowDir := filepath.Join(root, "workflows")

	wf, err := LoadWorkflow(workflowDir, workflowName)
	if err != nil {
		return nil, err
	}
	if len(wf.Steps) == 0 {
		return nil, fmt.Errorf("workflow %q declares no steps", workflowName)
	}

	wfName := wf.Name
	if wfName == "" {
		wfName = workflowName
	}
	run := &WorkflowRun{
		Workflow:  wfName,
		StartedAt: time.Now().UTC().Format(time.RFC3339),
	}
	start := time.Now()

	agents := e.agentIDs(root)

	for i, step := range wf.Steps {
		res := e.runStep(i, step, agents)
		run.Steps = append(run.Steps, res)
		e.recordStep(wfName, res)

		if res.Status == "error" && !onFailureContinues(step) {
			for _, rest := range wf.Steps[i+1:] {
				run.Steps = append(run.Steps, StepResult{
					StepID: rest.ID, Name: rest.Name, Agent: rest.Agent,
					Status: "skipped", Error: "skipped after failure of " + step.ID,
				})
			}
			break
		}
	}

	run.FinishedAt = time.Now().UTC().Format(time.RFC3339)
	run.DurationMs = time.Since(start).Milliseconds()

	errCount, skipped := 0, 0
	for _, s := range run.Steps {
		if s.Status == "error" {
			errCount++
		}
		if s.Status == "skipped" {
			skipped++
		}
	}
	switch {
	case errCount == 0 && skipped == 0:
		run.Status = "success"
	case errCount == len(run.Steps)-skipped:
		run.Status = "failed"
	default:
		run.Status = "partial"
	}
	return run, nil
}

func onFailureContinues(step core.WorkflowStep) bool {
	return step.OnFailure == "continue" || step.OnFailure == "skip"
}

// agentIDs indexes discoverable agent resource IDs/names for step resolution.
func (e *Executor) agentIDs(root string) map[string]bool {
	ids := map[string]bool{}
	for _, res := range parser.DiscoverResources(root) {
		if res.Type == core.ResourceTypeAgent {
			ids[res.ID] = true
		}
	}
	return ids
}

func (e *Executor) runStep(idx int, step core.WorkflowStep, agents map[string]bool) StepResult {
	start := time.Now()
	res := StepResult{
		StepID: step.ID,
		Name:   step.Name,
		Agent:  step.Agent,
		Skill:  step.Skill,
		Status: "prepared",
	}
	fail := func(msg string) StepResult {
		res.Status = "error"
		res.Error = msg
		res.DurationMs = time.Since(start).Milliseconds()
		return res
	}

	if step.Agent != "" && !agents[step.Agent] {
		return fail(fmt.Sprintf("unknown agent %q", step.Agent))
	}

	var skills []string
	if step.Skill != "" {
		skills = []string{step.Skill}
	}
	snap := e.runtime.GetContextManager().BuildContextSnapshot(context.ContextOptions{
		ActiveAgent:  step.Agent,
		ActiveSkills: skills,
	})
	res.Instruction = snap.SystemInstruction

	if e.runner != nil {
		out, rerr := e.runner(step, StepContext{
			Runtime:   e.runtime,
			Snapshot:  snap,
			StepIndex: idx,
			Inputs:    step.Inputs,
		})
		if rerr != nil {
			res.Status = "error"
			res.Error = rerr.Error()
		} else {
			res.Status = "success"
			res.Output = out
		}
	}

	res.DurationMs = time.Since(start).Milliseconds()
	return res
}

func (e *Executor) recordStep(workflowName string, res StepResult) {
	key := fmt.Sprintf("workflow:%s:%s:%d", workflowName, res.StepID, time.Now().UnixNano())
	e.runtime.GetMemoryManager().Set(key, res, core.MemoryScopeEpisodic,
		[]string{"workflow-run", workflowName, res.Status}, 0)
}
