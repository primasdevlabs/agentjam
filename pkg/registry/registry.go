package registry

import (
	"github.com/agentjam/agentjam/pkg/core"
	"github.com/agentjam/agentjam/pkg/parser"
)

// RegistryIndex holds repository resource stats.
type RegistryIndex struct {
	AgentsCount    int `json:"agentsCount"`
	SkillsCount    int `json:"skillsCount"`
	ToolsCount     int `json:"toolsCount"`
	WorkflowsCount int `json:"workflowsCount"`
	StacksCount    int `json:"stacksCount"`
	LanguagesCount int `json:"languagesCount"`
	PoliciesCount  int `json:"policiesCount"`
}

// BuildRegistryIndex scans workspace and returns structured registry metadata.
func BuildRegistryIndex(rootDir string) RegistryIndex {
	discovered := parser.DiscoverResources(rootDir)
	idx := RegistryIndex{}

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
		}
	}

	return idx
}
