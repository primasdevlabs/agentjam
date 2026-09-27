package runtime

import (
	"regexp"
	"strconv"
	"time"
)

// FreshnessResult contains documentation age calculation.
type FreshnessResult struct {
	IsFresh    bool   `json:"isFresh"`
	AgeDays    int    `json:"ageDays"`
	MaxDocAge  string `json:"maxDocAge"`
	LastUpdate string `json:"lastUpdate"`
}

var docAgePattern = regexp.MustCompile(`^(\d+)\s*([hdwmy])$`)

// maxDocAgeDays converts a max-age spec to days. Supported suffixes:
// h (hours, rounded up), d (days), w (weeks), m (~30 days), y (~365 days).
// Defaults to 7 days on unrecognized input.
func maxDocAgeDays(spec string) float64 {
	m := docAgePattern.FindStringSubmatch(spec)
	if m == nil {
		return 7
	}
	n, err := strconv.Atoi(m[1])
	if err != nil {
		return 7
	}
	switch m[2] {
	case "h":
		return float64(n) / 24.0
	case "d":
		return float64(n)
	case "w":
		return float64(n) * 7
	case "m":
		return float64(n) * 30
	case "y":
		return float64(n) * 365
	}
	return 7
}

// CheckFreshness evaluates documentation age against maximum age policy (e.g., "7d").
func CheckFreshness(lastUpdatedIso string, maxDocAge string) FreshnessResult {
	parsedTime, err := time.Parse(time.RFC3339, lastUpdatedIso)
	if err != nil {
		return FreshnessResult{IsFresh: true, AgeDays: 0, MaxDocAge: maxDocAge, LastUpdate: lastUpdatedIso}
	}

	maxDays := maxDocAgeDays(maxDocAge)
	ageHours := time.Since(parsedTime).Hours()
	ageDays := int(ageHours / 24.0)

	return FreshnessResult{
		IsFresh:    ageHours <= maxDays*24.0,
		AgeDays:    ageDays,
		MaxDocAge:  maxDocAge,
		LastUpdate: lastUpdatedIso,
	}
}
