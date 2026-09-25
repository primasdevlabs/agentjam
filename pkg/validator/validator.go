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

// ValidateRepository validates repository resources and cross-references.
func ValidateRepository(rootDir string) []ValidationError {
	errors := make([]ValidationError, 0)
	discovered := parser.DiscoverResources(rootDir)

	knownSkills := make(map[string]bool)
	knownTools := make(map[string]bool)

	for _, res := range discovered {
		if res.Type == core.ResourceTypeSkill {
			knownSkills[res.ID] = true
		} else if res.Type == core.ResourceTypeTool {
			knownTools[res.ID] = true
		}
	}

	for _, res := range discovered {
		if res.Type == core.ResourceTypeAgent {
			bundle, err := parser.ParseAgent(res.Path)
			if err == nil {
				for _, skillID := range bundle.Manifest.Skills {
					if !knownSkills[skillID] {
						errors = append(errors, ValidationError{
							ResourcePath: res.Path,
							ResourceID:   res.ID,
							Field:        "skills",
							Message:      fmt.Sprintf("Agent '%s' references unknown skill '%s'", res.ID, skillID),
							Severity:     "warning",
						})
					}
				}
			}
		}
	}

	return errors
}
