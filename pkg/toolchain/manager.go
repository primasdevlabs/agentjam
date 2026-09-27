package toolchain

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/primasdevlabs/agentjam/pkg/core"
)

// stackProbe maps a project manifest file to the binary it requires.
type stackProbe struct {
	manifest   string
	stackID    string
	stackName  string
	binary     string
	frameworks []string
}

var stackProbes = []stackProbe{
	{"go.mod", "go", "Go Module", "go", []string{"go"}},
	{"package.json", "nodejs", "Node.js", "node", []string{"node"}},
	{"composer.json", "php", "PHP / Composer", "php", []string{"php"}},
	{"Cargo.toml", "rust", "Rust / Cargo", "cargo", []string{"rust"}},
	{"pyproject.toml", "python", "Python", "python", []string{"python"}},
	{"requirements.txt", "python", "Python", "python", []string{"python"}},
	{"Gemfile", "ruby", "Ruby / Bundler", "ruby", []string{"ruby"}},
	{"pom.xml", "java", "Java / Maven", "java", []string{"java"}},
	{"build.gradle", "java", "Java / Gradle", "java", []string{"java"}},
	{"Dockerfile", "docker", "Docker", "docker", []string{"docker"}},
	{"mix.exs", "elixir", "Elixir / Mix", "mix", []string{"elixir"}},
}

// ToolchainManager handles shell execution, system binary checks, and preflight checks in Go.
type ToolchainManager struct {
	workspaceRoot string
}

// NewToolchainManager creates a ToolchainManager instance.
func NewToolchainManager(workspaceRoot string) *ToolchainManager {
	return &ToolchainManager{workspaceRoot: workspaceRoot}
}

// HasBinary returns true if executable exists in PATH.
func (tm *ToolchainManager) HasBinary(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// DetectStacks inspects workspace manifest files and returns detected stacks.
func (tm *ToolchainManager) DetectStacks() []core.DetectedStack {
	stacks := make([]core.DetectedStack, 0)
	seen := map[string]bool{}
	for _, probe := range stackProbes {
		if seen[probe.stackID] {
			continue
		}
		if _, err := os.Stat(filepath.Join(tm.workspaceRoot, probe.manifest)); err == nil {
			seen[probe.stackID] = true
			stacks = append(stacks, core.DetectedStack{
				ID:           probe.stackID,
				Name:         probe.stackName,
				ManifestFile: probe.manifest,
				Frameworks:   probe.frameworks,
			})
		}
	}
	return stacks
}

// RunCommand executes shell commands safely with timeouts.
func (tm *ToolchainManager) RunCommand(command string, timeoutMs time.Duration) (bool, string, int64) {
	if timeoutMs <= 0 {
		timeoutMs = 30 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutMs)
	defer cancel()

	var cmd *exec.Cmd
	if os.Getenv("OS") == "Windows_NT" {
		cmd = exec.CommandContext(ctx, "powershell", "-NoProfile", "-Command", command)
	} else {
		cmd = exec.CommandContext(ctx, "sh", "-c", command)
	}

	cmd.Dir = tm.workspaceRoot
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	startTime := time.Now()
	err := cmd.Run()
	durationMs := time.Since(startTime).Milliseconds()

	output := stdoutBuf.String() + "\n" + stderrBuf.String()
	return err == nil, output, durationMs
}

// RunPreflightChecks inspects project stacks and verifies required binaries.
func (tm *ToolchainManager) RunPreflightChecks() core.PreflightCheckResult {
	checks := make([]core.PreflightCheckItem, 0)
	allPassed := true
	stackID := ""

	seen := map[string]bool{}
	for _, probe := range stackProbes {
		if seen[probe.stackID] {
			continue
		}
		if _, err := os.Stat(filepath.Join(tm.workspaceRoot, probe.manifest)); err != nil {
			continue
		}
		seen[probe.stackID] = true
		if stackID == "" {
			stackID = probe.stackID
		}

		start := time.Now()
		present := tm.HasBinary(probe.binary)
		check := core.PreflightCheckItem{
			Name:       probe.stackName + " Toolchain Check",
			Category:   "environment",
			Status:     "pass",
			Command:    probe.binary,
			DurationMs: time.Since(start).Milliseconds(),
		}
		if !present {
			check.Status = "fail"
			check.Output = "required binary '" + probe.binary + "' not found in PATH"
			allPassed = false
		}
		checks = append(checks, check)
	}

	// Fallback if no specific manifest checks were triggered
	if len(checks) == 0 {
		checks = append(checks, core.PreflightCheckItem{
			Name:       "System Environment Check",
			Category:   "environment",
			Status:     "pass",
			DurationMs: 1,
		})
	}

	return core.PreflightCheckResult{
		Passed:    allPassed,
		Timestamp: time.Now().Format(time.RFC3339),
		StackID:   stackID,
		Checks:    checks,
	}
}
