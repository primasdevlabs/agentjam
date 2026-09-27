# Runtime Rules Engine

Evaluates non-negotiable policies prior to and during LLM execution.

## Responsibilities

- Loads policy manifests from `policies/` via `PolicyEngine.LoadPolicies`.
- Evaluates source and UI content through `EvaluateDesignRules` and `EvaluateSecurityRules`.
- Aggregates violations by enforcement level (`strict-block`, `warning`, `info`) into a `PolicyEngineSummary`.
- Blocks execution when any `strict-block` violation is present (`summary.Allowed == false`).

## Evaluation Path

1. `agentjam eval <file>` reads the target file.
2. The engine runs each built-in evaluator gated by loaded policy categories.
3. Violations are emitted with `policyId`, `ruleName`, `line`, and enforcement level.
4. The CLI exits non-zero when `Allowed` is false.

## Adding a Rule

Add a manifest under `policies/<category>/` (YAML for structured rules, Markdown for narrative policy). Set `enforcement` to control blocking behavior.
