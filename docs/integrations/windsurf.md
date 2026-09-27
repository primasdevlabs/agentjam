# Windsurf Integration

- **Type**: `ide`
- **Compatibility Level**: Level 5 (Validation)
- **Manifest**: `integrations/ide/windsurf/integration.yaml`
- **Output**: `.windsurfrules`

## Capabilities

Instructions, skills, agents, workflows, tools, filesystem, terminal, MCP, project rules, context files, and validation.

Not supported: hooks and lifecycle events.

## Usage

```bash
agentjam export --harness windsurf
```

Writes Cascade rules containing the active persona, stack profile, and enforced policy matrix.
