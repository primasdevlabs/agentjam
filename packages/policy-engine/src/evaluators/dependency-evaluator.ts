/**
 * @agentjam/policy-engine — Dependency Evaluator
 *
 * Evaluates dependency policy rules (prohibited packages, version caps, required packages).
 */

import type { PolicyManifest } from '@agentjam/core';

export interface EvaluationViolation {
  policyId: string;
  policyName: string;
  category: string;
  enforcement: 'strict-block' | 'warning' | 'info';
  message: string;
  rule?: string;
  context?: Record<string, unknown>;
}

export function evaluateDependencyRules(
  policies: PolicyManifest[],
  dependencyName: string,
  version?: string
): EvaluationViolation[] {
  const violations: EvaluationViolation[] = [];
  const lowerDep = dependencyName.toLowerCase();

  for (const policy of policies) {
    if (policy.category !== 'dependencies') continue;
    const rules = policy.rules || {};

    // Check forbidden dependencies
    if (Array.isArray(rules.forbidden)) {
      const forbidden = rules.forbidden.map((f: unknown) => String(f).toLowerCase());
      if (forbidden.includes(lowerDep)) {
        violations.push({
          policyId: policy.id,
          policyName: policy.name,
          category: policy.category,
          enforcement: policy.enforcement,
          message: `Dependency '${dependencyName}' violates policy '${policy.name}': prohibited dependency.`,
          rule: 'forbidden',
          context: { dependencyName, version, forbiddenList: rules.forbidden },
        });
      }
    }

    // Check allowed list if explicitly restricted
    if (Array.isArray(rules.allowedOnly) && rules.allowedOnly.length > 0) {
      const allowedOnly = rules.allowedOnly.map((a: unknown) => String(a).toLowerCase());
      if (!allowedOnly.includes(lowerDep)) {
        violations.push({
          policyId: policy.id,
          policyName: policy.name,
          category: policy.category,
          enforcement: policy.enforcement,
          message: `Dependency '${dependencyName}' is not in the allowed dependency list for policy '${policy.name}'.`,
          rule: 'allowedOnly',
          context: { dependencyName, version, allowedOnly: rules.allowedOnly },
        });
      }
    }
  }

  return violations;
}
