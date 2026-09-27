# Gemini CLI Integration

- **Type**: `cli`
- **Compatibility Level**: Level 5 (Validation)
- **Manifest**: `integrations/cli/gemini/integration.yaml`
- **Output**: `GEMINI.md`

## Capabilities

Instructions, tools, filesystem, terminal, MCP, context files, and validation.

Not supported: skills, agents, workflows, project rules, hooks, lifecycle events.

## Usage

```bash
agentjam export --harness gemini
```

Writes `GEMINI.md` with the active persona and policy matrix. Gemini CLI loads it as its system instruction context.
