package core

import "time"

// ResourceType defines canonical resource kinds.
type ResourceType string

const (
	ResourceTypeAgent       ResourceType = "agent"
	ResourceTypeSkill       ResourceType = "skill"
	ResourceTypeTool        ResourceType = "tool"
	ResourceTypeWorkflow    ResourceType = "workflow"
	ResourceTypePolicy      ResourceType = "policy"
	ResourceTypeStack       ResourceType = "stack"
	ResourceTypeLanguage    ResourceType = "language"
	ResourceTypeIntegration ResourceType = "integration"
	ResourceTypePrompt      ResourceType = "prompt"
	ResourceTypeTemplate    ResourceType = "template"
)

// DependencyReference models cross-resource dependencies.
type DependencyReference struct {
	Name     string `json:"name" yaml:"name"`
	Version  string `json:"version" yaml:"version"`
	Optional bool   `json:"optional,omitempty" yaml:"optional,omitempty"`
}

// MCPServerConfig defines Model Context Protocol server parameters.
type MCPServerConfig struct {
	Name      string            `json:"name" yaml:"name"`
	Command   string            `json:"command" yaml:"command"`
	Args      []string          `json:"args,omitempty" yaml:"args,omitempty"`
	Env       map[string]string `json:"env,omitempty" yaml:"env,omitempty"`
	Transport string            `json:"transport" yaml:"transport"`
}

