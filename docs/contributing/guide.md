# Contributor Guide

Thank you for contributing to AgentJam!

---

## Development Setup

1. **Clone & Install**:
   ```bash
   git clone https://github.com/primasdevlabs/agentjam.git
   cd agentjam
   npm install
   ```

2. **Build Packages**:
   ```bash
   npm run build
   ```

3. **Run Unit Tests**:
   ```bash
   npm test
   ```

4. **Run Repository Validation**:
   ```bash
   npm run validate
   ```

---

## Engineering Standards

- **Toolchain Preflight**: Verify that lint, format, typecheck, and test tools pass before committing changes.
- **Modularity**: Use public package APIs only. Never import private internal module files across package boundaries.
- **Security**: Never commit secrets or hardcoded API keys. Parameterize database queries.
- **Verification**: Always run `npm run build`, `npm test`, and `npm run validate` to ensure 100% clean passage.
