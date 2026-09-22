/**
 * @agentjam/validator — Entry Point
 *
 * Combines schema validation, cross-reference integrity checks, and instruction completeness checks.
 */

import { validateSchemas, type ValidationError } from './schema-validator.js';
import { validateCrossReferences } from './cross-ref-validator.js';
import { validateInstructions } from './instruction-validator.js';

export * from './schema-validator.js';
export * from './cross-ref-validator.js';
export * from './instruction-validator.js';

export interface ValidationResult {
  valid: boolean;
  errors: ValidationError[];
  summary: {
    total: number;
    errorCount: number;
    warningCount: number;
  };
}

/**
 * Perform a full repository validation scan.
 */
export function validateRepository(rootDir: string): ValidationResult {
  const schemaErrors = validateSchemas(rootDir);
  const crossRefErrors = validateCrossReferences(rootDir);
  const instructionErrors = validateInstructions(rootDir);

  const allErrors = [...schemaErrors, ...crossRefErrors, ...instructionErrors];
  const errorCount = allErrors.filter(e => e.severity === 'error').length;
  const warningCount = allErrors.filter(e => e.severity === 'warning').length;

  return {
    valid: errorCount === 0,
    errors: allErrors,
    summary: {
      total: allErrors.length,
      errorCount,
      warningCount,
    },
  };
}
