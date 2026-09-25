# Contributing to AgentJam

Thank you for helping build AgentJam! AgentJam is a harness-agnostic, open-source collection of AI agents, skills, tools, and workflows.

## How to Contribute a Resource

1. **Choose a Resource Type**:
   - `agents/`: AI personas with role definitions and required skills.
   - `skills/`: Reusable task-specific instruction modules.
   - `tools/`: Abstract tool capability interfaces.
   - `workflows/`: Orchestrated execution flows.

2. **Use Templates**:
   Copy a starter template from `templates/<type>/` to your new resource location (e.g. `skills/my-new-skill/`).

3. **Fill in Metadata & Instructions**:
   - Create `*.yaml` manifest adhering to the canonical schema.
   - Create modular Markdown instructions inside `instructions/`.

4. **Validate**:
   Run the validation script to verify schema compliance:
   ```bash
   go run ./cmd/agentjam validate
   ```

5. **Submit a Pull Request**:
   Push your changes and open a PR using our PR template.

## Code of Conduct

Please review and follow our [Code of Conduct](CODE_OF_CONDUCT.md).
