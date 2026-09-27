package policy

import (
	"bufio"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/primasdevlabs/agentjam/pkg/core"
)

// Structural policies — god-file limits and domain-driven layout — evaluated
// by Scan on top of the content evaluators.

// topLevelDecl matches common top-level declarations across languages:
// func, class, type/struct/interface, export const/let, def, public/private
// methods at column 0. Heuristic — line-anchored so nested decls don't count.
var topLevelDecl = regexp.MustCompile(`^(export\s+(default\s+)?(async\s+)?(func|function|class|const|let|var|type|interface|struct|enum)|func\s|class\s|def\s|public\s|private\s|protected\s|internal\s|fn\s|type\s|interface\s|struct\s|enum\s|impl\s|module\s|package\s)`)

// responsibilityMarkers detect mixed-concern files: a file touching UI markup,
// SQL/data access, and HTTP orchestration simultaneously is likely a god file.
var responsibilityMarkers = map[string]*regexp.Regexp{
	"ui":          regexp.MustCompile(`(?i)(className=|<div|<button|render\(|fmt\.Print|console\.log|print\()`),
	"data":        regexp.MustCompile(`(?i)(SELECT\s|INSERT\s|UPDATE\s|DELETE\s|\.Query\(|\.Exec\(|db\.|repository)`),
	"http":        regexp.MustCompile(`(?i)(http\.(Get|Post|Handle|ListenAndServe)|fetch\(|axios|app\.(get|post|put|delete)\(|router\.)`),
	"io":          regexp.MustCompile(`(?i)(os\.(ReadFile|WriteFile|Open)|fs\.|File\(|fopen|open\()`),
	"concurrency": regexp.MustCompile(`(?i)(go\s+func|goroutine|thread|async\s|await\s|Promise|Mutex|Channel)`),
	"validation":  regexp.MustCompile(`(?i)(valid|sanitize|escape|check\(|assert)`),
}

func ruleInt(p *core.PolicyManifest, key string, fallback int) int {
	if p == nil || p.Rules == nil {
		return fallback
	}
	if v, ok := p.Rules[key]; ok {
		switch n := v.(type) {
		case int:
			return n
		case int64:
			return int(n)
		case float64:
			return int(n)
		}
	}
	return fallback
}

func ruleFloat(p *core.PolicyManifest, key string, fallback float64) float64 {
	if p == nil || p.Rules == nil {
		return fallback
	}
	if v, ok := p.Rules[key]; ok {
		if f, ok := v.(float64); ok {
			return f
		}
	}
	return fallback
}

// policyByID finds a loaded policy manifest by ID suffix match (e.g.
// "god-files" matches "06-god-files").
func policyByID(policies []core.PolicyManifest, suffix string) *core.PolicyManifest {
	for i := range policies {
		if policies[i].ID == suffix || strings.HasSuffix(policies[i].ID, "-"+suffix) {
			return &policies[i]
		}
	}
	return nil
}

