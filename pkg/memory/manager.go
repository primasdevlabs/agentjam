package memory

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/primasdevlabs/agentjam/pkg/core"
)

// MemoryQuery configures search queries.
type MemoryQuery struct {
	Scope  core.MemoryScope `json:"scope,omitempty"`
	Search string           `json:"search,omitempty"`
	Tags   []string         `json:"tags,omitempty"`
	Limit  int              `json:"limit,omitempty"`
}

// MemoryManager manages working, episodic, and semantic memory in Go.
type MemoryManager struct {
	mu             sync.RWMutex
	workspaceRoot  string
	workingMemory  map[string]core.MemoryEntry
	episodicMemory map[string]core.MemoryEntry
	semanticMemory map[string]core.MemoryEntry
	memoryFilePath string
}

// NewMemoryManager initializes a MemoryManager.
func NewMemoryManager(workspaceRoot string) *MemoryManager {
	mm := &MemoryManager{
		workspaceRoot:  workspaceRoot,
		workingMemory:  make(map[string]core.MemoryEntry),
		episodicMemory: make(map[string]core.MemoryEntry),
		semanticMemory: make(map[string]core.MemoryEntry),
		memoryFilePath: filepath.Join(workspaceRoot, ".agentjam", "memory", "memory_store.json"),
	}
	mm.LoadFromDisk()
	return mm
}

func (mm *MemoryManager) getStore(scope core.MemoryScope) map[string]core.MemoryEntry {
	switch scope {
	case core.MemoryScopeEpisodic:
		return mm.episodicMemory
	case core.MemoryScopeSemantic:
		return mm.semanticMemory
	default:
		return mm.workingMemory
	}
}

// Set stores a key-value memory entry.
func (mm *MemoryManager) Set(key string, value interface{}, scope core.MemoryScope, tags []string, ttlMs int64) core.MemoryEntry {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	if scope == "" {
		scope = core.MemoryScopeWorking
	}
	if tags == nil {
		tags = []string{}
	}

	entry := core.MemoryEntry{
		ID:        fmt.Sprintf("%s:%s:%d", scope, key, time.Now().UnixNano()),
		Scope:     scope,
		Key:       key,
		Value:     value,
		Tags:      tags,
		Timestamp: time.Now().Format(time.RFC3339),
		TTLMs:     ttlMs,
	}

	store := mm.getStore(scope)
	store[key] = entry

	if scope == core.MemoryScopeEpisodic || scope == core.MemoryScopeSemantic {
		mm.saveToDiskLocked()
	}

	return entry
}

// Get retrieves a memory entry.
func (mm *MemoryManager) Get(key string, scope core.MemoryScope) (core.MemoryEntry, bool) {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	if scope == "" {
		scope = core.MemoryScopeWorking
	}

	store := mm.getStore(scope)
	entry, exists := store[key]
	if !exists {
		return core.MemoryEntry{}, false
	}

	if entry.TTLMs > 0 {
		createdTime, err := time.Parse(time.RFC3339, entry.Timestamp)
		if err == nil && time.Since(createdTime).Milliseconds() > entry.TTLMs {
			delete(store, key)
			return core.MemoryEntry{}, false
		}
	}

	return entry, true
}

// Query performs search filtering over memory stores.
func (mm *MemoryManager) Query(q MemoryQuery) []core.MemoryEntry {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	if q.Limit <= 0 {
		q.Limit = 50
	}

	var stores []map[string]core.MemoryEntry
	if q.Scope != "" {
		stores = []map[string]core.MemoryEntry{mm.getStore(q.Scope)}
	} else {
		stores = []map[string]core.MemoryEntry{mm.workingMemory, mm.episodicMemory, mm.semanticMemory}
	}

	results := make([]core.MemoryEntry, 0)
	searchLower := strings.ToLower(q.Search)

	for _, store := range stores {
		for _, entry := range store {
			// TTL Check
			if entry.TTLMs > 0 {
				createdTime, err := time.Parse(time.RFC3339, entry.Timestamp)
				if err == nil && time.Since(createdTime).Milliseconds() > entry.TTLMs {
					continue
				}
			}

			// Tag filter
			if len(q.Tags) > 0 {
				hasTag := false
				for _, reqTag := range q.Tags {
					for _, entryTag := range entry.Tags {
						if reqTag == entryTag {
							hasTag = true
							break
						}
					}
				}
				if !hasTag {
					continue
				}
			}

			// Search text filter
			if searchLower != "" {
				keyMatch := strings.Contains(strings.ToLower(entry.Key), searchLower)
				valBytes, _ := json.Marshal(entry.Value)
				valMatch := strings.Contains(strings.ToLower(string(valBytes)), searchLower)
				if !keyMatch && !valMatch {
					continue
				}
			}

			results = append(results, entry)
		}
	}

	// Sort descending by timestamp
	sort.Slice(results, func(i, j int) bool {
		return results[i].Timestamp > results[j].Timestamp
	})

	if len(results) > q.Limit {
		return results[:q.Limit]
	}
	return results
}

func (mm *MemoryManager) saveToDiskLocked() {
	dir := filepath.Dir(mm.memoryFilePath)
	_ = os.MkdirAll(dir, 0755)

	payload := map[string]interface{}{
		"episodic":  mm.episodicMemory,
		"semantic":  mm.semanticMemory,
		"updatedAt": time.Now().Format(time.RFC3339),
	}

	bytes, err := json.MarshalIndent(payload, "", "  ")
	if err == nil {
		_ = os.WriteFile(mm.memoryFilePath, bytes, 0644)
	}
}

// SaveToDisk writes memory state to disk.
func (mm *MemoryManager) SaveToDisk() {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	mm.saveToDiskLocked()
}

// LoadFromDisk loads memory state from disk.
func (mm *MemoryManager) LoadFromDisk() {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	bytes, err := os.ReadFile(mm.memoryFilePath)
	if err != nil {
		return
	}

	var payload struct {
		Episodic map[string]core.MemoryEntry `json:"episodic"`
		Semantic map[string]core.MemoryEntry `json:"semantic"`
	}

	if err := json.Unmarshal(bytes, &payload); err == nil {
		if payload.Episodic != nil {
			mm.episodicMemory = payload.Episodic
		}
		if payload.Semantic != nil {
			mm.semanticMemory = payload.Semantic
		}
	}
}
