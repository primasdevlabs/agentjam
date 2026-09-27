# Review Prompt Template

## Scope

Review the following changes: {{change_set}}

## Review Checklist

1. **Correctness** — logic errors, edge cases, off-by-one boundaries, nil handling.
2. **Security** — hardcoded secrets, unparameterized queries, missing input validation, authz gaps.
3. **Conventions** — consistency with existing codebase patterns and naming.
4. **Testing** — new behavior covered by tests; assertions meaningful, not smoke-only.
5. **Policy** — no violations of the enforced policy matrix.

## Output Format

For each finding: severity (critical | recommended | optional), file/line, issue, and suggested fix. End with an overall verdict: approve, approve-with-comments, or request-changes.
