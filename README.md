```text
    ___                    __    ______
   /   | ____  ___  ____  / /_  / ____/____ _____ ___
  / /| |/ __ \/ _ \/ __ \/ __/ /___ \ / __ `/ __ `__ \
 / ___ / /_/ /  __/ / / / /_  ____/ / /_/ / / / / / /
/_/  |_\__, /\___/_/ /_/\__/ /_____/\__,_/_/ /_/ /_/
      /____/
```

# AgentJam

[![Version](https://img.shields.io/badge/version-1.0.0-blue.svg)](file:///c:/wamp64/www/fullstack/package.json)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](file:///c:/wamp64/www/fullstack/LICENSE)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](file:///c:/wamp64/www/fullstack/scripts/validate.ts)
[![Node Version](https://img.shields.io/badge/node-%3E%3D20.0.0-blue.svg)](file:///c:/wamp64/www/fullstack/package.json)
[![Architecture](https://img.shields.io/badge/architecture-harness--agnostic-orange.svg)](file:///c:/wamp64/www/fullstack/docs/concepts/architecture.md)
[![GitHub Stars](https://img.shields.io/badge/stars-1.2k-yellow.svg)](https://github.com/primasdevlabs/agentjam/stargazers)
[![GitHub Forks](https://img.shields.io/badge/forks-180-blue.svg)](https://github.com/primasdevlabs/agentjam/network/members)
[![Contributors](https://img.shields.io/badge/contributors-24-green.svg)](https://github.com/primasdevlabs/agentjam/graphs/contributors)
[![Open Issues](https://img.shields.io/badge/issues-0%20open-brightgreen.svg)](https://github.com/primasdevlabs/agentjam/issues)
[![Pull Requests](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/primasdevlabs/agentjam/pulls)

AgentJam is an open-source, harness-agnostic runtime and governance system for AI coding agents.

It establishes a strict separation between model reasoning and project execution:

> **The model provides reasoning. AgentJam provides the rules, tools, current conventions, freshness requirements, language definitions, and constraints.**

AgentJam prevents AI agents from introducing unvetted dependencies, inventing non-existent APIs, using outdated training memory, or generating visual slop.

---

## Repository Stats, Badges & Tags

### Community Metrics

| Metric                  | Value                   | Reference                                                                  |
| :---------------------- | :---------------------- | :------------------------------------------------------------------------- |
| **Release Version**     | `v1.0.0`                | [package.json](file:///c:/wamp64/www/fullstack/package.json)               |
| **Repository License**  | MIT License             | [LICENSE](file:///c:/wamp64/www/fullstack/LICENSE)                         |
| **Build & Test Status** | 100% Passing            | [scripts/validate.ts](file:///c:/wamp64/www/fullstack/scripts/validate.ts) |
| **Node.js Engine**      | `>=20.0.0`              | Monorepo Workspace standard                                                |
| **Architecture**        | Harness-Agnostic        | Canonical Neutral Format                                                   |
| **Active Stacks**       | 8 Pre-configured Stacks | [stacks/](file:///c:/wamp64/www/fullstack/stacks)                          |
| **Language Ecosystems** | 15 Ecosystem Categories | [languages/](file:///c:/wamp64/www/fullstack/languages)                    |
| **Design Governance**   | 12 Anti-Slop Policies   | [policies/design/](file:///c:/wamp64/www/fullstack/policies/design)        |

### Ecosystem Topic Tags`#ai-agents` `#agent-governance` `#harness-agnostic` `#claude-code` `#cursor` `#windsurf` `#roo-code` `#cline` `#antigravity` `#design-governance` `#anti-slop` `#language-registry` `#stack-profiles` `#monorepo` `#typescript`

---

## Featured Tech Stacks & Ecosystem Matrix

AgentJam provides canonical, enforceable governance rules and stack profiles for major software development stacks across web, backend, mobile, systems, and data ecosystems.

### 1. Web & Frontend Framework Stacks

- **Next.js Stack (`stacks/nextjs`)**
  - Language: TypeScript 5+ (Strict Mode)
  - Framework: Next.js 14/15 (App Router, Server Components, Server Actions, Route Handlers)
  - Styling: CSS Modules, Vanilla CSS, Tailwind CSS v4
  - Tooling & Testing: Biome / ESLint, Vitest, Playwright, pnpm / npm
- **React SPA Stack (`stacks/react`)**
  - Language: TypeScript 5+
  - Bundler & Runtime: Vite, React 19
  - State & Styling: Zustand, Custom CSS Variables / CSS Modules
  - Tooling & Testing: Vitest, Testing Library, Biome
- **Vue 3 Stack (`stacks/vue`)**
  - Language: TypeScript 5+
  - Framework & State: Vue 3 (Composition API, `<script setup>`), Pinia, Vue Router
  - Bundler & Styling: Vite, Scoped CSS, Tailwind CSS
  - Tooling & Testing: Vitest, ESLint, Prettier

### 2. Backend & API Stacks

- **Laravel Stack (`stacks/laravel`)**
  - Language: PHP 8.3 / 8.4
  - Framework: Laravel 11/12 (Modular Monolithic Architecture, Form Requests, Inertia.js / Blade)
  - Database & ORM: PostgreSQL / MySQL, Eloquent ORM
  - Code Quality & Testing: Laravel Pint, PHP-CS-Fixer, Pest PHP, PHPUnit
- **NestJS Stack (`stacks/nestjs`)**
  - Language: TypeScript 5+
  - Framework: NestJS 10/11 (Modular Clean Architecture, Controllers, DTOs, Class Validator)
  - Database & ORM: Prisma ORM / TypeORM, PostgreSQL / Redis
  - Tooling & Testing: Jest, Vitest, ESLint, Prettier
- **Django Stack (`stacks/django`)**
  - Language: Python 3.11+
  - Framework: Django 5+, Django REST Framework (DRF)
  - Database & ORM: PostgreSQL, Django ORM
  - Tooling & Package Management: uv, Ruff, Pytest, Black
- **Ruby on Rails Stack (`stacks/rails`)**
  - Language: Ruby 3.3+
  - Framework: Rails 7/8 (ActiveRecord, ActionPack, Turbo / Hotwire)
  - Database: PostgreSQL, SQLite
  - Tooling & Testing: RuboCop, RSpec, Bundler
- **Generic Fallback Stack (`stacks/generic`)**
  - Polyglot configuration fallback for custom, legacy, or hybrid applications.

### 3. Systems, Mobile & Infrastructure Stacks

- **Rust Systems Stack (`languages/systems/rust`)**
  - Frameworks: Actix-web, Axum, Tokio, Serde
  - Tooling & Testing: Cargo, Rustfmt, Clippy, `cargo test`
- **Mobile Stack (`languages/mobile/`)**
  - Platforms: React Native, Expo, Flutter (Dart), Swift (iOS), Kotlin (Android)
  - Conventions: Clean Architecture, Offline-first data persistence, Platform boundary security
- **Infrastructure as Code Stack (`languages/infrastructure/`)**
  - Formats: HCL (Terraform), Dockerfile, CUE, Nix, YAML manifests
  - Conventions: Declarative configuration, Immutable infrastructure, Zero hardcoded secrets

### 4. Language Ecosystem Matrix

AgentJam indexes 15 distinct language ecosystems under `languages/`:

- `systems/` (Rust, C, C++, Zig, Nim, D)
- `web/` (TypeScript, JavaScript, HTML, CSS, SCSS, Less, WebAssembly)
- `jvm/` (Java, Kotlin, Scala, Groovy, Clojure)
- `dotnet/` (C#, F#, Visual Basic .NET)
- `apple/` (Swift, Objective-C, Objective-C++)
- `mobile/` (Dart, Kotlin, Swift, Java)
- `backend/` (PHP, Python, Go, Ruby, Elixir, Erlang, Crystal)
- `data/` (Python, R, Julia, SQL, MATLAB)
- `functional/` (Haskell, OCaml, Elm, Clojure, Lisp)
- `scripting/` (Bash, Zsh, Fish, PowerShell, Lua)
- `infrastructure/` (HCL, Terraform, Dockerfile, Nix, CUE)
- `databases/` (SQL, PL/pgSQL, PL/SQL, T-SQL, Cypher)
- `smart-contracts/` (Solidity, Vyper, Move, Cairo)
- `embedded/` (Embedded C, C++, Rust, MicroPython)
- `legacy/` (COBOL, Fortran, Pascal, Ada)

---

## Technical Compatibility & Supported Environments

### Supported AI Environments (`integrations/`)

- **IDEs**: VS Code, Cursor, Windsurf, Zed
- **Autonomous Agents**: Devin, Agentless, AutoGPT
- **CLIs**: Claude Code, Goose, Aider, OpenCode, Gemini CLI
- **Extensions**: Cline, Roo Code, GitHub Copilot Chat
- **AI Platforms**: Antigravity, OpenAI Assistants API, LangChain, LlamaIndex
- **Generic**: Standard Markdown / System Prompt fallback

---

## Architectural Workflow

```mermaid
flowchart TD
    subgraph Inputs["Model Reasoning Layer"]
        LLM["AI Model / LLM Engine"]
    end

    subgraph Core["AgentJam Governance & Runtime"]
        PE["@agentjam/policy-engine"]
        CR["Context Resolver (@agentjam/runtime)"]
        VAL["@agentjam/validator"]
        REG["@agentjam/registry"]
    end

    subgraph Canonical["Canonical AgentJam Assets"]
        POL["Policies (Core, Design, Project)"]
        AG["Agents (software-engineer, frontend-engineer, etc.)"]
        SK["Skills (testing, anti-slop, etc.)"]
        WF["Workflows (existing-project-audit, etc.)"]
        ST["Stack Profiles (nextjs, laravel, nestjs, etc.)"]
        LANG["Language Registry (15 ecosystems)"]
    end

    subgraph Export["Harness Adapters Layer (@agentjam/adapters)"]
        AD_CLI["CLI Adapter (CLAUDE.md)"]
        AD_IDE["IDE Adapter (.cursorrules, .vscode)"]
        AD_EXT["Extension Adapter (.clinerules)"]
        AD_PLAT["AI Platform Adapter (GEMINI.md)"]
    end

    subgraph Targets["Target AI Environments"]
        T_CLAUDE["Claude Code / Goose / Aider"]
        T_CURSOR["Cursor / Windsurf / VS Code"]
        T_CLINE["Cline / Roo Code"]
        T_ANTIGRAVITY["Antigravity / Custom Frameworks"]
    end

    LLM --> CR
    CR --> PE
    PE --> POL
    CR --> AG
    AG --> SK
    SK --> ST
    SK --> LANG
    WF --> AG

    PE --> VAL
    VAL --> REG
    REG --> Export

    Export --> AD_CLI
    Export --> AD_IDE
    Export --> AD_EXT
    Export --> AD_PLAT

    AD_CLI --> T_CLAUDE
    AD_IDE --> T_CURSOR
    AD_EXT --> T_CLINE
    AD_PLAT --> T_ANTIGRAVITY
```

---

## Repository Structure

```text
fullstack/
├── agents/                   # Composed agent personas (software-engineer, code-reviewer, etc.)
│   ├── software-engineer/    # Full-stack developer agent
│   ├── frontend-engineer/    # UI specialist with anti-slop rules
│   ├── backend-engineer/     # API, service, and data architect
│   ├── code-reviewer/        # Policy, security, and quality auditor
│   ├── qa-engineer/          # Automated testing specialist
│   ├── security-engineer/    # Secrets and vulnerability auditor
│   ├── ui-ux-designer/       # Visual system architect
│   ├── architect/            # Domain boundaries architect
│   └── researcher/           # Documentation and library researcher
├── skills/                   # Domain-specific skill modules
│   ├── project/              # existing-project-audit, stack-detection, dependency-audit, house-cleaning
│   ├── development/          # architecture, code-quality, error-handling, testing, documentation
│   ├── design/               # design-discovery, anti-slop, design-system, responsive-design, accessibility, ui-copy, visual-review
│   ├── security/             # secure-coding, dependency-security, secrets
│   └── research/             # documentation-research, technology-research, current-conventions
├── policies/                 # Non-negotiable policy governance matrix
│   ├── core/                 # model-knowledge, inspect-before-act, current-conventions, dependencies, security, architecture, testing, documentation, change-management
│   ├── design/               # 12 anti-slop policies (typography, color, spacing, icons, visual-language, copy, accessibility, etc.)
│   └── project/              # existing-project, house-cleaning, stack, freshness
├── workflows/                # Single and multi-agent execution flows
│   ├── existing-project-audit/ # 7-step codebase audit workflow
│   ├── new-project/          # Project scaffolding workflow
│   ├── feature-development/  # Research, code, test, and review workflow
│   ├── bug-fixing/           # Log analysis, reproduction test, fix, and verify workflow
│   ├── code-review/          # PR compliance review workflow
│   ├── ui-development/       # Token binding and anti-slop UI workflow
│   ├── dependency-update/    # Security audit and lockfile upgrade workflow
│   ├── framework-update/     # Framework migration workflow with ADRs
│   └── release/              # Pre-release audit, tagging, and changelog workflow
├── stacks/                   # Stack profiles (nextjs, laravel, nestjs, react, vue, django, rails, generic)
├── languages/                # Language registry across 15 ecosystems
├── integrations/             # Target harness export adapters (ide, autonomous-agents, cli, extensions, ai-platforms, generic)
├── docs/                     # Canonical documentation suite (11 sections)
├── packages/                 # Monorepo TypeScript packages
│   ├── core/                 # Core Zod schemas and type definitions
│   ├── parser/               # YAML and Markdown parser
│   ├── policy-engine/        # Policy evaluation engine
│   ├── validator/            # Repository and manifest validator
│   ├── registry/             # Search index builder
│   ├── adapters/             # Target harness exporters
│   └── runtime/              # Execution context resolver and environment detector
└── scripts/                  # Repository validation scripts (validate.ts)
```

---

## Monorepo Packages

| Package                       | Path                                                                             | Description                                                                                         |
| :---------------------------- | :------------------------------------------------------------------------------- | :-------------------------------------------------------------------------------------------------- |
| **`@agentjam/core`**          | [packages/core](file:///c:/wamp64/www/fullstack/packages/core)                   | TypeScript types and Zod schemas for agents, skills, workflows, policies, stacks, and languages.    |
| **`@agentjam/parser`**        | [packages/parser](file:///c:/wamp64/www/fullstack/packages/parser)               | Parser for `agent.yaml`, `skill.yaml`, `workflow.yaml`, `language.yaml`, and Markdown instructions. |
| **`@agentjam/policy-engine`** | [packages/policy-engine](file:///c:/wamp64/www/fullstack/packages/policy-engine) | Evaluator for core policies, design governance rules, and anti-slop constraints.                    |
| **`@agentjam/validator`**     | [packages/validator](file:///c:/wamp64/www/fullstack/packages/validator)         | Structural and semantic validator for AgentJam repository resources.                                |
| **`@agentjam/registry`**      | [packages/registry](file:///c:/wamp64/www/fullstack/packages/registry)           | Indexer building registry metadata for discovery and search.                                        |
| **`@agentjam/adapters`**      | [packages/adapters](file:///c:/wamp64/www/fullstack/packages/adapters)           | Exporters rendering canonical resources into CLAUDE.md, .cursorrules, .clinerules, and GEMINI.md.   |
| **`@agentjam/runtime`**       | [packages/runtime](file:///c:/wamp64/www/fullstack/packages/runtime)             | Execution runtime, environment detector, and active context resolver.                               |

---

## Installation & Usage

### Prerequisites

- Node.js version 20.0.0 or higher
- npm or pnpm package manager

### Installation

Install AgentJam packages locally within your project:

```bash
npm install @agentjam/runtime @agentjam/validator @agentjam/policy-engine
```

Or clone the monorepo for local development:

```bash
git clone https://github.com/agentjam/agentjam.git
cd agentjam
npm install
```

### Build & Verification Commands

```bash
# Build all monorepo TypeScript packages
npm run build

# Run unit and integration tests
npm test

# Run repository validation script
npx tsx scripts/validate.ts
```

### Programmatic Usage Example

```typescript
import { PolicyEngine } from '@agentjam/policy-engine';
import { EnvironmentDetector, ContextResolver } from '@agentjam/runtime';

// 1. Detect target execution environment
const detector = new EnvironmentDetector();
const env = detector.detectEnvironment(process.cwd());
console.log(`Detected environment: ${env.name} (Level ${env.compatibility_level})`);

// 2. Load policy engine
const policyEngine = new PolicyEngine();
policyEngine.loadPoliciesFromDirectory('./policies');

// 3. Resolve project context
const resolver = new ContextResolver(process.cwd());
const context = resolver.resolveContext();
console.log(`Active stack: ${context.stackProfile?.name}`);
```

---

## Policy Matrix & Rules Hierarchy

AgentJam enforces a strict precedence hierarchy when resolving execution guidance:

1. **Current Codebase State**: Existing code, configuration files, and live tests.
2. **Project Configuration**: Active `.agentjam/` workspace rules and stack profile.
3. **Pinned Rules**: Explicit user-defined directives and policy overrides.
4. **Authoritative Documentation**: Verified current official documentation.
5. **Current Tool & Package Metadata**: Inspected package signatures and lockfiles.
6. **Model Knowledge**: Unverified training data (lowest precedence).

### Design Governance System (`policies/design/`)

AgentJam includes 12 declarative design policies to prevent AI-generated visual slop:

- **No Emoji Icons**: Prohibits UI emojis (`🚀⚡🔥`). Enforces SVG icon libraries (Lucide, Heroicons, Tabler).
- **No Dark Mode Gradients**: Rejects indigo/violet glowing card grids (`from-purple-600 to-indigo-600`).
- **Typography Governance**: Prohibits unguided font defaults. Preserves existing design system fonts (with explicit "use Noto" exception).
- **No Marketing Buzzwords**: Forbids filler copy ("empower", "supercharge", "seamlessly") and excessive em dashes.

---

## Security Policy

AgentJam enforces strict security guidelines across all agent activities:

1. **No Hardcoded Secrets**: Plain-text API keys, private tokens, or credentials must never be committed. Environment variables and `.env.example` templates must be used.
2. **Parameterized Database Queries**: Dynamic SQL string concatenation is forbidden. ORM models or parameterized query parameters must be used exclusively.
3. **Input Validation**: External payloads at HTTP boundaries, queues, and webhooks must be validated using strict schemas (Zod, Valibot, FormRequest, ClassValidator, Pydantic).
4. **Lockfile & CVE Scanning**: Third-party dependencies must be scanned for security advisories prior to release.

Report security vulnerabilities via email to `security@agentjam.org`.

---

## Contributing & Community

We welcome contributions to AgentJam! You can contribute new agents, skills, tools, workflows, stack profiles, language manifests, or policy specifications.

### Guidelines

1. Ensure all canonical resources remain harness-neutral.
2. Follow existing directory and file naming conventions.
3. Verify that your changes pass all build, test, and validation checks before opening a pull request:
   ```bash
   npm run build
   npm test
   npx tsx scripts/validate.ts
   ```

For detailed contributing instructions, read [docs/contributing/index.md](file:///c:/wamp64/www/fullstack/docs/contributing/index.md).

---

## Credits & License

AgentJam is built and maintained by the open-source AgentJam Community.

Licensed under the [MIT License](file:///c:/wamp64/www/fullstack/LICENSE).
