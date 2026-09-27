package policy

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
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
}

// Scan walks root, evaluates each source file against the loaded policies,
// and returns an aggregate report. exit-fatal semantics are the caller's.
func (pe *PolicyEngine) Scan(root string, extensions []string, verbose bool) (*ScanReport, error) {
	if len(extensions) == 0 {
		extensions = DefaultScanExtensions
	}
	exts := map[string]bool{}
	for _, e := range extensions {
		exts[strings.ToLower(strings.TrimSpace(e))] = true
	}

	report := &ScanReport{Root: root, Allowed: true}
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
		if !exts[strings.ToLower(filepath.Ext(path))] {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			report.FilesSkipped++
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		summary := pe.EvaluateAll(string(data))
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
		if verbose || !summary.Allowed {
			report.Files = append(report.Files, FileResult{Path: rel, Summary: summary})
		}
		return nil
	})
	if err != nil {
		return report, err
	}
	sort.Slice(report.Files, func(i, j int) bool { return report.Files[i].Path < report.Files[j].Path })
	return report, nil
}
