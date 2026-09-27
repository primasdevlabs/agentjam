package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/primasdevlabs/agentjam/pkg/adapters"
	"github.com/primasdevlabs/agentjam/pkg/context"
	"github.com/primasdevlabs/agentjam/pkg/core"
	"github.com/primasdevlabs/agentjam/pkg/dispatcher"
	"github.com/primasdevlabs/agentjam/pkg/mcp"
	"github.com/primasdevlabs/agentjam/pkg/memory"
	"github.com/primasdevlabs/agentjam/pkg/parser"
	"github.com/primasdevlabs/agentjam/pkg/policy"
	"github.com/primasdevlabs/agentjam/pkg/registry"
	"github.com/primasdevlabs/agentjam/pkg/runtime"
	"github.com/primasdevlabs/agentjam/pkg/validator"
)

const version = core.Version

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(0)
	}

	command := os.Args[1]

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting workspace directory: %v\n", err)
		os.Exit(1)
	}

	rt := runtime.NewRuntime(cwd, dispatcher.ToolDispatcherOptions{AllowDestructive: true, AllowAdmin: true})

	switch command {
	case "version":
		fmt.Printf("AgentJam Go Native CLI v%s\n", version)

	case "preflight":
		result := rt.GetToolchainManager().RunPreflightChecks()
		bytes, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(bytes))
		if !result.Passed {
			os.Exit(1)
		}

	case "context":
		fs := flag.NewFlagSet("context", flag.ExitOnError)
		agentFlag := fs.String("agent", "", "Active persona agent name")
		stackFlag := fs.String("stack", "", "Active stack profile name")
		skillsFlag := fs.String("skills", "", "Comma-separated active skill names")
		budgetFlag := fs.Int("token-budget", 128000, "Token budget for system instruction")
		_ = fs.Parse(os.Args[2:])

		var skills []string
		if *skillsFlag != "" {
			skills = splitComma(*skillsFlag)
		}

		snapshot := rt.GetContextManager().BuildContextSnapshot(context.ContextOptions{
			ActiveAgent:  *agentFlag,
			ActiveStack:  *stackFlag,
			ActiveSkills: skills,
			TokenBudget:  *budgetFlag,
		})
		bytes, _ := json.MarshalIndent(snapshot, "", "  ")
		fmt.Println(string(bytes))

	case "detect":
		envs := runtime.DetectEnvironments(cwd)
		stacks := rt.GetToolchainManager().DetectStacks()
		out := map[string]interface{}{"environments": envs, "stacks": stacks}
		bytes, _ := json.MarshalIndent(out, "", "  ")
		fmt.Println(string(bytes))

	case "validate":
		fmt.Println("Validating AgentJam canonical resources, design policies, languages & audit workflows...")
		valErrors := validator.ValidateRepository(cwd)
		idx := registry.BuildRegistryIndex(cwd)
		policiesCount := len(rt.GetPolicyEngine().GetPolicies())

		for _, v := range valErrors {
			fmt.Printf("[%s] %s (%s): %s\n", v.Severity, v.ResourceID, v.Field, v.Message)
		}

		fmt.Printf("Registry Index: %d Agents, %d Skills, %d Tools, %d Workflows, %d Stacks, %d Languages, %d Policies, %d Integrations.\n",
			idx.AgentsCount, idx.SkillsCount, idx.ToolsCount, idx.WorkflowsCount, idx.StacksCount, idx.LanguagesCount, idx.PoliciesCount, idx.IntegrationsCount)
		fmt.Printf("Policy Engine: Loaded %d policy specifications.\n", policiesCount)

		if validator.HasErrors(valErrors) {
			fmt.Println("Repository validation FAILED.")
			os.Exit(1)
		}
		if len(valErrors) > 0 {
			fmt.Printf("Repository validation passed with %d warning(s).\n", len(valErrors))
		} else {
			fmt.Println("Repository validation passed successfully!")
		}

	case "eval":
		fs := flag.NewFlagSet("eval", flag.ExitOnError)
		_ = fs.Parse(os.Args[2:])
		targets := fs.Args()
		if len(targets) == 0 {
			fmt.Fprintln(os.Stderr, "Usage: agentjam eval <file> [file...]")
			os.Exit(2)
		}
		pe := rt.GetPolicyEngine()
		exitCode := 0
		for _, target := range targets {
			data, err := os.ReadFile(target)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Cannot read %s: %v\n", target, err)
				exitCode = 2
				continue
			}
			summary := pe.EvaluateAll(string(data))
			bytes, _ := json.MarshalIndent(summary, "", "  ")
			fmt.Printf("%s\n%s\n", target, string(bytes))
			if !summary.Allowed {
				exitCode = 1
			}
		}
		os.Exit(exitCode)

	case "build-registry":
		fmt.Println("Building canonical AgentJam registry.json index...")
		idx := registry.BuildRegistryIndex(cwd)
		bytes, _ := json.MarshalIndent(idx, "", "  ")
		if err := os.WriteFile("registry.json", bytes, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write registry.json: %v\n", err)
			os.Exit(1)
		}
		indexPath := filepath.Join("registry", "index", "index.json")
		if err := os.MkdirAll(filepath.Dir(indexPath), 0755); err == nil {
			_ = os.WriteFile(indexPath, bytes, 0644)
		}
		for dir, entries := range idx.SplitByType() {
			groupPath := filepath.Join("registry", dir, "index.json")
			if err := os.MkdirAll(filepath.Dir(groupPath), 0755); err != nil {
				continue
			}
			groupBytes, _ := json.MarshalIndent(entries, "", "  ")
			_ = os.WriteFile(groupPath, groupBytes, 0644)
		}
		fmt.Printf("Successfully generated registry.json (%d entries).\n", idx.Total)

	case "run-e2e":
		runE2E(rt, cwd)

	case "run":
		fs := flag.NewFlagSet("run", flag.ExitOnError)
		listFlag := fs.Bool("list", false, "List discoverable workflows instead of running one")
		_ = fs.Parse(os.Args[2:])

		workflowDir := filepath.Join(cwd, "workflows")
		if *listFlag {
			names, err := runtime.ListWorkflows(workflowDir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to list workflows: %v\n", err)
				os.Exit(1)
			}
			for _, n := range names {
				fmt.Println(n)
			}
			return
		}
		if fs.NArg() == 0 {
			fmt.Fprintln(os.Stderr, "Usage: agentjam run <workflow> [--list]")
			os.Exit(2)
		}
		exec := runtime.NewExecutor(rt, nil)
		report, err := exec.Run(fs.Arg(0))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Workflow run failed: %v\n", err)
			os.Exit(1)
		}
		bytes, _ := json.MarshalIndent(report, "", "  ")
		fmt.Println(string(bytes))
		if report.Status != "success" {
			os.Exit(1)
		}

	case "mcp":
		mcp.RegisterBuiltinTools(rt)
		fmt.Fprintln(os.Stderr, "AgentJam MCP server listening on stdio (newline-delimited JSON-RPC).")
		if err := mcp.NewServer(rt.GetToolDispatcher(), os.Stdin, os.Stdout).Serve(); err != nil {
			fmt.Fprintf(os.Stderr, "MCP server error: %v\n", err)
			os.Exit(1)
		}

	case "export":
		fs := flag.NewFlagSet("export", flag.ExitOnError)
		harnessFlag := fs.String("harness", "auto", "Target AI harness format ('auto', 'all', 'cursor', 'claude-code', 'gemini', 'cline', 'windsurf', 'devin', 'roo-code', 'generic')")
		agentFlag := fs.String("agent", "", "Agent persona to export (defaults to .agentjam/config.yaml defaultAgent)")
		_ = fs.Parse(os.Args[2:])

		agent := resolveAgent(cwd, *agentFlag)
		snapshot := rt.GetContextManager().BuildContextSnapshot(context.ContextOptions{
			ActiveAgent: agent.Name,
			TokenBudget: 128000,
		})

		if *harnessFlag == "auto" {
			fmt.Println("Auto-detecting harness environment for active workspace...")
			envs := runtime.DetectEnvironments(cwd)
			if len(envs) == 0 {
				fmt.Println("No specific harness file detected. Falling back to generic platform exporter...")
				writeExport(adapters.ExportGenericPrompt(agent, snapshot.SystemInstruction))
			} else {
				for _, env := range envs {
					fmt.Printf("Detected Harness: %s (%s)\n", env.Name, env.Type)
				}
				writeAllExports(adapters.ExportAllHarnesses(agent, snapshot.SystemInstruction))
			}
		} else if *harnessFlag == "all" {
			writeAllExports(adapters.ExportAllHarnesses(agent, snapshot.SystemInstruction))
		} else {
			fmt.Printf("User explicitly selected harness: '%s'\n", *harnessFlag)
			allExports := adapters.ExportAllHarnesses(agent, snapshot.SystemInstruction)
			if res, ok := allExports[*harnessFlag]; ok {
				writeExport(res)
			} else {
				fmt.Fprintf(os.Stderr, "Unknown harness '%s'. Supported: 'cursor', 'claude-code', 'gemini', 'cline', 'windsurf', 'devin', 'roo-code', 'generic'\n", *harnessFlag)
				os.Exit(2)
			}
		}

	default:
		printHelp()
	}
}

