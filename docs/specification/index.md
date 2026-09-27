# AgentJam Specification

The normative specification for the AgentJam canonical format and tooling.

- **[Canonical Format](canonical-format.md)** — resource types, manifest schemas, directory layout, and referential integrity rules.
- **[CLI Specification](cli-spec.md)** — the `agentjam` command surface: `validate`, `context`, `export`, `run`, `mcp`, `build-registry`, `eval`, `preflight`, `detect`, `run-e2e`.

## Invariants

- Every canonical resource is a directory containing a manifest file plus optional supporting content.
- Manifests are YAML; instruction content is Markdown (optionally frontmatter-delimited).
- Cross-resource references resolve by `name` and are enforced by the validator.
- The generated `registry.json` and `registry/<type>s/index.json` files are derived artifacts — never hand-edit them.
