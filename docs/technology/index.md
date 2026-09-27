# Technology Stacks & Languages

AgentJam ships canonical stack profiles (`stacks/`) and a language registry (`languages/`).

## Stack Profiles

Eight stack manifests describe conventions, policies, and documentation sources per ecosystem — e.g. Go, Node/React, Python, .NET. The active stack is selected via `.agentjam/config.yaml` (`stack:`) and surfaces in every context snapshot as the **Active Stack Profile**.

## Language Registry

18 language manifests across 15 ecosystems (web, mobile, backend, jvm, dotnet, apple, embedded, functional, legacy, smart-contracts, data, databases, scripting, and more). Each manifest captures extensions, aliases, preferred tooling, framework precedence, and canonical documentation URLs.

## Detection

```bash
agentjam detect      # detected stacks from manifest files (go.mod, package.json, Cargo.toml, ...)
agentjam preflight   # verifies required binaries for detected stacks
```

`DetectStacks` probes 11 manifest types, so preflight adapts to whatever stack the workspace actually uses.
