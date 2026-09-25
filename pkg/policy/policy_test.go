package policy_test

import (
	"testing"

	"github.com/agentjam/agentjam/pkg/core"
	"github.com/agentjam/agentjam/pkg/policy"
)

func TestPolicyEngineAndEvaluators(t *testing.T) {
	pe := policy.NewPolicyEngine()
	pe.AddPolicy(core.PolicyManifest{
		ID:          "sec-1",
		Name:        "Security",
		Enforcement: core.EnforceStrictBlock,
	})

	if len(pe.GetPolicies()) != 1 {
		t.Errorf("Expected 1 policy, got %d", len(pe.GetPolicies()))
	}

	// Test Design Evaluator (Emoji & Slop Gradient Detection)
	violations := policy.EvaluateDesignRules(pe.GetPolicies(), "Button with 🚀 and from-purple-600 to-indigo-600 gradient")
	if len(violations) < 2 {
		t.Errorf("Expected at least 2 design violations, got %d", len(violations))
	}

	// Test Security Evaluator (Hardcoded Secrets & Unparameterized SQL)
	secViolations := policy.EvaluateSecurityRules(pe.GetPolicies(), "api_key = 'abcdef1234567890123'\nSELECT * FROM users + id")
	if len(secViolations) < 2 {
		t.Errorf("Expected 2 security violations, got %d", len(secViolations))
	}
}
