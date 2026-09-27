package adapters_test

import (
	"strings"
	"testing"

	"github.com/primasdevlabs/agentjam/pkg/adapters"
	"github.com/primasdevlabs/agentjam/pkg/core"
)

func testAgent() core.AgentManifest {
	return core.AgentManifest{
		Name:        "software-engineer",
		Description: "Fullstack Engineer",
		Version:     "1.0.0",
	}
}

func TestExportAllHarnesses(t *testing.T) {
	results := adapters.ExportAllHarnesses(testAgent(), "INSTRUCTIONS-BODY")

	expected := map[string]string{
		"claude-code": "CLAUDE.md",
		"cursor":      ".cursorrules",
		"gemini":      "GEMINI.md",
		"cline":       ".clinerules",
		"windsurf":    ".windsurfrules",
		"devin":       ".devin/playbook.md",
		"roo-code":    ".roomodes",
		"generic":     "SYSTEM_PROMPT.md",
	}

	for harness, file := range expected {
		res, ok := results[harness]
		if !ok {
			t.Errorf("Missing harness export %q", harness)
			continue
		}
		content, ok := res.Files[file]
		if !ok {
			t.Errorf("Harness %q missing expected file %q", harness, file)
			continue
		}
		if !strings.Contains(content, "INSTRUCTIONS-BODY") {
			t.Errorf("Harness %q output missing instructions body", harness)
		}
		if !strings.Contains(content, "software-engineer") {
			t.Errorf("Harness %q output missing agent name", harness)
		}
	}
}

func TestExportIndividualAdapters(t *testing.T) {
	res := adapters.ExportClaudeCode(testAgent(), "body")
	if res.Harness != "claude-code" {
		t.Errorf("Unexpected harness name %q", res.Harness)
	}
	if !strings.Contains(res.Files["CLAUDE.md"], "Fullstack Engineer") {
		t.Error("CLAUDE.md should include agent description")
	}
}
