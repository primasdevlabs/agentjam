# Generic Fallback Integration

- **Type**: `generic`
- **Compatibility Level**: Level 1 (Instructions)
- **Manifest**: `integrations/generic/integration.yaml`
- **Output**: `AGENTJAM.md`

## Capabilities

Instructions, skills, filesystem, project rules, and context files.

## Usage

```bash
agentjam export --harness generic
```

The generic exporter writes a portable `AGENTJAM.md` system prompt suitable for any harness that accepts markdown instructions — the baseline for environments without a dedicated integration.
