# API Design Instructions

## Protocol
1. Model endpoints around resources; use nouns, not verbs, in paths.
2. Return consistent error envelopes with machine-readable codes.
3. Paginate collection endpoints; never return unbounded lists.
4. Make mutating endpoints idempotent where retries are possible.
5. Version breaking changes; evolve additive fields without a version bump.
