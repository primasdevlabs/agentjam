package parser

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/primasdevlabs/agentjam/pkg/core"
)

// DiscoveredResource represents a discovered file resource in workspace.
type DiscoveredResource struct {
	ID   string            `json:"id"`
	Type core.ResourceType `json:"type"`
	Path string            `json:"path"`
}

// ResourceBundle represents parsed resource with instruction files.
type ResourceBundle[T any] struct {
	Manifest     T                 `json:"manifest"`
	Instructions map[string]string `json:"instructions"`
	BasePath     string            `json:"basePath"`
}

// skippedDirs are directory names excluded from resource discovery.
var skippedDirs = map[string]bool{
	".git": true, ".agentjam": true, ".github": true, ".cursor": true,
	"node_modules": true, "vendor": true, "pkg": true, "cmd": true,
	"docs": true, "runtime": true, "registry": true,
}

// ParseFrontmatter splits Markdown frontmatter (YAML/JSON block) from body.
func ParseFrontmatter(content string) (string, string) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	var headLines []string
	var bodyLines []string

	inFrontmatter := false
	delimiterCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "---" {
			delimiterCount++
			if delimiterCount == 1 {
				inFrontmatter = true
				continue
			} else if delimiterCount == 2 {
				inFrontmatter = false
				continue
			}
		}

		if inFrontmatter {
			headLines = append(headLines, line)
		} else if delimiterCount >= 2 || delimiterCount == 0 {
			bodyLines = append(bodyLines, line)
		}
	}

	return strings.Join(headLines, "\n"), strings.Join(bodyLines, "\n")
}

// UnmarshalManifest decodes a manifest file. YAML is a superset of JSON, so a
// single YAML decoder handles both formats.
func UnmarshalManifest(data []byte, out interface{}) error {
	if err := yaml.Unmarshal(data, out); err != nil {
		return fmt.Errorf("%w: %v", core.ErrInvalidManifest, err)
	}
	return nil
}

// classifyResource determines the canonical resource type for a file based on
// its directory context and file name. Returns "" for unrecognized files.
func classifyResource(rootDir, path, base string) (core.ResourceType, string) {
	rel, err := filepath.Rel(rootDir, path)
	if err != nil {
		return "", ""
	}
	segments := strings.Split(filepath.ToSlash(rel), "/")
	top := segments[0]
	parentDir := filepath.Base(filepath.Dir(path))
	id := parentDir
	ext := filepath.Ext(base)
	stem := strings.TrimSuffix(base, ext)

	switch top {
	case "policies":
		if ext == ".yaml" || ext == ".yml" || ext == ".md" {
			return core.ResourceTypePolicy, stem
		}
		return "", ""
	case "prompts":
		if ext == ".md" {
			return core.ResourceTypePrompt, parentDir + "/" + stem
		}
		return "", ""
	case "templates":
		if stem == "agent" || stem == "manifest" {
			return core.ResourceTypeTemplate, parentDir + "/" + stem
		}
		if stem == "skill" || stem == "workflow" || stem == "tool" {
			return core.ResourceTypeTemplate, parentDir + "/" + stem
		}
		return "", ""
	}

	switch stem {
	case "agent":
		if ext == ".yaml" || ext == ".yml" {
			return core.ResourceTypeAgent, id
		}
	case "skill":
		if ext == ".yaml" || ext == ".yml" {
			return core.ResourceTypeSkill, id
		}
	case "SKILL":
		if ext == ".md" {
			return core.ResourceTypeSkill, id
		}
	case "workflow":
		if ext == ".yaml" || ext == ".yml" {
			return core.ResourceTypeWorkflow, id
		}
	case "stack":
		if ext == ".yaml" || ext == ".yml" {
			return core.ResourceTypeStack, id
		}
	case "language":
		if ext == ".yaml" || ext == ".yml" {
			return core.ResourceTypeLanguage, id
		}
	case "tool":
		if ext == ".yaml" || ext == ".yml" {
			return core.ResourceTypeTool, id
		}
	case "integration":
		if ext == ".yaml" || ext == ".yml" {
			return core.ResourceTypeIntegration, id
		}
	case "manifest":
		if ext != ".yaml" && ext != ".yml" {
			return "", ""
		}
		switch top {
		case "integrations":
			return core.ResourceTypeIntegration, id
		case "agents":
			return core.ResourceTypeAgent, id
		default:
			return "", ""
		}
	}

	return "", ""
}

