package validator

import (
	"fmt"

	"github.com/primasdevlabs/agentjam/pkg/core"
	"github.com/primasdevlabs/agentjam/pkg/parser"
)

// ValidationError represents a repository validation issue.
type ValidationError struct {
	ResourcePath string `json:"resourcePath"`
	ResourceID   string `json:"resourceId"`
	Field        string `json:"field"`
	Message      string `json:"message"`
	Severity     string `json:"severity"` // error, warning
}

func err(path, id, field, msg string) ValidationError {
	return ValidationError{ResourcePath: path, ResourceID: id, Field: field, Message: msg, Severity: "error"}
}

func warn(path, id, field, msg string) ValidationError {
	return ValidationError{ResourcePath: path, ResourceID: id, Field: field, Message: msg, Severity: "warning"}
}

// HasErrors reports whether the list contains any error-severity issues.
func HasErrors(errors []ValidationError) bool {
	for _, e := range errors {
		if e.Severity == "error" {
			return true
		}
	}
	return false
}

// ValidateRepository validates repository resources and cross-references.
func ValidateRepository(rootDir string) []ValidationError {
	errors := make([]ValidationError, 0)
	discovered := parser.DiscoverResources(rootDir)

	known := map[core.ResourceType]map[string]bool{}
	for _, res := range discovered {
		if known[res.Type] == nil {
			known[res.Type] = map[string]bool{}
		}
		known[res.Type][res.ID] = true
	}

	checkRefs := func(res parser.DiscoveredResource, field string, refs []string, target core.ResourceType) {
		for _, ref := range refs {
			if !known[target][ref] {
				errors = append(errors, warn(res.Path, res.ID, field,
					fmt.Sprintf("%s '%s' references unknown %s '%s'", res.Type, res.ID, target, ref)))
			}
		}
	}

	for _, res := range discovered {
		switch res.Type {
		case core.ResourceTypeAgent:
			bundle, perr := parser.ParseAgent(res.Path)
			if perr != nil {
				errors = append(errors, err(res.Path, res.ID, "manifest", perr.Error()))
				continue
			}
			m := bundle.Manifest
			if m.Description == "" {
				errors = append(errors, err(res.Path, res.ID, "description", "Agent manifest is missing a description"))
			}
			checkRefs(res, "skills", m.Skills, core.ResourceTypeSkill)
			checkRefs(res, "tools", m.Tools, core.ResourceTypeTool)

		case core.ResourceTypeWorkflow:
			bundle, perr := parser.ParseWorkflow(res.Path)
			if perr != nil {
				errors = append(errors, err(res.Path, res.ID, "manifest", perr.Error()))
				continue
			}
			m := bundle.Manifest
			if m.Description == "" {
				errors = append(errors, err(res.Path, res.ID, "description", "Workflow manifest is missing a description"))
			}
			if len(m.Steps) == 0 {
				errors = append(errors, warn(res.Path, res.ID, "steps", "Workflow defines no steps"))
			}
			checkRefs(res, "agents", m.Agents, core.ResourceTypeAgent)
			checkRefs(res, "skills", m.Skills, core.ResourceTypeSkill)
			seen := map[string]bool{}
			for _, step := range m.Steps {
				if step.ID == "" {
					errors = append(errors, err(res.Path, res.ID, "steps", "Workflow step is missing an id"))
				} else if seen[step.ID] {
					errors = append(errors, err(res.Path, res.ID, "steps",
						fmt.Sprintf("Duplicate workflow step id '%s'", step.ID)))
				}
				seen[step.ID] = true
			}

		case core.ResourceTypeSkill:
			bundle, perr := parser.ParseSkill(res.Path)
			if perr != nil {
				errors = append(errors, err(res.Path, res.ID, "manifest", perr.Error()))
				continue
			}
			if bundle.Manifest.Description == "" {
				errors = append(errors, warn(res.Path, res.ID, "description", "Skill manifest is missing a description"))
			}

		case core.ResourceTypeTool:
			m, perr := parser.ParseTool(res.Path)
			if perr != nil {
				errors = append(errors, err(res.Path, res.ID, "manifest", perr.Error()))
				continue
			}
			if m.Description == "" {
				errors = append(errors, err(res.Path, res.ID, "description", "Tool manifest is missing a description"))
			}
			switch m.SafetyLevel {
			case "", core.SafetyReadOnly, core.SafetySafeWrite, core.SafetyDestructive, core.SafetyAdmin:
			default:
				errors = append(errors, err(res.Path, res.ID, "safetyLevel",
					fmt.Sprintf("Tool '%s' has invalid safetyLevel '%s'", res.ID, m.SafetyLevel)))
			}

		case core.ResourceTypeStack:
			m, perr := parser.ParseStack(res.Path)
			if perr != nil {
				errors = append(errors, err(res.Path, res.ID, "manifest", perr.Error()))
				continue
			}
			if m.Description == "" {
				errors = append(errors, warn(res.Path, res.ID, "description", "Stack manifest is missing a description"))
			}

		case core.ResourceTypePolicy:
			p, perr := parser.ParsePolicy(res.Path)
			if perr != nil {
				errors = append(errors, err(res.Path, res.ID, "manifest", perr.Error()))
				continue
			}
			switch p.Enforcement {
			case "", core.EnforceStrictBlock, core.EnforceWarning, core.EnforceInfo:
			default:
				errors = append(errors, err(res.Path, res.ID, "enforcement",
					fmt.Sprintf("Policy '%s' has invalid enforcement '%s'", res.ID, p.Enforcement)))
			}
		}
	}

	return errors
}
