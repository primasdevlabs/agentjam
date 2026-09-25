# AgentJam Architecture (Go Native Engine)

AgentJam is built as a high-performance, harness-agnostic Go architecture separating canonical assets from export adapter implementations.

---

## Go Packages Architecture (`pkg/`)

1. **`github.com/primasdevlabs/agentjam/pkg/core`**: Go domain types, structs, constants, errors, and token estimation utilities.
2. **`github.com/primasdevlabs/agentjam/pkg/parser`**: YAML manifest loader, Markdown frontmatter parser, and directory resource discovery walker.
3. **`github.com/primasdevlabs/agentjam/pkg/policy`**: Governance evaluator for core policies, 12 anti-slop design rules, and security scans.
4. **`github.com/primasdevlabs/agentjam/pkg/validator`**: Repository linting, schema validation, and cross-reference integrity checker.
5. **`github.com/primasdevlabs/agentjam/pkg/registry`**: Search index builder and repository statistics calculator.
6. **`github.com/primasdevlabs/agentjam/pkg/adapters`**: Harness export engine converting canonical rules into target AI harness formats (Claude Code, Cursor, Gemini/Antigravity, Cline, Windsurf, Devin, Roo Code, Generic).
7. **`github.com/primasdevlabs/agentjam/pkg/context`**: Context Manager, token estimator, system prompt formatter, and window truncator.
8. **`github.com/primasdevlabs/agentjam/pkg/memory`**: Thread-safe Working, Episodic, and Semantic memory store.
9. **`github.com/primasdevlabs/agentjam/pkg/toolchain`**: Toolchain preflight check executor and PATH binary inspector.
10. **`github.com/primasdevlabs/agentjam/pkg/dispatcher`**: Tool dispatcher enforcing safety levels and workspace path boundary security.
11. **`github.com/primasdevlabs/agentjam/pkg/runtime`**: Main execution runtime orchestrator.

---

## Executable CLI (`cmd/agentjam`)

The native Go binary CLI (`agentjam`) provides subcommands:
- `agentjam validate`: Validates repository resources, policies, and stacks.
- `agentjam export`: Auto-detects or explicitly exports rules for AI harnesses (`--harness auto|all|cursor|claude-code|gemini|...`).
- `agentjam context`: Generates context snapshot and system instructions.
- `agentjam preflight`: Runs toolchain preflight checks.
- `agentjam build-registry`: Generates `registry.json` index.
- `agentjam run-e2e`: Runs full end-to-end verification suite.
