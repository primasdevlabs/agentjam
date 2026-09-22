import { z } from 'zod';

export const ResourceTypeSchema = z.enum(['agent', 'skill', 'tool', 'workflow', 'prompt', 'template', 'language']);
export type ResourceType = z.infer<typeof ResourceTypeSchema>;

export const SemVerSchema = z.string().regex(/^\^?\d+\.\d+\.\d+(-[0-9A-Za-z.-]+)?$/);

export const DependencyReferenceSchema = z.object({
  name: z.string(),
  version: SemVerSchema,
  optional: z.boolean().optional(),
});
export type DependencyReference = z.infer<typeof DependencyReferenceSchema>;

export const AgentManifestSchema = z.object({
  name: z.string(),
  version: SemVerSchema,
  type: z.literal('agent'),
  description: z.string(),
  author: z.string().optional(),
  license: z.string().optional(),
  skills: z.array(z.string()).default([]),
  tools: z.array(z.string()).default([]),
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

export const ToolManifestSchema = z.object({
  name: z.string(),
  version: SemVerSchema,
  type: z.literal('tool'),
  description: z.string(),
  category: z.string().optional(),
  capabilities: z.array(z.string()).default([]),
  safetyLevel: z.enum(['read-only', 'safe-write', 'destructive', 'admin']).default('read-only'),
  parameters: z.record(z.unknown()).optional(),
});
export type ToolManifest = z.infer<typeof ToolManifestSchema>;

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

export const IntegrationManifestSchema = z.object({
  name: z.string(),
  harness: z.string(),
  description: z.string(),
  version: SemVerSchema,
  capabilities: z.object({
    supportsAgents: z.boolean(),
    supportsSkills: z.boolean(),
    supportsWorkflows: z.boolean(),
    supportsTools: z.boolean(),
  }),
  mapping: z.object({
    instructionFormat: z.enum(['markdown', 'system-prompt', 'xml', 'json']),
    fileNamingConvention: z.string(),
    toolCallStyle: z.string().optional(),
  }),
  limitations: z.array(z.string()).optional(),
});
export type IntegrationManifest = z.infer<typeof IntegrationManifestSchema>;

export const EnvironmentTypeSchema = z.enum(['ide', 'autonomous-agents', 'cli', 'extensions', 'ai-platforms', 'generic']);
export type EnvironmentType = z.infer<typeof EnvironmentTypeSchema>;

export const EnvironmentCapabilitiesSchema = z.object({
  instructions: z.boolean().default(true),
  skills: z.boolean().default(true),
  agents: z.union([z.boolean(), z.literal('partial')]).default(true),
  workflows: z.union([z.boolean(), z.literal('partial')]).default(true),
  tools: z.boolean().default(true),
  filesystem: z.boolean().default(true),
  terminal: z.boolean().default(true),
  browser: z.boolean().default(false),
  mcp: z.boolean().default(false),
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
