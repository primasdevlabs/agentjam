package memory_test

import (
	"sync"
	"testing"
	"time"

	"github.com/primasdevlabs/agentjam/pkg/core"
	"github.com/primasdevlabs/agentjam/pkg/memory"
)

func newManager(t *testing.T) *memory.MemoryManager {
	t.Helper()
	return memory.NewMemoryManager(t.TempDir())
}

func TestSetGetAcrossScopes(t *testing.T) {
	mm := newManager(t)

	mm.Set("k1", "working-val", core.MemoryScopeWorking, nil, 0)
	mm.Set("k2", "episodic-val", core.MemoryScopeEpisodic, nil, 0)
	mm.Set("k3", "semantic-val", core.MemoryScopeSemantic, nil, 0)

	if e, ok := mm.Get("k1", core.MemoryScopeWorking); !ok || e.Value != "working-val" {
		t.Error("working memory get failed")
	}
	if e, ok := mm.Get("k2", core.MemoryScopeEpisodic); !ok || e.Value != "episodic-val" {
		t.Error("episodic memory get failed")
	}
	if e, ok := mm.Get("k3", core.MemoryScopeSemantic); !ok || e.Value != "semantic-val" {
		t.Error("semantic memory get failed")
	}
	if _, ok := mm.Get("k1", core.MemoryScopeEpisodic); ok {
		t.Error("key should not leak across scopes")
	}
}

func TestTTLExpiry(t *testing.T) {
	mm := newManager(t)
	mm.Set("short", "gone", core.MemoryScopeWorking, nil, 1) // 1ms TTL
	time.Sleep(10 * time.Millisecond)

	if _, ok := mm.Get("short", core.MemoryScopeWorking); ok {
		t.Error("Expected expired entry to be evicted")
	}
}

func TestDeleteClearCount(t *testing.T) {
	mm := newManager(t)
	mm.Set("a", 1, core.MemoryScopeWorking, nil, 0)
	mm.Set("b", 2, core.MemoryScopeWorking, nil, 0)

	if mm.Count(core.MemoryScopeWorking) != 2 {
		t.Errorf("Expected count 2, got %d", mm.Count(core.MemoryScopeWorking))
	}
	if !mm.Delete("a", core.MemoryScopeWorking) {
		t.Error("Delete should return true for existing key")
	}
	if mm.Delete("a", core.MemoryScopeWorking) {
		t.Error("Delete should return false for missing key")
	}
	mm.Clear(core.MemoryScopeWorking)
	if mm.Count(core.MemoryScopeWorking) != 0 {
		t.Error("Expected empty working store after Clear")
	}
}

func TestPrune(t *testing.T) {
	mm := newManager(t)
	mm.Set("old", "x", core.MemoryScopeEpisodic, nil, 1)
	mm.Set("fresh", "y", core.MemoryScopeEpisodic, nil, 0)
	time.Sleep(10 * time.Millisecond)

	if removed := mm.Prune(); removed != 1 {
		t.Errorf("Expected 1 pruned entry, got %d", removed)
	}
}

func TestQueryFilters(t *testing.T) {
	mm := newManager(t)
	mm.Set("auth-token", "abc", core.MemoryScopeSemantic, []string{"auth"}, 0)
	mm.Set("db-host", "localhost", core.MemoryScopeSemantic, []string{"db"}, 0)

	byTag := mm.Query(memory.MemoryQuery{Tags: []string{"auth"}})
	if len(byTag) != 1 || byTag[0].Key != "auth-token" {
		t.Errorf("Tag query failed: %+v", byTag)
	}

	bySearch := mm.Query(memory.MemoryQuery{Search: "localhost"})
	if len(bySearch) != 1 || bySearch[0].Key != "db-host" {
		t.Errorf("Search query failed: %+v", bySearch)
	}

	byScope := mm.Query(memory.MemoryQuery{Scope: core.MemoryScopeWorking})
	if len(byScope) != 0 {
		t.Error("Scope-filtered query should return no working entries")
	}
}

func TestPersistenceRoundTrip(t *testing.T) {
	dir := t.TempDir()
	mm1 := memory.NewMemoryManager(dir)
	mm1.Set("persisted", map[string]int{"n": 42}, core.MemoryScopeSemantic, []string{"p"}, 0)

	mm2 := memory.NewMemoryManager(dir)
	e, ok := mm2.Get("persisted", core.MemoryScopeSemantic)
	if !ok {
		t.Fatal("Expected semantic memory to persist across manager instances")
	}
	if m, isMap := e.Value.(map[string]interface{}); !isMap || m["n"] != 42.0 {
		t.Errorf("Persisted value mismatch: %+v", e.Value)
	}
	// Working memory must NOT persist
	mm1.Set("volatile", "x", core.MemoryScopeWorking, nil, 0)
	mm3 := memory.NewMemoryManager(dir)
	if _, ok := mm3.Get("volatile", core.MemoryScopeWorking); ok {
		t.Error("Working memory should not persist to disk")
	}
}

func TestConcurrentAccess(t *testing.T) {
	mm := newManager(t)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := "k"
			mm.Set(key, i, core.MemoryScopeWorking, []string{"c"}, 0)
			mm.Get(key, core.MemoryScopeWorking)
			mm.Query(memory.MemoryQuery{Tags: []string{"c"}})
		}(i)
	}
	wg.Wait()
	if mm.Count(core.MemoryScopeWorking) != 1 {
		t.Errorf("Expected 1 entry after concurrent writes to same key, got %d", mm.Count(core.MemoryScopeWorking))
	}
}
