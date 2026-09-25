package adapters

import (
	"fmt"
	"strings"

	"github.com/primasdevlabs/agentjam/pkg/core"
)

// ExportResult holds exported file content map.
type ExportResult struct {
	Harness  string            `json:"harness"`
	Files    map[string]string `json:"files"`
	Warnings []string          `json:"warnings,omitempty"`
}

// ExportClaudeCode renders CLAUDE.md for Claude Code CLI.
func ExportClaudeCode(agent core.AgentManifest, instructions string) ExportResult {
	files := make(map[string]string)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Claude Code Configuration — %s\n\n", agent.Name))
	sb.WriteString(fmt.Sprintf("Description: %s\n\n", agent.Description))
	sb.WriteString("## System Guidance & AgentJam Policies\n")
	sb.WriteString(instructions)

	files["CLAUDE.md"] = sb.String()
	return ExportResult{Harness: "claude-code", Files: files}
}

// ExportCursorRules renders .cursorrules for Cursor IDE.
func ExportCursorRules(agent core.AgentManifest, instructions string) ExportResult {
	files := make(map[string]string)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Cursor Project Rules — %s\n\n", agent.Name))
	sb.WriteString(instructions)

	files[".cursorrules"] = sb.String()
	return ExportResult{Harness: "cursor", Files: files}
}

// ExportGeminiMarkdown renders GEMINI.md for Antigravity & Gemini CLI.
func ExportGeminiMarkdown(agent core.AgentManifest, instructions string) ExportResult {
	files := make(map[string]string)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Antigravity / Gemini Instructions — %s\n\n", agent.Name))
	sb.WriteString(instructions)

	files["GEMINI.md"] = sb.String()
	return ExportResult{Harness: "gemini", Files: files}
}

// ExportClineRules renders .clinerules for Cline Extension.
func ExportClineRules(agent core.AgentManifest, instructions string) ExportResult {
	files := make(map[string]string)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Cline Custom System Prompt — %s\n\n", agent.Name))
	sb.WriteString(instructions)

	files[".clinerules"] = sb.String()
	return ExportResult{Harness: "cline", Files: files}
}

// ExportWindsurfRules renders .windsurfrules for Windsurf IDE.
func ExportWindsurfRules(agent core.AgentManifest, instructions string) ExportResult {
	files := make(map[string]string)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Windsurf Cascade Rules — %s\n\n", agent.Name))
	sb.WriteString(instructions)

	files[".windsurfrules"] = sb.String()
	return ExportResult{Harness: "windsurf", Files: files}
}

// ExportDevinConfig renders devin.json for Devin Autonomous Agent.
func ExportDevinConfig(agent core.AgentManifest, instructions string) ExportResult {
	files := make(map[string]string)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Devin Agent Playbook — %s\n\n", agent.Name))
	sb.WriteString(instructions)

	files[".devin/playbook.md"] = sb.String()
	return ExportResult{Harness: "devin", Files: files}
}

// ExportRooCodeModes renders .roomodes for Roo Code Extension.
func ExportRooCodeModes(agent core.AgentManifest, instructions string) ExportResult {
	files := make(map[string]string)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Roo Code Custom Mode — %s\n\n", agent.Name))
	sb.WriteString(instructions)

	files[".roomodes"] = sb.String()
	return ExportResult{Harness: "roo-code", Files: files}
}

// ExportGenericPrompt renders SYSTEM_PROMPT.md for generic platforms.
func ExportGenericPrompt(agent core.AgentManifest, instructions string) ExportResult {
	files := make(map[string]string)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Generic AI System Prompt — %s\n\n", agent.Name))
	sb.WriteString(instructions)

	files["SYSTEM_PROMPT.md"] = sb.String()
	return ExportResult{Harness: "generic", Files: files}
}

// ExportAllHarnesses generates configs across all supported AI harnesses.
func ExportAllHarnesses(agent core.AgentManifest, instructions string) map[string]ExportResult {
	return map[string]ExportResult{
		"claude-code": ExportClaudeCode(agent, instructions),
		"cursor":      ExportCursorRules(agent, instructions),
		"gemini":      ExportGeminiMarkdown(agent, instructions),
		"cline":       ExportClineRules(agent, instructions),
		"windsurf":    ExportWindsurfRules(agent, instructions),
		"devin":       ExportDevinConfig(agent, instructions),
		"roo-code":    ExportRooCodeModes(agent, instructions),
		"generic":     ExportGenericPrompt(agent, instructions),
	}
}
