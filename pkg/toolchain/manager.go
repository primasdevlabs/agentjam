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

// RunPreflightChecks inspects project stack and runs preflight verification.
func (tm *ToolchainManager) RunPreflightChecks() core.PreflightCheckResult {
	checks := make([]core.PreflightCheckItem, 0)
	allPassed := true

	// Check go.mod presence
	goMod := filepath.Join(tm.workspaceRoot, "go.mod")
	if _, err := os.Stat(goMod); err == nil {
		goCheck := core.PreflightCheckItem{
			Name:       "Go Engine Check",
			Category:   "environment",
			Status:     "fail",
			DurationMs: 5,
		}
		if tm.HasBinary("go") {
			goCheck.Status = "pass"
		} else {
			allPassed = false
		}
		checks = append(checks, goCheck)
	}

	// Check package.json presence
	pkgJson := filepath.Join(tm.workspaceRoot, "package.json")
	if _, err := os.Stat(pkgJson); err == nil {
		nodeCheck := core.PreflightCheckItem{
			Name:       "Node.js Engine Check",
			Category:   "environment",
			Status:     "fail",
			DurationMs: 5,
		}
		if tm.HasBinary("node") {
			nodeCheck.Status = "pass"
		}
		checks = append(checks, nodeCheck)
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
		Checks:    checks,
	}
}
