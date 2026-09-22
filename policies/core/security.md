# Policy: Security & Secret Management

## Policy Statement

Zero tolerance for committed secrets, unparameterized database queries, and unvalidated input boundaries.

## Rules

1. **Secrets Governance**:
   - Never commit or hardcode API keys, passwords, or credentials. Use environment variables.
2. **Parameterized Queries**:
   - Parameterize all SQL and ORM queries to prevent injection attacks.
3. **Boundary Input Validation**:
   - Validate and sanitize external input at HTTP endpoints, webhooks, and message queues.
