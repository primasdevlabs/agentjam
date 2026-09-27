package registry_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/primasdevlabs/agentjam/pkg/core"
	"github.com/primasdevlabs/agentjam/pkg/registry"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	return filepath.Clean(filepath.Join(cwd, "..", ".."))
}

func TestBuildRegistryIndex(t *testing.T) {
	idx := registry.BuildRegistryIndex(repoRoot(t))

	if idx.AgentsCount == 0 || idx.SkillsCount == 0 || idx.PoliciesCount == 0 {
		t.Errorf("Expected non-zero counts, got %+v", idx)
	}
	if idx.Total != len(idx.Entries) {
		t.Errorf("Total %d != entries %d", idx.Total, len(idx.Entries))
	}
	if idx.GeneratedAt == "" {
		t.Error("Expected generatedAt timestamp")
	}
}

func TestRegistryEntriesPopulated(t *testing.T) {
	idx := registry.BuildRegistryIndex(repoRoot(t))

	var found *registry.RegistryEntry
	for i, e := range idx.Entries {
		if e.Type == core.ResourceTypeAgent && e.ID == "software-engineer" {
			found = &idx.Entries[i]
		}
		if e.Path == "" || filepath.IsAbs(e.Path) {
			t.Errorf("Entry path must be relative and non-empty: %+v", e)
		}
	}
	if found == nil {
		t.Fatal("Expected software-engineer agent entry")
	}
	if found.Name != "software-engineer" || found.Description == "" {
		t.Errorf("Entry should carry parsed manifest fields: %+v", found)
	}
}

func TestRegistryIndexSerializable(t *testing.T) {
	idx := registry.BuildRegistryIndex(repoRoot(t))
	bytes, err := json.Marshal(idx)
	if err != nil {
		t.Fatalf("Index not serializable: %v", err)
	}
	var roundTrip registry.RegistryIndex
	if err := json.Unmarshal(bytes, &roundTrip); err != nil {
		t.Fatalf("Index not deserializable: %v", err)
	}
	if roundTrip.AgentsCount != idx.AgentsCount {
		t.Error("Round-trip mismatch")
	}
}
