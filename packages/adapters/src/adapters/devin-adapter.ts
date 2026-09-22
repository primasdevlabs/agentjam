/**
 * @agentjam/adapters — Devin Adapter
 *
 * Exports agents and skills to DEVIN.md and .devin/rules/ workspace configurations
 * compatible with Devin CLI and Devin Desktop (ACP).
 */

import { BaseAdapter, type ExportResult } from '../base-adapter.js';
import type { ResourceBundle } from '@agentjam/parser';
import type { AgentManifest, SkillManifest } from '@agentjam/core';

export class DevinAdapter extends BaseAdapter {
  harnessName = 'devin';

  exportAgent(bundle: ResourceBundle<AgentManifest>): ExportResult {
    const { manifest, instructions } = bundle;
    const body = this.concatInstructions(instructions);

    const content = `# ${manifest.name} (v${manifest.version})

> ${manifest.description}

## Context & Role Instructions
${body}
`;

    return {
      harness: this.harnessName,
      files: {
        'DEVIN.md': content,
        [`.devin/rules/${manifest.name}.md`]: content,
      },
    };
  }

  exportSkill(bundle: ResourceBundle<SkillManifest>): ExportResult {
    const { manifest, instructions } = bundle;
    const body = this.concatInstructions(instructions);

    return {
      harness: this.harnessName,
      files: {
        [`.devin/rules/skill-${manifest.name}.md`]: `# Skill: ${manifest.name}\n\n${body}`,
      },
    };
  }
}
