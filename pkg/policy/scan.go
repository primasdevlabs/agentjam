package policy

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/primasdevlabs/agentjam/pkg/core"
)

// DefaultScanExtensions are source file types scanned by `agentjam scan`.
var DefaultScanExtensions = []string{
	".go", ".ts", ".tsx", ".js", ".jsx", ".mjs",
	".py", ".rb", ".php", ".java", ".cs", ".kt", ".rs", ".swift",
	".html", ".htm", ".vue", ".svelte", ".css", ".scss", ".less",
	".sql", ".sh", ".ps1", ".yaml", ".yml", ".json",
}

// skippedDirs are never traversed by the scanner.
var skippedDirs = map[string]bool{
	".git": true, "node_modules": true, "vendor": true, "dist": true,
	"build": true, "out": true, ".next": true, ".nuxt": true,
	"target": true, "bin": true, "obj": true, "coverage": true,
}

// FileResult holds violations found in one file.
type FileResult struct {
	Path       string
	Summary    PolicyEngineSummary
	Skipped    bool
	SkipReason string
}

// ScanReport aggregates a project scan.
type ScanReport struct {
	Root          string       `json:"root"`
	FilesScanned  int          `json:"filesScanned"`
	FilesSkipped  int          `json:"filesSkipped"`
	FilesWithHits int          `json:"filesWithViolations"`
	Allowed       bool         `json:"allowed"`
	StrictBlocks  int          `json:"strictBlocks"`
	Warnings      int          `json:"warnings"`
	InfoCount     int          `json:"infoCount"`
	Files         []FileResult `json:"files,omitempty"`
	// ProjectViolations are structural findings that span the whole tree —
	// e.g. flat non-domain layout — rather than a single file.
	ProjectViolations []EvaluationViolation `json:"projectViolations,omitempty"`
}

// Scan walks root, evaluates each source file against the loaded policies,
// and returns an aggregate report. exit-fatal semantics are the caller's.
func (pe *PolicyEngine) Scan(root string, extensions []string, verbose, online bool) (*ScanReport, error) {
	if len(extensions) == 0 {
		extensions = DefaultScanExtensions
	}
	exts := map[string]bool{}
	for _, e := range extensions {
		exts[strings.ToLower(strings.TrimSpace(e))] = true
	}

	report := &ScanReport{Root: root, Allowed: true}
	var sourceFiles []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // unreadable entries are skipped, not fatal
		}
		if info.IsDir() {
			if skippedDirs[info.Name()] || strings.HasPrefix(info.Name(), ".") && info.Name() != "." {
				return filepath.SkipDir
			}
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		// Filename-level guard: committed secrets/keys are violations on any
		// file, regardless of extension.
		for _, sv := range pe.EvaluateSensitivePath(rel) {
			report.ProjectViolations = append(report.ProjectViolations, sv)
			report.StrictBlocks++
			report.Allowed = false
		}
		if !exts[strings.ToLower(filepath.Ext(path))] {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			report.FilesSkipped++
			return nil
		}
		sourceFiles = append(sourceFiles, rel)
		summary := pe.EvaluateAll(string(data))
		// Structural pass: god-file limits against the architecture policies.
		summary.Violations = append(summary.Violations, pe.EvaluateFileStructure(rel, string(data))...)
		summary.TotalViolations = len(summary.Violations)
		summary.StrictBlocks, summary.Warnings, summary.InfoCount, summary.Allowed = 0, 0, 0, true
		for _, v := range summary.Violations {
			switch v.Enforcement {
			case core.EnforceStrictBlock:
				summary.StrictBlocks++
				summary.Allowed = false
			case core.EnforceWarning:
				summary.Warnings++
			default:
				summary.InfoCount++
			}
		}
		report.FilesScanned++
		if summary.TotalViolations == 0 {
			return nil
		}
		report.FilesWithHits++
		report.StrictBlocks += summary.StrictBlocks
		report.Warnings += summary.Warnings
		report.InfoCount += summary.InfoCount
		if !summary.Allowed {
			report.Allowed = false
		}
		if summary.TotalViolations > 0 {
			report.Files = append(report.Files, FileResult{Path: rel, Summary: summary})
		}
		return nil
	})
	if err != nil {
		return report, err
	}

	// Project-level structural pass: domain layout against the architecture policies.
	for _, pv := range pe.EvaluateProjectLayout(root, sourceFiles) {
		report.ProjectViolations = append(report.ProjectViolations, pv)
		countSeverity(report, pv)
	}

	// Dependency pass: deprecated-package advisories and (when online)
	// package-manager freshness, grouped per manifest file.
	for rel, vls := range pe.EvaluateDependencies(root, online) {
		merged := false
		for i := range report.Files {
			if report.Files[i].Path == rel {
				report.Files[i].Summary.Violations = append(report.Files[i].Summary.Violations, vls...)
				report.Files[i].Summary.TotalViolations = len(report.Files[i].Summary.Violations)
				merged = true
				break
			}
		}
		if !merged {
			report.Files = append(report.Files, FileResult{Path: rel, Summary: PolicyEngineSummary{
				Violations:      vls,
				TotalViolations: len(vls),
			}})
			report.FilesWithHits++
		}
		for _, v := range vls {
			countSeverity(report, v)
		}
	}

	sort.Slice(report.Files, func(i, j int) bool { return report.Files[i].Path < report.Files[j].Path })
	return report, nil
}

// countSeverity folds one violation into the report totals.
func countSeverity(report *ScanReport, v EvaluationViolation) {
	switch v.Enforcement {
	case core.EnforceStrictBlock:
		report.StrictBlocks++
		report.Allowed = false
	case core.EnforceWarning:
		report.Warnings++
	default:
		report.InfoCount++
	}
}
