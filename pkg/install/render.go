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

// catPath joins a base dir, optional category, and file name.
func catPath(base, category, file string) string {
	if category == "" {
		return base + "/" + file
	}
	return base + "/" + category + "/" + file
}

// Render produces the full file map for one harness in a domain-driven
// layout: rules grouped under skills/<category>/, policies/<category>/,
// and agents/, plus a top-level rules file carrying the assembled system
// instruction and the AGENTJAM.md agent entry-point runbook.
func Render(rs *ResourceSet, harness, systemInstruction, agentName, agentDesc string) (map[string]string, error) {
	files := map[string]string{}

	skillBody := func(s SkillEntry) string {
		return fmt.Sprintf("# Skill: %s\n\n%s\n\n%s", s.Manifest.Name, s.Manifest.Description, instructionsText(s.Instructions))
	}
	agentBody := func(a AgentEntry) string {
		return fmt.Sprintf("# Agent: %s\n\n%s\n\nSkills: %s\n\n%s", a.Manifest.Name, a.Manifest.Description,
			strings.Join(a.Manifest.Skills, ", "), instructionsText(a.Instructions))
	}
	policyBody := func(p PolicyEntry) string {
		return fmt.Sprintf("# Policy: %s\n\n%s\n\n%s", p.Manifest.Name, p.Manifest.Description, strings.TrimSpace(p.Manifest.Instructions))
	}

	mdc := func(desc, body string, alwaysApply bool) string {
		return fmt.Sprintf("---\ndescription: %s\nalwaysApply: %t\n---\n\n%s\n", desc, alwaysApply, body)
	}

	switch harness {
	case "cursor":
		files[".cursorrules"] = "# Cursor Project Rules — " + agentName + "\n\n" + systemInstruction
		for _, s := range rs.Skills {
			files[catPath(".cursor/rules/skills", s.Category, s.Manifest.Name+".mdc")] =
				mdc(s.Manifest.Description, skillBody(s), false)
		}
		for _, a := range rs.Agents {
			files[".cursor/rules/agents/"+a.Manifest.Name+".mdc"] =
				mdc("agent: "+a.Manifest.Name, agentBody(a), false)
		}
		for _, p := range rs.Policies {
			files[catPath(".cursor/rules/policies", p.Category, p.Manifest.ID+".mdc")] =
				mdc("policy: "+p.Manifest.Name, policyBody(p), p.Manifest.Enforcement == core.EnforceStrictBlock)
		}

	case "claude-code":
		files["CLAUDE.md"] = "# Claude Code Configuration — " + agentName + "\n\n" + systemInstruction
		for _, s := range rs.Skills {
			files[catPath(".claude/skills", s.Category, s.Manifest.Name)+"/SKILL.md"] =
				fmt.Sprintf("---\nname: %s\ndescription: %s\n---\n\n%s\n", s.Manifest.Name, s.Manifest.Description, skillBody(s))
		}
		for _, a := range rs.Agents {
			files[".claude/agents/"+a.Manifest.Name+".md"] =
				fmt.Sprintf("---\nname: %s\ndescription: %s\n---\n\n%s\n", a.Manifest.Name, a.Manifest.Description, agentBody(a))
		}
		for _, p := range rs.Policies {
			files[catPath(".claude/policies", p.Category, p.Manifest.ID+".md")] = policyBody(p) + "\n"
		}

	case "windsurf":
		files[".windsurfrules"] = "# Windsurf Cascade Rules — " + agentName + "\n\n" + systemInstruction
		for _, s := range rs.Skills {
			files[catPath(".windsurf/rules/skills", s.Category, s.Manifest.Name+".md")] = skillBody(s) + "\n"
		}
		for _, a := range rs.Agents {
			files[".windsurf/rules/agents/"+a.Manifest.Name+".md"] = agentBody(a) + "\n"
		}
		for _, p := range rs.Policies {
			files[catPath(".windsurf/rules/policies", p.Category, p.Manifest.ID+".md")] = policyBody(p) + "\n"
		}

	case "cline":
		files[".clinerules"] = "# Cline Custom System Prompt — " + agentName + "\n\n" + systemInstruction
		for _, s := range rs.Skills {
			files[catPath(".clinerules/skills", s.Category, s.Manifest.Name+".md")] = skillBody(s) + "\n"
		}
		for _, a := range rs.Agents {
			files[".clinerules/agents/"+a.Manifest.Name+".md"] = agentBody(a) + "\n"
		}
		for _, p := range rs.Policies {
			files[catPath(".clinerules/policies", p.Category, p.Manifest.ID+".md")] = policyBody(p) + "\n"
		}

	case "roo-code":
		files[".roomodes"] = "# Roo Code Custom Mode — " + agentName + "\n\n" + systemInstruction
		for _, s := range rs.Skills {
			files[catPath(".roo/rules/skills", s.Category, s.Manifest.Name+".md")] = skillBody(s) + "\n"
		}
		for _, a := range rs.Agents {
			files[".roo/rules/agents/"+a.Manifest.Name+".md"] = agentBody(a) + "\n"
		}
		for _, p := range rs.Policies {
			files[catPath(".roo/rules/policies", p.Category, p.Manifest.ID+".md")] = policyBody(p) + "\n"
		}

	case "devin":
		files[".devin/playbook.md"] = "# Devin Agent Playbook — " + agentName + "\n\n" + systemInstruction
		files["AGENTS.md"] = agentsDoc(rs, agentName, agentDesc)
		for _, s := range rs.Skills {
			files[catPath(".devin/rules/skills", s.Category, s.Manifest.Name+".md")] = skillBody(s) + "\n"
		}
		for _, a := range rs.Agents {
			files[".devin/rules/agents/"+a.Manifest.Name+".md"] = agentBody(a) + "\n"
		}
		for _, p := range rs.Policies {
			files[catPath(".devin/rules/policies", p.Category, p.Manifest.ID+".md")] = policyBody(p) + "\n"
		}

	case "gemini":
		files["GEMINI.md"] = "# Gemini Instructions — " + agentName + "\n\n" + systemInstruction
		for _, s := range rs.Skills {
			files[catPath(".gemini/skills", s.Category, s.Manifest.Name+".md")] = skillBody(s) + "\n"
		}
		for _, a := range rs.Agents {
			files[".gemini/agents/"+a.Manifest.Name+".md"] = agentBody(a) + "\n"
		}
		for _, p := range rs.Policies {
			files[catPath(".gemini/policies", p.Category, p.Manifest.ID+".md")] = policyBody(p) + "\n"
		}

	case "generic":
		files["SYSTEM_PROMPT.md"] = "# Generic AI System Prompt — " + agentName + "\n\n" + systemInstruction
		files["AGENTS.md"] = agentsDoc(rs, agentName, agentDesc)
		for _, s := range rs.Skills {
			files[catPath(".agentjam/rules/skills", s.Category, s.Manifest.Name+".md")] = skillBody(s) + "\n"
		}
		for _, a := range rs.Agents {
			files[".agentjam/rules/agents/"+a.Manifest.Name+".md"] = agentBody(a) + "\n"
		}
		for _, p := range rs.Policies {
			files[catPath(".agentjam/rules/policies", p.Category, p.Manifest.ID+".md")] = policyBody(p) + "\n"
		}

	default:
		return nil, fmt.Errorf("unknown harness %q (supported: %s)", harness, strings.Join(SupportedHarnesses, ", "))
	}

	// Every harness gets the operational entry point: the runbook an agent
	// reads first — enforcement protocol, scan/fix loop, command map.
	files["AGENTJAM.md"] = Runbook(agentName)

	return files, nil
}

