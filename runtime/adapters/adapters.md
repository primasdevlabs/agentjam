# Environment Harness Adapters

Injects runtime context and non-negotiable policy rules into target harness instruction files.

## Supported Targets

| Harness | Output file | Function |
| --- | --- | --- |
| Claude Code | `CLAUDE.md` | `ExportClaudeCode` |
| Cursor | `.cursorrules` | `ExportCursorRules` |
| Gemini / Antigravity | `GEMINI.md` | `ExportGeminiMarkdown` |
| Cline | `.clinerules` | `ExportClineRules` |
| Windsurf | `.windsurfrules` | `ExportWindsurfRules` |
| Devin | `.devin/playbook.md` | `ExportDevinConfig` |
| Roo Code | `.roomodes` | `ExportRooCodeModes` |
| Generic | `SYSTEM_PROMPT.md` | `ExportGenericPrompt` |

## Export Flow

1. `ContextManager.BuildContextSnapshot` assembles the system instruction (policies, skill instructions, custom rules).
2. `adapters.ExportAllHarnesses` renders the instruction into every supported harness file.
3. `agentjam export --harness <name>` writes the selected files to the workspace root.

## Contract

Each exporter returns an `ExportResult` (`Harness`, `Files`, `Warnings`). Files are path → content maps; the CLI creates parent directories before writing.
