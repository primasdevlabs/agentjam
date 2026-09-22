/**
 * @agentjam/validator — Instruction Validator
 *
 * Checks completeness of instruction files across agents, skills, and workflows.
 */

import { discoverResources, parseAgent, parseSkill, parseWorkflow } from '@agentjam/parser';
import type { ValidationError } from './schema-validator.js';

export function validateInstructions(rootDir: string): ValidationError[] {
  const errors: ValidationError[] = [];
  const discovered = discoverResources(rootDir);

  for (const item of discovered) {
    if (item.type === 'agent') {
      try {
        const bundle = parseAgent(item.path);
        if (!bundle.instructions['role.md']) {
          errors.push({
            resourcePath: item.path,
            resourceId: item.id,
            field: 'instructions/role.md',
            message: `Agent '${item.id}' is missing role.md in instructions/`,
            severity: 'warning',
          });
        }
      } catch {
        // Handled by schema validator
      }
    } else if (item.type === 'workflow') {
      try {
        const bundle = parseWorkflow(item.path);
        if (Object.keys(bundle.instructions).length === 0) {
          errors.push({
            resourcePath: item.path,
            resourceId: item.id,
            field: 'instructions/',
            message: `Workflow '${item.id}' has no instruction markdown files in instructions/ or steps/`,
            severity: 'warning',
          });
        }
      } catch {
        // Handled by schema validator
      }
    }
  }

  return errors;
}
