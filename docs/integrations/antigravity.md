# Antigravity Integration

- **Type**: `ai-platforms`
- **Compatibility Level**: Level 6 (Lifecycle Integration)
- **Manifest**: `integrations/ai-platforms/antigravity/integration.yaml`
- **Output**: `.gemini/GEMINI.md`

## Capabilities

Full surface: instructions, skills, agents, workflows, tools, filesystem, terminal, MCP, project rules, context files, hooks, lifecycle events, and validation.

## Integration Path

Antigravity consumes the Gemini-format context file:

```bash
agentjam export --harness gemini
```

Tool access is provided through the MCP stdio server (`agentjam mcp`).
