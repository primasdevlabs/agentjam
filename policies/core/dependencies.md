# Policy: Dependency Governance

## Policy Statement

Agents must respect project dependency boundaries and prevent unapproved package inflation.

## Rules

1. **No Forbidden Dependencies**:
   - Check `stacks/<profile>/` and project manifests before introducing packages.
2. **No Duplicate Libraries**:
   - Do not install competing libraries (e.g. Axios when Fetch/Native HTTP is configured, second icon library).
3. **Audit First**:
   - Verify package compatibility and maintenance status before recommending updates.
