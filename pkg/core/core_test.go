package core_test

import (
	"testing"

	"github.com/agentjam/agentjam/pkg/core"
)

func TestCoreUtilsAndTypes(t *testing.T) {
	// Test EstimateTokens
	text := "AgentJam system instructions text"
	tokens := core.EstimateTokens(text)
	if tokens <= 0 {
		t.Errorf("Expected positive token estimate, got %d", tokens)
	}

	// Test MatchRegex
	matched := core.MatchRegex(`(?i)agentjam`, "Welcome to AgentJam")
	if !matched {
		t.Errorf("Regex matching failed")
	}

	// Test TimestampNow
	ts := core.TimestampNow()
	if ts == "" {
		t.Errorf("TimestampNow returned empty string")
	}
}