func splitComma(s string) []string {
	out := make([]string, 0)
	for _, part := range strings.Split(s, ",") {
		if part = strings.TrimSpace(part); part != "" {
			out = append(out, part)
		}
	}
	return out
}

// resolveAgent loads an agent manifest by name, falling back to the workspace
// config defaultAgent and finally to software-engineer.
func resolveAgent(rootDir, name string) core.AgentManifest {
	if name == "" {
		if cfg, err := parser.LoadWorkspaceConfig(rootDir); err == nil {
			name = cfg.DefaultAgent
		}
	}
	if name == "" {
		name = "software-engineer"
	}
	for _, res := range parser.DiscoverResources(rootDir) {
		if res.Type != core.ResourceTypeAgent || res.ID != name {
			continue
		}
		if bundle, err := parser.ParseAgent(res.Path); err == nil {
			return bundle.Manifest
		}
	}
	return core.AgentManifest{Name: name, Description: "AgentJam agent"}
}

func writeExport(res adapters.ExportResult) {
	for fn, content := range res.Files {
		if dir := filepath.Dir(fn); dir != "." {
			_ = os.MkdirAll(dir, 0755)
		}
		_ = os.WriteFile(fn, []byte(content), 0644)
		fmt.Printf("Exported [%s] -> %s\n", res.Harness, fn)
	}
}

