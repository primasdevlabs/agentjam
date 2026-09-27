// Package install bootstraps AgentJam into arbitrary projects: materializing
// the embedded canonical tree (init) and rendering per-resource harness files
// (skills, agents, policies) in each AI environment's native format.
package install

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/primasdevlabs/agentjam/pkg/core"
	"github.com/primasdevlabs/agentjam/pkg/parser"
)

// InitOptions controls `agentjam init` materialization.
type InitOptions struct {
	Bare  bool     // only .agentjam/config.yaml + top-level skeleton dirs
	Force bool     // overwrite existing files
	Only  []string // subset of canonical top-level dirs (empty = all)
}

// InitResult reports what init wrote.
type InitResult struct {
	Written []string `json:"written"`
	Skipped []string `json:"skipped"`
}

// canonicalDirs are the embedded top-level resource directories.
var canonicalDirs = []string{
	"agents", "skills", "tools", "workflows",
	"policies", "stacks", "languages", "integrations",
	"prompts", "templates",
}

// defaultWorkspaceConfig seeds .agentjam/config.yaml in new workspaces.
const defaultWorkspaceConfig = `# AgentJam workspace configuration.
versionPolicy: current-stable
freshnessRequired: true
maxDocAge: 7d
stack: generic
defaultAgent: software-engineer
`

// Init materializes the embedded canonical tree into targetRoot. Canonical
// resources land under targetRoot/.agentjam/ — keeping the project root
// clean — and every consumer resolves that location via
// parser.CanonicalRoot.
func Init(fsys fs.FS, targetRoot string, opts InitOptions) (*InitResult, error) {
	res := &InitResult{}
	canonRoot := filepath.Join(targetRoot, ".agentjam")
	write := func(rel string, data []byte) error {
		dst := filepath.Join(canonRoot, rel)
		if _, err := os.Stat(dst); err == nil && !opts.Force {
			res.Skipped = append(res.Skipped, rel)
			return nil
		}
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(dst, data, 0o644); err != nil {
			return err
		}
		res.Written = append(res.Written, rel)
		return nil
	}

	if err := write("config.yaml", []byte(defaultWorkspaceConfig)); err != nil {
		return res, err
	}
	if opts.Bare {
		for _, d := range canonicalDirs {
			if err := os.MkdirAll(filepath.Join(canonRoot, d), 0o755); err != nil {
				return res, err
			}
		}
		return res, nil
	}

	dirs := canonicalDirs
	if len(opts.Only) > 0 {
		dirs = opts.Only
	}
	for _, d := range dirs {
		err := fs.WalkDir(fsys, d, func(path string, de fs.DirEntry, err error) error {
			if err != nil || de.IsDir() {
				return err
			}
			data, err := fs.ReadFile(fsys, path)
			if err != nil {
				return err
			}
			return write(filepath.FromSlash(path), data)
		})
		if err != nil {
			return res, fmt.Errorf("materializing %s: %w", d, err)
		}
	}
	return res, nil
}

// ResourceSet holds everything a harness renderer needs. Entries carry the
// resource's domain category so installed trees follow DDD layout:
// skills/<category>/<id>, policies/<category>/<id>, agents/<id>.
type ResourceSet struct {
	Agents   []AgentEntry
	Skills   []SkillEntry
	Policies []PolicyEntry
}

type AgentEntry struct {
	Manifest     core.AgentManifest
	Instructions map[string]string
}

type SkillEntry struct {
	Manifest     core.SkillManifest
	Category     string
	Instructions map[string]string
}

type PolicyEntry struct {
	Manifest core.PolicyManifest
	Category string
}

// categoryOf extracts the domain segment between the type dir and the
// resource id: skills/<cat>/<id>/skill.yaml -> cat, policies/<cat>/<file> ->
// cat. Returns "" for flat layouts (agents/<id>).
func categoryOf(root, resPath, typeDir string) string {
	rel, err := filepath.Rel(root, resPath)
	if err != nil {
		return ""
	}
	rel = filepath.ToSlash(rel)
	parts := strings.Split(rel, "/")
	// expect: <typeDir>/<cat>/<id>/<file> or <typeDir>/<cat>/<file>
	// skills/<cat>/<id>/file needs 4 segments to have a category;
	// policies/<cat>/file needs 3. Anything shallower is flat (no category).
	minDepth := 3
	if typeDir == "skills" {
		minDepth = 4
	}
	for i, p := range parts {
		if p == typeDir && i+1 < len(parts)-1 && len(parts) >= minDepth {
			return parts[i+1]
		}
	}
	return ""
}

// Collect parses the canonical resource tree under root. When root lacks a
// canonical tree, callers should materialize the embedded FS to a temp dir
// and collect from there.
func Collect(root string) (*ResourceSet, error) {
	rs := &ResourceSet{}
	for _, res := range parser.DiscoverResources(root) {
		switch res.Type {
		case core.ResourceTypeAgent:
			if b, err := parser.ParseAgent(res.Path); err == nil {
				rs.Agents = append(rs.Agents, AgentEntry{Manifest: b.Manifest, Instructions: b.Instructions})
			}
		case core.ResourceTypeSkill:
			if b, err := parser.ParseSkill(res.Path); err == nil {
				rs.Skills = append(rs.Skills, SkillEntry{
					Manifest: b.Manifest, Category: categoryOf(root, res.Path, "skills"), Instructions: b.Instructions})
			}
		case core.ResourceTypePolicy:
			if p, err := parser.ParsePolicy(res.Path); err == nil {
				rs.Policies = append(rs.Policies, PolicyEntry{
					Manifest: p, Category: categoryOf(root, res.Path, "policies")})
			}
		}
	}
	sort.Slice(rs.Agents, func(i, j int) bool { return rs.Agents[i].Manifest.Name < rs.Agents[j].Manifest.Name })
	sort.Slice(rs.Skills, func(i, j int) bool { return rs.Skills[i].Manifest.Name < rs.Skills[j].Manifest.Name })
	sort.Slice(rs.Policies, func(i, j int) bool { return rs.Policies[i].Manifest.Name < rs.Policies[j].Manifest.Name })
	return rs, nil
}

// instructionsText concatenates an instruction map into one markdown body.
func instructionsText(instructions map[string]string) string {
	names := make([]string, 0, len(instructions))
	for n := range instructions {
		names = append(names, n)
	}
	sort.Strings(names)
	var sb strings.Builder
	for _, n := range names {
		sb.WriteString(strings.TrimSpace(instructions[n]))
		sb.WriteString("\n\n")
	}
	return strings.TrimSpace(sb.String())
}
