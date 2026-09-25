# Installation & Workspace Operational Behavior

This guide details how to install AgentJam Go Native Engine and **what AgentJam actively enforces inside a workspace where it is installed**.

---

## Installation

Install AgentJam CLI and runtime using Go:

```bash
# Clone repository
git clone https://github.com/agentjam/agentjam.git
cd agentjam

# Build native Go binary
go build -o agentjam ./cmd/agentjam
```

Or run directly using `go run`:

```bash
go run ./cmd/agentjam validate
```

---

## Harness Export & Rule Deployment

AgentJam supports both **auto-detection** and **explicit harness selection**:

```bash
# Auto-detect workspace AI environment (Cursor, Antigravity, Claude Code, etc.)
agentjam export

# Explicitly export rules for a specific AI harness
agentjam export --harness cursor
agentjam export --harness claude-code
agentjam export --harness gemini
agentjam export --harness all
```

This scaffolds `.agentjam/config.yaml`, `.agentjam/policies/`, and exports harness rule configurations (`AGENTS.md`, `.cursorrules`, `CLAUDE.md`, `GEMINI.md`, `.clinerules`) into your repository root.

---

## What AgentJam DOES in an Installed Repository

When an AI coding agent operates inside an AgentJam-enabled workspace, AgentJam actively enforces the following operational guardrails:

### 1. Toolchain Preflight (Before Writing Code)
- **Stack Detection**: Automatically detects workspace stack from manifests (`package.json`, `composer.json`, `pyproject.toml`, `go.mod`).
- **Tooling Verification**: Ensures `lint`, `format`, `typecheck`, and `test` scripts are present and passing before shipping changes.
- **Config Alignment**: Forces new code to conform to existing ESLint, Prettier, Pint, or TypeScript compiler configurations.

### 2. Modular Boundary Enforcement
- **Public API Isolation**: Restricts cross-feature imports to exported public APIs—prohibiting imports from other features' private internals (`/pkg/internal/`).
- **Single Responsibility**: Enforces small, composable units with colocated tests and type definitions.

### 3. Security Governance
- **Boundary Input Validation**: Enforces input validation schemas at HTTP, queue, and webhook boundaries.
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
- Never declares a task complete without executing empirical verification (`go build ./...`, `go test ./...`, `agentjam validate`).
