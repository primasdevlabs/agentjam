/**
 * @agentjam/parser — Workflow Parser
 */

import fs from 'node:fs';
import path from 'node:path';
import yaml from 'js-yaml';
import { WorkflowManifestSchema, type WorkflowManifest, type ResourceBundle, ParseError } from '@agentjam/core';
import { loadInstructionFiles } from '../markdown.js';

/**
 * Parse a workflow directory into a fully-resolved ResourceBundle.
 *
 * Reads `workflow.yaml`, loads instruction files from `instructions/` or `steps/`,
 * and validates that step IDs are unique.
 */
export function parseWorkflow(dirPath: string): ResourceBundle<WorkflowManifest> {
  const manifestPath = path.join(dirPath, 'workflow.yaml');

  if (!fs.existsSync(manifestPath)) {
    throw new ParseError(`workflow.yaml not found in ${dirPath}`, {
      resourcePath: dirPath,
      field: 'workflow.yaml',
    });
  }

  let raw: unknown;
  try {
    const content = fs.readFileSync(manifestPath, 'utf-8');
    raw = yaml.load(content);
  } catch (err: unknown) {
    const msg = err instanceof Error ? err.message : String(err);
    throw new ParseError(`Failed to parse workflow.yaml: ${msg}`, {
      resourcePath: manifestPath,
    });
  }

  const result = WorkflowManifestSchema.safeParse(raw);
  if (!result.success) {
    throw new ParseError(
      `Invalid workflow.yaml in ${dirPath}`,
      { resourcePath: manifestPath },
      result.error.issues.map((i) => ({ path: i.path.join('.'), message: i.message })),
    );
  }
  const manifest = result.data;

  // Validate unique step IDs
  const stepIds = new Set<string>();
  for (const step of manifest.steps) {
    if (stepIds.has(step.id)) {
      throw new ParseError(
        `Duplicate step ID '${step.id}' in workflow '${manifest.name}'`,
        { resourcePath: manifestPath, field: `steps.${step.id}` },
      );
    }
    stepIds.add(step.id);
  }

  // Load instruction files from `instructions/` or `steps/` directory
  let instructions: Record<string, string> = {};
  const instructionsDir = path.join(dirPath, 'instructions');
  const stepsDir = path.join(dirPath, 'steps');

  if (fs.existsSync(instructionsDir)) {
    instructions = { ...instructions, ...loadInstructionFiles(instructionsDir) };
  }
  if (fs.existsSync(stepsDir)) {
    instructions = { ...instructions, ...loadInstructionFiles(stepsDir) };
  }

  return {
    manifest,
    instructions,
    basePath: dirPath,
  };
}
