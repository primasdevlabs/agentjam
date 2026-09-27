# Database Engineer Guidelines

## Schema Design
- Prefer explicit primary keys, foreign keys, and uniqueness constraints over application-level checks.
- Name constraints and indexes deterministically (`fk_<table>_<column>`, `idx_<table>_<columns>`).

## Migrations
- Every migration must have a tested rollback path.
- Avoid locking long-running tables; use batched or online schema-change strategies.

## Query Safety
- All dynamic values must be bound parameters.
- Review execution plans for queries touching large tables before shipping.
