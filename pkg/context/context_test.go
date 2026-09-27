package context_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	agentctx "github.com/primasdevlabs/agentjam/pkg/context"
	"github.com/primasdevlabs/agentjam/pkg/policy"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	return filepath.Clean(filepath.Join(cwd, "..", ".."))
}

func TestEstimateTokenCount(t *testing.T) {
	cm := agentctx.NewContextManager(t.TempDir(), nil)
	if cm.EstimateTokenCount("") != 0 {
		t.Error("Empty string should estimate 0 tokens")
	}
	if cm.EstimateTokenCount("abcd") != 1 {
		t.Error("4 chars should estimate 1 token")
	}
	if cm.EstimateTokenCount("abcde") != 2 {
		t.Error("5 chars should ceil to 2 tokens")
	}
}

func TestBuildContextSnapshot(t *testing.T) {
	root := repoRoot(t)
	pe := policy.NewPolicyEngineFromDir(root)
	cm := agentctx.NewContextManager(root, pe)

	snap := cm.BuildContextSnapshot(agentctx.ContextOptions{
		ActiveAgent:  "software-engineer",
		ActiveStack:  "nextjs",
		ActiveSkills: []string{"testing"},
		CustomRules:  []string{"No hardcoded secrets"},
	})

	if snap.WorkspaceRoot != root {
		t.Errorf("WorkspaceRoot mismatch: %q", snap.WorkspaceRoot)
	}
	if !strings.Contains(snap.SystemInstruction, "software-engineer") {
		t.Error("System instruction should mention active agent")
	}
	if !strings.Contains(snap.SystemInstruction, "nextjs") {
		t.Error("System instruction should mention active stack")
	}
	if !strings.Contains(snap.SystemInstruction, "No hardcoded secrets") {
		t.Error("System instruction should include custom rules")
	}
	if snap.PoliciesCount != len(pe.GetPolicies()) {
		t.Error("PoliciesCount should match loaded policy count")
	}
	if snap.TokenCountEstimate <= 0 {
		t.Error("Expected positive token estimate")
	}
}

func TestPolicyInstructionsInjected(t *testing.T) {
	root := repoRoot(t)
	pe := policy.NewPolicyEngineFromDir(root)
	cm := agentctx.NewContextManager(root, pe)

	snap := cm.BuildContextSnapshot(agentctx.ContextOptions{})
	if !strings.Contains(snap.SystemInstruction, "Policy:") {
		t.Error("Expected policy instructions embedded in system instruction")
	}
}

func TestSkillInstructionsInjected(t *testing.T) {
	root := repoRoot(t)
	cm := agentctx.NewContextManager(root, policy.NewPolicyEngine())

	snap := cm.BuildContextSnapshot(agentctx.ContextOptions{
		ActiveSkills: []string{"anti-slop"},
	})
	if !strings.Contains(snap.SystemInstruction, "anti-slop") {
		t.Error("Expected skill name in system instruction")
	}
	if !strings.Contains(snap.SystemInstruction, "Skill `anti-slop`") {
		t.Error("Expected resolved skill instruction block in system instruction")
	}
}

func TestTokenBudgetTruncation(t *testing.T) {
	cm := agentctx.NewContextManager(t.TempDir(), nil)
	snap := cm.BuildContextSnapshot(agentctx.ContextOptions{
		TokenBudget: 10, // 10 tokens ≈ 40 chars
	})
	if !strings.Contains(snap.SystemInstruction, "Context Truncated") {
		t.Error("Expected truncation marker when exceeding token budget")
	}
}
