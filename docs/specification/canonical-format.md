# Canonical Manifest Specification

AgentJam defines strict YAML schemas for all resource types.

---

## 1. Agent Manifest (`agent.yaml`)
```yaml
name: software-engineer
version: 1.0.0
type: agent
description: Generalist software engineer for full-stack feature development.
skills: [architecture, testing, secure-coding]
tools: [filesystem, terminal]
inputs: [requirement-spec, codebase]
outputs: [source-code, pull-request]
```

## 2. Skill Manifest (`skill.yaml`)
```yaml
name: anti-slop
version: 1.0.0
type: skill
description: Audits UI implementations to eliminate generic AI slop.
category: design
triggers: [anti-slop, design-governance]
```

## 3. Workflow Manifest (`workflow.yaml`)
```yaml
name: feature-development
version: 1.0.0
type: workflow
description: End-to-end feature implementation flow.
mode: multi-agent
steps:
  - id: planning
    agent: architect
  - id: implementation
    agent: software-engineer
```

## 4. Policy Manifest (`policy.yaml`)
```yaml
id: design-anti-slop
name: Design Anti-Slop Policy
description: Prohibits emoji UI icons and generic purple gradients.
category: design
enforcement: strict-block
```
