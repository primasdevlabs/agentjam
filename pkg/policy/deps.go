package policy

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/primasdevlabs/agentjam/pkg/core"
)

// depDecl is one dependency pinned in a manifest.
type depDecl struct {
	Name    string
	Version string // raw constraint as written (e.g. "^4.17.0", ">=1.2")
	Eco     string // npm | go | pip | cargo | composer | gem
}

// DeprecatedPackage is a bundled advisory entry for a package that is
// deprecated, renamed, abandoned, or carries known vulnerabilities.
type DeprecatedPackage struct {
	Eco         string
	Name        string // exact name, or prefix when Prefix is true
	Prefix      bool
	Below       string // optional: only flag when pinned version < Below
	Severity    core.PolicyEnforcement
	Reason      string
	Replacement string
}

// deprecatedPackages is the built-in advisory DB — well-known deprecated,
// renamed, or known-vulnerable packages. Kept tight to avoid false
// positives; entries must be defensible.
var deprecatedPackages = []DeprecatedPackage{
	// npm
	{"npm", "request", false, "", core.EnforceWarning, "deprecated since 2020", "undici, node-fetch, or axios"},
	{"npm", "request-promise", false, "", core.EnforceWarning, "deprecated with request", "undici or axios"},
	{"npm", "request-promise-native", false, "", core.EnforceWarning, "deprecated with request", "undici or axios"},
	{"npm", "node-sass", false, "", core.EnforceWarning, "deprecated", "sass (dart-sass)"},
	{"npm", "tslint", false, "", core.EnforceWarning, "deprecated since 2019", "eslint + @typescript-eslint"},
	{"npm", "jade", false, "", core.EnforceWarning, "renamed", "pug"},
	{"npm", "bower", false, "", core.EnforceWarning, "deprecated package manager", "npm/pnpm/yarn"},
	{"npm", "gulp-util", false, "", core.EnforceWarning, "deprecated", "vinyl utilities / plugin-error"},
	{"npm", "phantomjs-prebuilt", false, "", core.EnforceWarning, "PhantomJS is abandoned", "playwright or puppeteer"},
	{"npm", "phantomjs", false, "", core.EnforceWarning, "PhantomJS is abandoned", "playwright or puppeteer"},
	{"npm", "faker", false, "", core.EnforceWarning, "original package abandoned", "@faker-js/faker"},
	{"npm", "moment", false, "", core.EnforceWarning, "maintenance mode per maintainers", "dayjs, date-fns, or luxon"},
	{"npm", "har-validator", false, "", core.EnforceWarning, "deprecated with request", "ajv"},
	{"npm", "core-js", false, "3.0.0", core.EnforceWarning, "core-js 2.x is unsupported", "core-js@3"},
	{"npm", "uuid", false, "4.0.0", core.EnforceWarning, "uuid v3 uses weak RNG", "uuid >= 9"},
	{"npm", "querystring", false, "", core.EnforceWarning, "deprecated Node shim", "URLSearchParams"},
	// Go modules
	{"go", "github.com/dgrijalva/jwt-go", false, "", core.EnforceStrictBlock, "archived, CVE-2020-26160", "github.com/golang-jwt/jwt/v5"},
	{"go", "github.com/satori/go.uuid", false, "", core.EnforceStrictBlock, "unmaintained, known RNG weakness", "github.com/gofrs/uuid"},
	{"go", "github.com/golang/protobuf", false, "", core.EnforceWarning, "superseded upstream", "google.golang.org/protobuf"},
	{"go", "github.com/golang/mock", false, "", core.EnforceWarning, "repository moved", "go.uber.org/mock"},
	{"go", "github.com/pkg/errors", false, "", core.EnforceWarning, "archived", "fmt.Errorf with %w"},
	{"go", "github.com/Sirupsen/logrus", false, "", core.EnforceWarning, "renamed (lowercase)", "github.com/sirupsen/logrus"},
	// pip
	{"pip", "pycrypto", false, "", core.EnforceStrictBlock, "dead since 2013, known CVEs", "pycryptodome"},
	{"pip", "sklearn", false, "", core.EnforceStrictBlock, "PyPI package is a deprecated shim that fails installs", "scikit-learn"},
	{"pip", "mysql-python", false, "", core.EnforceStrictBlock, "Python 2 only, unmaintained", "mysqlclient"},
	{"pip", "pil", false, "", core.EnforceWarning, "PIL is ancient; the maintained fork is Pillow", "Pillow"},
	{"pip", "nose", false, "", core.EnforceWarning, "unmaintained test runner", "pytest"},
	{"pip", "pep8", false, "", core.EnforceWarning, "renamed", "pycodestyle"},
	{"pip", "mock", false, "", core.EnforceWarning, "obsolete backport on Python 3", "unittest.mock"},
	{"pip", "enum34", false, "", core.EnforceWarning, "obsolete backport on Python 3.4+", "stdlib enum"},
	{"pip", "futures", false, "", core.EnforceWarning, "obsolete backport on Python 3", "concurrent.futures"},
	{"pip", "beautifulsoup", false, "", core.EnforceWarning, "superseded by beautifulsoup4", "beautifulsoup4"},
	{"pip", "pyyaml", false, "5.4", core.EnforceStrictBlock, "PyYAML < 5.4 has CVE-2020-14343", "pyyaml >= 5.4"},
	// cargo
	{"cargo", "failure", false, "", core.EnforceWarning, "deprecated", "anyhow / thiserror"},
	{"cargo", "rustc-serialize", false, "", core.EnforceWarning, "deprecated", "serde"},
	// composer
	{"composer", "zendframework/", true, "", core.EnforceWarning, "Zend Framework migrated to Laminas", "laminas/*"},
	// gem
	{"gem", "sass", false, "", core.EnforceWarning, "Ruby Sass is deprecated", "sassc / dartsass-rails"},
	{"gem", "therubyracer", false, "", core.EnforceWarning, "deprecated", "mini_racer"},
}

