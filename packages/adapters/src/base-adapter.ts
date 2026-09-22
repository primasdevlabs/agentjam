/**
 * @agentjam/adapters — Base Harness Adapter Interface
 */

import type { ResourceBundle } from '@agentjam/parser';
import type { AgentManifest, SkillManifest, WorkflowManifest } from '@agentjam/core';

export interface ExportResult {
  harness: string;
  files: Record<string, string>;
  warnings?: string[];
}

export interface HarnessAdapter {
  harnessName: string;
  exportAgent(bundle: ResourceBundle<AgentManifest>): ExportResult;
  exportSkill(bundle: ResourceBundle<SkillManifest>): ExportResult;
  exportWorkflow?(bundle: ResourceBundle<WorkflowManifest>): ExportResult;
}

export abstract class BaseAdapter implements HarnessAdapter {
  abstract harnessName: string;

  abstract exportAgent(bundle: ResourceBundle<AgentManifest>): ExportResult;
  abstract exportSkill(bundle: ResourceBundle<SkillManifest>): ExportResult;

  protected concatInstructions(instructions: Record<string, string>): string {
    return Object.entries(instructions)
      .map(([file, text]) => `<!-- Source: ${file} -->\n${text}`)
      .join('\n\n---\n\n');
  }
}