// DiscoverResources recursively scans root directory for agents, skills, tools,
// workflows, policies, stacks, languages, integrations, prompts, and templates.
func DiscoverResources(rootDir string) []DiscoveredResource {
	results := make([]DiscoveredResource, 0)

	_ = filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if skippedDirs[info.Name()] && path != rootDir {
				return filepath.SkipDir
			}
			if strings.HasPrefix(info.Name(), ".") && path != rootDir {
				return filepath.SkipDir
			}
			return nil
		}

		resType, id := classifyResource(rootDir, path, info.Name())
		if resType != "" {
			results = append(results, DiscoveredResource{
				ID:   id,
				Type: resType,
				Path: path,
			})
		}
		return nil
	})

	return results
}

// ReadInstructionFiles reads all markdown instruction files in directory.
func ReadInstructionFiles(dir string) map[string]string {
	instructions := make(map[string]string)
	instDir := filepath.Join(dir, "instructions")

	if _, err := os.Stat(instDir); err == nil {
		_ = filepath.Walk(instDir, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
				rel, _ := filepath.Rel(instDir, path)
				content, err := os.ReadFile(path)
				if err == nil {
					instructions[filepath.ToSlash(rel)] = string(content)
				}
			}
			return nil
		})
	}

	return instructions
}

// ParseAgent parses an agent resource bundle (agent.yaml / manifest.yaml).
func ParseAgent(manifestPath string) (ResourceBundle[core.AgentManifest], error) {
	dir := filepath.Dir(manifestPath)
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return ResourceBundle[core.AgentManifest]{}, err
	}

	var manifest core.AgentManifest
	if err := UnmarshalManifest(data, &manifest); err != nil {
		return ResourceBundle[core.AgentManifest]{}, err
	}
	if manifest.Name == "" {
		manifest.Name = filepath.Base(dir)
	}
	if manifest.Type == "" {
		manifest.Type = string(core.ResourceTypeAgent)
	}

	mergeAgentOverrides(dir, &manifest)

	return ResourceBundle[core.AgentManifest]{
		Manifest:     manifest,
		Instructions: ReadInstructionFiles(dir),
		BasePath:     dir,
	}, nil
}

// mergeAgentOverrides folds the optional skills.yaml / tools.yaml override
// files into the manifest's constraint maps.
func mergeAgentOverrides(dir string, manifest *core.AgentManifest) {
	var skillsFile struct {
		Skills []struct {
			Name    string `yaml:"name"`
			Version string `yaml:"version"`
		} `yaml:"skills"`
	}
	if data, err := os.ReadFile(filepath.Join(dir, "skills.yaml")); err == nil {
		if yaml.Unmarshal(data, &skillsFile) == nil && len(skillsFile.Skills) > 0 {
			manifest.SkillVersions = make(map[string]string, len(skillsFile.Skills))
			for _, s := range skillsFile.Skills {
				manifest.SkillVersions[s.Name] = s.Version
				if !containsStr(manifest.Skills, s.Name) {
					manifest.Skills = append(manifest.Skills, s.Name)
				}
			}
		}
	}

	var toolsFile struct {
		Tools []struct {
			Name         string   `yaml:"name"`
			Capabilities []string `yaml:"capabilities"`
		} `yaml:"tools"`
	}
	if data, err := os.ReadFile(filepath.Join(dir, "tools.yaml")); err == nil {
		if yaml.Unmarshal(data, &toolsFile) == nil && len(toolsFile.Tools) > 0 {
			manifest.ToolCapabilities = make(map[string][]string, len(toolsFile.Tools))
			for _, t := range toolsFile.Tools {
				manifest.ToolCapabilities[t.Name] = t.Capabilities
				if !containsStr(manifest.Tools, t.Name) {
					manifest.Tools = append(manifest.Tools, t.Name)
				}
			}
		}
	}
}

func containsStr(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

// ParseSkill parses a skill resource bundle (skill.yaml or SKILL.md).
func ParseSkill(manifestPath string) (ResourceBundle[core.SkillManifest], error) {
	dir := filepath.Dir(manifestPath)
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return ResourceBundle[core.SkillManifest]{}, err
	}

	var manifest core.SkillManifest
	if strings.HasSuffix(manifestPath, ".md") {
		fm, body := ParseFrontmatter(string(data))
		if err := UnmarshalManifest([]byte(fm), &manifest); err != nil && fm != "" {
			return ResourceBundle[core.SkillManifest]{}, err
		}
		if manifest.Description == "" {
			manifest.Description = strings.TrimSpace(body)
		}
	} else {
		if err := UnmarshalManifest(data, &manifest); err != nil {
			return ResourceBundle[core.SkillManifest]{}, err
		}
	}
	if manifest.Name == "" {
		manifest.Name = filepath.Base(dir)
	}
	if manifest.Type == "" {
		manifest.Type = string(core.ResourceTypeSkill)
	}

	return ResourceBundle[core.SkillManifest]{
		Manifest:     manifest,
		Instructions: ReadInstructionFiles(dir),
		BasePath:     dir,
	}, nil
}

