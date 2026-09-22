# Changelog

All notable changes to AgentJam will be documented in this file.

## [0.4.0] - 2026-09-22

### Added

- **AI Development Environment Compatibility (Section 22)**:
  - Environment compatibility taxonomy under `integrations/`: `ide/`, `autonomous-agents/`, `cli/`, `extensions/`, `ai-platforms/`, `generic/`.
  - Compatibility Level system (Levels 0 through 6: Documentation, Instructions, Skills, Tools, Workflows, Validation, Lifecycle Integration).
  - Capability-based capability schema (instructions, skills, agents, workflows, tools, filesystem, terminal, browser, mcp, project_rules, context_files, hooks, lifecycle_events, validation).
  - `EnvironmentDetector` and `PrecedenceOrder` evaluation engine in `@agentjam/runtime`.
  - Comprehensive documentation suite under `docs/integrations/` (`overview.md`, `compatibility.md`, `cursor.md`, `vscode.md`, `claude-code.md`, `antigravity.md`, `generic.md`, etc.).

## [0.3.0] - 2026-09-22

### Added

- **Language Registry & Policy Layer (`languages/`)**: Organized language specifications across 15 ecosystem categories.
- **Hierarchy Model**: Defined `Language` -> `Framework` -> `Ecosystem` -> `Project Stack` precedence ordering.

## [0.2.0] - 2026-09-22

### Added

- **Design Governance System (`policies/design/`)**: Added 12 non-negotiable design governance policies.
- **Existing Project Audit Workflow (`workflows/existing-project-audit/`)**: Mandatory house-cleaning workflow.

## [0.1.0] - 2026-09-22

### Added

- Initial harness-agnostic repository architecture and monorepo TypeScript packages.
