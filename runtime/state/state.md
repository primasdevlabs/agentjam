# Runtime State & Memory

Three-tier memory store backing agent sessions.

## Scopes

- **Working** — ephemeral, in-memory only. Cleared when the process exits.
- **Episodic** — persisted session events (e.g. tool-call results recorded by the dispatcher).
- **Semantic** — persisted durable facts and decisions.

Episodic and semantic scopes persist to `.agentjam/memory/memory_store.json`.

## Operations

| Method | Behavior |
| --- | --- |
| `Set(key, value, scope, tags, ttlMs)` | Upsert entry; persists for non-working scopes |
| `Get(key, scope)` | Retrieve; lazily evicts expired entries |
| `Delete(key, scope)` | Remove a single entry |
| `Query(MemoryQuery)` | Filter by scope, tags, and text search; newest first |
| `Count(scope)` | Live entry count (empty scope counts all stores) |
| `Prune()` | Remove all expired entries |
| `Clear(scope)` | Empty one scope |
| `SaveToDisk` / `LoadFromDisk` | Explicit persistence control |

All operations are guarded by `sync.RWMutex` and safe for concurrent use.
