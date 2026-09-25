package policy

import (
	"regexp"

	"github.com/primasdevlabs/agentjam/pkg/core"
)

// EvaluateDesignRules evaluates design content against anti-slop rules.
func EvaluateDesignRules(policies []core.PolicyManifest, content string) []EvaluationViolation {
	violations := make([]EvaluationViolation, 0)

	// 1. Prohibited Emoji Icons Rule
	emojiRegex := regexp.MustCompile(`[\x{1F600}-\x{1F64F}\x{1F300}-\x{1F5FF}\x{1F680}-\x{1F6FF}\x{1F700}-\x{1F77F}\x{2600}-\x{26FF}\x{2700}-\x{27BF}]`)
	if emojiRegex.MatchString(content) {
		violations = append(violations, EvaluationViolation{
			PolicyID:    "design-no-emoji-icons",
			PolicyName:  "Prohibit UI Emojis",
			Enforcement: core.EnforceStrictBlock,
			RuleName:    "no-emojis-in-ui",
			Message:     "UI contains prohibited emoji characters. Use SVG icon libraries (Lucide, Heroicons) instead.",
		})
	}

	// 2. Prohibited Dark Mode Gradient Cards Rule
	gradientRegex := regexp.MustCompile(`from-(purple|indigo|violet)-[0-9]+.*to-(indigo|blue|purple)-[0-9]+`)
	if gradientRegex.MatchString(content) {
		violations = append(violations, EvaluationViolation{
			PolicyID:    "design-no-slop-gradients",
			PolicyName:  "Prohibit Slop Gradients",
			Enforcement: core.EnforceStrictBlock,
			RuleName:    "no-purple-indigo-gradients",
			Message:     "UI contains generic glowing dark mode card gradients. Use curated HSL color tokens.",
		})
	}

	// 3. Prohibited Marketing Buzzwords Rule
	buzzwordsRegex := regexp.MustCompile(`(?i)\b(empower|supercharge|seamlessly|next-generation|game-changer)\b`)
	if buzzwordsRegex.MatchString(content) {
		violations = append(violations, EvaluationViolation{
			PolicyID:    "design-copy-governance",
			PolicyName:  "No Marketing Buzzwords",
			Enforcement: core.EnforceWarning,
			RuleName:    "no-filler-marketing-copy",
			Message:     "Copy contains filler marketing buzzwords. Use concise, human-centric product copy.",
		})
	}

	return violations
}
