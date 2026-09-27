# OpenCode Integration

- **Type**: `cli`
- **Compatibility Level**: Level 4 (Workflows)
- **Manifest**: `integrations/cli/opencode/integration.yaml`
- **Output**: `opencode.json`

## Capabilities

Instructions, skills, tools, filesystem, terminal, MCP, project rules, and context files.

Not supported: agents, workflows, hooks, lifecycle events, validation.

## Integration Path

Register the AgentJam MCP server in `opencode.json` under `mcp.servers` so OpenCode can call workspace tools:

```bash
agentjam mcp
```

Project rules are exported via `agentjam export --harness generic`.
