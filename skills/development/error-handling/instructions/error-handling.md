# Error Handling Instructions

## Protocol
1. No empty catch blocks or silently swallowed errors.
2. Provide explicit error types with rich context metadata (correlation ID, operation name, module).
3. Validate all input boundaries upfront.
