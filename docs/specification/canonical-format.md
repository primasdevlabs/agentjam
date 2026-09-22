# Canonical Resource Specification

Every AgentJam resource requires machine-readable metadata.

## Agent Manifest (`agent.yaml`)

```yaml
name: code-reviewer
version: 1.0.0
type: agent
description: Reviews code changes.
skills:
  - code-review
tools:
  - filesystem
```

## Skill Manifest (`skill.yaml`)

```yaml
name: code-review
version: 1.0.0
type: skill
description: Code review skill instructions.
```
