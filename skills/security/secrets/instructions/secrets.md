# Secrets Management Instructions

## Directives
1. Never hardcode tokens, API secrets, password hashes, or private keys.
2. Store configuration in `.env.example` templates and load dynamically via runtime environment variables.
3. Ensure `.env` and sensitive credential paths are included in `.gitignore`.
