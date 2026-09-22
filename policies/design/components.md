# Component & Design Token Policy

## No Inline Design Tokens

Design tokens MUST NOT be hardcoded inside components.

**Forbidden**:
```tsx
<div style={{ color: "#344C36" }} />
```
or
```tsx
className="text-[#344C36]"
```

All reusable colors, spacing, border radii, typography values, shadows, and animation durations must reference central design tokens (`theme/`, `styles/`, `tokens.css`).

## Balanced Component Architecture

- **No God Components**: Do not place entire pages or 2000+ lines of logic inside a single component file (`Everything.tsx`).
- **No Component Explosion**: Do not split simple markup into hundreds of single-line micro-wrapper files.
- **Component Decomposition Criteria**: Create reusable components based on single responsibility, domain domain meaning, clear prop interfaces, and maintainability.
