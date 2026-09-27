package install_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	agentjam "github.com/primasdevlabs/agentjam"
	"github.com/primasdevlabs/agentjam/pkg/install"
)

func TestInitMaterializesEmbeddedTree(t *testing.T) {
	dir := t.TempDir()
	res, err := install.Init(agentjam.CanonicalFS, dir, install.InitOptions{})
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	if len(res.Written) < 100 {
		t.Fatalf("expected >=100 files materialized, got %d", len(res.Written))
	}
	for _, rel := range []string{
		".agentjam/config.yaml",
		"agents/software-engineer/agent.yaml",
		"policies/security/04-security.yaml",
		"skills/development/testing/skill.yaml",
	} {
		if _, err := os.Stat(filepath.Join(dir, filepath.FromSlash(rel))); err != nil {
			t.Errorf("expected %s to exist: %v", rel, err)
		}
	}
}

func TestInitSkipsExistingAndForceOverwrites(t *testing.T) {
	dir := t.TempDir()
	if _, err := install.Init(agentjam.CanonicalFS, dir, install.InitOptions{}); err != nil {
		t.Fatal(err)
	}
	cfgPath := filepath.Join(dir, ".agentjam", "config.yaml")
	if err := os.WriteFile(cfgPath, []byte("custom: true\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	res, err := install.Init(agentjam.CanonicalFS, dir, install.InitOptions{})
	if err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(cfgPath)
	if !strings.Contains(string(data), "custom: true") {
		t.Error("existing config overwritten without --force")
	}
	if len(res.Skipped) == 0 {
		t.Error("expected skipped files for existing tree")
	}

	if _, err := install.Init(agentjam.CanonicalFS, dir, install.InitOptions{Force: true}); err != nil {
		t.Fatal(err)
	}
	data, _ = os.ReadFile(cfgPath)
	if strings.Contains(string(data), "custom: true") {
		t.Error("--force should overwrite existing config")
	}
}

func TestInitBare(t *testing.T) {
	dir := t.TempDir()
	if _, err := install.Init(agentjam.CanonicalFS, dir, install.InitOptions{Bare: true}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, ".agentjam", "config.yaml")); err != nil {
		t.Error("bare init should write config.yaml")
	}
	if _, err := os.Stat(filepath.Join(dir, "agents", "software-engineer")); !os.IsNotExist(err) {
		t.Error("bare init should not materialize resources")
	}
	for _, d := range []string{"agents", "skills", "policies"} {
		if info, err := os.Stat(filepath.Join(dir, d)); err != nil || !info.IsDir() {
			t.Errorf("bare init should create %s skeleton dir", d)
		}
	}
}

func TestCollectFromEmbeddedTree(t *testing.T) {
	dir := t.TempDir()
	if _, err := install.Init(agentjam.CanonicalFS, dir, install.InitOptions{}); err != nil {
		t.Fatal(err)
	}
	rs, err := install.Collect(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(rs.Agents) == 0 || len(rs.Skills) == 0 || len(rs.Policies) == 0 {
		t.Fatalf("empty ResourceSet: %+v agents=%d skills=%d policies=%d",
			rs, len(rs.Agents), len(rs.Skills), len(rs.Policies))
	}
}

func TestRenderCursorEmitsMDC(t *testing.T) {
	dir := t.TempDir()
	if _, err := install.Init(agentjam.CanonicalFS, dir, install.InitOptions{}); err != nil {
		t.Fatal(err)
	}
	rs, err := install.Collect(dir)
	if err != nil {
		t.Fatal(err)
	}
	files, err := install.Render(rs, "cursor", "INSTRUCTIONS", "software-engineer", "desc")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := files[".cursorrules"]; !ok {
		t.Fatal("expected .cursorrules")
	}
	if _, ok := files["AGENTJAM.md"]; !ok {
		t.Fatal("expected AGENTJAM.md runbook entry point")
	}
	mdcFound := 0
	for name, content := range files {
		if strings.HasPrefix(name, ".cursor/rules/skills/") && strings.HasSuffix(name, ".mdc") {
			mdcFound++
			if !strings.Contains(content, "alwaysApply:") {
				t.Errorf(".mdc missing frontmatter: %s", name)
			}
		}
	}
	if mdcFound == 0 {
		t.Error("expected per-skill .mdc files")
	}
}

func TestRenderClaudeCodeEmitsSkillDirs(t *testing.T) {
	dir := t.TempDir()
	install.Init(agentjam.CanonicalFS, dir, install.InitOptions{})
	rs, _ := install.Collect(dir)
	files, err := install.Render(rs, "claude-code", "INSTRUCTIONS", "a", "d")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := files["CLAUDE.md"]; !ok {
		t.Fatal("expected CLAUDE.md")
	}
	found := false
	for name := range files {
		if strings.HasPrefix(name, ".claude/skills/") && strings.HasSuffix(name, "/SKILL.md") {
			found = true
		}
	}
	if !found {
		t.Error("expected .claude/skills/<name>/SKILL.md files")
	}
}

func TestRenderUnknownHarness(t *testing.T) {
	if _, err := install.Render(&install.ResourceSet{}, "bogus", "", "a", "d"); err == nil {
		t.Fatal("expected error for unknown harness")
	}
}
