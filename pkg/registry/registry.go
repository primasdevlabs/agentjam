package registry

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/primasdevlabs/agentjam/pkg/core"
	"github.com/primasdevlabs/agentjam/pkg/parser"
)

// RegistryEntry describes a single indexed canonical resource.
type RegistryEntry struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Version     string            `json:"version,omitempty"`
	Type        core.ResourceType `json:"type"`
	Description string            `json:"description,omitempty"`
	Path        string            `json:"path"`
}

// RegistryIndex holds repository resource stats and indexed entries.
type RegistryIndex struct {
	GeneratedAt       string          `json:"generatedAt"`
	AgentsCount       int             `json:"agentsCount"`
	SkillsCount       int             `json:"skillsCount"`
	ToolsCount        int             `json:"toolsCount"`
	WorkflowsCount    int             `json:"workflowsCount"`
	StacksCount       int             `json:"stacksCount"`
	LanguagesCount    int             `json:"languagesCount"`
	PoliciesCount     int             `json:"policiesCount"`
	IntegrationsCount int             `json:"integrationsCount"`
	Total             int             `json:"total"`
	Entries           []RegistryEntry `json:"entries"`
}

// manifestSummary extracts name/version/description from a parsed resource.
func manifestSummary(res parser.DiscoveredResource) (name, version, desc string) {
	switch res.Type {
	case core.ResourceTypeAgent:
		if b, err := parser.ParseAgent(res.Path); err == nil {
			return b.Manifest.Name, b.Manifest.Version, b.Manifest.Description
		}
	case core.ResourceTypeSkill:
		if b, err := parser.ParseSkill(res.Path); err == nil {
			return b.Manifest.Name, b.Manifest.Version, b.Manifest.Description
		}
	case core.ResourceTypeWorkflow:
		if b, err := parser.ParseWorkflow(res.Path); err == nil {
			return b.Manifest.Name, b.Manifest.Version, b.Manifest.Description
		}
	case core.ResourceTypeTool:
		if m, err := parser.ParseTool(res.Path); err == nil {
			return m.Name, m.Version, m.Description
		}
	case core.ResourceTypeStack:
		if m, err := parser.ParseStack(res.Path); err == nil {
			return m.Name, m.Version, m.Description
		}
	case core.ResourceTypeLanguage:
		if m, err := parser.ParseLanguage(res.Path); err == nil {
			return m.Name, "", ""
		}
	case core.ResourceTypePolicy:
		if m, err := parser.ParsePolicy(res.Path); err == nil {
			return m.Name, "", m.Description
		}
	case core.ResourceTypeIntegration:
		if m, err := parser.ParseIntegration(res.Path); err == nil {
			return m.Name, "", ""
		}
	}
	return res.ID, "", ""
}

// BuildRegistryIndex scans workspace and returns structured registry metadata.
func BuildRegistryIndex(rootDir string) RegistryIndex {
	discovered := parser.DiscoverResources(rootDir)
	idx := RegistryIndex{
		GeneratedAt: core.TimestampNow(),
		Entries:     make([]RegistryEntry, 0, len(discovered)),
	}

	for _, res := range discovered {
		switch res.Type {
		case core.ResourceTypeAgent:
			idx.AgentsCount++
		case core.ResourceTypeSkill:
			idx.SkillsCount++
		case core.ResourceTypeTool:
			idx.ToolsCount++
		case core.ResourceTypeWorkflow:
			idx.WorkflowsCount++
		case core.ResourceTypeStack:
			idx.StacksCount++
		case core.ResourceTypeLanguage:
			idx.LanguagesCount++
		case core.ResourceTypePolicy:
			idx.PoliciesCount++
		case core.ResourceTypeIntegration:
			idx.IntegrationsCount++
		}

		name, version, desc := manifestSummary(res)
		rel, err := filepath.Rel(rootDir, res.Path)
		if err != nil {
			rel = res.Path
		}
		idx.Entries = append(idx.Entries, RegistryEntry{
			ID:          res.ID,
			Name:        name,
			Version:     version,
			Type:        res.Type,
			Description: desc,
			Path:        filepath.ToSlash(rel),
		})
	}

	sort.Slice(idx.Entries, func(i, j int) bool {
		if idx.Entries[i].Type != idx.Entries[j].Type {
			return idx.Entries[i].Type < idx.Entries[j].Type
		}
		return idx.Entries[i].ID < idx.Entries[j].ID
	})

	idx.Total = len(idx.Entries)
	return idx
}

// SplitByType groups index entries by resource type, keyed on the plural
// registry directory name (e.g. "agents", "skills").
func (idx RegistryIndex) SplitByType() map[string][]RegistryEntry {
	groups := make(map[string][]RegistryEntry)
	for _, e := range idx.Entries {
		t := string(e.Type)
		dir := t + "s"
		if strings.HasSuffix(t, "y") && !strings.HasSuffix(t, "ay") && !strings.HasSuffix(t, "ey") {
			dir = t[:len(t)-1] + "ies"
		}
		groups[dir] = append(groups[dir], e)
	}
	return groups
}