// AgentManifest defines an Agent persona manifest.
type AgentManifest struct {
	Name          string                 `json:"name" yaml:"name"`
	Version       string                 `json:"version" yaml:"version"`
	Type          string                 `json:"type" yaml:"type"`
	Description   string                 `json:"description" yaml:"description"`
	Author        string                 `json:"author,omitempty" yaml:"author,omitempty"`
	License       string                 `json:"license,omitempty" yaml:"license,omitempty"`
	Skills        []string               `json:"skills" yaml:"skills"`
	Tools         []string               `json:"tools" yaml:"tools"`
	MCPServers    []MCPServerConfig      `json:"mcpServers,omitempty" yaml:"mcpServers,omitempty"`
	Inputs        []string               `json:"inputs" yaml:"inputs"`
	Outputs       []string               `json:"outputs" yaml:"outputs"`
	Compatibility []string               `json:"compatibility,omitempty" yaml:"compatibility,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

// SkillManifest defines a domain skill manifest.
type SkillManifest struct {
	Name          string                 `json:"name" yaml:"name"`
	Version       string                 `json:"version" yaml:"version"`
	Type          string                 `json:"type" yaml:"type"`
	Description   string                 `json:"description" yaml:"description"`
	Author        string                 `json:"author,omitempty" yaml:"author,omitempty"`
	Category      string                 `json:"category,omitempty" yaml:"category,omitempty"`
	Triggers      []string               `json:"triggers,omitempty" yaml:"triggers,omitempty"`
	Compatibility []string               `json:"compatibility,omitempty" yaml:"compatibility,omitempty"`
	Parameters    map[string]interface{} `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

// ToolSafetyLevel defines risk level for tools.
type ToolSafetyLevel string

const (
	SafetyReadOnly    ToolSafetyLevel = "read-only"
	SafetySafeWrite   ToolSafetyLevel = "safe-write"
	SafetyDestructive ToolSafetyLevel = "destructive"
	SafetyAdmin       ToolSafetyLevel = "admin"
)

// ToolManifest defines an executable tool specification.
type ToolManifest struct {
	Name         string                 `json:"name" yaml:"name"`
	Version      string                 `json:"version" yaml:"version"`
	Type         string                 `json:"type" yaml:"type"`
	Description  string                 `json:"description" yaml:"description"`
	Category     string                 `json:"category,omitempty" yaml:"category,omitempty"`
	Capabilities []string               `json:"capabilities" yaml:"capabilities"`
	SafetyLevel  ToolSafetyLevel        `json:"safetyLevel" yaml:"safetyLevel"`
	MCPServer    *MCPServerConfig       `json:"mcpServer,omitempty" yaml:"mcpServer,omitempty"`
	Parameters   map[string]interface{} `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

// WorkflowStep defines a single execution step in a workflow.
type WorkflowStep struct {
	ID          string                 `json:"id" yaml:"id"`
	Name        string                 `json:"name" yaml:"name"`
	Agent       string                 `json:"agent,omitempty" yaml:"agent,omitempty"`
	Skill       string                 `json:"skill,omitempty" yaml:"skill,omitempty"`
	Description string                 `json:"description,omitempty" yaml:"description,omitempty"`
	Inputs      map[string]interface{} `json:"inputs,omitempty" yaml:"inputs,omitempty"`
	OnFailure   string                 `json:"onFailure,omitempty" yaml:"onFailure,omitempty"`
}

// WorkflowManifest defines a workflow manifest.
type WorkflowManifest struct {
	Name        string         `json:"name" yaml:"name"`
	Version     string         `json:"version" yaml:"version"`
	Type        string         `json:"type" yaml:"type"`
	Description string         `json:"description" yaml:"description"`
	Mode        string         `json:"mode" yaml:"mode"`
	Agents      []string       `json:"agents,omitempty" yaml:"agents,omitempty"`
	Skills      []string       `json:"skills,omitempty" yaml:"skills,omitempty"`
	Steps       []WorkflowStep `json:"steps" yaml:"steps"`
	Inputs      []string       `json:"inputs" yaml:"inputs"`
	Outputs     []string       `json:"outputs" yaml:"outputs"`
}

// StackManifest defines a stack profile manifest.
type StackManifest struct {
	Name          string            `json:"name" yaml:"name"`
	Version       string            `json:"version" yaml:"version"`
	Type          string            `json:"type" yaml:"type"`
	Description   string            `json:"description" yaml:"description"`
	Framework     string            `json:"framework" yaml:"framework"`
	Ecosystem     string            `json:"ecosystem" yaml:"ecosystem"`
	Conventions   map[string]string `json:"conventions,omitempty" yaml:"conventions,omitempty"`
	Policies      []string          `json:"policies,omitempty" yaml:"policies,omitempty"`
	Documentation map[string]string `json:"documentation,omitempty" yaml:"documentation,omitempty"`
}

// PolicyEnforcement defines enforcement level.
type PolicyEnforcement string

const (
	EnforceStrictBlock PolicyEnforcement = "strict-block"
	EnforceWarning     PolicyEnforcement = "warning"
	EnforceInfo        PolicyEnforcement = "info"
)

// PolicyManifest defines policy rule specifications.
type PolicyManifest struct {
	ID           string                 `json:"id" yaml:"id"`
	Name         string                 `json:"name" yaml:"name"`
	Description  string                 `json:"description" yaml:"description"`
	Category     string                 `json:"category" yaml:"category"`
	Enforcement  PolicyEnforcement      `json:"enforcement" yaml:"enforcement"`
	Scope        string                 `json:"scope" yaml:"scope"`
	AppliesTo    []string               `json:"appliesTo,omitempty" yaml:"appliesTo,omitempty"`
	Rules        map[string]interface{} `json:"rules,omitempty" yaml:"rules,omitempty"`
	Instructions string                 `json:"instructions,omitempty" yaml:"instructions,omitempty"`
}

// ContextSnapshot holds active snapshot context.
type ContextSnapshot struct {
	WorkspaceRoot      string                 `json:"workspaceRoot"`
	SystemInstruction  string                 `json:"systemInstruction"`
	TokenCountEstimate int                    `json:"tokenCountEstimate"`
	PoliciesCount      int                    `json:"policiesCount"`
	ActiveFiles        []string               `json:"activeFiles"`
	Timestamp          string                 `json:"timestamp"`
	Metadata           map[string]interface{} `json:"metadata,omitempty"`
}

// MemoryScope defines memory storage scope.
type MemoryScope string

const (
	MemoryScopeWorking  MemoryScope = "working"
	MemoryScopeEpisodic MemoryScope = "episodic"
	MemoryScopeSemantic MemoryScope = "semantic"
)

// MemoryEntry represents a memory item.
type MemoryEntry struct {
	ID        string                 `json:"id"`
	Scope     MemoryScope            `json:"scope"`
	Key       string                 `json:"key"`
	Value     interface{}            `json:"value"`
	Tags      []string               `json:"tags"`
	Timestamp string                 `json:"timestamp"`
	TTLMs     int64                  `json:"ttlMs,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// PreflightCheckItem represents a single check outcome.
type PreflightCheckItem struct {
	Name       string `json:"name"`
	Category   string `json:"category"`
	Status     string `json:"status"` // pass, fail, warn, skip
	Command    string `json:"command,omitempty"`
	Output     string `json:"output,omitempty"`
	DurationMs int64  `json:"durationMs"`
}

// PreflightCheckResult represents full preflight evaluation.
type PreflightCheckResult struct {
	Passed    bool                 `json:"passed"`
	Timestamp string               `json:"timestamp"`
	StackID   string               `json:"stackId,omitempty"`
	Checks    []PreflightCheckItem `json:"checks"`
}

// ToolCall represents an incoming tool invocation.
type ToolCall struct {
	ID        string                 `json:"id"`
	ToolName  string                 `json:"toolName"`
	Arguments map[string]interface{} `json:"arguments"`
}

// ToolResult represents the outcome of a tool execution.
type ToolResult struct {
	ID         string      `json:"id"`
	ToolName   string      `json:"toolName"`
	Status     string      `json:"status"` // success, error, blocked
	Output     interface{} `json:"output"`
	Error      string      `json:"error,omitempty"`
	DurationMs int64       `json:"durationMs"`
}

// DetectedStack holds detected tech stack details.
type DetectedStack struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	ManifestFile string   `json:"manifestFile"`
	Frameworks   []string `json:"frameworks"`
}

// TimestampNow returns current RFC3339 timestamp.
func TimestampNow() string {
	return time.Now().Format(time.RFC3339)
}
