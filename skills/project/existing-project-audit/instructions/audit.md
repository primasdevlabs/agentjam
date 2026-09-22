# Existing Project Audit Skill

This skill provides step-by-step guidance for auditing an existing codebase before attempting modification or feature creation.

## Objectives
1. Inspect project directory layout, configuration manifests, and package managers.
2. Identify existing architectural patterns, module boundaries, and conventions.
3. Discover design tokens, UI component structure, typography, and styling setup.
4. Uncover security vulnerabilities, hardcoded secrets, and legacy code anti-patterns.
5. Produce a comprehensive Project Profile (`.agentjam/project-profile.yaml`) and Audit Report.

## Audit Workflow
1. **Manifest Inspection**: Read `package.json`, `composer.json`, `pyproject.toml`, `Cargo.toml`, etc.
2. **Directory & Structure Analysis**: Map route structure, component directories, domain models, and asset organization.
3. **Conventions Audit**: Identify linting rules (`.eslintrc`, `biome.json`, `pint.json`), code styling conventions, and TypeScript strictness.
4. **Design Governance Audit**: Detect design tokens, Tailwind / CSS variables, icon libraries, and visual pattern consistency.
5. **Security & Secrets Scan**: Verify absence of plain-text API keys, missing parameterization in database queries, or loose CORS/auth policies.
