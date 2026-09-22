import { z } from 'zod';

// ---------------------------------------------------------------------------
// Shared Primitives
// ---------------------------------------------------------------------------

export const ResourceTypeSchema = z.enum([
  'agent', 'skill', 'tool', 'workflow', 'policy',
  'stack', 'language', 'integration', 'prompt', 'template',
]);
export type ResourceType = z.infer<typeof ResourceTypeSchema>;

export const SemVerSchema = z.string().regex(/^\^?\d+\.\d+\.\d+(-[0-9A-Za-z.-]+)?$/);

export const DependencyReferenceSchema = z.object({
  name: z.string(),
  version: SemVerSchema,
  optional: z.boolean().optional(),
});
export type DependencyReference = z.infer<typeof DependencyReferenceSchema>;

export const MCPServerConfigSchema = z.object({
  name: z.string(),
  command: z.string(),
  args: z.array(z.string()).optional(),
  env: z.record(z.string()).optional(),
  transport: z.enum(['stdio', 'sse']).default('stdio'),
});
export type MCPServerConfig = z.infer<typeof MCPServerConfigSchema>;

export interface ResourceBundle<T = unknown> {
  manifest: T;
  instructions: Record<string, string>;
  basePath: string;
  metadata?: Record<string, unknown>;
}


// ---------------------------------------------------------------------------
// Agent Manifest
// ---------------------------------------------------------------------------

export const AgentManifestSchema = z.object({
  name: z.string(),
  version: SemVerSchema,
  type: z.literal('agent'),
  description: z.string(),
  author: z.string().optional(),
  license: z.string().optional(),
  skills: z.array(z.string()).default([]),
  tools: z.array(z.string()).default([]),
  mcpServers: z.array(MCPServerConfigSchema).optional(),
  dependencies: z.object({
    skills: z.array(DependencyReferenceSchema).optional(),
    tools: z.array(DependencyReferenceSchema).optional(),
  }).optional(),
  inputs: z.array(z.string()).default([]),
  outputs: z.array(z.string()).default([]),
  compatibility: z.array(z.string()).optional(),
  metadata: z.record(z.unknown()).optional(),
});
export type AgentManifest = z.infer<typeof AgentManifestSchema>;

// ---------------------------------------------------------------------------
// Skill Manifest
// ---------------------------------------------------------------------------

export const SkillManifestSchema = z.object({
  name: z.string(),
  version: SemVerSchema,
  type: z.literal('skill'),
  description: z.string(),
  author: z.string().optional(),
  category: z.string().optional(),
  triggers: z.array(z.string()).optional(),
  parameters: z.record(z.object({
    type: z.string(),
    description: z.string(),
    required: z.boolean().optional(),
    default: z.unknown().optional(),
  })).optional(),
  dependencies: z.object({
    skills: z.array(DependencyReferenceSchema).optional(),
  }).optional(),
  compatibility: z.array(z.string()).optional(),
});
export type SkillManifest = z.infer<typeof SkillManifestSchema>;

// ---------------------------------------------------------------------------
// Tool Manifest
// ---------------------------------------------------------------------------

export const ToolManifestSchema = z.object({
  name: z.string(),
  version: SemVerSchema,
  type: z.literal('tool'),
  description: z.string(),
  category: z.string().optional(),
  capabilities: z.array(z.string()).default([]),
  safetyLevel: z.enum(['read-only', 'safe-write', 'destructive', 'admin']).default('read-only'),
  mcpServer: MCPServerConfigSchema.optional(),
  parameters: z.record(z.unknown()).optional(),
});
export type ToolManifest = z.infer<typeof ToolManifestSchema>;

// ---------------------------------------------------------------------------
// Workflow Manifest
// ---------------------------------------------------------------------------

export const WorkflowStepSchema = z.object({
  id: z.string(),
  name: z.string(),
  agent: z.string().optional(),
  skill: z.string().optional(),
  description: z.string().optional(),
  inputs: z.record(z.unknown()).optional(),
  onFailure: z.enum(['abort', 'continue', 'retry']).optional(),
});

export const WorkflowManifestSchema = z.object({
  name: z.string(),
  version: SemVerSchema,
  type: z.literal('workflow'),
  description: z.string(),
  mode: z.enum(['single-agent', 'multi-agent']).default('single-agent'),
  agents: z.array(z.string()).optional(),
  skills: z.array(z.string()).optional(),
  steps: z.array(WorkflowStepSchema),
  inputs: z.array(z.string()).default([]),
  outputs: z.array(z.string()).default([]),
});
export type WorkflowManifest = z.infer<typeof WorkflowManifestSchema>;