// ParseWorkflow parses a workflow resource bundle.
func ParseWorkflow(manifestPath string) (ResourceBundle[core.WorkflowManifest], error) {
	dir := filepath.Dir(manifestPath)
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return ResourceBundle[core.WorkflowManifest]{}, err
	}

	var manifest core.WorkflowManifest
	if err := UnmarshalManifest(data, &manifest); err != nil {
		return ResourceBundle[core.WorkflowManifest]{}, err
	}
	if manifest.Name == "" {
		manifest.Name = filepath.Base(dir)
	}
	if manifest.Type == "" {
		manifest.Type = string(core.ResourceTypeWorkflow)
	}

	return ResourceBundle[core.WorkflowManifest]{
		Manifest:     manifest,
		Instructions: ReadInstructionFiles(dir),
		BasePath:     dir,
	}, nil
}

// ParseTool parses a tool manifest.
func ParseTool(manifestPath string) (core.ToolManifest, error) {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return core.ToolManifest{}, err
	}
	var manifest core.ToolManifest
	if err := UnmarshalManifest(data, &manifest); err != nil {
		return core.ToolManifest{}, err
	}
	if manifest.Name == "" {
		manifest.Name = filepath.Base(filepath.Dir(manifestPath))
	}
	return manifest, nil
}

// ParseStack parses a stack profile manifest.
func ParseStack(manifestPath string) (core.StackManifest, error) {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return core.StackManifest{}, err
	}
	var manifest core.StackManifest
	if err := UnmarshalManifest(data, &manifest); err != nil {
		return core.StackManifest{}, err
	}
	if manifest.Name == "" {
		manifest.Name = filepath.Base(filepath.Dir(manifestPath))
	}
	return manifest, nil
}

// ParseLanguage parses a language registry manifest.
func ParseLanguage(manifestPath string) (core.LanguageManifest, error) {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return core.LanguageManifest{}, err
	}
	var manifest core.LanguageManifest
	if err := UnmarshalManifest(data, &manifest); err != nil {
		return core.LanguageManifest{}, err
	}
	if manifest.ID == "" {
		manifest.ID = filepath.Base(filepath.Dir(manifestPath))
	}
	return manifest, nil
}

// ParseIntegration parses an integration manifest (integration.yaml / manifest.yaml).
func ParseIntegration(manifestPath string) (core.IntegrationManifest, error) {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return core.IntegrationManifest{}, err
	}
	var manifest core.IntegrationManifest
	if err := UnmarshalManifest(data, &manifest); err != nil {
		return core.IntegrationManifest{}, err
	}
	if manifest.Name == "" {
		manifest.Name = filepath.Base(filepath.Dir(manifestPath))
	}
	return manifest, nil
}

// ParsePolicy parses a policy manifest (YAML or Markdown with frontmatter).
func ParsePolicy(policyPath string) (core.PolicyManifest, error) {
	data, err := os.ReadFile(policyPath)
	if err != nil {
		return core.PolicyManifest{}, err
	}

	var p core.PolicyManifest
	if strings.HasSuffix(policyPath, ".md") {
		fm, body := ParseFrontmatter(string(data))
		p.ID = strings.TrimSuffix(filepath.Base(policyPath), ".md")
		p.Name = p.ID
		p.Enforcement = core.EnforceStrictBlock
		p.Instructions = body
		if err := UnmarshalManifest([]byte(fm), &p); err != nil && fm != "" {
			return core.PolicyManifest{}, err
		}
	} else {
		if err := UnmarshalManifest(data, &p); err != nil {
			return core.PolicyManifest{}, err
		}
		if p.ID == "" {
			p.ID = strings.TrimSuffix(filepath.Base(policyPath), filepath.Ext(policyPath))
		}
		if p.Name == "" {
			p.Name = p.ID
		}
	}
	return p, nil
}

// LoadWorkspaceConfig reads .agentjam/config.yaml if present.
func LoadWorkspaceConfig(rootDir string) (core.WorkspaceConfig, error) {
	cfgPath := filepath.Join(rootDir, core.AgentJamDir, core.ConfigFileName)
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return core.WorkspaceConfig{}, err
	}
	var cfg core.WorkspaceConfig
	if err := UnmarshalManifest(data, &cfg); err != nil {
		return core.WorkspaceConfig{}, err
	}
	return cfg, nil
}
