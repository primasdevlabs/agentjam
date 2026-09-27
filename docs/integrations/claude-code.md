# Claude Code Integration

- **Type**: `cli`
- **Compatibility Level**: Level 6 (Lifecycle Integration)
- **Manifest**: `integrations/cli/claude-code/integration.yaml`
- **Output**: `CLAUDE.md`

## Capabilities

Full surface: instructions, skills, agents, workflows, tools, filesystem, terminal, MCP, project rules, context files, hooks, lifecycle events, and validation.

## Usage

```bash
agentjam export --harness claude-code
```

Writes `CLAUDE.md` at the workspace root with the active persona (from `.agentjam/config.yaml` `defaultAgent` or `--agent`), the active stack profile, and all enforced policy instructions. Claude Code picks it up automatically at session start.
