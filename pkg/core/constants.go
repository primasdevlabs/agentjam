package core

const (
	Version       = "1.0.0-go"
	ConfigFileName = "config.yaml"
	AgentJamDir   = ".agentjam"
)

// SupportedHarnesses lists exported AI targets.
var SupportedHarnesses = []string{
	"claude-code",
	"cursor",
	"cline",
	"gemini",
	"devin",
	"windsurf",
	"generic",
}
