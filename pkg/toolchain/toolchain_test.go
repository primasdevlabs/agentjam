package toolchain_test

import (
	"os"
	"testing"
	"time"

	"github.com/agentjam/agentjam/pkg/toolchain"
)

func TestToolchainManagerModule(t *testing.T) {
	cwd, _ := os.Getwd()
	tm := toolchain.NewToolchainManager(cwd)

	// RunCommand
	ok, out, dur := tm.RunCommand("echo Hello", 5*time.Second)
	if !ok || dur < 0 {
		t.Errorf("RunCommand failed: ok=%v, out=%s", ok, out)
	}

	// RunPreflightChecks
	res := tm.RunPreflightChecks()
	if len(res.Checks) == 0 {
		t.Errorf("Preflight checks returned no items")
	}
}
