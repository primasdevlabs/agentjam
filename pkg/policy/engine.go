package policy

import (
	"path/filepath"
	"strings"

	"github.com/primasdevlabs/agentjam/pkg/core"
	"github.com/primasdevlabs/agentjam/pkg/parser"
)

// PolicyEngineSummary summarizes evaluation results.
type PolicyEngineSummary struct {
	Allowed         bool                  `json:"allowed"`
	TotalViolations int                   `json:"totalViolations"`
	StrictBlocks    int                   `json:"strictBlocks"`
	Warnings        int                   `json:"warnings"`
	InfoCount       int                   `json:"infoCount"`
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

// HasPolicy reports whether a policy with the given ID is loaded.
func (pe *PolicyEngine) HasPolicy(id string) bool {
	for _, p := range pe.policies {
		if p.ID == id || p.Name == id {
			return true
		}
	}
	return false
}

// LoadPolicies discovers policy manifests under rootDir/policies and registers
// them in the engine. Returns the number of policies loaded.
func (pe *PolicyEngine) LoadPolicies(rootDir string) int {
	discovered := parser.DiscoverResources(rootDir)
	loaded := 0
	for _, res := range discovered {
		if res.Type != core.ResourceTypePolicy {
			continue
		}
		p, err := parser.ParsePolicy(res.Path)
		if err != nil {
			continue
		}
		// Carry the policy category from its directory when unset.
		if p.Category == "" {
			p.Category = filepath.Base(filepath.Dir(res.Path))
		}
		pe.AddPolicy(p)
		loaded++
	}
	return loaded
}

// NewPolicyEngineFromDir creates a PolicyEngine pre-loaded with all policies
// discovered under rootDir.
func NewPolicyEngineFromDir(rootDir string) *PolicyEngine {
	pe := NewPolicyEngine()
	pe.LoadPolicies(rootDir)
	return pe
}

// policyActive reports whether an evaluation rule should run. When the caller
// supplies a non-empty policy list, a rule only runs if its policy ID is
// present in that list; an empty list runs all built-in rules.
func policyActive(policies []core.PolicyManifest, policyID string) bool {
	if len(policies) == 0 {
		return true
	}
	for _, p := range policies {
		if p.ID == policyID || p.Name == policyID ||
			strings.EqualFold(p.Category, policyID) {
			return true
		}
	}
	return false
}

// EvaluateAll runs every built-in evaluator against content and returns the
// aggregated summary.
func (pe *PolicyEngine) EvaluateAll(content string) PolicyEngineSummary {
	policies := pe.policies
	violations := append(
		EvaluateDesignRules(policies, content),
		EvaluateSecurityRules(policies, content)...,
	)
	return pe.Summarize(violations)
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
