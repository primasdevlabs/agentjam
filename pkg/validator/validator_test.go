package validator_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/primasdevlabs/agentjam/pkg/validator"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	return filepath.Clean(filepath.Join(cwd, "..", ".."))
}

func TestValidateRepositoryRealRepo(t *testing.T) {
	errs := validator.ValidateRepository(repoRoot(t))
	if validator.HasErrors(errs) {
		for _, e := range errs {
			if e.Severity == "error" {
				t.Errorf("Unexpected validation error: %+v", e)
			}
		}
	}
}

func TestValidateAgentCrossReferences(t *testing.T) {
	dir := t.TempDir()

	skillDir := filepath.Join(dir, "skills", "testing")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(skillDir, "skill.yaml"),
		[]byte("name: testing\nversion: 1.0.0\ntype: skill\ndescription: Testing skill\n"), 0644)

	agentDir := filepath.Join(dir, "agents", "coder")
	if err := os.MkdirAll(agentDir, 0755); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(agentDir, "agent.yaml"),
		[]byte("name: coder\nversion: 1.0.0\ntype: agent\ndescription: Test agent\nskills:\n  - testing\n  - nonexistent\ntools:\n  - ghost-tool\n"), 0644)

	errs := validator.ValidateRepository(dir)

	var missing []string
	for _, e := range errs {
		missing = append(missing, e.Message)
	}
	foundSkill, foundTool := false, false
	for _, m := range missing {
		if contains(m, "unknown skill 'nonexistent'") {
			foundSkill = true
		}
		if contains(m, "unknown tool 'ghost-tool'") {
			foundTool = true
		}
	}
	if !foundSkill {
		t.Errorf("Expected unknown skill warning, got: %v", missing)
	}
	if !foundTool {
		t.Errorf("Expected unknown tool warning, got: %v", missing)
	}
	// 'testing' should not be flagged
	for _, m := range missing {
		if contains(m, "unknown skill 'testing'") {
			t.Errorf("Known skill incorrectly flagged: %s", m)
		}
	}
}

func TestValidateMalformedManifest(t *testing.T) {
	dir := t.TempDir()
	agentDir := filepath.Join(dir, "agents", "broken")
	if err := os.MkdirAll(agentDir, 0755); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(agentDir, "agent.yaml"), []byte("name: [unclosed"), 0644)

	errs := validator.ValidateRepository(dir)
	if !validator.HasErrors(errs) {
		t.Error("Expected error severity for malformed manifest")
	}
}

func TestValidateWorkflowDuplicateSteps(t *testing.T) {
	dir := t.TempDir()
	wfDir := filepath.Join(dir, "workflows", "bad-flow")
	if err := os.MkdirAll(wfDir, 0755); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(wfDir, "workflow.yaml"),
		[]byte("name: bad-flow\ntype: workflow\ndescription: dup steps\nsteps:\n  - id: s1\n    name: a\n  - id: s1\n    name: b\n"), 0644)

	errs := validator.ValidateRepository(dir)
	found := false
	for _, e := range errs {
		if contains(e.Message, "Duplicate workflow step id 's1'") {
			found = true
		}
	}
	if !found {
		t.Errorf("Expected duplicate step id error, got: %+v", errs)
	}
}

func contains(s, sub string) bool {
	return strings.Contains(s, sub)
}
