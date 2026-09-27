package policy_test

import (
	"strings"
	"testing"

	"github.com/primasdevlabs/agentjam/pkg/core"
	"github.com/primasdevlabs/agentjam/pkg/policy"
)

func structureEngine(t *testing.T, maxLines int) *policy.PolicyEngine {
	t.Helper()
	pe := policy.NewPolicyEngine()
	pe.AddPolicy(core.PolicyManifest{
		ID:          "06-god-files",
		Name:        "God File Limit",
		Category:    "architecture",
		Enforcement: core.EnforceWarning,
		Rules:       map[string]interface{}{"max_file_lines": maxLines, "max_top_level_decls": 40, "max_responsibility_markers": 8},
	})
	pe.AddPolicy(core.PolicyManifest{
		ID:          "07-domain-layout",
		Name:        "Domain-Driven Layout",
		Category:    "architecture",
		Enforcement: core.EnforceWarning,
		Rules:       map[string]interface{}{"flat_source_threshold": 5, "max_files_per_dir": 4, "min_domain_ratio": 0.6},
	})
	return pe
}

func TestGodFileDetection(t *testing.T) {
	pe := structureEngine(t, 10)
	big := strings.Repeat("// filler\n", 50)
	v := pe.EvaluateFileStructure("src/mega.go", big)
	if len(v) == 0 {
		t.Fatal("expected god-file violation for oversized file")
	}
	if v[0].RuleName != "max_file_lines" {
		t.Fatalf("expected max_file_lines violation, got %s", v[0].RuleName)
	}
}

func TestCleanFileNoViolation(t *testing.T) {
	pe := structureEngine(t, 500)
	v := pe.EvaluateFileStructure("src/small.go", "package x\n\nfunc a() {}\n")
	if len(v) != 0 {
		t.Fatalf("expected no violations, got %+v", v)
	}
}

func TestFlatLayoutViolation(t *testing.T) {
	pe := structureEngine(t, 500)
	files := []string{"a.go", "b.go", "c.go", "d.go", "e.go", "f.go", "g.go"}
	v := pe.EvaluateProjectLayout(".", files)
	if len(v) == 0 {
		t.Fatal("expected flat-layout violation")
	}
	if v[0].RuleName != "flat_layout" {
		t.Fatalf("expected flat_layout rule, got %s", v[0].RuleName)
	}
}

func TestDomainLayoutPasses(t *testing.T) {
	pe := structureEngine(t, 500)
	files := []string{
		"main.go",
		"auth/handler.go", "auth/store.go", "auth/token.go",
		"billing/invoice.go", "billing/charge.go",
		"inventory/stock.go", "inventory/sku.go",
	}
	v := pe.EvaluateProjectLayout(".", files)
	if len(v) != 0 {
		t.Fatalf("expected clean domain layout, got %+v", v)
	}
}

func TestOversizedDirectoryViolation(t *testing.T) {
	pe := structureEngine(t, 500)
	files := []string{
		"auth/a.go", "auth/b.go", "auth/c.go", "auth/d.go", "auth/e.go", "auth/f.go",
	}
	v := pe.EvaluateProjectLayout(".", files)
	found := false
	for _, vi := range v {
		if vi.RuleName == "oversized_directory" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected oversized_directory violation, got %+v", v)
	}
}

func TestNoStructurePoliciesNoViolations(t *testing.T) {
	pe := policy.NewPolicyEngine()
	if v := pe.EvaluateFileStructure("x.go", strings.Repeat("x\n", 9999)); len(v) != 0 {
		t.Fatal("expected no violations when god-file policy absent")
	}
	if v := pe.EvaluateProjectLayout(".", []string{"a.go"}); len(v) != 0 {
		t.Fatal("expected no violations when domain-layout policy absent")
	}
}
