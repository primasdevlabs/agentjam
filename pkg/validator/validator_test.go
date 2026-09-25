package validator_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/agentjam/agentjam/pkg/validator"
)

func TestValidatorModule(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	rootDir := filepath.Clean(filepath.Join(cwd, "..", ".."))
	errors := validator.ValidateRepository(rootDir)
	// Should execute without panicking
	_ = errors
}
