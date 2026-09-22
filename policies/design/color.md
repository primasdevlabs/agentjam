# Color Policy

## Forbidden Color Defaults

Agents MUST NOT default to:
- Neon colors
- Indigo (`#4F46E5` / Tailwind indigo)
- Generic electric blue (`#3B82F6`)
- Generic purple
- Blue/purple gradients
- Neon gradients or arbitrary rainbow background fills

## Palette Derivation Rules

1. **Derive from Established Direction**: Colors must be derived from the project's design system tokens (`tokens.css`, `tailwind.config`, CSS variables).
2. **No Invented Palettes**: Do not invent a new color palette simply because one is missing. Ask for direction or inspect existing brand assets first.
3. **Functional Color Application**: Use color purposefully for state (success, warning, error, info) and primary action contrast, not for background noise.
