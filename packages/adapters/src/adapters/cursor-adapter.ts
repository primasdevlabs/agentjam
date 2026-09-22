/**
 * @agentjam/adapters — Cursor Harness Adapter
 *
 * Exports agents and skills to .cursorrules and .cursor/rules/*.mdc format.
 */

import { BaseAdapter, type ExportResult } from '../base-adapter.js';
import type { ResourceBundle } from '@agentjam/parser';
import type { AgentManifest, SkillManifest } from '@agentjam/core';

export class CursorAdapter extends BaseAdapter {
  harnessName = 'cursor';

  exportAgent(bundle: ResourceBundle<AgentManifest>): ExportResult {
    const { manifest, instructions } = bundle;
    const body = this.concatInstructions(instructions);

    const cursorRulesContent = `# ${manifest.name} (v${manifest.version})
# ${manifest.description}

${body}
`;

    return {
      harness: this.harnessName,
      files: {
        '.cursorrules': cursorRulesContent,
        [`.cursor/rules/${manifest.name}.mdc`]: `---
description: ${manifest.description}
globs: *
---

${body}`,
      },
    };
  }

  exportSkill(bundle: ResourceBundle<SkillManifest>): ExportResult {
    const { manifest, instructions } = bundle;
    const body = this.concatInstructions(instructions);

    return {
      harness: this.harnessName,
      files: {
        [`.cursor/rules/skill-${manifest.name}.mdc`]: `---
description: ${manifest.description}
globs: *
---

# Skill: ${manifest.name}

${body}`,
      },
    };
  }
}
