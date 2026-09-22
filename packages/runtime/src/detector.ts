/**
 * @agentjam/runtime — Environment Detector
 *
 * Scans a project directory to identify active AI coding environments, IDEs, and CLI tools.
 */

import fs from 'node:fs';
import path from 'node:path';

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
      { name: 'cursorrules', category: 'ide', marker: '.cursorrules' },
      { name: 'vscode', category: 'ide', marker: '.vscode' },
      { name: 'windsurf', category: 'ide', marker: '.windsurf' },
      { name: 'windsurfrules', category: 'ide', marker: '.windsurfrules' },
      { name: 'claude-code', category: 'cli', marker: 'CLAUDE.md' },
      { name: 'codex', category: 'cli', marker: 'codex.config.json' },
      { name: 'gemini', category: 'cli', marker: '.gemini' },
      { name: 'geminimd', category: 'cli', marker: 'GEMINI.md' },
      { name: 'opencode', category: 'cli', marker: 'opencode.json' },
      { name: 'cline', category: 'extensions', marker: '.cline' },
      { name: 'clinerules', category: 'extensions', marker: '.clinerules' },
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
