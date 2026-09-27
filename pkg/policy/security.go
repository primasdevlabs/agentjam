package policy

import (
	"regexp"

	"github.com/primasdevlabs/agentjam/pkg/core"
)

var (
	secretPattern    = regexp.MustCompile(`(?i)(api_key|secret_key|private_key|password|access_token)\s*[:=]\s*['"][A-Za-z0-9_\-\.]{16,}['"]`)
	sqlConcatPattern = regexp.MustCompile(`(?i)(SELECT|INSERT INTO|UPDATE|DELETE FROM)\s+[^;\n]*\+\s*`)
)

// EvaluateSecurityRules evaluates source code for hardcoded secrets and
// unparameterized SQL. When a non-empty policy list is supplied, security rules
// only run if a "security" category policy is loaded.
func EvaluateSecurityRules(policies []core.PolicyManifest, content string) []EvaluationViolation {
	violations := make([]EvaluationViolation, 0)
	if !policyActive(policies, "security") {
		return violations
	}

	if secretPattern.MatchString(content) {
		violations = append(violations, EvaluationViolation{
			PolicyID:    "security-no-secrets",
			PolicyName:  "No Hardcoded Secrets",
			Enforcement: core.EnforceStrictBlock,
			RuleName:    "no-plaintext-credentials",
			Message:     "Source code contains hardcoded secrets or API keys. Use environment variables.",
			Line:        firstMatchLine(secretPattern, content),
		})
	}

	if sqlConcatPattern.MatchString(content) {
		violations = append(violations, EvaluationViolation{
			PolicyID:    "security-sql-parameterization",
			PolicyName:  "Parameterized Database Queries",
			Enforcement: core.EnforceStrictBlock,
			RuleName:    "parameterized-sql-only",
			Message:     "Dynamic SQL string concatenation detected. Parameterize database queries.",
			Line:        firstMatchLine(sqlConcatPattern, content),
		})
	}

	return violations
}