// Runbook renders the agent entry point installed into every target
// project: enforcement protocol, the scan→fix loop, and the command map.
func Runbook(agentName string) string {
	return `# AgentJam Runbook — Operational Entry Point for AI Agents

You are operating in an AgentJam-managed workspace as **` + agentName + `**.
This file is your entry point. Follow the protocol below exactly.

## Enforcement Protocol (before you act)

1. Read the installed rules for your harness — they are domain-organized:
   ` + "```" + `
   skills/<category>/<skill>     — capability instructions
   policies/<category>/<policy>  — enforced rules (strict-block = must obey)
   agents/<id>                   — persona definitions
   ` + "```" + `
2. strict-block policies are non-negotiable. A response that violates one is
   a failed response.
3. Do not produce output that ` + "`agentjam scan`" + ` would flag.

## Scan & Fix Loop

` + "```bash" + `
agentjam scan              # policy-scan all source files; exit 1 on strict blocks
agentjam eval <file>       # evaluate one file with line numbers
agentjam validate          # canonical resource integrity (errors must be 0)
agentjam run <workflow>    # structured multi-step workflows (e.g. bug-fixing)
` + "```" + `

For every task:
1. **Scan first** — ` + "`agentjam scan`" + ` on the touched area before and after edits.
2. **Fix violations** by severity: strict-block → warning → info.
3. **Re-scan** until clean. ` + "`agentjam scan`" + ` exit code 0 is the gate.

## Tooling

` + "```bash" + `
agentjam mcp               # serve workspace tools (filesystem/git/terminal/run_workflow) over MCP stdio
agentjam preflight         # verify toolchain binaries for the detected stack
agentjam detect            # show detected harnesses and stacks
` + "```" + `

## Empty Projects

If this project was just initialized, the canonical tree under agents/,
skills/, policies/, tools/, workflows/, stacks/, languages/ is the source of
truth — extend it rather than working around it, and keep
` + "`agentjam validate`" + ` green.
`
}

// agentsDoc renders the AGENTS.md standard file: skills, agents, and
// policies grouped by domain category.
func agentsDoc(rs *ResourceSet, agentName, agentDesc string) string {
	var sb strings.Builder
	sb.WriteString("# AGENTS.md — AgentJam Workspace Instructions\n\n")
	fmt.Fprintf(&sb, "Active agent: **%s** — %s\n\n", agentName, agentDesc)
	sb.WriteString("Operational protocol: see AGENTJAM.md (scan → fix → validate).\n\n")

	sb.WriteString("## Policies\n\n")
	cat := ""
	for _, p := range rs.Policies {
		if p.Category != cat {
			cat = p.Category
			fmt.Fprintf(&sb, "#### %s\n\n", orCore(cat))
		}
		fmt.Fprintf(&sb, "- **[%s] %s** — %s\n", strings.ToUpper(string(p.Manifest.Enforcement)), p.Manifest.Name, p.Manifest.Description)
	}
	sb.WriteString("\n")

	sb.WriteString("## Skills\n\n")
	cat = ""
	for _, s := range rs.Skills {
		if s.Category != cat {
			cat = s.Category
			fmt.Fprintf(&sb, "#### %s\n\n", orCore(cat))
		}
		fmt.Fprintf(&sb, "- **%s** — %s\n", s.Manifest.Name, s.Manifest.Description)
	}
	sb.WriteString("\n")

	sb.WriteString("## Agents\n\n")
	for _, a := range rs.Agents {
		fmt.Fprintf(&sb, "- **%s** — %s (skills: %s)\n", a.Manifest.Name, a.Manifest.Description, strings.Join(a.Manifest.Skills, ", "))
	}
	return sb.String()
}

func orCore(cat string) string {
	if cat == "" {
		return "core"
	}
	return cat
}
