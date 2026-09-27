# Visual Studio Code Integration

- **Type**: `ide`
- **Compatibility Level**: Level 5 (Validation)
- **Manifest**: `integrations/ide/vscode/integration.yaml`
- **Output**: `.vscode/settings.json`

## Capabilities

Instructions, skills, agents, workflows, tools, filesystem, terminal, MCP, project rules, context files, and validation.

Not supported: hooks and lifecycle events (VS Code relies on extension-provided hooks).

## Integration Path

```bash
agentjam export --harness generic   # write AGENTJAM.md workspace rules
agentjam mcp                        # expose workspace tools over MCP stdio
```

Pair with an MCP-capable extension (e.g. Cline or Roo Code) to consume tools — see the [extensions catalog](index.md).
