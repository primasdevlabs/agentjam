package core

// Version is the CLI/engine version. It is a var (not const) so release
// builds can inject the git tag via -ldflags "-X .../pkg/core.Version=vX.Y.Z".
var Version = "1.0.0-go"

const (
	ConfigFileName = "config.yaml"
	AgentJamDir    = ".agentjam"
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
