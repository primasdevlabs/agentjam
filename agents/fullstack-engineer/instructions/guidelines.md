# Full-Stack Engineer Guidelines

## API Boundaries
- Validate input at the HTTP boundary; never trust client payloads.
- Keep serialization formats consistent with existing endpoints.

## Frontend
- Reuse existing components and design tokens before introducing new ones.

## Data
- New fields and tables go through migrations, never ad-hoc queries.
