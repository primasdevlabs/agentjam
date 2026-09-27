# Runtime Validators

Repository and manifest cross-reference validation.

## Checks Performed

- **Manifest syntax**: malformed YAML surfaces as `error` severity.
- **Required fields**: agents, workflows, and tools must declare a `description`.
- **Cross-references**: agent `skills`/`tools` and workflow `agents`/`skills` must resolve to discovered resources (`warning` severity).
- **Workflow integrity**: steps require unique non-empty `id`s.
- **Enum fields**: tool `safetyLevel` and policy `enforcement` must match the declared vocabularies.

## Usage

```bash
go run ./cmd/agentjam validate
```

`validator.HasErrors` reports whether any `error`-severity issue exists; the CLI exits non-zero in that case. Warnings are printed but non-blocking.
