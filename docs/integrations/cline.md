# Cline Integration

AgentJam provides complete integration support for **Cline** (formerly Claude Dev), the autonomous coding extension for VS Code.

---

## Capabilities & Mapping

| Feature | Support | Details |
| :--- | :---: | :--- |
| Integration Type | Extension | VS Code / VSCodium |
| Configuration Target | `.clinerules` | Root workspace rule injection |
| Level Compatibility | Level 5 | Full instructions, skills, tools, and validation |

---

## Setup & Rules Export

Export canonical AgentJam instructions to `.clinerules` format:

```bash
npx @agentjam/cli export --agent software-engineer --harness cline
```

This generates `.clinerules` at the project root, defining system role boundaries, code standards, and execution rules.
