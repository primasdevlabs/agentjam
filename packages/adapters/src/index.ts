import { ResourceBundle } from '@agentjam/parser';
import { AgentManifest, SkillManifest } from '@agentjam/core';

export interface ExportResult {
  harness: string;
  files: Record<string, string>;
}

export interface HarnessAdapter {
  harnessName: string;
  exportAgent(bundle: ResourceBundle<AgentManifest>): ExportResult;
  exportSkill(bundle: ResourceBundle<SkillManifest>): ExportResult;
}

export class GenericAdapter implements HarnessAdapter {
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
    for (const [filename, text] of Object.entries(instructions)) {
      content += `### ${filename}\n\n${text}\n\n`;
    }

    return {
      harness: 'generic',
      files: {
        [`${manifest.name}.md`]: content,
      },
    };
  }

  exportSkill(bundle: ResourceBundle<SkillManifest>): ExportResult {
    const { manifest, instructions } = bundle;
    let content = `# Skill: ${manifest.name} (v${manifest.version})\n\n`;
    content += `> ${manifest.description}\n\n`;

    for (const [filename, text] of Object.entries(instructions)) {
      content += `### ${filename}\n\n${text}\n\n`;
    }

    return {
      harness: 'generic',
      files: {
        [`skills/${manifest.name}.md`]: content,
      },
    };
  }
}
