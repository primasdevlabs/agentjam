package context_test

import (
	"os"
	"testing"

	"github.com/agentjam/agentjam/pkg/context"
	"github.com/agentjam/agentjam/pkg/policy"
)

func TestContextManagerModule(t *testing.T) {
	cwd, _ := os.Getwd()
	pe := policy.NewPolicyEngine()
	cm := context.NewContextManager(cwd, pe)

	snapshot := cm.BuildContextSnapshot(context.ContextOptions{
		TokenBudget:  50,
		ActiveAgent:  "backend-engineer",
		ActiveSkills: []string{"testing"},
		CustomRules:  []string{"Rule 1"},
	})

	if snapshot.TokenCountEstimate <= 0 {
		t.Errorf("Expected positive token count estimate")
	}
}
