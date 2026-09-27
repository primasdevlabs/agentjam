# AgentJam Tools

Canonical tools define the execution capabilities agents may invoke. Each lives under `tools/<id>/tool.yaml` with a declared `safetyLevel`.

## Catalog

| Tool | Safety | Live Handler (`agentjam mcp`) |
|---|---|---|
| `filesystem` | safe-write | yes — read / write / list within workspace root |
| `git` | safe-write | yes — status / diff / log / branch |
| `terminal` | destructive | yes — shell commands with denylist + timeout |
| `github` | safe-write | yes — `gh` CLI passthrough (pr/issue/repo view) |
| `run_workflow` | read-only | yes — meta-tool driving the workflow executor |
| `browser` | read-only | manifest only |
| `database` | safe-write | manifest only |
| `communication` | safe-write | manifest only |
| `mcp` | safe-write | manifest only |

- **[MCP Tools Catalog & Safety Levels](mcp-tools.md)** — capability and schema details.

## Safety Levels

`read-only` → `safe-write` → `destructive` → `admin`. The dispatcher blocks any tool above its configured permission ceiling, enforces workspace path containment, applies per-call timeouts, and recovers handler panics.

## Serving Tools

```bash
agentjam mcp   # newline-delimited JSON-RPC on stdio
```

Supports `initialize`, `tools/list`, `tools/call`, `ping`. Any MCP client (Claude Code, OpenCode, Zed, Cline, Roo Code) can attach.