// manifestParsers map manifest basenames to their dependency parser.
var manifestParsers = map[string]func(data []byte) []depDecl{
	"go.mod":           parseGoMod,
	"package.json":     parsePackageJSON,
	"requirements.txt": parseRequirements,
	"pyproject.toml":   parsePyproject,
	"Cargo.toml":       parseCargo,
	"composer.json":    parseComposer,
	"Gemfile":          parseGemfile,
}

// FindManifests locates dependency manifests under root, honoring the same
// directory exclusions as the file scanner.
func FindManifests(root string) []string {
	var found []string
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if skippedDirs[info.Name()] || strings.HasPrefix(info.Name(), ".") && info.Name() != "." {
				return filepath.SkipDir
			}
			return nil
		}
		if _, ok := manifestParsers[info.Name()]; ok {
			found = append(found, path)
		}
		return nil
	})
	return found
}

// EvaluateDependencies checks dependency manifests against the bundled
// deprecated-package DB and, when online is true, shells out to the
// ecosystem package manager for freshness data. Returns violations grouped
// by manifest path (relative to root).
func (pe *PolicyEngine) EvaluateDependencies(root string, online bool) map[string][]EvaluationViolation {
	out := map[string][]EvaluationViolation{}
	deprecatedOK := policyByID(pe.policies, "deprecated-packages") != nil
	freshOK := policyByID(pe.policies, "dependency-freshness") != nil

	for _, mpath := range FindManifests(root) {
		parse := manifestParsers[filepath.Base(mpath)]
		data, err := os.ReadFile(mpath)
		if err != nil {
			continue
		}
		rel, _ := filepath.Rel(root, mpath)
		deps := parse(data)

		if deprecatedOK {
			for _, d := range deps {
				for _, adv := range deprecatedPackages {
					if !adv.matches(d) {
						continue
					}
					msg := adv.Reason + " — use " + adv.Replacement + " instead."
					out[rel] = append(out[rel], EvaluationViolation{
						PolicyID:    "07-deprecated-packages",
						PolicyName:  "Deprecated Package Guard",
						Enforcement: adv.Severity,
						RuleName:    "no-deprecated-packages",
						Message:     messagef("%s %s is flagged: %s", d.Name, d.Version, msg),
					})
				}
			}
		}
	}

	if freshOK && online {
		for mpath, vls := range pmFreshness(root) {
			out[mpath] = append(out[mpath], vls...)
		}
	}
	return out
}

// matches reports whether the advisory applies to the declared dependency.
func (adv DeprecatedPackage) matches(d depDecl) bool {
	if d.Eco != adv.Eco {
		return false
	}
	name := strings.ToLower(d.Name)
	target := strings.ToLower(adv.Name)
	if adv.Prefix {
		if !strings.HasPrefix(name, target) {
			return false
		}
	} else if name != target {
		return false
	}
	if adv.Below != "" {
		ver := normalizeVersion(d.Version)
		if ver == "" || !versionLess(ver, adv.Below) {
			return false
		}
	}
	return true
}

// normalizeVersion extracts a leading semver-ish "X.Y.Z" from a constraint
// like "^4.17.21" or ">=1.2,<2".
var semverLeading = regexp.MustCompile(`(\d+(?:\.\d+){0,3})`)

func normalizeVersion(raw string) string {
	m := semverLeading.FindString(raw)
	return m
}

