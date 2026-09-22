# Secure Coding Instructions

## Directives
1. Parameterized SQL queries or ORM models only. Zero string concatenation in database queries.
2. Validate external inputs at boundary endpoints using schema validators (Zod, Valibot, FormRequest, Pydantic).
3. Sanitize HTML/rendered outputs to prevent XSS.
