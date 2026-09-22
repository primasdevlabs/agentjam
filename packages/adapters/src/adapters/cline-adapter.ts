/**
 * @agentjam/adapters — Cline Adapter
 *
 * Exports agents to .clinerules format.
 */

import { BaseAdapter, type ExportResult } from '../base-adapter.js';
import type { ResourceBundle } from '@agentjam/parser';
import type { AgentManifest, SkillManifest } from '@agentjam/core';

export class ClineAdapter extends BaseAdapter {
  harnessName = 'cline';

  exportAgent(bundle: ResourceBundle<AgentManifest>): ExportResult {
    const { manifest, instructions } = bundle;
    const body = this.concatInstructions(instructions);

    return {
      harness: this.harnessName,
      files: {
        '.clinerules': `# ${manifest.name}\n\n${body}`,
      },
    };
  }

  exportSkill(bundle: ResourceBundle<SkillManifest>): ExportResult {
    const { manifest, instructions } = bundle;
    const body = this.concatInstructions(instructions);

    return {
      harness: this.harnessName,
      files: {
        [`.cline/skills/${manifest.name}.md`]: `# ${manifest.name}\n\n${body}`,
      },
    };
  }
}
