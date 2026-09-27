# Roo Code Integration

- **Type**: `extensions`
- **Compatibility Level**: Level 5 (Validation)
- **Manifest**: `integrations/extensions/roo-code/integration.yaml`
- **Output**: `.roo/.roomodes`

## Capabilities

Instructions, skills, agents, workflows, tools, filesystem, terminal, MCP, project rules, context files, and validation.

Not supported: hooks and lifecycle events.

## Usage

```bash
agentjam export --harness roo-code
```

Writes a Roo Code custom-modes file carrying the active agent persona and enforced policy matrix.
