# AgentJam Workspace Guidance (AGENTS.md)

## Operational Rules & Guardrails
1. **Toolchain Preflight**: Run lint, format, typecheck, and unit tests before shipping changes.
2. **Modular Boundaries**: Use public package APIs only. Never import private internal paths.
3. **Security Standards**: Parameterize SQL queries, sanitize inputs, and use secret managers.
4. **Design Governance**: Follow HSL color tokens, use SVG icons (no emojis), and write human copy.
5. **Empirical Verification**: Verify changes with test and build commands before declaring task completion.
