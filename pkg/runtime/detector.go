package runtime

import (
	"os"
	"path/filepath"
)

// DetectedEnvironment represents detected IDE or AI harness.
type DetectedEnvironment struct {
	Name               string   `json:"name"`
	Type               string   `json:"type"`
	CompatibilityLevel int      `json:"compatibilityLevel"`
	DetectedFiles      []string `json:"detectedFiles"`
}

// DetectEnvironments inspects workspace configuration files.
func DetectEnvironments(projectRoot string) []DetectedEnvironment {
	envs := make([]DetectedEnvironment, 0)

	// Cursor check
	if exists(filepath.Join(projectRoot, ".cursorrules")) || exists(filepath.Join(projectRoot, ".cursor")) {
		envs = append(envs, DetectedEnvironment{
			Name:               "Cursor IDE",
			Type:               "ide",
			CompatibilityLevel: 5,
			DetectedFiles:      []string{".cursorrules"},
		})
	}

	// Antigravity / Gemini check
	if exists(filepath.Join(projectRoot, "GEMINI.md")) || exists(filepath.Join(projectRoot, ".gemini")) {
		envs = append(envs, DetectedEnvironment{
			Name:               "Antigravity / Gemini Platform",
			Type:               "ai-platforms",
			CompatibilityLevel: 6,
			DetectedFiles:      []string{"GEMINI.md"},
		})
	}

	// Claude Code check
	if exists(filepath.Join(projectRoot, "CLAUDE.md")) {
		envs = append(envs, DetectedEnvironment{
			Name:               "Claude Code CLI",
			Type:               "cli",
			CompatibilityLevel: 5,
			DetectedFiles:      []string{"CLAUDE.md"},
		})
	}

	return envs
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
