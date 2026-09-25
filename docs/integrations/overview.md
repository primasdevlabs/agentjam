# AgentJam AI Development Environment Compatibility

> **Core Principle**: AgentJam is environment-neutral. One project standard, many AI environments. The developer chooses their AI environment; AgentJam provides the common operating rules.

## Environment Compatibility Taxonomy

AgentJam categorizes AI development environments into 6 categories under `integrations/`:

1. **IDEs and Editors (`integrations/ide/`)**: Visual Studio Code, Cursor, Windsurf, Zed, Void.
2. **Autonomous Coding Agents (`integrations/autonomous-agents/`)**: Devin, hosted coding agents.
3. **CLI / Terminal Agents (`integrations/cli/`)**: Claude Code, Codex, Gemini CLI, OpenCode.
4. **IDE Extensions (`integrations/extensions/`)**: Cline, Roo Code.
5. **AI Platforms (`integrations/ai-platforms/`)**: Antigravity, emerging AI platforms.
6. **Generic Fallback (`integrations/generic/`)**: Markdown instructions, standard project configuration, CLI fallback.

## Compatibility Levels (0 to 6)

- **Level 0 — Documentation**: Manual consumption of AgentJam documentation.
- **Level 1 — Instructions**: Project rules & markdown instructions (`.cursorrules`, `CLAUDE.md`, `GEMINI.md`).
- **Level 2 — Skills**: Reusable task-specific AgentJam skill modules.
- **Level 3 — Tools**: Connection to AgentJam abstract tools or MCP servers.
- **Level 4 — Workflows**: Single and multi-agent execution flows.
- **Level 5 — Validation**: AgentJam post-execution validation engine (`agentjam validate`).
- **Level 6 — Lifecycle Integration**: Participating in lifecycle events (`before task`, `before commit`, etc.).
