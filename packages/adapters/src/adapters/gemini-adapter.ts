/**
 * @agentjam/adapters — Gemini CLI Adapter
 *
 * Exports agents and skills to GEMINI.md and .gemini/ directory configs.
 */

import { BaseAdapter, type ExportResult } from '../base-adapter.js';
import type { ResourceBundle } from '@agentjam/parser';
import type { AgentManifest, SkillManifest } from '@agentjam/core';

export class GeminiAdapter extends BaseAdapter {
  harnessName = 'gemini';

  exportAgent(bundle: ResourceBundle<AgentManifest>): ExportResult {
    const { manifest, instructions } = bundle;
    const body = this.concatInstructions(instructions);

    const content = `# ${manifest.name}

> ${manifest.description}

${body}
`;

    return {
      harness: this.harnessName,
      files: {
        'GEMINI.md': content,
        [`.gemini/agents/${manifest.name}.md`]: content,
      },
    };
  }

  exportSkill(bundle: ResourceBundle<SkillManifest>): ExportResult {
    const { manifest, instructions } = bundle;
    const body = this.concatInstructions(instructions);

    return {
      harness: this.harnessName,
      files: {
        [`.gemini/skills/${manifest.name}/SKILL.md`]: `---
name: ${manifest.name}
description: ${manifest.description}
---

# ${manifest.name}

${body}`,
      },
    };
  }
}