// EvaluateFileStructure checks a single file against the god-file policy.
// relPath is used only for reporting context.
func (pe *PolicyEngine) EvaluateFileStructure(relPath, content string) []EvaluationViolation {
	p := policyByID(pe.policies, "god-files")
	if p == nil {
		return nil
	}
	maxLines := ruleInt(p, "max_file_lines", 500)
	maxDecls := ruleInt(p, "max_top_level_decls", 40)
	maxResp := ruleInt(p, "max_responsibility_markers", 8)

	var violations []EvaluationViolation
	v := func(rule, msg string) {
		violations = append(violations, EvaluationViolation{
			PolicyID: p.ID, PolicyName: p.Name, Enforcement: p.Enforcement,
			RuleName: rule, Message: msg,
		})
	}

	lines := strings.Count(content, "\n") + 1
	if lines > maxLines {
		v("max_file_lines", messagef("%s has %d lines (limit %d) — split by responsibility.", filepath.Base(relPath), lines, maxLines))
	}

	decls := 0
	sc := bufio.NewScanner(strings.NewReader(content))
	for sc.Scan() {
		if topLevelDecl.MatchString(sc.Text()) {
			decls++
		}
	}
	if decls > maxDecls {
		v("max_top_level_decls", messagef("%s declares %d top-level symbols (limit %d) — extract cohesive units.", filepath.Base(relPath), decls, maxDecls))
	}

	respHits := 0
	for _, re := range responsibilityMarkers {
		if re.MatchString(content) {
			respHits++
		}
	}
	if respHits > maxResp {
		v("mixed_responsibilities", messagef("%s mixes %d concern domains (limit %d) — decompose into per-domain units.", filepath.Base(relPath), respHits, maxResp))
	}

	// Function-level segmentation: lines between top-level declarations.
	maxFnLines := ruleInt(p, "max_function_lines", 80)
	maxComplexity := ruleInt(p, "max_complexity_score", 12)
	type seg struct{ lines, complexity int }
	var segs []seg
	cur := seg{}
	for sc := bufio.NewScanner(strings.NewReader(content)); sc.Scan(); {
		line := sc.Text()
		if topLevelDecl.MatchString(line) {
			segs = append(segs, cur)
			cur = seg{}
		}
		cur.lines++
		cur.complexity += len(branchKeywords.FindAllString(line, -1))
	}
	segs = append(segs, cur)

	var longest seg
	for _, s := range segs {
		if s.lines > longest.lines {
			longest = s
		}
		if s.lines > maxFnLines {
			v("max_function_lines", messagef("%s contains a ~%d-line unit (limit %d) — extract named helpers.", filepath.Base(relPath), s.lines, maxFnLines))
			break
		}
	}
	if longest.complexity > maxComplexity {
		v("max_complexity", messagef("%s has a unit with branching score %d (limit %d) — decompose conditionals.", filepath.Base(relPath), longest.complexity, maxComplexity))
	}

	return violations
}

var branchKeywords = regexp.MustCompile(`\b(if|elif|else\s+if|for|while|case|when|catch|except|and|or)\b|&&|\|\||\?`)

// EvaluateProjectLayout checks the collected source-file list against the
// domain-layout policy. sourceFiles must be workspace-relative paths.
func (pe *PolicyEngine) EvaluateProjectLayout(root string, sourceFiles []string) []EvaluationViolation {
	p := policyByID(pe.policies, "domain-layout")
	if p == nil || len(sourceFiles) == 0 {
		return nil
	}
	flatThreshold := ruleInt(p, "flat_source_threshold", 12)
	maxPerDir := ruleInt(p, "max_files_per_dir", 15)
	minRatio := ruleFloat(p, "min_domain_ratio", 0.6)

	var violations []EvaluationViolation
	v := func(rule, msg string) {
		violations = append(violations, EvaluationViolation{
			PolicyID: p.ID, PolicyName: p.Name, Enforcement: p.Enforcement,
			RuleName: rule, Message: msg,
		})
	}

	// Flat ratio: files at the scan root (no domain subdirectory).
	flat := 0
	perDir := map[string]int{}
	for _, f := range sourceFiles {
		dir := filepath.ToSlash(filepath.Dir(f))
		perDir[dir]++
		if dir == "." || dir == "" {
			flat++
		}
	}
	nested := len(sourceFiles) - flat
	if len(sourceFiles) > flatThreshold && float64(nested)/float64(len(sourceFiles)) < minRatio {
		v("flat_layout", messagef("Only %d/%d source files live in domain subdirectories (need %.0f%%) — group by bounded context.", nested, len(sourceFiles), minRatio*100))
	}

	for dir, count := range perDir {
		if count > maxPerDir {
			v("oversized_directory", messagef("Directory %q holds %d source files (limit %d) — split into narrower domains.", dir, count, maxPerDir))
		}
	}
	return violations
}

func messagef(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
