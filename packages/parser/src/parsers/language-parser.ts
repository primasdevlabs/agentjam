/**
 * @agentjam/parser — Language Parser
 */

import fs from 'node:fs';
import path from 'node:path';
import yaml from 'js-yaml';
import { LanguageManifestSchema, type LanguageManifest, type ResourceBundle, ParseError } from '@agentjam/core';

export function parseLanguage(dirPath: string): ResourceBundle<LanguageManifest> {
  const manifestPath = path.join(dirPath, 'language.yaml');

  if (!fs.existsSync(manifestPath)) {
    throw new ParseError(`Language manifest file not found in ${dirPath}`, { resourcePath: dirPath });
  }

  try {
    const raw = yaml.load(fs.readFileSync(manifestPath, 'utf-8'));
    const manifest = LanguageManifestSchema.parse(raw);
    return { manifest, instructions: {}, basePath: dirPath };
  } catch (err: unknown) {
    throw new ParseError(
      `Failed to parse language manifest at ${manifestPath}: ${err instanceof Error ? err.message : String(err)}`,
      { resourcePath: manifestPath }
    );
  }
}
