# Getting Started with AgentJam

AgentJam is an open-source governance framework, canonical resource library, and runtime for AI-powered software engineering.

---

## Quick Navigation

- **[Installation & Operational Behavior](installation.md)**: How to install AgentJam and what it does inside a workspace where it is installed.
- **[System Architecture](../concepts/architecture.md)**: Monorepo package structure and system data flow.
- **[Design Governance](../concepts/design-governance.md)**: The 12 design policies and anti-slop rules.
- **[Precedence Pipeline](../concepts/precedence-pipeline.md)**: 7-level rule resolution hierarchy.
- **[Canonical Agent Catalog](../agents/catalog.md)**: Overview of the 13 canonical agent roles.

---

## The Core AgentJam Workflow

1. **Initialize Workspace**: Run `npx agentjam init` in your project root to generate `.agentjam/config.yaml` and export harness configurations (`AGENTS.md`, `.cursorrules`, `CLAUDE.md`, `GEMINI.md`).
2. **Toolchain Preflight**: AgentJam inspects your project stack manifests (`package.json`, `composer.json`, etc.) and verifies linting, formatting, typechecking, and test scripts.
3. **Execute AI Development**: AI assistants operate under AgentJam precedence rules (User Explicit Directives -> Project Rules -> Global Policies).
4. **Empirical Verification**: Run `npm run validate` to verify schema compliance, cross-reference integrity, security rules, and design governance.
