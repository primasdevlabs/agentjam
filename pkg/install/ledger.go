package install

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// LedgerFile records every file `agentjam install` wrote, per harness, so
// `agentjam uninstall` removes exactly what it created — never user files.
const LedgerFile = ".agentjam/installed.json"

// Ledger maps harness name -> relative paths written for that harness.
type Ledger struct {
	Version   string              `json:"version"`
	Harnesses map[string][]string `json:"harnesses"`
}

// WriteLedger records files written per harness under root. Entries are
// merged with any existing ledger so repeated installs accumulate.
func WriteLedger(root string, harness string, files []string) error {
	l, _ := ReadLedger(root)
	if l.Harnesses == nil {
		l.Harnesses = map[string][]string{}
	}
	seen := map[string]bool{}
	merged := append(l.Harnesses[harness], files...)
	var out []string
	for _, f := range merged {
		if !seen[f] {
			seen[f] = true
			out = append(out, f)
		}
	}
	sort.Strings(out)
	l.Version = "1"
	l.Harnesses[harness] = out
	data, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(root, filepath.Dir(LedgerFile)), 0755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(root, LedgerFile), data, 0644)
}

// ReadLedger loads the install ledger, or returns an empty one.
func ReadLedger(root string) (*Ledger, error) {
	data, err := os.ReadFile(filepath.Join(root, LedgerFile))
	if err != nil {
		return &Ledger{Harnesses: map[string][]string{}}, err
	}
	var l Ledger
	if err := json.Unmarshal(data, &l); err != nil {
		return &Ledger{Harnesses: map[string][]string{}}, err
	}
	if l.Harnesses == nil {
		l.Harnesses = map[string][]string{}
	}
	return &l, nil
}

// UninstallResult summarizes an uninstall pass.
type UninstallResult struct {
	Removed  []string
	Missing  []string
	KeptUser []string
	Purged   []string
}

// harnessOwnedPaths are AgentJam-owned paths removed on uninstall when no
// ledger exists (pre-ledger installs). Keyed by harness; "shared" applies
// to every harness. Single files listed here are 100% AgentJam-generated.
var harnessOwnedPaths = map[string][]string{
	"shared":      {"AGENTJAM.md"},
	"cursor":      {".cursor/rules/skills", ".cursor/rules/policies", ".cursor/rules/agents"},
	"claude-code": {".claude/skills", ".claude/agents", ".claude/policies"},
	"windsurf":    {".windsurf/rules"},
	"cline":       {".clinerules/skills", ".clinerules/policies", ".clinerules/agents"},
	"roo-code":    {".roo/rules"},
	"devin":       {".devin/playbook.md", ".devin/rules"},
	"generic":     {".agentjam/rules"},
	"gemini":      {".gemini/skills", ".gemini/policies", ".gemini/agents"},
}

// rootRuleFiles are top-level harness files install overwrites; removed
// only on --purge since a project may have had its own before install.
var rootRuleFiles = []string{
	".cursorrules", ".windsurfrules", ".clinerules", ".roomodes",
	"CLAUDE.md", "GEMINI.md", "AGENTS.md", "SYSTEM_PROMPT.md",
}

// Uninstall removes AgentJam-installed files from root. With a ledger it
// deletes exactly the recorded files (plus ledger itself); without one it
// removes only wellKnownPaths. purge additionally removes rootRuleFiles
// and the .agentjam directory.
func Uninstall(root string, harnesses []string, purge bool) (*UninstallResult, error) {
	res := &UninstallResult{}
	l, lerr := ReadLedger(root)

	targets := map[string]bool{}
	wantAll := len(harnesses) == 0
	if len(l.Harnesses) > 0 {
		for h, files := range l.Harnesses {
			if !wantAll && !contains(harnesses, h) {
				continue
			}
			for _, f := range files {
				targets[f] = true
			}
		}
		// Keep the shared runbook while other harnesses remain installed.
		if !wantAll {
			othersRemain := false
			for h := range l.Harnesses {
				if !contains(harnesses, h) {
					othersRemain = true
				}
			}
			if othersRemain {
				delete(targets, "AGENTJAM.md")
			}
		}
	} else {
		if wantAll {
			for _, p := range harnessOwnedPaths["shared"] {
				targets[p] = true
			}
		}
		for h, paths := range harnessOwnedPaths {
			if h == "shared" || (!wantAll && !contains(harnesses, h)) {
				continue
			}
			for _, p := range paths {
				targets[p] = true
			}
		}
	}

	for f := range targets {
		full := filepath.Join(root, filepath.FromSlash(f))
		info, err := os.Stat(full)
		if err != nil {
			res.Missing = append(res.Missing, f)
			continue
		}
		if info.IsDir() {
			if err := os.RemoveAll(full); err != nil {
				return res, fmt.Errorf("uninstall %s: %w", f, err)
			}
		} else if err := os.Remove(full); err != nil {
			return res, fmt.Errorf("uninstall %s: %w", f, err)
		}
		res.Removed = append(res.Removed, f)
		pruneEmptyDirs(root, filepath.Dir(f))
	}

	if purge {
		for _, f := range append(append([]string{}, rootRuleFiles...), ".agentjam") {
			full := filepath.Join(root, filepath.FromSlash(f))
			if _, err := os.Stat(full); err != nil {
				continue
			}
			if err := os.RemoveAll(full); err != nil {
				return res, fmt.Errorf("purge %s: %w", f, err)
			}
			res.Purged = append(res.Purged, f)
		}
	} else if len(l.Harnesses) > 0 {
		// Update ledger: drop uninstalled harnesses; remove ledger file
		// when nothing remains.
		wantAll := len(harnesses) == 0
		for h := range l.Harnesses {
			if wantAll || contains(harnesses, h) {
				delete(l.Harnesses, h)
			}
		}
		if len(l.Harnesses) == 0 {
			_ = os.Remove(filepath.Join(root, LedgerFile))
			pruneEmptyDirs(root, filepath.Dir(LedgerFile))
		} else {
			data, _ := json.MarshalIndent(l, "", "  ")
			_ = os.WriteFile(filepath.Join(root, LedgerFile), data, 0644)
		}
	}

	if lerr != nil && len(res.Removed) == 0 {
		return res, errors.New("nothing installed: no AgentJam ledger or known paths found")
	}
	sort.Strings(res.Removed)
	sort.Strings(res.Purged)
	return res, nil
}

// pruneEmptyDirs removes empty parent directories of a deleted file,
// stopping at root.
func pruneEmptyDirs(root, dir string) {
	for dir != "." && dir != "" && !strings.HasPrefix(dir, "..") {
		full := filepath.Join(root, dir)
		entries, err := os.ReadDir(full)
		if err != nil || len(entries) > 0 {
			return
		}
		_ = os.Remove(full)
		dir = filepath.Dir(dir)
	}
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
