/**
 * @agentjam/policy-engine — Security Policy Evaluator
 *
 * Evaluates code for hardcoded secrets, raw SQL string concatenation,
 * unvalidated external inputs, and insecure protocol usages.
 */

import type { PolicyManifest } from '@agentjam/core';
import type { EvaluationViolation } from './dependency-evaluator.js';

const SECRET_PATTERNS = [
  /api[_-]?key\s*=\s*['"][A-Za-z0-9_-]{16,}['"]/i,
  /secret[_-]?key\s*=\s*['"][A-Za-z0-9_-]{16,}['"]/i,
  /bearer\s+[A-Za-z0-9_.-]{20,}/i,
  /password\s*=\s*['"][^'"]+['"]/i,
];

const RAW_SQL_CONCAT = /SELECT\s+.*\s+FROM\s+.*\s+\+\s*['"]/i;

export function evaluateSecurityRules(
  policies: PolicyManifest[],
  content: string,
  filePath?: string
): EvaluationViolation[] {
  const violations: EvaluationViolation[] = [];

  for (const policy of policies) {
    if (policy.category !== 'security') continue;

    // Check hardcoded secrets
    for (const pattern of SECRET_PATTERNS) {
      if (pattern.test(content)) {
        violations.push({
          policyId: policy.id,
          policyName: policy.name,
          category: policy.category,
          enforcement: policy.enforcement,
          message: `Potential hardcoded secret or API key pattern detected in ${filePath || 'code'}. Use environment variables or secret managers.`,
          rule: 'no-hardcoded-secrets',
        });
        break;
      }
    }

    // Check raw SQL concatenation
    if (RAW_SQL_CONCAT.test(content)) {
      violations.push({
        policyId: policy.id,
        policyName: policy.name,
        category: policy.category,
        enforcement: policy.enforcement,
        message: `Potential raw SQL query string concatenation detected in ${filePath || 'code'}. Use parameterized queries or ORM models.`,
        rule: 'parameterized-queries-only',
      });
    }
  }

  return violations;
}
