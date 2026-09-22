/**
 * @agentjam/core — Canonical Constants
 *
 * Single source of truth for all enum-like values, identifiers,
 * and configuration constants used across AgentJam packages.
 */

// ---------------------------------------------------------------------------
// Resource Types
// ---------------------------------------------------------------------------

/** All canonical resource types in the AgentJam system. */
export const RESOURCE_TYPES = [
  'agent',
  'skill',
  'tool',
  'workflow',
  'policy',
  'stack',
  'language',
  'integration',
  'prompt',
  'template',
] as const;
export type ResourceTypeName = (typeof RESOURCE_TYPES)[number];

// ---------------------------------------------------------------------------
// Policy Categories
// ---------------------------------------------------------------------------

/** All supported policy governance categories. */
export const POLICY_CATEGORIES = [
  'technology',
  'versions',
  'architecture',
  'security',
  'dependencies',
  'coding-standards',
  'freshness',
  'design',
  'project',
  'core',
] as const;
export type PolicyCategoryName = (typeof POLICY_CATEGORIES)[number];

// ---------------------------------------------------------------------------
// Design Policy Identifiers
// ---------------------------------------------------------------------------

/** Canonical IDs for each of the 12 design governance policies. */
export const DESIGN_POLICY_IDS = [
  'design-anti-slop',
  'design-typography',
  'design-color',
  'design-spacing',
  'design-icons',
  'design-components',
  'design-animation',
  'design-responsive',
  'design-accessibility',
  'design-copy',
  'design-visual-language',
  'design-design-system',
] as const;
export type DesignPolicyId = (typeof DESIGN_POLICY_IDS)[number];

// ---------------------------------------------------------------------------
// Environment Types
// ---------------------------------------------------------------------------

/** AI development environment categories. */
export const ENVIRONMENT_TYPES = [
  'ide',
  'autonomous-agents',
  'cli',
  'extensions',
  'ai-platforms',
  'generic',
] as const;
export type EnvironmentTypeName = (typeof ENVIRONMENT_TYPES)[number];

// ---------------------------------------------------------------------------
// Compatibility Levels
// ---------------------------------------------------------------------------

/**
 * Compatibility levels (0–6) describing the depth of AgentJam integration
 * an AI environment supports.
 */
export const COMPATIBILITY_LEVELS = {
  0: { name: 'Documentation', description: 'Manual consumption of AgentJam docs only.' },
  1: { name: 'Instructions', description: 'Project rules and markdown instructions (.cursorrules, CLAUDE.md, GEMINI.md).' },
  2: { name: 'Skills', description: 'Reusable task-specific AgentJam skill modules.' },
  3: { name: 'Tools', description: 'Connection to AgentJam abstract tools or MCP servers.' },
  4: { name: 'Workflows', description: 'Single and multi-agent execution flows.' },
  5: { name: 'Validation', description: 'AgentJam post-execution validation engine (npm run validate).' },
  6: { name: 'Lifecycle Integration', description: 'Participating in lifecycle events (before task, before commit, etc.).' },
} as const;
export type CompatibilityLevel = keyof typeof COMPATIBILITY_LEVELS;

// ---------------------------------------------------------------------------
// Precedence Order
// ---------------------------------------------------------------------------

/**
 * AgentJam enforces a strict precedence hierarchy when resolving
 * execution guidance.  Index 0 is highest priority.
 */
export const PRECEDENCE_ORDER = [
  'User explicit instruction',
  'Current project state (codebase, config, lockfiles)',
  'Project configuration (.agentjam/config.yaml)',
  'AgentJam project policy (.agentjam/policies/)',
  'AgentJam stack policy (stacks/<profile>/)',
  'AgentJam global policy (policies/)',
  'Environment defaults',
  'Model knowledge (reasoning fallback only)',
] as const;
export type PrecedenceLevel = (typeof PRECEDENCE_ORDER)[number];


// ---------------------------------------------------------------------------
// Manifest Filenames
// ---------------------------------------------------------------------------

/** Maps each resource type to its expected YAML manifest filename. */
export const MANIFEST_FILENAMES: Record<string, string> = {
  agent: 'agent.yaml',
  skill: 'skill.yaml',
  tool: 'tool.yaml',
  workflow: 'workflow.yaml',
  policy: 'policy.yaml',
  stack: 'stack.yaml',
  language: 'language.yaml',
  integration: 'manifest.yaml',
  prompt: 'prompt.yaml',
  template: 'template.yaml',
} as const;

// ---------------------------------------------------------------------------
// Policy Enforcement Levels
// ---------------------------------------------------------------------------

/** How a policy violation is treated at runtime. */
export const ENFORCEMENT_LEVELS = [
  'strict-block',
  'warning',
  'info',
] as const;
export type EnforcementLevel = (typeof ENFORCEMENT_LEVELS)[number];

// ---------------------------------------------------------------------------
// Safety Levels for Tools
// ---------------------------------------------------------------------------

/** Classifies the risk of a tool invocation. */
export const SAFETY_LEVELS = [
  'read-only',
  'safe-write',
  'destructive',
  'admin',
] as const;
export type SafetyLevel = (typeof SAFETY_LEVELS)[number];

// ---------------------------------------------------------------------------
// Workflow Modes
// ---------------------------------------------------------------------------

export const WORKFLOW_MODES = ['single-agent', 'multi-agent'] as const;
export type WorkflowMode = (typeof WORKFLOW_MODES)[number];

// ---------------------------------------------------------------------------
// Stack Ecosystem Tags
// ---------------------------------------------------------------------------

export const STACK_ECOSYSTEMS = [
  'web',
  'backend',
  'mobile',
  'systems',
  'data',
  'infrastructure',
] as const;
export type StackEcosystem = (typeof STACK_ECOSYSTEMS)[number];