// ---------------------------------------------------------------------------
// Stack Manifest (NEW)
// ---------------------------------------------------------------------------

export const StackDefaultsSchema = z.object({
  language: z.string(),
  package_manager: z.string().optional(),
  linter: z.string().optional(),
  formatter: z.string().optional(),
  test_runner: z.string().optional(),
  style_system: z.string().optional(),
});
export type StackDefaults = z.infer<typeof StackDefaultsSchema>;

export const StackManifestSchema = z.object({
  name: z.string(),
  version: SemVerSchema,
  type: z.literal('stack'),
  description: z.string(),
  framework: z.string(),
  ecosystem: z.string(),
  defaults: StackDefaultsSchema,
  conventions: z.record(z.string()).optional(),
  policies: z.array(z.string()).optional(),
  documentation: z.record(z.string()).optional(),
});
export type StackManifest = z.infer<typeof StackManifestSchema>;

// ---------------------------------------------------------------------------
// Policy Manifest (canonical home — previously in policy-engine)
// ---------------------------------------------------------------------------

export const PolicyEnforcementSchema = z.enum(['strict-block', 'warning', 'info']);
export type PolicyEnforcement = z.infer<typeof PolicyEnforcementSchema>;

export const PolicyCategorySchema = z.enum([
  'technology', 'versions', 'architecture', 'security',
  'dependencies', 'coding-standards', 'freshness', 'design',
  'project', 'core',
]);
export type PolicyCategory = z.infer<typeof PolicyCategorySchema>;

export const PolicyManifestSchema = z.object({
  id: z.string(),
  name: z.string(),
  description: z.string(),
  category: PolicyCategorySchema,
  enforcement: PolicyEnforcementSchema.default('strict-block'),
  scope: z.enum(['global', 'stack', 'project']).default('global'),
  appliesTo: z.array(z.string()).optional(),
  rules: z.record(z.unknown()).optional(),
  examples: z.array(z.object({
    title: z.string(),
    violation: z.string().optional(),
    compliant: z.string().optional(),
  })).optional(),
  instructions: z.string().optional(),
});
export type PolicyManifest = z.infer<typeof PolicyManifestSchema>;

// ---------------------------------------------------------------------------
// Design Policy Rule (structured pattern-matching for design governance)
// ---------------------------------------------------------------------------

export const DesignPolicyRuleSchema = z.object({
  /** Human-readable rule name. */
  name: z.string(),
  /** Description of what this rule enforces. */
  description: z.string(),
  /** Regex patterns to flag as violations. */
  prohibitedPatterns: z.array(z.string()).optional(),
  /** Literal strings to flag as violations. */
  prohibitedLiterals: z.array(z.string()).optional(),
  /** Required patterns that must be present. */
  requiredPatterns: z.array(z.string()).optional(),
  /** File glob patterns this rule applies to (e.g. '*.tsx', '*.css'). */
  fileGlobs: z.array(z.string()).optional(),
  /** Exemptions — patterns or files that are allowed despite the rule. */
  exemptions: z.array(z.string()).optional(),
  /** Suggested remediation. */
  suggestion: z.string().optional(),
});
export type DesignPolicyRule = z.infer<typeof DesignPolicyRuleSchema>;

// ---------------------------------------------------------------------------
// Environment Schemas
// ---------------------------------------------------------------------------

export const EnvironmentTypeSchema = z.enum(['ide', 'autonomous-agents', 'cli', 'extensions', 'ai-platforms', 'generic']);
export type EnvironmentType = z.infer<typeof EnvironmentTypeSchema>;

// ---------------------------------------------------------------------------
// Integration Manifest
// ---------------------------------------------------------------------------

export const IntegrationManifestSchema = z.object({
  name: z.string(),
  type: EnvironmentTypeSchema.optional(),
  harness: z.string().optional(),
  description: z.string().optional(),
  version: SemVerSchema.optional(),
  status: z.enum(['supported', 'partial', 'experimental', 'planned']).default('supported'),
  compatibility_level: z.number().min(0).max(6).default(5),
  capabilities: z.record(z.unknown()),
  mapping: z.object({
    instructionFormat: z.enum(['markdown', 'system-prompt', 'xml', 'json']).default('markdown'),
    fileNamingConvention: z.string(),
    outputDir: z.string().optional(),
    toolCallStyle: z.string().optional(),
  }),
  limitations: z.array(z.string()).optional(),
});
export type IntegrationManifest = z.infer<typeof IntegrationManifestSchema>;

