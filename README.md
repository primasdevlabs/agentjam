# AgentJam

> Open-source runtime, configuration layer, agents, skills, tools, workflows, and language registry for AI-powered development.

AgentJam is a community-driven, harness-agnostic ecosystem built on a core principle:

> **The model provides reasoning. AgentJam provides the rules, tools, current conventions, freshness requirements, language definitions, and constraints.**

AgentJam prevents AI coding agents from relying on outdated training memory or casually introducing unapproved dependencies. It establishes non-negotiable policies, language standards, and stack profiles that agents must obey.

---

## Language & Layer Hierarchy

```text
Language (e.g., TypeScript, PHP)
   ↓
Framework (e.g., Next.js, Laravel)
   ↓
Ecosystem (e.g., Web, Backend)
   ↓
Project Stack (Active configuration)
```

*The language policy establishes language conventions, the framework skill establishes framework conventions, and the project stack determines active constraints. Framework and project conventions override generic language defaults (e.g., Laravel conventions take precedence over generic PHP PSR defaults).*

---

## Core Principles

- **`TRAINING KNOWLEDGE ≠ SOURCE OF TRUTH`**: Enforces documentation freshness rules requiring agents to verify current official documentation.
- **Enforced Stack Boundaries**: Project stack profiles (`stacks/nextjs`, `stacks/laravel`, etc.) prevent agents from introducing unapproved libraries (e.g. Express into NestJS, MongoDB into PostgreSQL).
- **Design Governance System**: Declarative anti-slop rules (`policies/design/`) prohibiting generic SaaS cards, emoji UI icons (`🚀⚡🔥`), neon/indigo defaults, font slop, and AI marketing buzzwords.
- **Structured Language Registry**: Categorized language ecosystem specifications under `languages/` (15 ecosystems: `systems/`, `web/`, `jvm/`, `dotnet/`, `apple/`, `mobile/`, `backend/`, `data/`, `functional/`, `scripting/`, `infrastructure/`, `databases/`, `smart-contracts/`, `embedded/`, `legacy/`).
- **Harness-Agnostic Engine**: Export canonical rules and resources to Claude Code, Cursor, Codex, Gemini, OpenCode, Cline, Roo Code, or Generic setups.

---

## Repository Structure

```text
agentjam/
├── languages/                # Language registry organized across 15 ecosystems
├── policies/                 # Non-negotiable rules (stack, freshness, security, design, dependencies)
├── runtime/                  # Context resolvers, rule evaluators, post-generation validators
├── stacks/                   # Stack profiles (nextjs, laravel, nestjs, react, vue, django, rails)
├── agents/                   # Canonical AI agent personas (code-reviewer, software-engineer, etc.)
├── skills/                   # Domain-specific skill modules (code-review, testing, etc.)
├── tools/                    # Abstract tool capability interfaces (filesystem, git, etc.)
├── workflows/                # Single and multi-agent execution flows (including existing-project-audit)
├── prompts/                  # Modular prompt templates
├── templates/                # Community contributor starters
├── integrations/             # Target harness export adapters
├── registry/                 # Search indices and metadata catalogs
└── packages/                 # Monorepo TypeScript core libraries
```

---

## Monorepo Packages (`packages/`)

- **`@agentjam/core`**: Domain models, Zod schemas, and LanguageManifest definitions.
- **`@agentjam/parser`**: YAML manifest, Markdown instruction, and language parser.
- **`@agentjam/policy-engine`**: Non-negotiable policy and design slop evaluation engine.
- **`@agentjam/runtime`**: Project context resolver and rule validator.
- **`@agentjam/validator`**: Structural resource, policy, and language compliance checker.
- **`@agentjam/registry`**: Search index builder.
- **`@agentjam/adapters`**: Harness export engine.

---

## Quick Start

### Installation & Validation

```bash
# Validate canonical resources, policies, languages, and stack profiles
npx tsx scripts/validate.ts
```

---

## Contributing

We welcome contributions of new stack profiles, language definitions, policies, agents, skills, and harness adapters! Please read our [Contributing Guide](CONTRIBUTING.md) and check out [templates/](templates/) to get started.

---

## License

[MIT](LICENSE) © AgentJam Community
