/**
 * @agentjam/parser — Agent Parser
 */

import fs from 'node:fs';
import path from 'node:path';
import yaml from 'js-yaml';
import { AgentManifestSchema, type AgentManifest, type ResourceBundle, ParseError } from '@agentjam/core';
import { loadInstructionFiles } from '../markdown.js';

/**
 * Parse an agent directory into a fully-resolved ResourceBundle.
 *
 * Reads `agent.yaml`, loads all instruction markdown files from `instructions/`,
 * and optionally loads `skills.yaml` and `tools.yaml` sidecar files.
 */
export function parseAgent(dirPath: string): ResourceBundle<AgentManifest> {
  const manifestPath = path.join(dirPath, 'agent.yaml');

  if (!fs.existsSync(manifestPath)) {
    throw new ParseError(`agent.yaml not found in ${dirPath}`, {
      resourcePath: dirPath,
      field: 'agent.yaml',
    });
  }

  let raw: unknown;
  try {
    const content = fs.readFileSync(manifestPath, 'utf-8');
    raw = yaml.load(content);
  } catch (err: unknown) {
    const msg = err instanceof Error ? err.message : String(err);
    throw new ParseError(`Failed to parse agent.yaml: ${msg}`, {
      resourcePath: manifestPath,
    });
  }

  const result = AgentManifestSchema.safeParse(raw);
  if (!result.success) {
    throw new ParseError(
      `Invalid agent.yaml in ${dirPath}`,
      { resourcePath: manifestPath },
      result.error.issues.map((i) => ({ path: i.path.join('.'), message: i.message })),
    );
  }
  const manifest = result.data;

  // Load instruction files
  const instructions = loadInstructionFiles(path.join(dirPath, 'instructions'));

  // Load optional sidecar files (skills.yaml, tools.yaml)
  const sidecars: Record<string, unknown> = {};
  for (const sidecar of ['skills.yaml', 'tools.yaml']) {
    const sidecarPath = path.join(dirPath, sidecar);
    if (fs.existsSync(sidecarPath)) {
      try {
        const content = fs.readFileSync(sidecarPath, 'utf-8');
        sidecars[sidecar] = yaml.load(content);
      } catch {
        // Non-critical: skip malformed sidecars
      }
    }
  }

  return {
    manifest,
    instructions,
    basePath: dirPath,
    metadata: Object.keys(sidecars).length > 0 ? sidecars : undefined,
  };
}
