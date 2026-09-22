/**
 * @agentjam/policy-engine — Architecture Policy Evaluator
 *
 * Evaluates module boundaries, private import violations, and layering rules.
 */

import type { PolicyManifest } from '@agentjam/core';
import type { EvaluationViolation } from './dependency-evaluator.js';

export function evaluateArchitectureRules(
  policies: PolicyManifest[],
  filePath: string,
  importedModule: string
): EvaluationViolation[] {
  const violations: EvaluationViolation[] = [];

  for (const policy of policies) {
    if (policy.category !== 'architecture') continue;

    // Prohibit accessing private internal paths of other modules
    if (importedModule.includes('/src/internal/') || importedModule.includes('/private/')) {
      violations.push({
        policyId: policy.id,
        policyName: policy.name,
        category: policy.category,
        enforcement: policy.enforcement,
        message: `Import '${importedModule}' in '${filePath}' accesses private/internal module paths. Use public module APIs only.`,
        rule: 'public-apis-only',
        context: { filePath, importedModule },
      });
    }
  }

  return violations;
}
