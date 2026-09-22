/**
 * @agentjam/policy-engine — Main Policy Engine
 *
 * Loads, stores, and executes evaluation rules across all policy categories.
 */

import fs from 'node:fs';
import path from 'node:path';
import { type PolicyManifest, ParseError } from '@agentjam/core';
import { parsePolicy } from '@agentjam/parser';
import { evaluateDependencyRules, type EvaluationViolation } from './evaluators/dependency-evaluator.js';
import { evaluateDesignRules } from './evaluators/design-evaluator.js';
import { evaluateSecurityRules } from './evaluators/security-evaluator.js';
import { evaluateArchitectureRules } from './evaluators/architecture-evaluator.js';

export interface PolicyEngineSummary {
  allowed: boolean;
  totalViolations: number;
  strictBlocks: number;
  warnings: number;
  infoCount: number;
  violations: EvaluationViolation[];
}

export class PolicyEngine {
  private policies: PolicyManifest[] = [];

  /**
   * Add a policy manifest directly to the engine.
   */
  addPolicy(policy: PolicyManifest): void {
    this.policies.push(policy);
  }

  /**
   * Recursively load policies from a directory (supporting both YAML and Markdown design policies).
   */
  loadPoliciesFromDirectory(policiesDir: string): number {
    if (!fs.existsSync(policiesDir)) return 0;
    let count = 0;

    const walk = (dir: string) => {
      const entries = fs.readdirSync(dir, { withFileTypes: true });
      for (const entry of entries) {
        const fullPath = path.join(dir, entry.name);
        if (entry.isDirectory()) {
          walk(fullPath);
        } else if (
          entry.name.endsWith('.yaml') ||
          entry.name.endsWith('.yml') ||
          entry.name.endsWith('.md')
        ) {
          try {
            const parsed = parsePolicy(fullPath);
            const policy = 'manifest' in parsed ? (parsed.manifest as PolicyManifest) : (parsed as PolicyManifest);
            this.addPolicy(policy);
            count++;
          } catch {
            // Ignore non-policy files cleanly
          }
        }
      }
    };

    walk(policiesDir);
    return count;
  }

  /**
   * Retrieve all loaded policies.
   */
  getPolicies(): PolicyManifest[] {
    return [...this.policies];
  }

  /**
   * Evaluate a dependency against all dependency policies.
   */
  evaluateDependency(dependencyName: string, version?: string): PolicyEngineSummary {
    const rawViolations = evaluateDependencyRules(this.policies, dependencyName, version);
    return this.summarize(rawViolations);
  }

  /**
   * Evaluate design content (HTML/JSX/TSX/Vue/Svelte/CSS) against design policies.
   */
  evaluateDesignContent(content: string, filePath?: string): PolicyEngineSummary {
    const rawViolations = evaluateDesignRules(this.policies, content, filePath);
    return this.summarize(rawViolations);
  }

  /**
   * Evaluate source code content for security policy violations.
   */
  evaluateSecurity(content: string, filePath?: string): PolicyEngineSummary {
    const rawViolations = evaluateSecurityRules(this.policies, content, filePath);
    return this.summarize(rawViolations);
  }

  /**
   * Evaluate module import path for architectural compliance.
   */
  evaluateArchitecture(filePath: string, importedModule: string): PolicyEngineSummary {
    const rawViolations = evaluateArchitectureRules(this.policies, filePath, importedModule);
    return this.summarize(rawViolations);
  }

  /**
   * Helper to structure raw violations into a PolicyEngineSummary.
   */
  private summarize(violations: EvaluationViolation[]): PolicyEngineSummary {
    const strictBlocks = violations.filter(v => v.enforcement === 'strict-block').length;
    const warnings = violations.filter(v => v.enforcement === 'warning').length;
    const infoCount = violations.filter(v => v.enforcement === 'info').length;

    return {
      allowed: strictBlocks === 0,
      totalViolations: violations.length,
      strictBlocks,
      warnings,
      infoCount,
      violations,
    };
  }
}
