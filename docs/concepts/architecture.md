# AgentJam Architecture

AgentJam is built as a harness-agnostic monorepo separating canonical assets from adapter implementations.

## Monorepo Packages

1. `@agentjam/core`: TypeScript type definitions and Zod validation schemas.
2. `@agentjam/parser`: YAML & markdown manifest loader.
3. `@agentjam/validator`: Repository linting and structural validation engine.
4. `@agentjam/registry`: Local and remote resource indexer.
5. `@agentjam/policy-engine`: Governance and precedence evaluator.
6. `@agentjam/adapters`: Export engine converting canonical rules into target harness formats (Claude Code, Cursor, Windsurf, Roo Code, etc.).
7. `@agentjam/runtime`: AgentJam execution runtime and environment compatibility detector.
