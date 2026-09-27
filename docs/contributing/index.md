# Contributing to AgentJam

See the [full contribution guide](guide.md) for the complete process.

## Quick Start

1. **Add a resource** — copy the matching `templates/<type>/` manifest into the canonical directory tree and fill in the fields.
2. **Reference existing resources** — agents list skill and tool names; workflows list agents and skills. `agentjam validate` catches dangling references.
3. **Verify before shipping**:

```bash
go build ./...
go test ./...
go run ./cmd/agentjam validate
go run ./cmd/agentjam build-registry
```

## Rules

- Keep resource IDs unique — duplicates are validation errors.
- Regenerate the registry (`agentjam build-registry`) after adding or removing resources.
- Follow the repository's `AGENTS.md` guardrails: package boundaries, parameterized queries, sanitized inputs, no hardcoded secrets.
