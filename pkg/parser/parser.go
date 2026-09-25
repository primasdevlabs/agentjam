package parser

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/agentjam/agentjam/pkg/core"
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

// DiscoverResources recursively scans root directory for agents, skills, workflows, policies, stacks, and languages.
func DiscoverResources(rootDir string) []DiscoveredResource {
	results := make([]DiscoveredResource, 0)

	_ = filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		base := info.Name()
		dir := filepath.Dir(path)
		parentDir := filepath.Base(dir)

		var resType core.ResourceType
		var id string

		if base == "agent.yaml" || base == "agent.yml" || base == "manifest.yaml" {
			resType = core.ResourceTypeAgent
			id = parentDir
		} else if base == "skill.yaml" || base == "skill.yml" || base == "SKILL.md" {
			resType = core.ResourceTypeSkill
			id = parentDir
		} else if base == "workflow.yaml" || base == "workflow.yml" {
			resType = core.ResourceTypeWorkflow
			id = parentDir
		} else if base == "stack.yaml" || base == "stack.yml" {
			resType = core.ResourceTypeStack
			id = parentDir
		} else if base == "language.yaml" || base == "language.yml" {
			resType = core.ResourceTypeLanguage
			id = parentDir
		} else if strings.HasPrefix(filepath.ToSlash(path), "policies/") && (strings.HasSuffix(base, ".yaml") || strings.HasSuffix(base, ".md")) {
			resType = core.ResourceTypePolicy
			id = strings.TrimSuffix(strings.TrimSuffix(base, ".yaml"), ".md")
		}

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

// ParseAgent parses an agent resource bundle.
func ParseAgent(manifestPath string) (ResourceBundle[core.AgentManifest], error) {
	dir := filepath.Dir(manifestPath)
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return ResourceBundle[core.AgentManifest]{}, err
	}

	var manifest core.AgentManifest
	// Attempt JSON / YAML parsing
	if err := json.Unmarshal(data, &manifest); err != nil {
		manifest = core.AgentManifest{
			Name:        filepath.Base(dir),
			Version:     "1.0.0",
			Type:        "agent",
			Description: "Agent persona",
			Skills:      []string{},
			Tools:       []string{},
		}
	}

	instructions := ReadInstructionFiles(dir)
	return ResourceBundle[core.AgentManifest]{
		Manifest:     manifest,
		Instructions: instructions,
		BasePath:     dir,
	}, nil
}

// ParseSkill parses a skill resource bundle.
func ParseSkill(manifestPath string) (ResourceBundle[core.SkillManifest], error) {
	dir := filepath.Dir(manifestPath)
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return ResourceBundle[core.SkillManifest]{}, err
	}

	var manifest core.SkillManifest
	if strings.HasSuffix(manifestPath, ".md") {
		fm, body := ParseFrontmatter(string(data))
		manifest.Name = filepath.Base(dir)
		manifest.Version = "1.0.0"
		manifest.Type = "skill"
		manifest.Description = body
		_ = json.Unmarshal([]byte(fm), &manifest)
	} else {
		_ = json.Unmarshal(data, &manifest)
	}

	instructions := ReadInstructionFiles(dir)
	return ResourceBundle[core.SkillManifest]{
		Manifest:     manifest,
		Instructions: instructions,
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
	_ = json.Unmarshal(data, &manifest)
	if manifest.Name == "" {
		manifest.Name = filepath.Base(dir)
	}

	instructions := ReadInstructionFiles(dir)
	return ResourceBundle[core.WorkflowManifest]{
		Manifest:     manifest,
		Instructions: instructions,
		BasePath:     dir,
	}, nil
}

// ParsePolicy parses a policy manifest (YAML or Markdown).
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
		_ = json.Unmarshal([]byte(fm), &p)
	} else {
		_ = json.Unmarshal(data, &p)
		if p.ID == "" {
			p.ID = strings.TrimSuffix(filepath.Base(policyPath), filepath.Ext(policyPath))
		}
	}
	return p, nil
}
