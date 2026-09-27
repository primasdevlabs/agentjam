# Runtime Context Manager

Assembles the system instruction snapshot injected into agent harnesses.

## Snapshot Contents

- Workspace root and generation timestamp.
- Active agent persona, stack profile, and environment.
- Active skills with their resolved `instructions/*.md` bodies.
- The enforced policy matrix (ID, name, enforcement level, description) and policy instruction bodies.
- Project-specific custom rules.

## Token Budgeting

`ContextOptions.TokenBudget` (default 128000) bounds the instruction. Token counts are estimated at ~4 characters per token; oversized snapshots are truncated with a `Context Truncated` marker.

## API

```go
cm := rt.GetContextManager()
snapshot := cm.BuildContextSnapshot(context.ContextOptions{
    ActiveAgent:  "software-engineer",
    ActiveStack:  "nextjs",
    ActiveSkills: []string{"testing"},
    TokenBudget:  128000,
})
```