// versionLess compares dotted numeric versions ("1.2.3" < "5.4").
func versionLess(a, b string) bool {
	as := strings.Split(a, ".")
	bs := strings.Split(b, ".")
	for i := 0; i < len(as) || i < len(bs); i++ {
		var av, bv int
		if i < len(as) {
			av, _ = strconv.Atoi(as[i])
		}
		if i < len(bs) {
			bv, _ = strconv.Atoi(bs[i])
		}
		if av != bv {
			return av < bv
		}
	}
	return false
}

// --- manifest parsers ---

// parseGoMod reads `require` directives (single-line and block form).
func parseGoMod(data []byte) []depDecl {
	var deps []depDecl
	inBlock := false
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "//") || line == "" {
			continue
		}
		if strings.HasPrefix(line, "require (") {
			inBlock = true
			continue
		}
		if inBlock && line == ")" {
			inBlock = false
			continue
		}
		fields := strings.Fields(line)
		if inBlock && len(fields) >= 2 {
			deps = append(deps, depDecl{fields[0], fields[1], "go"})
		} else if strings.HasPrefix(line, "require ") && len(fields) >= 3 {
			deps = append(deps, depDecl{fields[1], fields[2], "go"})
		}
	}
	return deps
}

// parsePackageJSON reads dependencies/devDependencies/optionalDependencies.
func parsePackageJSON(data []byte) []depDecl {
	var doc struct {
		Deps     map[string]string `json:"dependencies"`
		DevDeps  map[string]string `json:"devDependencies"`
		OptDeps  map[string]string `json:"optionalDependencies"`
		PeerDeps map[string]string `json:"peerDependencies"`
	}
	if json.Unmarshal(data, &doc) != nil {
		return nil
	}
	var deps []depDecl
	for _, m := range []map[string]string{doc.Deps, doc.DevDeps, doc.OptDeps, doc.PeerDeps} {
		for name, ver := range m {
			deps = append(deps, depDecl{name, ver, "npm"})
		}
	}
	return deps
}

// requirementLine matches pip requirement lines: name[extra]==1.2 / >=1.2.
var requirementLine = regexp.MustCompile(`^\s*([A-Za-z0-9_\-.]+)\s*(\[[^\]]*\])?\s*(==|>=|<=|~=|!=|>|<)?\s*([A-Za-z0-9_.\-*!+]+)?`)

func parseRequirements(data []byte) []depDecl {
	var deps []depDecl
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "-") {
			continue
		}
		if m := requirementLine.FindStringSubmatch(line); m != nil && m[1] != "" {
			deps = append(deps, depDecl{strings.ToLower(m[1]), m[3] + m[4], "pip"})
		}
	}
	return deps
}

// parsePyproject handles project.dependencies items and poetry-style
// `name = "^x"` tables (heuristic — TOML lite).
var (
	pyprojectDepItem  = regexp.MustCompile(`^\s*["']([A-Za-z0-9_\-.]+)\s*[<>=!~]`)
	poetryDepLine     = regexp.MustCompile(`^\s*([A-Za-z0-9_\-.]+)\s*=\s*["']([^{]*)`)
	poetryInlineTable = regexp.MustCompile(`^\s*([A-Za-z0-9_\-.]+)\s*=\s*\{\s*version\s*=\s*["']([^"']+)`)
)

func parsePyproject(data []byte) []depDecl {
	var deps []depDecl
	for _, line := range strings.Split(string(data), "\n") {
		if m := pyprojectDepItem.FindStringSubmatch(line); m != nil {
			deps = append(deps, depDecl{strings.ToLower(m[1]), "", "pip"})
		} else if m := poetryInlineTable.FindStringSubmatch(line); m != nil {
			deps = append(deps, depDecl{strings.ToLower(m[1]), m[2], "pip"})
		} else if m := poetryDepLine.FindStringSubmatch(line); m != nil {
			name := strings.ToLower(m[1])
			if name != "python" && name != "name" && name != "version" {
				deps = append(deps, depDecl{name, m[2], "pip"})
			}
		}
	}
	return deps
}

// parseCargo reads [dependencies] / [dev-dependencies] `name = "x"` entries.
var (
	cargoSection = regexp.MustCompile(`^\s*\[(?:dev-|build-)?dependencies\]`)
	cargoOther   = regexp.MustCompile(`^\s*\[`)
	cargoEntry   = regexp.MustCompile(`^\s*([A-Za-z0-9_\-]+)\s*=\s*(?:["']([^"']+)["']|\{\s*version\s*=\s*["']([^"']+))`)
)

func parseCargo(data []byte) []depDecl {
	var deps []depDecl
	inDeps := false
	for _, line := range strings.Split(string(data), "\n") {
		if cargoOther.MatchString(line) {
			inDeps = cargoSection.MatchString(line)
			continue
		}
		if !inDeps {
			continue
		}
		if m := cargoEntry.FindStringSubmatch(line); m != nil {
			ver := m[2]
			if ver == "" {
				ver = m[3]
			}
			deps = append(deps, depDecl{m[1], ver, "cargo"})
		}
	}
	return deps
}

