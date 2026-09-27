package toolchain_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/primasdevlabs/agentjam/pkg/toolchain"
)

func TestHasBinary(t *testing.T) {
	tm := toolchain.NewToolchainManager(t.TempDir())
	if !tm.HasBinary("go") {
		t.Skip("go binary not in PATH; skipping")
	}
	if tm.HasBinary("definitely-not-a-real-binary-agentjam") {
		t.Error("HasBinary returned true for nonexistent binary")
	}
}

func TestDetectStacks(t *testing.T) {
	dir := t.TempDir()
	for _, f := range []string{"go.mod", "package.json", "Cargo.toml"} {
		if err := os.WriteFile(filepath.Join(dir, f), []byte("{}"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	tm := toolchain.NewToolchainManager(dir)
	stacks := tm.DetectStacks()

	ids := map[string]bool{}
	for _, s := range stacks {
		ids[s.ID] = true
	}
	for _, want := range []string{"go", "nodejs", "rust"} {
		if !ids[want] {
			t.Errorf("Expected stack %q to be detected, got %+v", want, stacks)
		}
	}
}

func TestRunPreflightChecks(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module x\n"), 0644); err != nil {
		t.Fatal(err)
	}

	tm := toolchain.NewToolchainManager(dir)
	result := tm.RunPreflightChecks()

	if result.StackID != "go" {
		t.Errorf("Expected stackId 'go', got %q", result.StackID)
	}
	if len(result.Checks) == 0 {
		t.Fatal("Expected at least one preflight check")
	}
	found := false
	for _, c := range result.Checks {
		if c.Command == "go" {
			found = true
			if c.Status != "pass" && c.Status != "fail" {
				t.Errorf("Unexpected status %q", c.Status)
			}
		}
	}
	if !found {
		t.Error("Expected go binary check in preflight results")
	}
}

func TestRunPreflightEmptyDir(t *testing.T) {
	tm := toolchain.NewToolchainManager(t.TempDir())
	result := tm.RunPreflightChecks()
	if !result.Passed || len(result.Checks) != 1 {
		t.Errorf("Expected fallback pass check for empty dir, got %+v", result)
	}
}

func TestRunCommand(t *testing.T) {
	tm := toolchain.NewToolchainManager(t.TempDir())
	ok, out, _ := tm.RunCommand("echo hello-agentjam", 5*time.Second)
	if !ok {
		t.Fatalf("RunCommand failed: %s", out)
	}
	if !containsStr(out, "hello-agentjam") {
		t.Errorf("Expected output to contain echo text, got %q", out)
	}
}

func containsStr(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
