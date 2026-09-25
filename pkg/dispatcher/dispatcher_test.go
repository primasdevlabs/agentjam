package dispatcher_test

import (
	"os"
	"testing"

	"github.com/primasdevlabs/agentjam/pkg/core"
	"github.com/primasdevlabs/agentjam/pkg/dispatcher"
)

func TestToolDispatcherModule(t *testing.T) {
	cwd, _ := os.Getwd()
	td := dispatcher.NewToolDispatcher(cwd, dispatcher.ToolDispatcherOptions{AllowDestructive: false}, nil)

	td.RegisterTool(
		core.ToolManifest{
			Name:        "safe_tool",
			Version:     "1.0.0",
			Type:        "tool",
			Description: "Safe execution",
			SafetyLevel: core.SafetyReadOnly,
		},
		func(args map[string]interface{}) (interface{}, error) {
			return "success", nil
		},
	)

	res := td.Dispatch(core.ToolCall{
		ID:        "c1",
		ToolName:  "safe_tool",
		Arguments: map[string]interface{}{},
	})

	if res.Status != "success" || res.Output != "success" {
		t.Errorf("Dispatch failed: %v", res)
	}
}
