package memory_test

import (
	"os"
	"testing"

	"github.com/agentjam/agentjam/pkg/core"
	"github.com/agentjam/agentjam/pkg/memory"
)

func TestMemoryManagerModule(t *testing.T) {
	cwd, _ := os.Getwd()
	mm := memory.NewMemoryManager(cwd)

	// Set & Get Working Memory
	mm.Set("key1", "val1", core.MemoryScopeWorking, []string{"tag1"}, 0)
	entry, ok := mm.Get("key1", core.MemoryScopeWorking)
	if !ok || entry.Value != "val1" {
		t.Errorf("Memory Get failed")
	}

	// Query Memory
	results := mm.Query(memory.MemoryQuery{Search: "val1"})
	if len(results) == 0 {
		t.Errorf("Memory Query failed")
	}
}
