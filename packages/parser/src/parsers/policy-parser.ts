/**
 * @agentjam/parser — Policy Parser
 *
 * Supports both YAML policy manifests and markdown design policies.
 */

import fs from 'node:fs';
import path from 'node:path';
import yaml from 'js-yaml';
import { PolicyManifestSchema, type PolicyManifest, ParseError } from '@agentjam/core';

/**
 * Parse a single policy file.
 *
 * - `.yaml` / `.yml` files are parsed against `PolicyManifestSchema`.
 * - `.md` files are treated as markdown design policies: the filename
 *   becomes the policy ID, and the full content becomes `instructions`.
 */
export function parsePolicy(filePath: string): PolicyManifest {
  if (!fs.existsSync(filePath)) {
    throw new ParseError(`Policy file not found: ${filePath}`, {
      resourcePath: filePath,
    });
  }

  const content = fs.readFileSync(filePath, 'utf-8');

  if (filePath.endsWith('.md')) {
    return parseMarkdownPolicy(filePath, content);
  }

  return parseYamlPolicy(filePath, content);
}

/**
 * Parse a YAML policy manifest.
 */
function parseYamlPolicy(filePath: string, content: string): PolicyManifest {
  let raw: unknown;
  try {
    raw = yaml.load(content);
  } catch (err: unknown) {
    const msg = err instanceof Error ? err.message : String(err);
    throw new ParseError(`Failed to parse policy YAML: ${msg}`, {
      resourcePath: filePath,
    });
  }

  const result = PolicyManifestSchema.safeParse(raw);
  if (!result.success) {
    throw new ParseError(
      `Invalid policy manifest in ${filePath}`,
      { resourcePath: filePath },
      result.error.issues.map((i) => ({ path: i.path.join('.'), message: i.message })),
    );
  }

  return result.data;
}

/**
 * Parse a markdown file as a design policy.
 *
 * Derives the policy ID from the parent directory and filename:
 *   `policies/design/anti-slop.md` → `design-anti-slop`
 *   `policies/core/security.md`    → `core-security`
 */
function parseMarkdownPolicy(filePath: string, content: string): PolicyManifest {
  const filename = path.basename(filePath, '.md');
  const parentDir = path.basename(path.dirname(filePath));

  // Determine category from parent directory name
  const validCategories = [
    'technology', 'versions', 'architecture', 'security',
    'dependencies', 'coding-standards', 'freshness', 'design',
    'project', 'core',
  ];
  const category = validCategories.includes(parentDir) ? parentDir : 'design';

  return {
    id: `${category}-${filename}`,
    name: `${capitalize(category)} Policy: ${capitalize(filename.replace(/-/g, ' '))}`,
    description: `${capitalize(category)} governance rule for ${filename.replace(/-/g, ' ')}`,
    category: category as PolicyManifest['category'],
    enforcement: 'strict-block',
    scope: 'global',
    instructions: content,
  };
}

function capitalize(str: string): string {
  return str.charAt(0).toUpperCase() + str.slice(1);
}

/**
 * Load all policies from a directory tree.
 * Scans for `.yaml`, `.yml`, and `.md` files recursively.
 */
export function loadPoliciesFromDirectory(policiesDir: string): PolicyManifest[] {
  const policies: PolicyManifest[] = [];

  if (!fs.existsSync(policiesDir)) return policies;

  function walk(dir: string) {
    const entries = fs.readdirSync(dir, { withFileTypes: true });
    for (const entry of entries) {
      const fullPath = path.join(dir, entry.name);
      if (entry.isDirectory()) {
        walk(fullPath);
      } else if (
        entry.isFile() &&
        (entry.name.endsWith('.yaml') || entry.name.endsWith('.yml') || entry.name.endsWith('.md'))
      ) {
        try {
          policies.push(parsePolicy(fullPath));
        } catch {
          // Skip files that don't conform to policy structure
        }
      }
    }
  }

  walk(policiesDir);
  return policies;
}
