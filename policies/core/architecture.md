# Policy: Architectural Governance

## Policy Statement

Agents must enforce modular boundaries, single responsibility, and prevent God files or component explosion.

## Rules

1. **No God Files**:
   - Avoid creating 1000+ line files containing mixed business and presentation logic.
2. **No Component Explosion**:
   - Do not split small 3-line snippets into unnecessary micro-components.
3. **Layering Integrity**:
   - Keep business and domain logic out of UI view components when the project architecture provides a domain layer.
