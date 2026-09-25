package core

import "errors"

// Common domain errors for AgentJam Go implementation.
var (
	ErrResourceNotFound = errors.New("agentjam: resource not found")
	ErrInvalidManifest  = errors.New("agentjam: invalid resource manifest")
	ErrPolicyViolation  = errors.New("agentjam: policy violation encountered")
	ErrToolExecution    = errors.New("agentjam: tool execution error")
	ErrSecurityBlock    = errors.New("agentjam: execution blocked by security policy")
)
