package context

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/primasdevlabs/agentjam/pkg/core"
	"github.com/primasdevlabs/agentjam/pkg/policy"
)

// ContextOptions configures snapshot generation.
type ContextOptions struct {
	TokenBudget   int      `json:"tokenBudget"`
	ActiveAgent   string   `json:"activeAgent,omitempty"`
	ActiveStack   string   `json:"activeStack,omitempty"`
	ActiveSkills  []string `json:"activeSkills,omitempty"`
	CustomRules   []string `json:"customRules,omitempty"`
	Environment   string   `json:"environment,omitempty"`
}

// ContextManager manages context assembly, token budget estimation, and snapshot generation.
type ContextManager struct {
	workspaceRoot string
	policyEngine  *policy.PolicyEngine
}

// NewContextManager creates a new ContextManager instance.
func NewContextManager(workspaceRoot string, pe *policy.PolicyEngine) *ContextManager {
	if pe == nil {
		pe = policy.NewPolicyEngine()
	}
	return &ContextManager{
		workspaceRoot: workspaceRoot,
		policyEngine:  pe,
	}
}

// EstimateTokenCount estimates token count for a text string (~4 chars per token).
func (cm *ContextManager) EstimateTokenCount(text string) int {
	if len(text) == 0 {
		return 0
	}
	return int(math.Ceil(float64(len(text)) / 4.0))
}

// BuildContextSnapshot builds a full system context snapshot.
func (cm *ContextManager) BuildContextSnapshot(opts ContextOptions) core.ContextSnapshot {
	if opts.TokenBudget <= 0 {
		opts.TokenBudget = 128000
	}

	policies := cm.policyEngine.GetPolicies()
	var sb strings.Builder

	sb.WriteString("# AgentJam System Instructions & Active Context (Go Engine)\n")
	sb.WriteString(fmt.Sprintf("**Workspace**: `%s`\n", cm.workspaceRoot))
	sb.WriteString(fmt.Sprintf("**Timestamp**: %s\n", time.Now().Format(time.RFC3339)))

	if opts.ActiveAgent != "" {
		sb.WriteString(fmt.Sprintf("**Active Persona/Agent**: %s\n", opts.ActiveAgent))
	}
	if opts.Environment != "" {
		sb.WriteString(fmt.Sprintf("**Environment**: %s\n", opts.Environment))
	}

	stackName := opts.ActiveStack
	if stackName == "" {
		stackName = "Generic Polyglot"
	}
	sb.WriteString(fmt.Sprintf("\n## Active Stack Profile: %s\n", stackName))

	if len(opts.ActiveSkills) > 0 {
		sb.WriteString("\n## Active Skills\n")
		for _, skill := range opts.ActiveSkills {
			sb.WriteString(fmt.Sprintf("- %s\n", skill))
		}
	}

	sb.WriteString(fmt.Sprintf("\n## Enforced Policy Matrix (%d Policies Loaded)\n", len(policies)))
	for _, p := range policies {
		sb.WriteString(fmt.Sprintf("- **[%s] %s** (`%s`): %s\n", strings.ToUpper(string(p.Enforcement)), p.Name, p.ID, p.Description))
	}

	if len(opts.CustomRules) > 0 {
		sb.WriteString("\n## Project Specific Rules\n")
		for _, rule := range opts.CustomRules {
			sb.WriteString(fmt.Sprintf("- %s\n", rule))
		}
	}

	systemInstruction := sb.String()
	estimatedTokens := cm.EstimateTokenCount(systemInstruction)

	if estimatedTokens > opts.TokenBudget {
		maxChars := opts.TokenBudget * 4
		if maxChars < len(systemInstruction) {
			systemInstruction = systemInstruction[:maxChars] + "\n\n[Context Truncated to fit Token Budget]"
			estimatedTokens = cm.EstimateTokenCount(systemInstruction)
		}
	}

	return core.ContextSnapshot{
		WorkspaceRoot:      cm.workspaceRoot,
		SystemInstruction:  systemInstruction,
		TokenCountEstimate: estimatedTokens,
		PoliciesCount:      len(policies),
		ActiveFiles:        []string{},
		Timestamp:          time.Now().Format(time.RFC3339),
		Metadata: map[string]interface{}{
			"activeAgent": opts.ActiveAgent,
			"engine":      "Go-Native",
		},
	}
}
