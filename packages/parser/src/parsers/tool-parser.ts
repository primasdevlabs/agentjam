/**
 * @agentjam/parser — Tool Parser
 */

import fs from 'node:fs';
import path from 'node:path';
import yaml from 'js-yaml';
import { ToolManifestSchema, type ToolManifest, type ResourceBundle, ParseError } from '@agentjam/core';

/**
 * Parse a tool directory into a ResourceBundle.
 * Tools typically have only `tool.yaml` and no instruction files.
 */
export function parseTool(dirPath: string): ResourceBundle<ToolManifest> {
  const manifestPath = path.join(dirPath, 'tool.yaml');

  if (!fs.existsSync(manifestPath)) {
    throw new ParseError(`tool.yaml not found in ${dirPath}`, {
      resourcePath: dirPath,
      field: 'tool.yaml',
    });
  }

  let raw: unknown;
  try {
    const content = fs.readFileSync(manifestPath, 'utf-8');
    raw = yaml.load(content);
  } catch (err: unknown) {
    const msg = err instanceof Error ? err.message : String(err);
    throw new ParseError(`Failed to parse tool.yaml: ${msg}`, {
      resourcePath: manifestPath,
    });
  }

  const result = ToolManifestSchema.safeParse(raw);
  if (!result.success) {
    throw new ParseError(
      `Invalid tool.yaml in ${dirPath}`,
      { resourcePath: manifestPath },
      result.error.issues.map((i) => ({ path: i.path.join('.'), message: i.message })),
    );
  }

  return {
    manifest: result.data,
    instructions: {},
    basePath: dirPath,
  };
}
