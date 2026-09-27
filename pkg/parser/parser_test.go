package parser_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/primasdevlabs/agentjam/pkg/core"
	"github.com/primasdevlabs/agentjam/pkg/parser"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	return filepath.Clean(filepath.Join(cwd, "..", ".."))
}

func TestParseFrontmatter(t *testing.T) {
	content := "---\nname: test-agent\nversion: 1.0.0\n---\n# Agent Role\nRole instructions"
	fm, body := parser.ParseFrontmatter(content)
	if fm == "" || body == "" {
		t.Fatalf("Frontmatter parsing failed: fm=%q, body=%q", fm, body)
	}
	if body == fm {
		t.Fatalf("Frontmatter and body must differ")
	}

	// No frontmatter: whole content is body
	fm2, body2 := parser.ParseFrontmatter("# Just markdown")
	if fm2 != "" {
		t.Errorf("Expected empty frontmatter, got %q", fm2)
	}
	if body2 != "# Just markdown" {
		t.Errorf("Expected body to equal input, got %q", body2)
	}
}

func TestDiscoverResourcesByType(t *testing.T) {
	resources := parser.DiscoverResources(repoRoot(t))
	if len(resources) == 0 {
		t.Fatal("Expected discovered resources, got none")
	}

	counts := map[core.ResourceType]int{}
	for _, r := range resources {
		counts[r.Type]++
		if r.ID == "" || r.Path == "" {
			t.Errorf("Discovered resource missing id/path: %+v", r)
		}
	}

	for _, rt := range []core.ResourceType{
		core.ResourceTypeAgent, core.ResourceTypeSkill, core.ResourceTypeTool,
		core.ResourceTypeWorkflow, core.ResourceTypePolicy, core.ResourceTypeStack,
		core.ResourceTypeLanguage, core.ResourceTypeIntegration,
	} {
		if counts[rt] == 0 {
			t.Errorf("Expected at least one %s resource, got 0", rt)
		}
	}
}

func TestParseAgentYAML(t *testing.T) {
	path := filepath.Join(repoRoot(t), "agents", "software-engineer", "agent.yaml")
	bundle, err := parser.ParseAgent(path)
	if err != nil {
		t.Fatalf("ParseAgent failed: %v", err)
	}
	m := bundle.Manifest
	if m.Name != "software-engineer" {
		t.Errorf("Expected name software-engineer, got %q", m.Name)
	}
	if len(m.Skills) == 0 {
		t.Error("Expected skills list to be populated from YAML")
	}
	if len(m.Tools) == 0 {
		t.Error("Expected tools list to be populated from YAML")
	}
	if len(bundle.Instructions) == 0 {
		t.Error("Expected instruction files to be loaded")
	}
}

func TestParseWorkflowYAML(t *testing.T) {
	path := filepath.Join(repoRoot(t), "workflows", "existing-project-audit", "workflow.yaml")
	bundle, err := parser.ParseWorkflow(path)
	if err != nil {
		t.Fatalf("ParseWorkflow failed: %v", err)
	}
	if bundle.Manifest.Name != "existing-project-audit" {
		t.Errorf("Unexpected workflow name %q", bundle.Manifest.Name)
	}
	if len(bundle.Manifest.Steps) != 7 {
		t.Errorf("Expected 7 workflow steps, got %d", len(bundle.Manifest.Steps))
	}
}

func TestParseToolYAML(t *testing.T) {
	path := filepath.Join(repoRoot(t), "tools", "filesystem", "tool.yaml")
	m, err := parser.ParseTool(path)
	if err != nil {
		t.Fatalf("ParseTool failed: %v", err)
	}
	if m.Name != "filesystem" {
		t.Errorf("Expected tool name filesystem, got %q", m.Name)
	}
	if m.SafetyLevel != core.SafetySafeWrite {
		t.Errorf("Expected safetyLevel safe-write, got %q", m.SafetyLevel)
	}
}

func TestParseStackYAML(t *testing.T) {
	path := filepath.Join(repoRoot(t), "stacks", "nextjs", "stack.yaml")
	m, err := parser.ParseStack(path)
	if err != nil {
		t.Fatalf("ParseStack failed: %v", err)
	}
	if m.Name != "nextjs" || m.Framework != "Next.js" {
		t.Errorf("Unexpected stack manifest: %+v", m)
	}
}

func TestParseLanguageYAML(t *testing.T) {
	path := filepath.Join(repoRoot(t), "languages", "web", "typescript", "language.yaml")
	m, err := parser.ParseLanguage(path)
	if err != nil {
		t.Fatalf("ParseLanguage failed: %v", err)
	}
	if m.ID != "lang-typescript" || m.Ecosystem != "web" {
		t.Errorf("Unexpected language manifest: %+v", m)
	}
	if len(m.Extensions) == 0 {
		t.Error("Expected language extensions to be parsed")
	}
}

func TestParsePolicyYAMLAndMarkdown(t *testing.T) {
	root := repoRoot(t)

	p, err := parser.ParsePolicy(filepath.Join(root, "policies", "security", "04-security.yaml"))
	if err != nil {
		t.Fatalf("ParsePolicy yaml failed: %v", err)
	}
	if p.ID != "04-security" || p.Enforcement != core.EnforceStrictBlock {
		t.Errorf("Unexpected policy manifest: %+v", p)
	}
	if p.Instructions == "" {
		t.Error("Expected policy instructions to be populated")
	}

	pm, err := parser.ParsePolicy(filepath.Join(root, "policies", "design", "anti-slop.md"))
	if err != nil {
		t.Fatalf("ParsePolicy markdown failed: %v", err)
	}
	if pm.ID != "anti-slop" || pm.Instructions == "" {
		t.Errorf("Unexpected markdown policy: %+v", pm)
	}
}

func TestParseAgentBadYAML(t *testing.T) {
	dir := t.TempDir()
	bad := filepath.Join(dir, "agent.yaml")
	if err := os.WriteFile(bad, []byte("name: [unclosed"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := parser.ParseAgent(bad); err == nil {
		t.Error("Expected parse error for malformed YAML, got nil")
	}
}

func TestLoadWorkspaceConfig(t *testing.T) {
	// Self-contained: .agentjam/ is gitignored, so never rely on the repo's copy.
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".agentjam"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".agentjam", "config.yaml"), []byte(`versionPolicy: current-stable
freshnessRequired: true
maxDocAge: 7d
stack: generic
defaultAgent: software-engineer
`), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := parser.LoadWorkspaceConfig(dir)
	if err != nil {
		t.Fatalf("LoadWorkspaceConfig failed: %v", err)
	}
	if cfg.MaxDocAge != "7d" {
		t.Errorf("Expected maxDocAge 7d, got %q", cfg.MaxDocAge)
	}
	if cfg.DefaultAgent != "software-engineer" || cfg.Stack != "generic" {
		t.Errorf("Expected stack/agent defaults, got %+v", cfg)
	}

	if _, err := parser.LoadWorkspaceConfig(t.TempDir()); err == nil {
		t.Error("Expected error for missing config.yaml")
	}
}
