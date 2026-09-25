package policy

import "github.com/primasdevlabs/agentjam/pkg/core"

// PolicyEngineSummary summarizes evaluation results.
type PolicyEngineSummary struct {
	Allowed         bool                   `json:"allowed"`
	TotalViolations int                    `json:"totalViolations"`
	StrictBlocks    int                    `json:"strictBlocks"`
	Warnings        int                    `json:"warnings"`
	InfoCount       int                    `json:"infoCount"`
	Violations      []EvaluationViolation `json:"violations"`
}

// EvaluationViolation represents a single policy violation.
type EvaluationViolation struct {
	PolicyID    string                 `json:"policyId"`
	PolicyName  string                 `json:"policyName"`
	Enforcement core.PolicyEnforcement `json:"enforcement"`
	RuleName    string                 `json:"ruleName"`
	Message     string                 `json:"message"`
	Line        int                    `json:"line,omitempty"`
}

// PolicyEngine manages and evaluates policy rules.
type PolicyEngine struct {
	policies []core.PolicyManifest
}

// NewPolicyEngine instantiates a PolicyEngine.
func NewPolicyEngine() *PolicyEngine {
	return &PolicyEngine{
		policies: make([]core.PolicyManifest, 0),
	}
}

// AddPolicy registers a policy manifest.
func (pe *PolicyEngine) AddPolicy(p core.PolicyManifest) {
	pe.policies = append(pe.policies, p)
}

// GetPolicies returns all loaded policies.
func (pe *PolicyEngine) GetPolicies() []core.PolicyManifest {
	return pe.policies
}

// Summarize aggregates raw violations.
func (pe *PolicyEngine) Summarize(violations []EvaluationViolation) PolicyEngineSummary {
	strictBlocks := 0
	warnings := 0
	infoCount := 0

	for _, v := range violations {
		switch v.Enforcement {
		case core.EnforceStrictBlock:
			strictBlocks++
		case core.EnforceWarning:
			warnings++
		case core.EnforceInfo:
			infoCount++
		}
	}

	return PolicyEngineSummary{
		Allowed:         strictBlocks == 0,
		TotalViolations: len(violations),
		StrictBlocks:    strictBlocks,
		Warnings:        warnings,
		InfoCount:       infoCount,
		Violations:      violations,
	}
}
