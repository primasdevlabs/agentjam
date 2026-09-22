# Policy: Model Knowledge Precedence

## Policy Statement

Model training data is a reasoning material, NOT a source of truth for framework APIs, package versions, recommended tooling, or current security practices.

## Rules

1. **Precedence Enforcement**:
   - `Current project state` > `Project config` > `Pinned AgentJam rules` > `Authoritative documentation` > `Model knowledge`.
2. **Fallback Verification**:
   - If an API or library convention cannot be verified through current documentation or project manifests, the agent must report that verification is unavailable.
3. **No Deprecated API Introduction**:
   - Do not introduce APIs or syntax marked deprecated in current official documentation.
