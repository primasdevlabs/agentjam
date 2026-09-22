/**
 * @agentjam/parser — Stack Parser
 */

import fs from 'node:fs';
import path from 'node:path';
import yaml from 'js-yaml';
import { StackManifestSchema, type StackManifest, type ResourceBundle, ParseError } from '@agentjam/core';
import { loadInstructionFiles } from '../markdown.js';

/**
 * Parse a stack profile directory.
 *
 * Reads `stack.yaml` and optionally loads instruction/convention files
 * from `instructions/`.
 */
export function parseStack(dirPath: string): ResourceBundle<StackManifest> {
  const manifestPath = path.join(dirPath, 'stack.yaml');

  if (!fs.existsSync(manifestPath)) {
    throw new ParseError(`stack.yaml not found in ${dirPath}`, {
      resourcePath: dirPath,
      field: 'stack.yaml',
    });
  }

  let raw: unknown;
  try {
    const content = fs.readFileSync(manifestPath, 'utf-8');
    raw = yaml.load(content);
  } catch (err: unknown) {
    const msg = err instanceof Error ? err.message : String(err);
    throw new ParseError(`Failed to parse stack.yaml: ${msg}`, {
      resourcePath: manifestPath,
    });
  }

  const result = StackManifestSchema.safeParse(raw);
  if (!result.success) {
    throw new ParseError(
      `Invalid stack.yaml in ${dirPath}`,
      { resourcePath: manifestPath },
      result.error.issues.map((i) => ({ path: i.path.join('.'), message: i.message })),
    );
  }

  const instructions = loadInstructionFiles(path.join(dirPath, 'instructions'));

  return {
    manifest: result.data,
    instructions,
    basePath: dirPath,
  };
}
