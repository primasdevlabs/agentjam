# Contributor Guide

Thank you for contributing to AgentJam!

---

## Development Setup

1. **Clone & Setup**:
   ```bash
   git clone https://github.com/primasdevlabs/agentjam.git
   cd agentjam
   ```

2. **Build Go Binaries**:
   ```bash
   go build ./...
   ```

3. **Run Unit Tests**:
   ```bash
   go test ./...
   ```

4. **Run Repository Validation**:
   ```bash
   go run ./cmd/agentjam validate
   ```

---

## Engineering Standards

- **Toolchain Preflight**: Verify that build, test, and validation tools pass before committing changes.
- **Modularity**: Use public Go package APIs only (`pkg/*`). Never import internal private packages across module boundaries.
- **Security**: Never commit secrets or hardcoded API keys. Parameterize database queries.
- **Verification**: Always run `go build ./...`, `go test ./...`, and `go run ./cmd/agentjam validate` to ensure 100% clean passage.
