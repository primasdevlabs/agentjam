# Mandatory Design Direction Request Policy

## Design Discovery Requirement

When an agent is asked to create or substantially redesign a UI and no established design system exists:

1. The agent MUST NOT proceed by inventing generic SaaS fallback styles.
2. The agent MUST request or determine:
   - **Design Direction**: Visual style, brand personality, color direction, typography, density.
   - **Design Tools / Assets**: Figma files, design token files, icon libraries, component libraries, brand guidelines.

## Established Project Priority

If the project ALREADY possesses an established design system (brand colors, tokens, component library, font choices):
- The agent MUST inspect and follow the existing design system as the primary source of truth.
- Do NOT replace an existing coherent design system with generic defaults.
