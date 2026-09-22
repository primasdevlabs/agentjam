import fs from 'node:fs';
import path from 'node:path';
import yaml from 'js-yaml';
import { z } from 'zod';

export const PolicyManifestSchema = z.object({
  id: z.string(),
  name: z.string(),
  description: z.string(),
  category: z.enum(['technology', 'versions', 'architecture', 'security', 'dependencies', 'coding-standards', 'freshness', 'design']),
  enforcement: z.enum(['strict-block', 'warning', 'info']).default('strict-block'),
  rules: z.record(z.unknown()).optional(),
  instructions: z.string().optional(),
});
export type PolicyManifest = z.infer<typeof PolicyManifestSchema>;

export interface PolicyEvaluationResult {
  allowed: boolean;
  violations: Array<{
    policyId: string;
    message: string;
  }>;
}

export function parsePolicy(filePath: string): PolicyManifest {
  const content = fs.readFileSync(filePath, 'utf-8');
  if (filePath.endsWith('.md')) {
    // Markdown design policy
    const filename = path.basename(filePath, '.md');
    return {
      id: `design-${filename}`,
      name: `Design Policy: ${filename}`,
      description: `Markdown design governance rule for ${filename}`,
      category: 'design',
      enforcement: 'strict-block',
      instructions: content,
    };
  }
  const raw = yaml.load(content);
  return PolicyManifestSchema.parse(raw);
}

export class PolicyEngine {
  private policies: PolicyManifest[] = [];

  loadPoliciesFromDirectory(policiesDir: string): void {
    if (!fs.existsSync(policiesDir)) return;
    const files = fs.readdirSync(policiesDir, { recursive: true });
    for (const file of files) {
      if (typeof file === 'string') {
        const fullPath = path.join(policiesDir, file);
        if (file.endsWith('.yaml') || file.endsWith('.yml') || file.endsWith('.md')) {
          try {
            const policy = parsePolicy(fullPath);
            this.policies.push(policy);
          } catch {
            // Ignore non-policy files cleanly
          }
        }
      }
    }
  }

  getPolicies(): PolicyManifest[] {
    return this.policies;
  }

  evaluateDependency(dependencyName: string): PolicyEvaluationResult {
    const violations: Array<{ policyId: string; message: string }> = [];

    for (const policy of this.policies) {
      if (policy.category === 'dependencies' && policy.rules?.forbidden) {
        const forbidden = policy.rules.forbidden as string[];
        if (forbidden.includes(dependencyName.toLowerCase())) {
          violations.push({
            policyId: policy.id,
            message: `Dependency '${dependencyName}' violates policy '${policy.name}': prohibited stack package.`,
          });
        }
      }
    }

    return {
      allowed: violations.length === 0,
      violations,
    };
  }

  evaluateDesignContent(content: string): PolicyEvaluationResult {
    const violations: Array<{ policyId: string; message: string }> = [];

    // 1. Emoji icon check
    const emojiRegex = /[\u{1F300}-\u{1F6FF}\u{1F900}-\u{1F9FF}\u{2600}-\u{26FF}\u{2700}-\u{27BF}]/u;
    if (emojiRegex.test(content)) {
      violations.push({
        policyId: 'design-icons',
        message: 'Emoji used as UI icons or interface elements is strictly prohibited by Design Governance Policy.',
      });
    }

    // 2. Prohibited AI marketing buzzwords check
    const forbiddenBuzzwords = [
      'empower your business',
      'unlock the power',
      'transform your workflow',
      'seamlessly connect',
      'next-generation platform',
      'cutting-edge solutions',
      'supercharge your',
      'elevate your',
    ];
    const lowerContent = content.toLowerCase();
    for (const buzzword of forbiddenBuzzwords) {
      if (lowerContent.includes(buzzword)) {
        violations.push({
          policyId: 'design-copy',
          message: `Forbidden AI marketing buzzword detected: "${buzzword}". Use direct, human product copy.`,
        });
      }
    }

    return {
      allowed: violations.length === 0,
      violations,
    };
  }
}
