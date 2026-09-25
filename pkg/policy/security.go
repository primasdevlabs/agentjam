package policy

import (
	"regexp"

	"github.com/primasdevlabs/agentjam/pkg/core"
)

// EvaluateSecurityRules evaluates source code for hardcoded secrets and unparameterized SQL.
func EvaluateSecurityRules(policies []core.PolicyManifest, content string) []EvaluationViolation {
	violations := make([]EvaluationViolation, 0)

	// Secret key detection
	secretRegex := regexp.MustCompile(`(?i)(api_key|secret_key|private_key|password)\s*=\s*['"][A-Za-z0-9_\-]{16,}['"]`)
	if secretRegex.MatchString(content) {
		violations = append(violations, EvaluationViolation{
			PolicyID:    "security-no-secrets",
			PolicyName:  "No Hardcoded Secrets",
			Enforcement: core.EnforceStrictBlock,
			RuleName:    "no-plaintext-credentials",
			Message:     "Source code contains hardcoded secrets or API keys. Use environment variables.",
		})
	}

	// SQL string concatenation check
	sqlConcatRegex := regexp.MustCompile(`(?i)(SELECT|INSERT|UPDATE|DELETE).*\+.*`)
	if sqlConcatRegex.MatchString(content) {
		violations = append(violations, EvaluationViolation{
			PolicyID:    "security-sql-parameterization",
			PolicyName:  "Parameterized Database Queries",
			Enforcement: core.EnforceStrictBlock,
			RuleName:    "parameterized-sql-only",
			Message:     "Dynamic SQL string concatenation detected. Parameterize database queries.",
		})
	}

	return violations
}
