/**
 * @agentjam/validator — Schema Validator
 *
 * Validates all discovered manifests against their canonical Zod schemas.
 */

import { discoverResources, parseResource } from '@agentjam/parser';

export interface ValidationError {
  resourcePath: string;
  resourceId?: string;
  field?: string;
  message: string;
  severity: 'error' | 'warning';
}

export function validateSchemas(rootDir: string): ValidationError[] {
  const errors: ValidationError[] = [];
  const discovered = discoverResources(rootDir);

  for (const item of discovered) {
    try {
      parseResource(item.path, item.type);
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      errors.push({
        resourcePath: item.path,
        message: `Schema validation failed for ${item.type} at ${item.path}: ${msg}`,
        severity: 'error',
      });
    }
  }

  return errors;
}
