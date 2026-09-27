# Core Concepts

AgentJam defines a canonical, environment-neutral standard for agentic development workspaces.

- **[Architecture](architecture.md)** — how the Go runtime composes parser, policy engine, validator, registry, memory, dispatcher, toolchain, and context manager.
- **[Philosophy](philosophy.md)** — the design rationale: one project standard, many AI environments, freshness-verified knowledge.

## The Canonical Model

Everything in the repository is a typed resource discovered from a conventional directory layout:

| Resource | Directory | Manifest |
|---|---|---|
| Agent personas | `agents/<id>/` | `agent.yaml` + `instructions/` |
| Skills | `skills/<category>/<id>/` | `skill.yaml` + `instructions/` |
| Tools | `tools/<id>/` | `tool.yaml` |
| Workflows | `workflows/<id>/` | `workflow.yaml` |
| Policies | `policies/<category>/<id>/` | `policy.yaml` |
| Stacks | `stacks/<id>/` | `stack.yaml` |
| Languages | `languages/<ecosystem>/<id>/` | `language.yaml` |
| Integrations | `integrations/<category>/<id>/` | `integration.yaml` |

`agentjam validate` enforces referential integrity across all of them.
