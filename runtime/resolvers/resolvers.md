# Runtime Resolvers

Resolves precedence between conflicting sources of guidance and checks documentation freshness.

## Precedence Resolution

`ResolveHighest` picks the winning `PrecedenceCandidate` using `PrecedenceOrder` weights (lower wins):

1. Current Codebase State
2. Project Configuration (`.agentjam/config.yaml`)
3. User explicit instruction
4. Pinned Rules
5. Authoritative Documentation
6. Environment defaults
7. Model Knowledge (unmapped levels resolve last)

## Freshness Checks

`CheckFreshness(lastUpdatedIso, maxDocAge)` compares a document timestamp against the max-age policy. Accepted units: `h`, `d`, `w`, `m` (30d), `y` (365d). Unparseable timestamps default to fresh; unrecognized specs default to 7 days.
