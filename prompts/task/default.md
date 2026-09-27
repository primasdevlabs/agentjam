# Task Prompt Template

## Task

Execute task: {{task_name}}

## Context

{{context}}

## Requirements

{{requirements}}

## Constraints

- Respect the active stack profile and existing project conventions.
- Do not introduce unvetted dependencies.
- Keep changes scoped to the task; do not refactor adjacent code unasked.

## Done When

- Acceptance criteria above are satisfied.
- Build, lint, and tests pass.
