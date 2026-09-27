package core_test

import (
	"testing"
	"time"

	"github.com/primasdevlabs/agentjam/pkg/core"
)

func TestEstimateTokens(t *testing.T) {
	if core.EstimateTokens("") != 0 {
		t.Error("Empty string must estimate 0 tokens")
	}
	if core.EstimateTokens("1234") != 1 {
		t.Error("4 chars must estimate 1 token")
	}
	if core.EstimateTokens("12345") != 2 {
		t.Error("5 chars must round up to 2 tokens")
	}
}

func TestSanitizeString(t *testing.T) {
	if got := core.SanitizeString("  hello\n"); got != "hello" {
		t.Errorf("Expected trimmed string, got %q", got)
	}
}

func TestMatchRegex(t *testing.T) {
	if !core.MatchRegex(`\d+`, "abc123") {
		t.Error("Expected regex match")
	}
	if core.MatchRegex(`[invalid`, "x") {
		t.Error("Invalid regex must return false, not panic")
	}
}

func TestTimestampNow(t *testing.T) {
	ts := core.TimestampNow()
	if _, err := time.Parse(time.RFC3339, ts); err != nil {
		t.Errorf("TimestampNow not RFC3339: %q", ts)
	}
}

func TestSafetyLevels(t *testing.T) {
	// Ensure declared safety constants serialize as expected
	levels := []core.ToolSafetyLevel{
		core.SafetyReadOnly, core.SafetySafeWrite, core.SafetyDestructive, core.SafetyAdmin,
	}
	want := []string{"read-only", "safe-write", "destructive", "admin"}
	for i, l := range levels {
		if string(l) != want[i] {
			t.Errorf("Safety level mismatch: %q != %q", l, want[i])
		}
	}
}
