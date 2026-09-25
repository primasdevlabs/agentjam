package runtime

import (
	"time"
)

// FreshnessResult contains documentation age calculation.
type FreshnessResult struct {
	IsFresh    bool   `json:"isFresh"`
	AgeDays    int    `json:"ageDays"`
	MaxDocAge  string `json:"maxDocAge"`
	LastUpdate string `json:"lastUpdate"`
}

// CheckFreshness evaluates documentation age against maximum age policy (e.g., "7d").
func CheckFreshness(lastUpdatedIso string, maxDocAge string) FreshnessResult {
	parsedTime, err := time.Parse(time.RFC3339, lastUpdatedIso)
	if err != nil {
		return FreshnessResult{IsFresh: true, AgeDays: 0, MaxDocAge: maxDocAge, LastUpdate: lastUpdatedIso}
	}

	ageDays := int(time.Since(parsedTime).Hours() / 24.0)
	maxDays := 7 // default
	if maxDocAge == "30d" {
		maxDays = 30
	} else if maxDocAge == "1d" {
		maxDays = 1
	}

	return FreshnessResult{
		IsFresh:    ageDays <= maxDays,
		AgeDays:    ageDays,
		MaxDocAge:  maxDocAge,
		LastUpdate: lastUpdatedIso,
	}
}
