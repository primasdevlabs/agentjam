package mcp

import (
	"bufio"
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/primasdevlabs/agentjam/pkg/core"
	"github.com/primasdevlabs/agentjam/pkg/dispatcher"
)

func newTestServer(t *testing.T, requests ...string) ([]map[string]interface{}, error) {
	t.Helper()
	td := dispatcher.NewToolDispatcher(t.TempDir(), dispatcher.ToolDispatcherOptions{}, nil)
	td.RegisterTool(core.ToolManifest{
		Name:        "echo",
		Version:     "1.0.0",
		Type:        "tool",
		Description: "Echoes arguments.",
		SafetyLevel: core.SafetyReadOnly,
	}, func(args map[string]interface{}) (interface{}, error) {
		return args, nil
	})

	in := strings.NewReader(strings.Join(requests, "\n") + "\n")
	var out bytes.Buffer
	srv := NewServer(td, in, &out)
	if err := srv.Serve(); err != nil {
		return nil, err
	}

	var responses []map[string]interface{}
	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		var m map[string]interface{}
		if err := json.Unmarshal(scanner.Bytes(), &m); err != nil {
			t.Fatalf("invalid JSON response line: %v", err)
		}
		responses = append(responses, m)
	}
	return responses, nil
}

func TestMCPInitialize(t *testing.T) {
	resps, err := newTestServer(t,
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"0"}}}`,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(resps) != 1 {
		t.Fatalf("expected 1 response, got %d", len(resps))
	}
	result := resps[0]["result"].(map[string]interface{})
	if result["protocolVersion"] != ProtocolVersion {
		t.Fatalf("bad protocol version: %v", result["protocolVersion"])
	}
	info := result["serverInfo"].(map[string]interface{})
	if info["name"] != "agentjam" {
		t.Fatalf("bad serverInfo: %v", info)
	}
}

func TestMCPToolsListAndCall(t *testing.T) {
	resps, err := newTestServer(t,
		`{"jsonrpc":"2.0","id":1,"method":"tools/list"}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"echo","arguments":{"msg":"hi"}}}`,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(resps) != 2 {
		t.Fatalf("expected 2 responses, got %d", len(resps))
	}

	tools := resps[0]["result"].(map[string]interface{})["tools"].([]interface{})
	if len(tools) != 1 {
		t.Fatalf("expected 1 tool, got %d", len(tools))
	}
	if tools[0].(map[string]interface{})["name"] != "echo" {
		t.Fatalf("bad tool entry: %v", tools[0])
	}

	callRes := resps[1]["result"].(map[string]interface{})
	if callRes["isError"] != false {
		t.Fatalf("call should succeed: %v", callRes)
	}
	content := callRes["content"].([]interface{})
	text := content[0].(map[string]interface{})["text"].(string)
	if !strings.Contains(text, "hi") {
		t.Fatalf("echo output missing arg: %s", text)
	}
}

func TestMCPUnknownTool(t *testing.T) {
	resps, err := newTestServer(t,
		`{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"ghost","arguments":{}}}`,
	)
	if err != nil {
		t.Fatal(err)
	}
	result := resps[0]["result"].(map[string]interface{})
	if result["isError"] != true {
		t.Fatalf("expected isError for unknown tool: %v", result)
	}
}

func TestMCPUnknownMethod(t *testing.T) {
	resps, err := newTestServer(t,
		`{"jsonrpc":"2.0","id":1,"method":"resources/list"}`,
	)
	if err != nil {
		t.Fatal(err)
	}
	rpcErr := resps[0]["error"].(map[string]interface{})
	if rpcErr["code"].(float64) != -32601 {
		t.Fatalf("expected -32601, got %v", rpcErr["code"])
	}
}

func TestMCPNotificationNoResponse(t *testing.T) {
	resps, err := newTestServer(t,
		`{"jsonrpc":"2.0","method":"notifications/initialized"}`,
		`{"jsonrpc":"2.0","id":1,"method":"ping"}`,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(resps) != 1 {
		t.Fatalf("notification should not produce a response, got %d", len(resps))
	}
}
