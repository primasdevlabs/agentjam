# AgentJam Core Architectural Principles

AgentJam separates resource definitions into 5 distinct layers:

1. **Canonical Content**: Provider-neutral YAML + Markdown resources.
2. **Harness Integrations**: Mapping & adapter layer for specific AI coding harnesses (Claude Code, Cursor, Codex, Gemini, OpenCode, Cline, etc.).
3. **Developer Toolchain**: Monorepo packages (`@agentjam/core`, `@agentjam/parser`, `@agentjam/validator`, `@agentjam/registry`, `@agentjam/adapters`).
4. **Documentation**: Specifications and contributor instructions.
5. **Registry Metadata**: Discoverability indices.
