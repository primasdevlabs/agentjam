# Model Context Protocol (MCP) Tools & Safety Levels

AgentJam tools connect AI agents to system environments via the **Model Context Protocol (MCP)** or native abstraction wrappers.

---

## Canonical Tool Catalog & Safety Levels

| Tool ID | Safety Level | Description & Operational Behavior |
| :--- | :---: | :--- |
| `filesystem` | `safe-write` | File viewing, directory walking, file replacement, multi-replace, and file creation with path validation. |
| `terminal` | `destructive` | Command line execution (`powershell`, `bash`, `zsh`) with async background task management. |
| `browser` | `read-only` | Web page inspection, HTML element extraction, and browser screenshot capture for UI validation. |
| `git` | `safe-write` | Git status inspection, branch verification, diff analysis, and commit generation. |
| `database` | `safe-write` | Schema inspection, migration execution, and parameterized SQL query execution. |
| `http-client` | `read-only` | External API HTTP requests with rate limiting, retries, and backoff handling. |
| `ast-parser` | `read-only` | Static code parsing, AST symbol extraction, and private import verification. |
| `docker` | `safe-write` | Container lifecycle management, image building, and containerized test execution. |

---

## Tool Execution Rules

1. **Safety Enclosures**: Destructive commands (`terminal`, `docker`) require user approval or confirmation flags in strict mode.
2. **Path Boundaries**: All filesystem tools validate paths to prevent directory traversal outside workspace boundaries.
3. **Structured Context**: Tool outputs are sanitized and formatted as markdown for agent context windows.
