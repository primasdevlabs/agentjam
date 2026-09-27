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

type envProbe struct {
	name               string
	envType            string
	compatibilityLevel int
	files              []string
}

var envProbes = []envProbe{
	{"Cursor IDE", "ide", 5, []string{".cursorrules", ".cursor"}},
	{"Windsurf IDE", "ide", 5, []string{".windsurfrules", ".codeium"}},
	{"VS Code", "ide", 4, []string{".vscode"}},
	{"Zed IDE", "ide", 3, []string{".zed"}},
	{"Claude Code CLI", "cli", 5, []string{"CLAUDE.md", ".claude"}},
	{"Antigravity / Gemini Platform", "ai-platforms", 6, []string{"GEMINI.md", ".gemini"}},
	{"Cline Extension", "extensions", 4, []string{".clinerules"}},
	{"Roo Code Extension", "extensions", 4, []string{".roomodes"}},
	{"Devin Autonomous Agent", "autonomous-agents", 5, []string{".devin"}},
	{"Generic Agents File", "generic", 3, []string{"AGENTS.md"}},
}

// DetectEnvironments inspects workspace configuration files.
func DetectEnvironments(projectRoot string) []DetectedEnvironment {
	envs := make([]DetectedEnvironment, 0)

	for _, probe := range envProbes {
		detected := make([]string, 0, len(probe.files))
		for _, f := range probe.files {
			if exists(filepath.Join(projectRoot, f)) {
				detected = append(detected, f)
			}
		}
		if len(detected) > 0 {
			envs = append(envs, DetectedEnvironment{
				Name:               probe.name,
				Type:               probe.envType,
				CompatibilityLevel: probe.compatibilityLevel,
				DetectedFiles:      detected,
			})
		}
	}

	return envs
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
