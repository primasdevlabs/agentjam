# AgentJam Integrations

AgentJam is environment-neutral: one canonical project standard, exported to many AI development environments.

## Integration Catalog

| Integration | Category | Level | Output |
|---|---|---|---|
| [Claude Code](claude-code.md) | CLI | 6 — Lifecycle | `CLAUDE.md` |
| [Codex](codex.md) | CLI | 4 — Workflows | `codex.config.json` |
| [Gemini CLI](gemini.md) | CLI | 5 — Validation | `GEMINI.md` |
| [OpenCode](opencode.md) | CLI | 4 — Workflows | `opencode.json` |
| [Cursor](cursor.md) | IDE | 5 — Validation | `.cursorrules`, `.cursor/rules/` |
| [VS Code](vscode.md) | IDE | 5 — Validation | `.vscode/settings.json` |
| [Windsurf](windsurf.md) | IDE | 5 — Validation | `.windsurfrules` |
| [Zed](zed.md) | IDE | 4 — Workflows | `.zed/zed.json` |
| [Void](overview.md) | IDE | 4 — Workflows | `.void/void.json` |
| [Cline](cline.md) | Extension | 5 — Validation | `.clinerules` |
| [Roo Code](roo-code.md) | Extension | 5 — Validation | `.roo/.roomodes` |
| [Devin](devin.md) | Autonomous Agent | 5 — Validation | `.devin/DEVIN.md` |
| [Antigravity](antigravity.md) | AI Platform | 6 — Lifecycle | `.gemini/GEMINI.md` |
| [Generic](generic.md) | Fallback | 1 — Instructions | `AGENTJAM.md` |

## Further Reading

- [Environment Compatibility Overview](overview.md) — taxonomy and compatibility levels 0–6.
- [Compatibility Matrix](compatibility.md) — per-capability breakdown across environments.

## Usage

```bash
agentjam export --harness auto      # detect environment and export matching rules
agentjam export --harness all       # write every supported harness file
agentjam export --harness windsurf  # target one environment
agentjam mcp                        # serve workspace tools over MCP stdio
```
