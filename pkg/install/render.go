package install

import (
	"fmt"
	"strings"

	"github.com/primasdevlabs/agentjam/pkg/core"
)

// SupportedHarnesses lists every render target for `agentjam install`.
var SupportedHarnesses = []string{
	"cursor", "claude-code", "gemini", "cline", "windsurf",
	"devin", "roo-code", "generic",
}

// Render produces the full file map for one harness: a top-level rules file
// carrying the assembled system instruction plus one file per skill, agent,
// and enforced policy in the harness's native format.
func Render(rs *ResourceSet, harness, systemInstruction, agentName, agentDesc string) (map[string]string, error) {
	files := map[string]string{}

	skillBody := func(s SkillEntry) string {
		return fmt.Sprintf("# Skill: %s\n\n%s\n\n%s", s.Manifest.Name, s.Manifest.Description, instructionsText(s.Instructions))
	}
	agentBody := func(a AgentEntry) string {
		return fmt.Sprintf("# Agent: %s\n\n%s\n\nSkills: %s\n\n%s", a.Manifest.Name, a.Manifest.Description,
			strings.Join(a.Manifest.Skills, ", "), instructionsText(a.Instructions))
	}
	policyBody := func(p core.PolicyManifest) string {
		return fmt.Sprintf("# Policy: %s\n\n%s\n\n%s", p.Name, p.Description, strings.TrimSpace(p.Instructions))
	}

	ruleFile := func(name, desc, body string, alwaysApply bool) string {
		return fmt.Sprintf("---\ndescription: %s\nalwaysApply: %t\n---\n\n%s\n", desc, alwaysApply, body)
	}

	switch harness {
	case "cursor":
		files[".cursorrules"] = "# Cursor Project Rules — " + agentName + "\n\n" + systemInstruction
		for _, s := range rs.Skills {
			files[".cursor/rules/skill-"+s.Manifest.Name+".mdc"] = ruleFile(s.Manifest.Name, s.Manifest.Description, skillBody(s), false)
		}
		for _, a := range rs.Agents {
			files[".cursor/rules/agent-"+a.Manifest.Name+".mdc"] = ruleFile("agent: "+a.Manifest.Name, a.Manifest.Description, agentBody(a), false)
		}
		for _, p := range rs.Policies {
			files[".cursor/rules/policy-"+p.ID+".mdc"] = ruleFile("policy: "+p.Name, p.Description, policyBody(p), p.Enforcement == core.EnforceStrictBlock)
		}

	case "claude-code":
		files["CLAUDE.md"] = "# Claude Code Configuration — " + agentName + "\n\n" + systemInstruction
		for _, s := range rs.Skills {
			files[".claude/skills/"+s.Manifest.Name+"/SKILL.md"] =
				fmt.Sprintf("---\nname: %s\ndescription: %s\n---\n\n%s\n", s.Manifest.Name, s.Manifest.Description, skillBody(s))
		}
		for _, a := range rs.Agents {
			files[".claude/agents/"+a.Manifest.Name+".md"] =
				fmt.Sprintf("---\nname: %s\ndescription: %s\n---\n\n%s\n", a.Manifest.Name, a.Manifest.Description, agentBody(a))
		}
		for _, p := range rs.Policies {
			files[".claude/policies/"+p.ID+".md"] = policyBody(p) + "\n"
		}

	case "windsurf":
		files[".windsurfrules"] = "# Windsurf Cascade Rules — " + agentName + "\n\n" + systemInstruction
		for _, s := range rs.Skills {
			files[".windsurf/rules/"+s.Manifest.Name+".md"] = skillBody(s) + "\n"
		}
		for _, a := range rs.Agents {
			files[".windsurf/rules/agent-"+a.Manifest.Name+".md"] = agentBody(a) + "\n"
		}

	case "cline":
		files[".clinerules"] = "# Cline Custom System Prompt — " + agentName + "\n\n" + systemInstruction
		for _, s := range rs.Skills {
			files[".clinerules/rules/"+s.Manifest.Name+".md"] = skillBody(s) + "\n"
		}
		for _, p := range rs.Policies {
			files[".clinerules/rules/policy-"+p.ID+".md"] = policyBody(p) + "\n"
		}

	case "roo-code":
		files[".roomodes"] = "# Roo Code Custom Mode — " + agentName + "\n\n" + systemInstruction
		for _, s := range rs.Skills {
			files[".roo/rules/"+s.Manifest.Name+".md"] = skillBody(s) + "\n"
		}

	case "devin":
		files[".devin/playbook.md"] = "# Devin Agent Playbook — " + agentName + "\n\n" + systemInstruction
		files["AGENTS.md"] = agentsDoc(rs, agentName, agentDesc)

	case "gemini":
		files["GEMINI.md"] = "# Gemini Instructions — " + agentName + "\n\n" + systemInstruction + "\n\n" + skillSections(rs)

	case "generic":
		files["SYSTEM_PROMPT.md"] = "# Generic AI System Prompt — " + agentName + "\n\n" + systemInstruction
		files["AGENTS.md"] = agentsDoc(rs, agentName, agentDesc)

	default:
		return nil, fmt.Errorf("unknown harness %q (supported: %s)", harness, strings.Join(SupportedHarnesses, ", "))
	}
	return files, nil
}

// agentsDoc renders the emerging AGENTS.md standard: skills, agents, and
// policies as markdown sections any agent can read.
func agentsDoc(rs *ResourceSet, agentName, agentDesc string) string {
	var sb strings.Builder
	sb.WriteString("# AGENTS.md — AgentJam Workspace Instructions\n\n")
	fmt.Fprintf(&sb, "Active agent: **%s** — %s\n\n", agentName, agentDesc)

	sb.WriteString("## Policies\n\n")
	for _, p := range rs.Policies {
		fmt.Fprintf(&sb, "### [%s] %s\n\n%s\n\n%s\n\n", strings.ToUpper(string(p.Enforcement)), p.Name, p.Description, strings.TrimSpace(p.Instructions))
	}

	sb.WriteString("## Skills\n\n")
	for _, s := range rs.Skills {
		fmt.Fprintf(&sb, "### %s\n\n%s\n\n%s\n\n", s.Manifest.Name, s.Manifest.Description, instructionsText(s.Instructions))
	}

	sb.WriteString("## Agents\n\n")
	for _, a := range rs.Agents {
		fmt.Fprintf(&sb, "### %s\n\n%s\n\n", a.Manifest.Name, a.Manifest.Description)
	}
	return sb.String()
}

func skillSections(rs *ResourceSet) string {
	var sb strings.Builder
	sb.WriteString("## Skills\n\n")
	for _, s := range rs.Skills {
		fmt.Fprintf(&sb, "### %s\n\n%s\n\n%s\n\n", s.Manifest.Name, s.Manifest.Description, instructionsText(s.Instructions))
	}
	return sb.String()
}
