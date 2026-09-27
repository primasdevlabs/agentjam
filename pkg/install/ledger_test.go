package install

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, root, rel, content string) {
	t.Helper()
	full := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestLedgerRoundTrip(t *testing.T) {
	dir := t.TempDir()
	if err := WriteLedger(dir, "cursor", []string{".cursorrules", ".cursor/rules/skills/dev/a.mdc"}); err != nil {
		t.Fatal(err)
	}
	if err := WriteLedger(dir, "cursor", []string{".cursor/rules/agents/x.mdc"}); err != nil {
		t.Fatal(err)
	}
	l, err := ReadLedger(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(l.Harnesses["cursor"]) != 3 {
		t.Fatalf("expected 3 merged entries, got %v", l.Harnesses["cursor"])
	}
}

func TestUninstallRemovesLedgerFiles(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, ".cursorrules", "rules")
	writeFile(t, dir, ".cursor/rules/skills/dev/a.mdc", "skill")
	writeFile(t, dir, "AGENTJAM.md", "runbook")
	writeFile(t, dir, "user-owned.txt", "mine")
	WriteLedger(dir, "cursor", []string{".cursorrules", ".cursor/rules/skills/dev/a.mdc", "AGENTJAM.md"})

	res, err := Uninstall(dir, nil, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Removed) != 3 {
		t.Fatalf("expected 3 removed, got %v", res.Removed)
	}
	for _, gone := range []string{".cursorrules", ".cursor/rules/skills/dev/a.mdc", "AGENTJAM.md", LedgerFile} {
		if _, err := os.Stat(filepath.Join(dir, filepath.FromSlash(gone))); !os.IsNotExist(err) {
			t.Fatalf("%s should be gone", gone)
		}
	}
	// User file untouched; empty .cursor/rules tree pruned.
	if _, err := os.Stat(filepath.Join(dir, "user-owned.txt")); err != nil {
		t.Fatal("user file was deleted")
	}
	if _, err := os.Stat(filepath.Join(dir, ".cursor", "rules")); !os.IsNotExist(err) {
		t.Fatal("empty rules dir should be pruned")
	}
}

func TestUninstallWithoutLedger(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, ".claude/skills/dev/x/SKILL.md", "s")
	writeFile(t, dir, "AGENTJAM.md", "r")
	writeFile(t, dir, "keep.go", "package main")
	res, err := Uninstall(dir, nil, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Removed) == 0 {
		t.Fatal("expected fallback removal of well-known paths")
	}
	if _, err := os.Stat(filepath.Join(dir, "keep.go")); err != nil {
		t.Fatal("keep.go deleted")
	}
}

func TestUninstallPurge(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "AGENTJAM.md", "r")
	writeFile(t, dir, "CLAUDE.md", "generated")
	writeFile(t, dir, ".agentjam/config.yaml", "x")
	res, err := Uninstall(dir, nil, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Purged) == 0 {
		t.Fatal("expected purged paths")
	}
	if _, err := os.Stat(filepath.Join(dir, ".agentjam")); !os.IsNotExist(err) {
		t.Fatal(".agentjam should be purged")
	}
	if _, err := os.Stat(filepath.Join(dir, "CLAUDE.md")); !os.IsNotExist(err) {
		t.Fatal("CLAUDE.md should be purged")
	}
}
