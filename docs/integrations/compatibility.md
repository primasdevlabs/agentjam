# Compatibility Matrix & Capability Schema

## Capability Schema

```yaml
capabilities:
  instructions: true
  skills: true
  agents: true
  workflows: true
  tools: true
  filesystem: true
  terminal: true
  browser: false
  mcp: true
  project_rules: true
  context_files: true
  hooks: false
  lifecycle_events: false
  validation: true
```

## Precedence Model

```text
User explicit instruction
        ↓
Project configuration (.agentjam/config.yaml)
        ↓
AgentJam project policy (.agentjam/policies/)
        ↓
AgentJam stack policy (stacks/<profile>/)
        ↓
AgentJam global policy (policies/)
        ↓
Environment defaults
        ↓
Model knowledge (Reasoning fallback only)
```
