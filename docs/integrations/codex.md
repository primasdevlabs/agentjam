# OpenAI Codex CLI Integration

- **Type**: `cli`
- **Compatibility Level**: Level 4 (Workflows)
- **Manifest**: `integrations/cli/codex/integration.yaml`
- **Output**: `codex.config.json`

## Capabilities

Instructions, skills, tools, filesystem, terminal, project rules, and context files.

Not supported: agents, workflows, MCP, hooks, lifecycle events, validation.

## Integration Path

Codex consumes project rules through its config file and workspace tools through the MCP stdio server:

```bash
agentjam mcp
```