export const EnvironmentCapabilitiesSchema = z.object({
  instructions: z.boolean().default(true),
  skills: z.boolean().default(true),
  agents: z.union([z.boolean(), z.literal('partial')]).default(true),
  workflows: z.union([z.boolean(), z.literal('partial')]).default(true),
  tools: z.boolean().default(true),
  filesystem: z.boolean().default(true),
  terminal: z.boolean().default(true),
  browser: z.boolean().default(false),
  mcp: z.boolean().default(true),
  project_rules: z.boolean().default(true),
  context_files: z.boolean().default(true),
  hooks: z.boolean().default(false),
  lifecycle_events: z.boolean().default(false),
  validation: z.union([z.boolean(), z.literal('partial')]).default(false),
});
export type EnvironmentCapabilities = z.infer<typeof EnvironmentCapabilitiesSchema>;

export const EnvironmentIntegrationManifestSchema = z.object({
  name: z.string(),
  type: EnvironmentTypeSchema,
  status: z.enum(['supported', 'partial', 'experimental', 'planned']).default('supported'),
  compatibility_level: z.number().min(0).max(6).default(5),
  capabilities: EnvironmentCapabilitiesSchema,
  limitations: z.array(z.string()).optional(),
  mapping: z.object({
    instructionFormat: z.enum(['markdown', 'system-prompt', 'xml', 'json']).default('markdown'),
    fileNamingConvention: z.string(),
    outputDir: z.string().optional(),
  }),
});
export type EnvironmentIntegrationManifest = z.infer<typeof EnvironmentIntegrationManifestSchema>;

// ---------------------------------------------------------------------------
// Language Manifest
// ---------------------------------------------------------------------------

export const LanguageManifestSchema = z.object({
  id: z.string(),
  name: z.string(),
  ecosystem: z.string(),
  aliases: z.array(z.string()).optional(),
  extensions: z.array(z.string()),
  tooling: z.object({
    package_manager: z.string().optional(),
    linter: z.string().optional(),
    formatter: z.string().optional(),
    type_checker: z.string().optional(),
  }).optional(),
  framework_precedence: z.array(z.string()).optional(),
  documentation: z.record(z.string()).optional(),
});
export type LanguageManifest = z.infer<typeof LanguageManifestSchema>;

// ---------------------------------------------------------------------------
// Prompt Manifest (NEW)
// ---------------------------------------------------------------------------

export const PromptManifestSchema = z.object({
  name: z.string(),
  version: SemVerSchema,
  type: z.literal('prompt'),
  description: z.string(),
  category: z.string().optional(),
  variables: z.array(z.object({
    name: z.string(),
    description: z.string(),
    required: z.boolean().default(true),
    default: z.string().optional(),
  })).optional(),
  template: z.string(),
});
export type PromptManifest = z.infer<typeof PromptManifestSchema>;

// ---------------------------------------------------------------------------
// Template Manifest (NEW)
// ---------------------------------------------------------------------------

export const TemplateManifestSchema = z.object({
  name: z.string(),
  version: SemVerSchema,
  type: z.literal('template'),
  description: z.string(),
  resourceType: ResourceTypeSchema,
  files: z.array(z.object({
    path: z.string(),
    description: z.string().optional(),
  })),
});
export type TemplateManifest = z.infer<typeof TemplateManifestSchema>;

// ---------------------------------------------------------------------------
// Unified Resource Manifest (discriminated union)
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Project Configuration Schema
// ---------------------------------------------------------------------------

export const AgentJamProjectConfigSchema = z.object({
  stackProfile: z.string().optional(),
  versionPolicy: z.enum(['current-stable', 'project-compatible', 'pinned']).default('current-stable'),
  freshnessRequired: z.boolean().default(true),
  maxDocAge: z.string().default('7d'),
});
export type AgentJamProjectConfig = z.infer<typeof AgentJamProjectConfigSchema>;

export type ResourceManifest =
  | AgentManifest
  | SkillManifest
  | ToolManifest
  | WorkflowManifest
  | StackManifest
  | PolicyManifest
  | LanguageManifest
  | IntegrationManifest
  | PromptManifest
  | TemplateManifest;

