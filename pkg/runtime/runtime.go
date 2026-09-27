package runtime

import (
	"github.com/primasdevlabs/agentjam/pkg/context"
	"github.com/primasdevlabs/agentjam/pkg/dispatcher"
	"github.com/primasdevlabs/agentjam/pkg/memory"
	"github.com/primasdevlabs/agentjam/pkg/policy"
	"github.com/primasdevlabs/agentjam/pkg/toolchain"
)

// AgentJamRuntime orchestrates all AgentJam subsystems.
type AgentJamRuntime struct {
	projectRoot      string
	policyEngine     *policy.PolicyEngine
	contextManager   *context.ContextManager
	memoryManager    *memory.MemoryManager
	toolchainManager *toolchain.ToolchainManager
	toolDispatcher   *dispatcher.ToolDispatcher
}

// NewRuntime initializes AgentJamRuntime in Go and loads workspace policies.
func NewRuntime(projectRoot string, opts dispatcher.ToolDispatcherOptions) *AgentJamRuntime {
	pe := policy.NewPolicyEngineFromDir(projectRoot)
	cm := context.NewContextManager(projectRoot, pe)
	mm := memory.NewMemoryManager(projectRoot)
	tm := toolchain.NewToolchainManager(projectRoot)
	td := dispatcher.NewToolDispatcher(projectRoot, opts, mm)

	return &AgentJamRuntime{
		projectRoot:      projectRoot,
		policyEngine:     pe,
		contextManager:   cm,
		memoryManager:    mm,
		toolchainManager: tm,
		toolDispatcher:   td,
	}
}

// RootPath returns the workspace root the runtime was initialized with.
func (r *AgentJamRuntime) RootPath() string {
	return r.projectRoot
}

func (r *AgentJamRuntime) GetPolicyEngine() *policy.PolicyEngine {
	return r.policyEngine
}

func (r *AgentJamRuntime) GetContextManager() *context.ContextManager {
	return r.contextManager
}

func (r *AgentJamRuntime) GetMemoryManager() *memory.MemoryManager {
	return r.memoryManager
}

func (r *AgentJamRuntime) GetToolchainManager() *toolchain.ToolchainManager {
	return r.toolchainManager
}

func (r *AgentJamRuntime) GetToolDispatcher() *dispatcher.ToolDispatcher {
	return r.toolDispatcher
}
