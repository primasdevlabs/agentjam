# Cursor Integration

- **Type**: `ide`
- **Compatibility Level**: Level 5 (Validation)
- **Manifest**: `integrations/ide/cursor/integration.yaml`
- **Output**: `.cursorrules` and `.cursor/rules/`

## Capabilities

Instructions, skills, agents, workflows, tools, filesystem, terminal, MCP, project rules, context files, and validation.

Not supported: hooks and lifecycle events.

## Usage

```bash
agentjam export --harness cursor
```

This writes the active agent's system instruction and the enforced policy matrix into Cursor's rules locations. Re-run after changing `.agentjam/config.yaml`, policies, or the active agent.
