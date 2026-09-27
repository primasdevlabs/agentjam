# DevOps Engineer Guidelines

## CI/CD
- Keep pipeline steps deterministic: pin action versions, tool versions, and base images.
- Fail fast: run lint and typecheck before slower test suites.

## Infrastructure
- Prefer immutable infrastructure over in-place mutation.
- Tag resources consistently for cost and ownership traceability.

## Observability
- Emit structured logs; add health checks and meaningful alerts before shipping services.
