# Query Optimization Instructions

## Protocol
1. Measure first: inspect the execution plan before rewriting a query.
2. Eliminate N+1 patterns with joins or batched loading.
3. Add composite indexes in column order matching the query's filter and sort.
4. Prefer keyset pagination over OFFSET for large tables.
