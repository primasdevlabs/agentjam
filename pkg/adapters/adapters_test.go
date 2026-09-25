package adapters_test

import (
	"testing"

	"github.com/agentjam/agentjam/pkg/adapters"
	"github.com/agentjam/agentjam/pkg/core"
)

func TestAdaptersModuleAllHarnesses(t *testing.T) {
	agent := core.AgentManifest{
		Name:        "test-engineer",
		Description: "Fullstack Engineer",
	}

	results := adapters.ExportAllHarnesses(agent, "System instructions")

	expectedHarnesses := []string{
		"claude-code", "cursor", "gemini", "cline",
		"windsurf", "devin", "roo-code", "generic",
	}

	for _, harness := range expectedHarnesses {
		res, ok := results[harness]
		if !ok {
			t.Errorf("Missing export for harness %s", harness)
			continue
		}
		if len(res.Files) == 0 {
			t.Errorf("Empty files map for harness %s", harness)
		}
	}
}
