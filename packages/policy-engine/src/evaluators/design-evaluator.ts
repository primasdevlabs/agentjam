/**
 * @agentjam/policy-engine — Design Governance Evaluator
 *
 * Evaluates code and content against UI design policies, icon standards,
 * copywriting guidelines, and visual design tokens.
 */

import type { PolicyManifest } from '@agentjam/core';
import type { EvaluationViolation } from './dependency-evaluator.js';

const EMOJI_REGEX = /[\u{1F300}-\u{1F6FF}\u{1F900}-\u{1F9FF}\u{2600}-\u{26FF}\u{2700}-\u{27BF}]/u;

const FORBIDDEN_BUZZWORDS = [
  'empower your business',
  'unlock the power',
  'transform your workflow',
  'seamlessly connect',
  'next-generation platform',
  'cutting-edge solutions',
  'supercharge your',
  'elevate your',
  'game-changer',
  'synergy',
  'disruptive',
  'operating-system-like',
  'productive-by-design',

];

export function evaluateDesignRules(
  policies: PolicyManifest[],
  content: string,
  filePath?: string
): EvaluationViolation[] {
  const violations: EvaluationViolation[] = [];
  const lowerContent = content.toLowerCase();

  for (const policy of policies) {
    if (policy.category !== 'design') continue;
    const rules = policy.rules || {};

    // 1. Emoji check
    if (rules.prohibitEmojis !== false && EMOJI_REGEX.test(content)) {
      violations.push({
        policyId: policy.id,
        policyName: policy.name,
        category: policy.category,
        enforcement: policy.enforcement,
        message: `Emoji icon detected in ${filePath || 'content'}: prohibited by Design Governance Policy. Use SVGs or Icon components.`,
        rule: 'prohibitEmojis',
      });
    }

    // 2. Buzzwords check
    if (rules.prohibitMarketingBuzzwords !== false) {
      for (const buzzword of FORBIDDEN_BUZZWORDS) {
        if (lowerContent.includes(buzzword)) {
          violations.push({
            policyId: policy.id,
            policyName: policy.name,
            category: policy.category,
            enforcement: policy.enforcement,
            message: `Forbidden marketing buzzword detected ("${buzzword}") in ${filePath || 'content'}: write direct, human product copy.`,
            rule: 'prohibitMarketingBuzzwords',
            context: { buzzword },
          });
        }
      }
    }

    // 3. Generic color check (e.g. hardcoded #ff0000 or plain 'red' inline styles)
    if (rules.prohibitGenericColors === true) {
      const genericColorRegex = /color:\s*(red|blue|green|yellow|purple);/i;
      if (genericColorRegex.test(content)) {
        violations.push({
          policyId: policy.id,
          policyName: policy.name,
          category: policy.category,
          enforcement: policy.enforcement,
          message: `Generic plain CSS color name detected in ${filePath || 'content'}: use design system color tokens or HSL palette variables.`,
          rule: 'prohibitGenericColors',
        });
      }
    }
  }

  return violations;
}
