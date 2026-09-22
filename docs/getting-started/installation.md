# Installation & Workspace Operational Behavior

This guide details how to install AgentJam and **what AgentJam actively enforces and does inside a workspace where it is installed**.

---

## Installation

Install AgentJam packages into your project workspace:

```bash
npm install @agentjam/runtime @agentjam/validator @agentjam/policy-engine
```

Or initialize AgentJam governance rules in any target repository:

```bash
npx @agentjam/cli init
```

This scaffolds `.agentjam/config.yaml`, `.agentjam/policies/`, and exports harness rules (`AGENTS.md`, `.cursorrules`, `CLAUDE.md`, `GEMINI.md`, `.clinerules`) into your repository root.

---

## What AgentJam DOES in an Installed Repository

When an AI coding agent operates inside an AgentJam-enabled workspace, AgentJam actively enforces the following operational guardrails:

### 1. Toolchain Preflight (Before Writing Code)
- **Stack Detection**: Automatically detects workspace stack from manifests (`package.json`, `composer.json`, `pyproject.toml`, `go.mod`).
- **Tooling Verification**: Ensures `lint`, `format`, `typecheck`, and `test` scripts are present and passing before shipping changes.
- **Config Alignment**: Forces new code to conform to existing ESLint, Prettier, Pint, or TypeScript compiler configurations.

### 2. Modular Boundary Enforcement
- **Public API Isolation**: Restricts cross-feature imports to exported public APIs—prohibiting imports from other features' private internals (`/src/internal/`).
- **Single Responsibility**: Enforces small, composable units with colocated tests and type definitions.

### 3. Security Governance
- **Boundary Input Validation**: Enforces input validation schemas (Zod, Valibot, standard validator) at HTTP, queue, and webhook boundaries.
- **Parameterized Database Queries**: Prohibits string-concatenated raw SQL queries; enforces parameterized queries or ORM models.
- **Secret Protection**: Scans for hardcoded secrets and API keys, requiring environment variables or secret managers.

### 4. Quality & Error Handling
- **Explicit Error Context**: Prohibits empty catch blocks; requires structured errors with machine-readable codes and context bags (`module`, `operation`, `correlationId`).
- **Performance Optimization**: Flags N+1 database queries on hot paths and requires pagination/indexing.

### 5. Anti-Slop Design Governance
- **Typography**: Replaces unstyled browser defaults with curated design fonts.
- **Color System**: Prohibits plain CSS color names (`red`, `blue`); enforces HSL design system tokens (`var(--color-primary-500)`).
- **SVG Icons Only**: Prohibits emoji characters as UI icons (`🚀⚡🔥`); enforces SVG icon libraries (Lucide, Heroicons).
- **Human Product Copy**: Prohibits AI marketing buzzwords ("empower your business", "unlock the power").

### 6. Verification Before Task Completion
- Never declares a task complete without executing empirical verification (`npm run build`, `npm test`, `npm run validate`).
