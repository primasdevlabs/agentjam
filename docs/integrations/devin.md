# Devin & Devin Desktop Integration

AgentJam provides complete Level 5 integration for Cognition AI's **Devin** (CLI, Autonomous Cloud Agent, and Devin Desktop).

---

## Integration Overview

| Feature | Support | Output Target |
| :--- | :---: | :--- |
| Environment Detection | ✅ Yes | Marker: `.devin` directory |
| Instruction Format | ✅ Markdown | Root `DEVIN.md` & `.devin/rules/*.md` |
| Agent Client Protocol (ACP) | ✅ Supported | Native sub-process execution |
| Model Routing & SWE-2 | ✅ Compatible | Harness configuration pass-through |

---

## Configuration & Usage

### 1. Project Configuration
AgentJam exports canonical rules to root `DEVIN.md` and project-level `.devin/rules/`:

```bash
# Export canonical agent rules for Devin / Devin Desktop
npx @agentjam/cli export --agent software-engineer --harness devin
```

This generates:
- `DEVIN.md` (root project instructions ingested automatically by Devin)
- `.devin/rules/software-engineer.md` (role-specific guidelines)

### 2. Devin Desktop Integration
Devin Desktop ingests rules automatically from `.devin/` and root `DEVIN.md` during session initialization.

Through the **Agent Client Protocol (ACP)**, Devin Desktop can also execute AgentJam agents as registered local sub-processes inside the Devin Desktop command center.
