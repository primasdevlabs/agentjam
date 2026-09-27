package policy

import (
	"regexp"
	"strings"

	"github.com/primasdevlabs/agentjam/pkg/core"
)

var (
	emojiPattern    = regexp.MustCompile(`[\x{1F600}-\x{1F64F}\x{1F300}-\x{1F5FF}\x{1F680}-\x{1F6FF}\x{1F700}-\x{1F77F}\x{2600}-\x{26FF}\x{2700}-\x{27BF}]`)
	gradientPattern = regexp.MustCompile(`from-(purple|indigo|violet)-[0-9]+.*to-(indigo|blue|purple)-[0-9]+`)
	buzzwordPattern = regexp.MustCompile(`(?i)\b(empower|supercharge|seamlessly|next-generation|game-changer)\b`)
)

// firstMatchLine returns the 1-based line number of the first regex match.
func firstMatchLine(re *regexp.Regexp, content string) int {
	loc := re.FindStringIndex(content)
	if loc == nil {
		return 0
	}
	return strings.Count(content[:loc[0]], "\n") + 1
}

// EvaluateDesignRules evaluates design content against anti-slop rules. When a
// non-empty policy list is supplied, design rules only run if a "design"
// category policy is loaded.
func EvaluateDesignRules(policies []core.PolicyManifest, content string) []EvaluationViolation {
	violations := make([]EvaluationViolation, 0)
	if !policyActive(policies, "design") {
		return violations
	}

	if emojiPattern.MatchString(content) {
		violations = append(violations, EvaluationViolation{
			PolicyID:    "design-no-emoji-icons",
			PolicyName:  "Prohibit UI Emojis",
			Enforcement: core.EnforceStrictBlock,
			RuleName:    "no-emojis-in-ui",
			Message:     "UI contains prohibited emoji characters. Use SVG icon libraries (Lucide, Heroicons) instead.",
			Line:        firstMatchLine(emojiPattern, content),
		})
	}

	if gradientPattern.MatchString(content) {
		violations = append(violations, EvaluationViolation{
			PolicyID:    "design-no-slop-gradients",
			PolicyName:  "Prohibit Slop Gradients",
			Enforcement: core.EnforceStrictBlock,
			RuleName:    "no-purple-indigo-gradients",
			Message:     "UI contains generic glowing dark mode card gradients. Use curated HSL color tokens.",
			Line:        firstMatchLine(gradientPattern, content),
		})
	}

	if buzzwordPattern.MatchString(content) {
		violations = append(violations, EvaluationViolation{
			PolicyID:    "design-copy-governance",
			PolicyName:  "No Marketing Buzzwords",
			Enforcement: core.EnforceWarning,
			RuleName:    "no-filler-marketing-copy",
			Message:     "Copy contains filler marketing buzzwords. Use concise, human-centric product copy.",
			Line:        firstMatchLine(buzzwordPattern, content),
		})
	}

	return violations
}
