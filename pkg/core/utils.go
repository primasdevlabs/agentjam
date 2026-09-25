package core

import (
	"math"
	"regexp"
	"strings"
)

// EstimateTokens calculates estimated token count (~4 chars/token).
func EstimateTokens(text string) int {
	if len(text) == 0 {
		return 0
	}
	return int(math.Ceil(float64(len(text)) / 4.0))
}

// SanitizeString cleans string for markdown output.
func SanitizeString(s string) string {
	return strings.TrimSpace(s)
}

// MatchRegex reports whether string matches regular expression.
func MatchRegex(pattern string, text string) bool {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return re.MatchString(text)
}
