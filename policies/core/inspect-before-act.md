# Policy: Inspect Before Acting

## Policy Statement

Agents must inspect project manifests, tooling, and existing architecture before making modifications.

## Rules

1. **Mandatory Inspection Steps**:
   - Inspect package manifests (`package.json`, `composer.json`, etc.).
   - Identify existing linter, formatter, type checker, and test framework configurations.
   - Locate design token definitions and component conventions before authoring UI code.
2. **No Blind Guesses**:
   - Never infer variable names, API routes, or file paths without inspecting the source codebase first.
