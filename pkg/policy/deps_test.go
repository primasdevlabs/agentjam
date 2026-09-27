package policy

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/primasdevlabs/agentjam/pkg/core"
)

func testPolicies() []core.PolicyManifest {
	return []core.PolicyManifest{
		{ID: "07-deprecated-packages", Category: "security", Enforcement: core.EnforceStrictBlock},
		{ID: "10-dependency-freshness", Category: "coding-standards", Enforcement: core.EnforceWarning},
	}
}

func TestGoModParser(t *testing.T) {
	data := []byte(`module example.com/app

go 1.22

require (
	github.com/dgrijalva/jwt-go v3.2.0
	github.com/pkg/errors v0.9.1 // indirect
)

require github.com/golang/protobuf v1.5.3
`)
	deps := parseGoMod(data)
	if len(deps) != 3 {
		t.Fatalf("expected 3 deps, got %v", deps)
	}
	if deps[0].Name != "github.com/dgrijalva/jwt-go" || deps[0].Version != "v3.2.0" {
		t.Fatalf("bad parse: %+v", deps[0])
	}
}

func TestPackageJSONParser(t *testing.T) {
	data := []byte(`{"dependencies":{"request":"^2.88.0","react":"^18.0.0"},"devDependencies":{"tslint":"^6.1.0"}}`)
	deps := parsePackageJSON(data)
	if len(deps) != 3 {
		t.Fatalf("expected 3 deps, got %v", deps)
	}
}

func TestRequirementsParser(t *testing.T) {
	data := []byte("# comment\npycrypto==2.6.1\nflask>=2.0\n-e .\n")
	deps := parseRequirements(data)
	if len(deps) != 2 {
		t.Fatalf("expected 2 deps, got %v", deps)
	}
	if deps[0].Name != "pycrypto" {
		t.Fatalf("expected pycrypto, got %s", deps[0].Name)
	}
}

func TestDeprecatedDetectionStrict(t *testing.T) {
	pe := &PolicyEngine{policies: testPolicies()}
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"),
		[]byte("module x\n\ngo 1.22\n\nrequire github.com/dgrijalva/jwt-go v3.2.0\n"), 0644); err != nil {
		t.Fatal(err)
	}
	out := pe.EvaluateDependencies(dir, false)
	vls := out["go.mod"]
	if len(vls) == 0 {
		t.Fatal("expected violation for jwt-go")
	}
	if vls[0].Enforcement != core.EnforceStrictBlock {
		t.Fatalf("expected strict-block, got %s", vls[0].Enforcement)
	}
}

func TestDeprecatedVersionGated(t *testing.T) {
	pe := &PolicyEngine{policies: testPolicies()}
	dir := t.TempDir()
	// core-js ^2.6 -> flagged; ^3.30 -> not flagged
	pkg := `{"dependencies":{"core-js":"^2.6.0"}}`
	if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte(pkg), 0644); err != nil {
		t.Fatal(err)
	}
	out := pe.EvaluateDependencies(dir, false)
	if len(out["package.json"]) == 0 {
		t.Fatal("expected core-js 2.x to be flagged")
	}
	pkg = `{"dependencies":{"core-js":"^3.30.0"}}`
	if err := os.WriteFile(filepath.Join(dir, "package.json"), []byte(pkg), 0644); err != nil {
		t.Fatal(err)
	}
	out = pe.EvaluateDependencies(dir, false)
	if len(out["package.json"]) != 0 {
		t.Fatalf("core-js 3.x should not be flagged, got %v", out["package.json"])
	}
}

func TestNoPoliciesNoDeps(t *testing.T) {
	// Engine with no policies → policyActive returns true (all run), but
	// deprecated/freshness still need their categories... verify gating:
	pe := &PolicyEngine{policies: []core.PolicyManifest{
		{ID: "01-x", Category: "design"},
	}}
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"),
		[]byte("module x\n\nrequire github.com/dgrijalva/jwt-go v3.2.0\n"), 0644)
	out := pe.EvaluateDependencies(dir, false)
	if len(out) != 0 {
		t.Fatalf("expected no violations without security policy, got %v", out)
	}
}

func TestCleanManifestNoViolation(t *testing.T) {
	pe := &PolicyEngine{policies: testPolicies()}
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "package.json"),
		[]byte(`{"dependencies":{"react":"^18.2.0","zod":"^3.22.0"}}`), 0644)
	out := pe.EvaluateDependencies(dir, false)
	if len(out) != 0 {
		t.Fatalf("expected clean, got %v", out)
	}
}

func TestDeprecatedAPISource(t *testing.T) {
	vls := EvaluateDeprecatedAPIs(testPolicies(), "package x\n\nimport \"io/ioutil\"\n")
	if len(vls) == 0 {
		t.Fatal("expected io/ioutil violation")
	}
	vls = EvaluateDeprecatedAPIs(testPolicies(), "import distutils\n")
	if len(vls) == 0 {
		t.Fatal("expected distutils violation")
	}
	vls = EvaluateDeprecatedAPIs(nil, "clean content")
	if len(vls) != 0 {
		t.Fatal("expected clean")
	}
}
