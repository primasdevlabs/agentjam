# AgentJam CLI Specification

The AgentJam CLI (`@agentjam/cli`) provides command-line tools for workspace management, validation, and harness rule export.

---

## Commands

### `agentjam init`
Scaffolds `.agentjam/config.yaml`, `.agentjam/policies/`, and exports default harness configurations into current workspace.

### `agentjam validate`
Runs three-tier validation (`schema-validator`, `cross-ref-validator`, `instruction-validator`) over current workspace.

### `agentjam export --agent <id> --harness <name>`
Exports canonical agent instructions into target harness configuration format (`cursor`, `claude`, `gemini`, `devin`, `cline`, `windsurf`).

### `agentjam build-registry`
Scans current repository and generates root `registry.json` index.
