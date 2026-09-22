import fs from 'node:fs';
import path from 'node:path';
import yaml from 'js-yaml';
import { PolicyEngine } from '@agentjam/policy-engine';
import { EnvironmentIntegrationManifest } from '@agentjam/core';

export interface AgentJamProjectConfig {
  stackProfile?: string;
  versionPolicy: 'current-stable' | 'project-compatible' | 'pinned';
  freshnessRequired: boolean;
  maxDocAge: string;
}

export interface DetectedEnvironment {
  name: string;
  category: 'ide' | 'autonomous-agents' | 'cli' | 'extensions' | 'ai-platforms' | 'generic';
  status: 'detected' | 'unknown';
  configMarker?: string;
}

export class EnvironmentDetector {
  static detectEnvironments(projectRoot: string): DetectedEnvironment[] {
    const detected: DetectedEnvironment[] = [];

    const signals: Array<{ name: string; category: DetectedEnvironment['category']; marker: string }> = [
      { name: 'cursor', category: 'ide', marker: '.cursor' },
      { name: 'vscode', category: 'ide', marker: '.vscode' },
      { name: 'windsurf', category: 'ide', marker: '.windsurf' },
      { name: 'claude-code', category: 'cli', marker: 'CLAUDE.md' },
      { name: 'codex', category: 'cli', marker: 'codex.config.json' },
      { name: 'gemini', category: 'cli', marker: '.gemini' },
      { name: 'opencode', category: 'cli', marker: 'opencode.json' },
      { name: 'cline', category: 'extensions', marker: '.cline' },
      { name: 'roo-code', category: 'extensions', marker: '.roo' },
      { name: 'devin', category: 'autonomous-agents', marker: '.devin' },
      { name: 'antigravity', category: 'ai-platforms', marker: '.antigravity' },
    ];

    for (const signal of signals) {
      const p = path.join(projectRoot, signal.marker);
      if (fs.existsSync(p)) {
        detected.push({
          name: signal.name,
          category: signal.category,
          status: 'detected',
          configMarker: signal.marker,
        });
      }
    }

    if (detected.length === 0) {
      detected.push({
        name: 'generic',
        category: 'generic',
        status: 'detected',
        configMarker: 'fallback',
      });
    }

    return detected;
  }
}

export const PrecedenceOrder = [
  'User explicit instruction',
  'Project configuration (.agentjam/config.yaml)',
  'AgentJam project policy (.agentjam/policies/)',
  'AgentJam stack policy (stacks/<profile>/)',
  'AgentJam global policy (policies/)',
  'Environment defaults',
  'Model knowledge (Reasoning fallback only)',
] as const;

export class AgentJamRuntime {
  private policyEngine: PolicyEngine;
  private projectRoot: string;

  constructor(projectRoot: string) {
    this.projectRoot = projectRoot;
    this.policyEngine = new PolicyEngine();
    this.policyEngine.loadPoliciesFromDirectory(path.join(projectRoot, 'policies'));
  }

  getPolicyEngine(): PolicyEngine {
    return this.policyEngine;
  }

  detectEnvironments(): DetectedEnvironment[] {
    return EnvironmentDetector.detectEnvironments(this.projectRoot);
  }

  loadConfig(): AgentJamProjectConfig {
    const configPath = path.join(this.projectRoot, '.agentjam', 'config.yaml');
    if (fs.existsSync(configPath)) {
      const content = fs.readFileSync(configPath, 'utf-8');
      return yaml.load(content) as AgentJamProjectConfig;
    }
    return {
      versionPolicy: 'current-stable',
      freshnessRequired: true,
      maxDocAge: '7d',
    };
  }
}
