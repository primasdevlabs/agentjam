/**
 * @agentjam/runtime — Main AgentJam Runtime
 *
 * Orchestrates environment detection, stack detection, configuration loading,
 * and policy engine evaluation for a workspace directory.
 */

import fs from 'node:fs';
import path from 'node:path';
import yaml from 'js-yaml';
import { PolicyEngine } from '@agentjam/policy-engine';
import { EnvironmentDetector, type DetectedEnvironment } from './detector.js';
import { StackDetector, type DetectedStack } from './stack-detector.js';
import type { AgentJamProjectConfig } from '@agentjam/core';

export class AgentJamRuntime {
  private policyEngine: PolicyEngine;
  private projectRoot: string;

  constructor(projectRoot: string) {
    this.projectRoot = projectRoot;
    this.policyEngine = new PolicyEngine();
    this.policyEngine.loadPoliciesFromDirectory(path.join(projectRoot, 'policies'));
    this.policyEngine.loadPoliciesFromDirectory(path.join(projectRoot, '.agentjam', 'policies'));
  }

  getPolicyEngine(): PolicyEngine {
    return this.policyEngine;
  }

  detectEnvironments(): DetectedEnvironment[] {
    return EnvironmentDetector.detectEnvironments(this.projectRoot);
  }

  detectStack(): DetectedStack | undefined {
    return StackDetector.detectStack(this.projectRoot);
  }

  loadConfig(): AgentJamProjectConfig {
    const configPath = path.join(this.projectRoot, '.agentjam', 'config.yaml');
    if (fs.existsSync(configPath)) {
      try {
        const content = fs.readFileSync(configPath, 'utf-8');
        return yaml.load(content) as AgentJamProjectConfig;
      } catch {
        // Fallback below
      }
    }

    return {
      versionPolicy: 'current-stable',
      freshnessRequired: true,
      maxDocAge: '7d',
    };
  }
}
