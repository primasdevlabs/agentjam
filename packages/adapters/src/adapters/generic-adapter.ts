/**
 * @agentjam/adapters — Generic Harness Adapter
 *
 * Markdown generator fallback for generic AI coding assistants.
 */

import { BaseAdapter, type ExportResult } from '../base-adapter.js';
import type { ResourceBundle } from '@agentjam/parser';
import type { AgentManifest, SkillManifest } from '@agentjam/core';

export class GenericAdapter extends BaseAdapter {
  harnessName = 'generic';

  exportAgent(bundle: ResourceBundle<AgentManifest>): ExportResult {
    const { manifest, instructions } = bundle;
    let content = `# Agent: ${manifest.name} (v${manifest.version})\n\n`;
    content += `> ${manifest.description}\n\n`;

    if (manifest.skills && manifest.skills.length > 0) {
      content += `## Skills\n` + manifest.skills.map((s: string) => `- ${s}`).join('\n') + '\n\n';
    }

    if (manifest.tools && manifest.tools.length > 0) {
      content += `## Tools\n` + manifest.tools.map((t: string) => `- ${t}`).join('\n') + '\n\n';
    }

    content += `## Instructions\n\n`;
    content += this.concatInstructions(instructions);

    return {
      harness: this.harnessName,
      files: {
        [`${manifest.name}.md`]: content,
      },
    };
  }

  exportSkill(bundle: ResourceBundle<SkillManifest>): ExportResult {
    const { manifest, instructions } = bundle;
    let content = `# Skill: ${manifest.name} (v${manifest.version})\n\n`;
    content += `> ${manifest.description}\n\n`;
    content += this.concatInstructions(instructions);

    return {
      harness: this.harnessName,
      files: {
        [`skills/${manifest.name}.md`]: content,
      },
    };
  }
}
