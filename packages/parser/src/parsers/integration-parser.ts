/**
 * @agentjam/parser — Integration Parser
 */

import fs from 'node:fs';
import path from 'node:path';
import yaml from 'js-yaml';
import { IntegrationManifestSchema, type IntegrationManifest, type ResourceBundle, ParseError } from '@agentjam/core';

export function parseIntegration(dirPath: string): ResourceBundle<IntegrationManifest> {
  let manifestPath = path.join(dirPath, 'integration.yaml');
  if (!fs.existsSync(manifestPath)) {
    manifestPath = path.join(dirPath, 'manifest.yaml');
  }

  if (!fs.existsSync(manifestPath)) {
    throw new ParseError(`Integration manifest file not found in ${dirPath}`, { resourcePath: dirPath });
  }

  try {
    const raw = yaml.load(fs.readFileSync(manifestPath, 'utf-8'));
    const manifest = IntegrationManifestSchema.parse(raw);
    return { manifest, instructions: {}, basePath: dirPath };
  } catch (err: unknown) {
    throw new ParseError(
      `Failed to parse integration manifest at ${manifestPath}: ${err instanceof Error ? err.message : String(err)}`,
      { resourcePath: manifestPath }
    );
  }
}
