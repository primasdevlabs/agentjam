/**
 * @agentjam/validator — Cross-Reference Validator
 *
 * Validates links between agents, skills, workflows, tools, and stacks across the repository.
 */

import { discoverResources, parseAgent, parseWorkflow, parseStack, type ResourceBundle } from '@agentjam/parser';
import type { AgentManifest, WorkflowManifest, StackManifest } from '@agentjam/core';
import type { ValidationError } from './schema-validator.js';

export function validateCrossReferences(rootDir: string): ValidationError[] {
  const errors: ValidationError[] = [];
  const discovered = discoverResources(rootDir);

  // Set of known skill IDs, tool IDs, and workflow IDs
  const knownSkillIds = new Set<string>();
  const knownToolIds = new Set<string>();
  const knownWorkflowIds = new Set<string>();

  const agentBundles: ResourceBundle<AgentManifest>[] = [];
  const workflowBundles: ResourceBundle<WorkflowManifest>[] = [];
  const stackBundles: ResourceBundle<StackManifest>[] = [];

  for (const item of discovered) {
    try {
      if (item.type === 'skill') {
        knownSkillIds.add(item.id);
      } else if (item.type === 'tool') {
        knownToolIds.add(item.id);
      } else if (item.type === 'workflow') {
        knownWorkflowIds.add(item.id);
        workflowBundles.push(parseWorkflow(item.path));
      } else if (item.type === 'agent') {
        agentBundles.push(parseAgent(item.path));
      } else if (item.type === 'stack') {
        stackBundles.push(parseStack(item.path));
      }
    } catch {
      // Ignore parse errors here as schema-validator handles them
    }
  }

  // 1. Check agent declared skill references
  for (const agentBundle of agentBundles) {
    const { manifest, basePath } = agentBundle;
    if (Array.isArray(manifest.skills)) {
      for (const skillId of manifest.skills) {
        if (!knownSkillIds.has(skillId)) {
          errors.push({
            resourcePath: basePath,
            resourceId: manifest.name,
            field: 'skills',
            message: `Agent '${manifest.name}' references unknown skill '${skillId}'`,
            severity: 'warning',
          });
        }
      }
    }
  }

  // 2. Check stack profile policy references
  for (const stackBundle of stackBundles) {
    const { manifest, basePath } = stackBundle;
    if (Array.isArray(manifest.policies)) {
      for (const policyId of manifest.policies) {
        // Stack policy checks can be added here
      }
    }
  }

  return errors;
}
