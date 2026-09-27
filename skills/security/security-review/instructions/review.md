# Security Review Instructions

## Protocol
1. Scan diffs for hardcoded secrets, tokens, and credentials before anything else.
2. Trace external input (HTTP, queues, webhooks) to sinks; verify validation and parameterization.
3. Check authorization checks exist at every privileged operation, not just at routing.
4. Flag new dependencies that are unmaintained or carry known advisories.
