# AgentJam Policy Matrix

Policies establish non-negotiable governance rules:

## Core Policies (`policies/core/`)
* `model-knowledge.md`: Model training data is not authoritative.
* `inspect-before-act.md`: Inspect existing project state before editing.
* `current-conventions.md`: Align with active project style.
* `dependencies.md`: Dependency security and pin requirements.
* `security.md`: Zero hardcoded secrets or raw string SQL.
* `architecture.md`: Modular boundaries and clean APIs.
* `testing.md`: Colocated tests and mandatory build verification.
* `documentation.md`: Direct technical documentation without marketing hype.
* `change-management.md`: Scoped edits and rollback readiness.

## Design Governance (`policies/design/`)
12 policies banning visual slop (`anti-slop.md`, `visual-language.md`, `typography.md`, `color.md`, `spacing.md`, `components.md`, `icons.md`, `animation.md`, `responsive.md`, `accessibility.md`, `copy.md`, `design-system.md`).

## Project Policies (`policies/project/`)
* `existing-project.md`: Existing project rules & overrides.
* `house-cleaning.md`: Cleanup and dead-code removal protocol.
* `stack.md`: Stack binding conventions.
* `freshness.md`: Freshness requirements for documentation & packages.
