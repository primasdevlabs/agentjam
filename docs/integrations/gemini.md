# Gemini CLI Integration

AgentJam provides Level 5 integration for **Gemini CLI**.

---

## Configuration & Export

Export workspace rules for Gemini:

```bash
npx @agentjam/cli export --agent software-engineer --harness gemini
```

This generates:
- `GEMINI.md` at project root
- `.gemini/skills/` skill instruction modules
