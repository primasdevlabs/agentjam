package registry_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/agentjam/agentjam/pkg/registry"
)

func TestRegistryModule(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	rootDir := filepath.Clean(filepath.Join(cwd, "..", ".."))
	index := registry.BuildRegistryIndex(rootDir)
	if index.AgentsCount < 0 {
		t.Errorf("BuildRegistryIndex returned invalid stats")
	}
}
