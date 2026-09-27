// Package mcp implements a minimal MCP (Model Context Protocol) server that
// exposes the AgentJam tool dispatcher over newline-delimited JSON-RPC on
// stdio, per the MCP stdio transport.
package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"

	"github.com/primasdevlabs/agentjam/pkg/core"
	"github.com/primasdevlabs/agentjam/pkg/dispatcher"
)

// ProtocolVersion is the MCP revision this server negotiates.
const ProtocolVersion = "2024-11-05"

const serverVersion = "1.0.0"

type rpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type rpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  interface{}     `json:"result,omitempty"`
	Error   *rpcError       `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Server bridges MCP stdio requests to a ToolDispatcher.
type Server struct {
	dispatcher *dispatcher.ToolDispatcher
	in         io.Reader
	out        io.Writer
}

// NewServer creates an MCP server reading requests from in and writing
// responses to out (typically os.Stdin / os.Stdout).
func NewServer(td *dispatcher.ToolDispatcher, in io.Reader, out io.Writer) *Server {
	return &Server{dispatcher: td, in: in, out: out}
}

// Serve processes newline-delimited JSON-RPC messages until EOF or a read error.
func (s *Server) Serve() error {
	scanner := bufio.NewScanner(s.in)
	scanner.Buffer(make([]byte, 0, 64*1024), 8*1024*1024)
	writer := bufio.NewWriter(s.out)
	defer writer.Flush()

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var req rpcRequest
		if err := json.Unmarshal(line, &req); err != nil {
			s.write(writer, rpcResponse{
				JSONRPC: "2.0", ID: json.RawMessage("null"),
				Error: &rpcError{Code: -32700, Message: "parse error"},
			})
			continue
		}
		resp, ok := s.handle(req)
		if !ok {
			continue // notification — no response
		}
		s.write(writer, resp)
	}
	return scanner.Err()
}

func (s *Server) write(w *bufio.Writer, resp rpcResponse) {
	data, err := json.Marshal(resp)
	if err != nil {
		return
	}
	w.Write(data)
	w.WriteByte('\n')
	w.Flush()
}

func (s *Server) handle(req rpcRequest) (rpcResponse, bool) {
	resp := rpcResponse{JSONRPC: "2.0", ID: req.ID}
	fail := func(code int, msg string) (rpcResponse, bool) {
		resp.Error = &rpcError{Code: code, Message: msg}
		return resp, len(req.ID) > 0
	}

	switch req.Method {
	case "initialize":
		resp.Result = map[string]interface{}{
			"protocolVersion": ProtocolVersion,
			"capabilities":    map[string]interface{}{"tools": map[string]interface{}{"listChanged": false}},
			"serverInfo":      map[string]interface{}{"name": "agentjam", "version": serverVersion},
		}
	case "ping":
		resp.Result = map[string]interface{}{}
	case "tools/list":
		resp.Result = map[string]interface{}{"tools": s.toolList()}
	case "tools/call":
		result, err := s.toolCall(req.Params)
		if err != nil {
			return fail(-32602, err.Error())
		}
		resp.Result = result
	default:
		if len(req.ID) == 0 {
			return resp, false
		}
		return fail(-32601, fmt.Sprintf("method not found: %s", req.Method))
	}
	return resp, len(req.ID) > 0
}

func (s *Server) toolList() []map[string]interface{} {
	tools := make([]map[string]interface{}, 0)
	for _, m := range s.dispatcher.ListTools() {
		schema := map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
		}
		if len(m.Parameters) > 0 {
			if props, ok := m.Parameters["properties"].(map[string]interface{}); ok {
				schema["properties"] = props
			} else {
				schema["properties"] = m.Parameters
			}
			if req, ok := m.Parameters["required"]; ok {
				schema["required"] = req
			}
		}
		tools = append(tools, map[string]interface{}{
			"name":        m.Name,
			"description": m.Description,
			"inputSchema": schema,
		})
	}
	return tools
}

func (s *Server) toolCall(paramsJSON json.RawMessage) (map[string]interface{}, error) {
	var params struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	}
	if err := json.Unmarshal(paramsJSON, &params); err != nil {
		return nil, fmt.Errorf("invalid tools/call params: %w", err)
	}
	if params.Name == "" {
		return nil, fmt.Errorf("tools/call requires a tool name")
	}

	result := s.dispatcher.Dispatch(core.ToolCall{
		ID:        "mcp-" + params.Name,
		ToolName:  params.Name,
		Arguments: params.Arguments,
	})

	text, _ := json.MarshalIndent(result.Output, "", "  ")
	payload := map[string]interface{}{
		"content": []map[string]interface{}{
			{"type": "text", "text": string(text)},
		},
		"isError": result.Status != "success",
	}
	if result.Status != "success" {
		payload["content"] = []map[string]interface{}{
			{"type": "text", "text": result.Error},
		}
	}
	return payload, nil
}
