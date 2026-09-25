# AgentJam Workspace Guidance (AGENTS.md)

## Operational Rules & Guardrails
1. **Toolchain Preflight**: Run Go build (`go build ./...`) and Go unit tests (`go test ./...`) before shipping changes.
2. **Modular Boundaries**: Respect Go package boundaries (`pkg/core`, `pkg/policy`, `pkg/runtime`, `pkg/context`, `pkg/memory`, `pkg/dispatcher`). Use public package functions only.
3. **Security Standards**: Parameterize SQL queries, sanitize inputs, and use secret managers.
4. **Design Governance**: Follow HSL color tokens, use SVG icons (no emojis), and write human copy.
5. **Empirical Verification**: Verify changes with Go build and test commands before declaring task completion.