// parseComposer reads require / require-dev objects.
func parseComposer(data []byte) []depDecl {
	var doc struct {
		Require    map[string]string `json:"require"`
		RequireDev map[string]string `json:"require-dev"`
	}
	if json.Unmarshal(data, &doc) != nil {
		return nil
	}
	var deps []depDecl
	for _, m := range []map[string]string{doc.Require, doc.RequireDev} {
		for name, ver := range m {
			if name == "php" || strings.HasPrefix(name, "ext-") {
				continue
			}
			deps = append(deps, depDecl{strings.ToLower(name), ver, "composer"})
		}
	}
	return deps
}

// gemLine matches `gem 'name', '~> 1.2'`.
var gemLine = regexp.MustCompile(`^\s*gem\s+["']([^"']+)["']\s*(?:,\s*["']([^"']+)["'])?`)

func parseGemfile(data []byte) []depDecl {
	var deps []depDecl
	for _, line := range strings.Split(string(data), "\n") {
		if m := gemLine.FindStringSubmatch(line); m != nil {
			deps = append(deps, depDecl{strings.ToLower(m[1]), m[2], "gem"})
		}
	}
	return deps
}

// --- package-manager freshness (online, best-effort) ---

// pmFreshness shells out to the ecosystem package manager for outdated-
// dependency data. Failures (no binary, no network, nonzero exit without
// output) produce an info-level note, never a hard error.
func pmFreshness(root string) map[string][]EvaluationViolation {
	out := map[string][]EvaluationViolation{}
	for _, mpath := range FindManifests(root) {
		dir := filepath.Dir(mpath)
		rel, _ := filepath.Rel(root, mpath)
		switch filepath.Base(mpath) {
		case "package.json":
			out[rel] = append(out[rel], npmOutdated(dir)...)
		case "go.mod":
			out[rel] = append(out[rel], goOutdated(dir)...)
		}
	}
	return out
}

func pmNote(manifest, eco, detail string) EvaluationViolation {
	return EvaluationViolation{
		PolicyID:    "10-dependency-freshness",
		PolicyName:  "Dependency Freshness",
		Enforcement: core.EnforceInfo,
		RuleName:    "freshness-check",
		Message:     messagef("%s freshness check skipped for %s: %s", eco, manifest, detail),
	}
}

func pmWarning(eco, name, current, latest string) EvaluationViolation {
	return EvaluationViolation{
		PolicyID:    "10-dependency-freshness",
		PolicyName:  "Dependency Freshness",
		Enforcement: core.EnforceWarning,
		RuleName:    "outdated-dependency",
		Message:     messagef("%s is outdated: %s -> %s (npm update / go get -u to upgrade)", name, current, latest),
	}
}

// npmOutdated parses `npm outdated --json` (exit 1 is normal when
// packages are outdated).
func npmOutdated(dir string) []EvaluationViolation {
	if _, err := exec.LookPath("npm"); err != nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "npm", "outdated", "--json")
	cmd.Dir = dir
	out, err := cmd.Output()
	if len(bytes.TrimSpace(out)) == 0 {
		if err != nil {
			return []EvaluationViolation{pmNote(dir, "npm", "registry unreachable or npm error")}
		}
		return nil
	}
	var rows map[string]struct {
		Current string `json:"current"`
		Latest  string `json:"latest"`
	}
	if json.Unmarshal(out, &rows) != nil {
		return nil
	}
	var vls []EvaluationViolation
	for name, r := range rows {
		vls = append(vls, pmWarning("npm", name, r.Current, r.Latest))
	}
	return vls
}

// goOutdated parses `go list -m -u` update lines.
func goOutdated(dir string) []EvaluationViolation {
	if _, err := exec.LookPath("go"); err != nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", "list", "-m", "-u",
		"-f", "{{if not .Indirect}}{{with .Update}}{{$.Path}}|{{$.Version}}|{{.Version}}{{end}}{{end}}",
		"all")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return []EvaluationViolation{pmNote(dir, "go", "module fetch failed (offline?)")}
	}
	var vls []EvaluationViolation
	for _, line := range strings.Split(string(out), "\n") {
		parts := strings.Split(strings.TrimSpace(line), "|")
		if len(parts) != 3 {
			continue
		}
		v := pmWarning("go", parts[0], parts[1], parts[2])
		v.Message = messagef("%s has an update: %s -> %s (go get -u %s)", parts[0], parts[1], parts[2], parts[0])
		vls = append(vls, v)
	}
	return vls
}
