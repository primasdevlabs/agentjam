package parser_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/agentjam/agentjam/pkg/parser"
)

func TestParserModule(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Navigate to root workspace if inside pkg/parser
	rootDir := filepath.Clean(filepath.Join(cwd, "..", ".."))

	// Test DiscoverResources
	resources := parser.DiscoverResources(rootDir)
	if len(resources) == 0 {
		t.Logf("No resources discovered in %s, testing frontmatter fallback", rootDir)
	}

	// Test ParseFrontmatter
	content := "---\nname: test-agent\nversion: 1.0.0\n---\n# Agent Role\nRole instructions"
	fm, body := parser.ParseFrontmatter(content)
	if fm == "" || body == "" {
		t.Errorf("Frontmatter parsing failed: fm=%q, body=%q", fm, body)
	}
}
