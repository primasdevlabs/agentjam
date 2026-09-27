# CI/CD Instructions

## Protocol
1. Order stages fast-to-slow: lint, typecheck, unit tests, integration tests, build, deploy.
2. Pin every external action, image, and tool version.
3. Cache dependency installs keyed on lockfile hashes.
4. Gate deployments on green tests and require explicit approval for production.
