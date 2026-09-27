package policy_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/primasdevlabs/agentjam/pkg/core"
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

func TestPolicyEngineLoadFromRepo(t *testing.T) {
	pe := policy.NewPolicyEngineFromDir(repoRoot(t))
	policies := pe.GetPolicies()
	if len(policies) == 0 {
		t.Fatal("Expected policies to be loaded from policies/ directory")
	}
	if !pe.HasPolicy("04-security") {
		t.Error("Expected policy '04-security' to be loaded")
	}
}

func TestSummarize(t *testing.T) {
	pe := policy.NewPolicyEngine()
	summary := pe.Summarize([]policy.EvaluationViolation{
		{PolicyID: "a", Enforcement: core.EnforceStrictBlock},
		{PolicyID: "b", Enforcement: core.EnforceWarning},
		{PolicyID: "c", Enforcement: core.EnforceInfo},
	})
	if summary.Allowed {
		t.Error("Expected Allowed=false when strict-block violations exist")
	}
	if summary.TotalViolations != 3 || summary.StrictBlocks != 1 || summary.Warnings != 1 || summary.InfoCount != 1 {
		t.Errorf("Unexpected summary counts: %+v", summary)
	}

	ok := pe.Summarize([]policy.EvaluationViolation{{Enforcement: core.EnforceWarning}})
	if !ok.Allowed {
		t.Error("Expected Allowed=true with only warnings")
	}
}

func TestEvaluateDesignRules(t *testing.T) {
	cases := []struct {
		name        string
		content     string
		wantRule    string
		wantViolate bool
	}{
		{"emoji", "<button>\xF0\x9F\x9A\x80 Go</button>", "no-emojis-in-ui", true},
		{"gradient", `<div class="from-purple-600 to-indigo-600">`, "no-purple-indigo-gradients", true},
		{"buzzword", "We seamlessly empower teams", "no-filler-marketing-copy", true},
		{"clean", "<button>Save changes</button>", "", false},
	}
	for _, tc := range cases {
		violations := policy.EvaluateDesignRules(nil, tc.content)
		if tc.wantViolate {
			found := false
			for _, v := range violations {
				if v.RuleName == tc.wantRule {
					found = true
					if v.Line < 1 {
						t.Errorf("%s: expected line number >= 1, got %d", tc.name, v.Line)
					}
				}
			}
			if !found {
				t.Errorf("%s: expected violation %q, got %+v", tc.name, tc.wantRule, violations)
			}
		} else if len(violations) != 0 {
			t.Errorf("%s: expected no violations, got %+v", tc.name, violations)
		}
	}
}

func TestEvaluateSecurityRules(t *testing.T) {
	violations := policy.EvaluateSecurityRules(nil, `const api_key = "abcdef1234567890abcd";`)
	found := false
	for _, v := range violations {
		if v.RuleName == "no-plaintext-credentials" {
			found = true
		}
	}
	if !found {
		t.Error("Expected hardcoded secret violation")
	}

	sqlViolations := policy.EvaluateSecurityRules(nil, `db.query("SELECT * FROM users WHERE id = " + userId)`)
	found = false
	for _, v := range sqlViolations {
		if v.RuleName == "parameterized-sql-only" {
			found = true
		}
	}
	if !found {
		t.Error("Expected SQL concatenation violation")
	}

	clean := policy.EvaluateSecurityRules(nil, `db.query("SELECT * FROM users WHERE id = ?", [userId])`)
	for _, v := range clean {
		if v.RuleName == "parameterized-sql-only" {
			t.Error("False positive on parameterized query")
		}
	}
}

func TestEvaluateAllWithLoadedPolicies(t *testing.T) {
	pe := policy.NewPolicyEngineFromDir(repoRoot(t))
	summary := pe.EvaluateAll("<div class=\"from-purple-500 to-indigo-700\">empower</div>")
	if summary.TotalViolations == 0 {
		t.Error("Expected EvaluateAll to flag violations with loaded repo policies")
	}
}