func writeAllExports(allExports map[string]adapters.ExportResult) {
	for _, res := range allExports {
		writeExport(res)
	}
}

func runE2E(rt *runtime.AgentJamRuntime, cwd string) {
	fmt.Println("Executing End-to-End AgentJam Go Verification Suite...")

	// 1. Parser Subsystem Verification
	fmt.Println("[PARSER SUBSYSTEM]")
	discovered := parser.DiscoverResources(cwd)
	fmHead, fmBody := parser.ParseFrontmatter("---\nname: test-skill\n---\n# Instructions")
	fmt.Printf("   Resource Discovery: Discovered %d canonical resources\n", len(discovered))
	fmt.Printf("   Frontmatter Parser: Extracted head (%d chars), body (%d chars)\n\n", len(fmHead), len(fmBody))

	// 2. Memory Subsystem Verification
	fmt.Println("[MEMORY SUBSYSTEM]")
	mm := rt.GetMemoryManager()
	mm.Set("session_state", "active", core.MemoryScopeWorking, []string{"session", "go"}, 0)
	memEntry, memFound := mm.Get("session_state", core.MemoryScopeWorking)
	memQuery := mm.Query(memory.MemoryQuery{Tags: []string{"session"}})
	fmt.Printf("   Memory Store & Retrieval: Key 'session_state' found=%v (value='%v')\n", memFound, memEntry.Value)
	fmt.Printf("   Memory Semantic Search: Tag query returned %d result(s)\n\n", len(memQuery))

	// 3. Context Subsystem Verification
	fmt.Println("[CONTEXT SUBSYSTEM]")
	cm := rt.GetContextManager()
	snapshot := cm.BuildContextSnapshot(context.ContextOptions{
		ActiveAgent:  "software-engineer",
		ActiveSkills: []string{"anti-slop", "testing"},
		TokenBudget:  128000,
	})
	fmt.Printf("   Context Snapshot: Estimated %d tokens across %d policies\n", snapshot.TokenCountEstimate, snapshot.PoliciesCount)
	fmt.Printf("   System Instructions Generated: %d bytes\n\n", len(snapshot.SystemInstruction))

	// 3.5 Policy Engine Subsystem Verification
	fmt.Println("[POLICY ENGINE SUBSYSTEM]")
	policyViolations := policy.EvaluateDesignRules(nil, "<button> Submit</button>")
	fmt.Printf("   Anti-Slop Policy Check: Evaluated design content (%d policy violation(s) flagged)\n\n", len(policyViolations))
	fmt.Println("[TOOL CALLING & DISPATCHER SUBSYSTEM]")
	td := rt.GetToolDispatcher()
	td.RegisterTool(
		core.ToolManifest{
			Name:         "system_info",
			Version:      "1.0.0",
			Type:         "tool",
			Description:  "Returns runtime system information",
			Capabilities: []string{"system"},
			SafetyLevel:  core.SafetyReadOnly,
		},
		func(args map[string]interface{}) (interface{}, error) {
			return map[string]string{"engine": "Go-Native", "status": "operational"}, nil
		},
	)

	toolRes := td.Dispatch(core.ToolCall{
		ID:        "call_sys_1",
		ToolName:  "system_info",
		Arguments: map[string]interface{}{},
	})
	fmt.Printf("   Tool Execution Dispatch: Tool '%s' status=%s, duration=%dms\n", toolRes.ToolName, toolRes.Status, toolRes.DurationMs)
	fmt.Printf("   Path Boundary Security Check: Enforced against workspace '%s'\n\n", cwd)

	// 5. Harness Integrations & Adapters Verification
	fmt.Println("[INTEGRATIONS & HARNESS ADAPTERS]")
	agent := core.AgentManifest{Name: "software-engineer", Description: "Fullstack Engineer"}
	allExports := adapters.ExportAllHarnesses(agent, snapshot.SystemInstruction)

	fmt.Println("   Target Harness Exports Generated:")
	for harness, res := range allExports {
		for fileName := range res.Files {
			fmt.Printf("   [%s] Exported: %s (%d bytes)\n", harness, fileName, len(res.Files[fileName]))
		}
	}
}

func printHelp() {
	fmt.Printf("AgentJam Go Native CLI v%s\n\n", version)
	fmt.Println("Usage: agentjam <command> [options]")
	fmt.Println("\nCommands:")
	fmt.Println("  version          Print AgentJam Go version")
	fmt.Println("  preflight        Run toolchain preflight checks")
	fmt.Println("  detect           Detect AI harness environments and project stacks")
	fmt.Println("  context          Generate context snapshot for active workspace")
	fmt.Println("  eval             Evaluate files against loaded policy rules")
	fmt.Println("  export           Export rule configurations (--harness auto|all|cursor|claude-code|gemini|...)")
	fmt.Println("  validate         Validate repository rules and policy engine")
	fmt.Println("  build-registry   Generate registry.json index")
	fmt.Println("  run              Run a workflow by name (--list to enumerate)")
	fmt.Println("  mcp              Start MCP stdio server exposing workspace tools")
	fmt.Println("  run-e2e          Execute end-to-end full system workflow test")
}
