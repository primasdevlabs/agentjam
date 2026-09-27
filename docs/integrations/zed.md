# Zed Editor Integration

- **Type**: `ide`
- **Compatibility Level**: Level 4 (Workflows)
- **Manifest**: `integrations/ide/zed/integration.yaml`
- **Output**: `.zed/zed.json`

## Capabilities

Instructions, skills, tools, filesystem, terminal, MCP, project rules, and context files.

Not supported: agents, workflows, hooks, lifecycle events, validation.

## Integration Path

Zed consumes AgentJam tooling through the MCP stdio server:

```bash
agentjam mcp
```

Register the command in Zed's assistant configuration as a context server. Workspace rules are exported via the generic exporter (`agentjam export --harness generic`).
