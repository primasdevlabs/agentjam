# Schema Design Instructions

## Protocol
1. Normalize to 3NF by default; document any deliberate denormalization.
2. Enforce integrity in the database: foreign keys, uniqueness, and check constraints.
3. Every schema change ships as a reversible migration.
4. Name indexes and constraints deterministically.
