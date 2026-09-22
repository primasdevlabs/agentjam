# AgentJam Core Philosophy

AgentJam is an open-source governance and runtime environment that establishes the rules, tools, current conventions, freshness requirements, and constraints AI agents operate within.

> **The model provides reasoning. AgentJam provides the rules, tools, current conventions, freshness requirements, and constraints.**

---

## Principle 1: Model Knowledge is Not Authoritative

The model's training memory must NEVER be treated as the source of truth for current APIs, package versions, framework conventions, security practices, or CLI commands.

AgentJam enforces the following precedence hierarchy:

```text
Current project state
        ↓
Project configuration
        ↓
Pinned AgentJam rules
        ↓
Authoritative current documentation
        ↓
Current package/tool metadata
        ↓
Model knowledge (Reasoning material only)
```

If current information cannot be verified, the agent must report that verification is unavailable rather than presenting potentially outdated training memory as current truth.

---

## Principle 2: Inspect Before Acting

An agent MUST understand its environment before making significant code or architecture modifications.

Before writing code, inspect:
- Project structure and domain boundaries
- Active package manager and manifest versions
- Configured framework, runtime, and language standards
- Existing design tokens and component library
- Test framework, linter, and formatter configurations
- CI/CD environment and available tools

Never make assumptions when the project can be inspected directly.

---

## Principle 3: Existing Projects Get an Audit First

When AgentJam is installed into an existing project, agents must NOT immediately start writing new feature code.

The mandatory workflow sequence is:

```text
Inspect
   ↓
Audit (Structure, Tech, Design, Architecture, Security)
   ↓
Recommend House Cleaning (Categorized Critical / Recommended / Optional)
   ↓
Establish Project Context & Stack Rules
   ↓
Resolve Current Conventions
   ↓
Implement Feature
   ↓
Validate (Build, Lint, Typecheck, Test)
```

---

## Principle 4: Follow the Project's Established System

If an existing codebase has an established, coherent architecture, design token system, or coding convention, preserve it.

AgentJam is a governance layer designed to protect project integrity, NOT an excuse to rewrite existing codebases to suit arbitrary AI preferences.

---

## Principle 5: Prefer Composition

AgentJam capabilities are modular and composable:

```text
Agent
  ↓
Skills
  ↓
Tools
  ↓
Policies
  ↓
Workflows
```

Agents compose reusable skills and policies rather than duplicating instruction blocks across personas.
