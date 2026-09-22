/**
 * @agentjam/adapters — Claude Code Adapter
 *
 * Exports agents and skills to CLAUDE.md and .claude/ workspace configurations.
 */

import { BaseAdapter, type ExportResult } from '../base-adapter.js';
import type { ResourceBundle } from '@agentjam/parser';
import type { AgentManifest, SkillManifest } from '@agentjam/core';

export class ClaudeCodeAdapter extends BaseAdapter {
  harnessName = 'claude-code';

  exportAgent(bundle: ResourceBundle<AgentManifest>): ExportResult {
    const { manifest, instructions } = bundle;
    const body = this.concatInstructions(instructions);

    const content = `# ${manifest.name}

> ${manifest.description}

## Context & Role
${body}
`;

    return {
      harness: this.harnessName,
      files: {
        'CLAUDE.md': content,
        [`.claude/agents/${manifest.name}.md`]: content,
      },
    };
  }

  exportSkill(bundle: ResourceBundle<SkillManifest>): ExportResult {
    const { manifest, instructions } = bundle;
    const body = this.concatInstructions(instructions);

    return {
      harness: this.harnessName,
      files: {
        [`.claude/skills/${manifest.name}.md`]: `# ${manifest.name}\n\n${body}`,
      },
    };
  }
}
