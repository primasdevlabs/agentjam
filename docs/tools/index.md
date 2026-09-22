# AgentJam Tool Specifications

AgentJam tools define capabilities provided by the target harness or external MCP servers to agents:

* **filesystem**: Read, write, list, and replace content in the project directory.
* **terminal**: Execute shell commands, run test runners, linters, and build scripts.
* **browser**: Interact with web applications, verify DOM state, and capture visual rendering.
* **mcp**: Connect to Model Context Protocol (MCP) servers to invoke external tools, stream prompts, and retrieve context resources via stdio or SSE transports.
