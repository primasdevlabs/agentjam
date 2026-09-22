/**
 * @agentjam/parser — Skill Parser
 */

import fs from 'node:fs';
import path from 'node:path';
import yaml from 'js-yaml';
import { SkillManifestSchema, type SkillManifest, type ResourceBundle, ParseError } from '@agentjam/core';
import { loadInstructionFiles } from '../markdown.js';

/**
 * Parse a skill directory into a fully-resolved ResourceBundle.
 *
 * Reads `skill.yaml` and loads all instruction markdown files from `instructions/`.
 */
export function parseSkill(dirPath: string): ResourceBundle<SkillManifest> {
  const manifestPath = path.join(dirPath, 'skill.yaml');

  if (!fs.existsSync(manifestPath)) {
    throw new ParseError(`skill.yaml not found in ${dirPath}`, {
      resourcePath: dirPath,
      field: 'skill.yaml',
    });
  }

  let raw: unknown;
  try {
    const content = fs.readFileSync(manifestPath, 'utf-8');
    raw = yaml.load(content);
  } catch (err: unknown) {
    const msg = err instanceof Error ? err.message : String(err);
    throw new ParseError(`Failed to parse skill.yaml: ${msg}`, {
      resourcePath: manifestPath,
    });
  }

  const result = SkillManifestSchema.safeParse(raw);
  if (!result.success) {
    throw new ParseError(
      `Invalid skill.yaml in ${dirPath}`,
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
