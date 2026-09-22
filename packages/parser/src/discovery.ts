/**
 * @agentjam/parser — Discovery Utilities
 *
 * File-system walking and manifest discovery used by all parsers,
 * the validator, and the registry builder.
 */

import fs from 'node:fs';
import path from 'node:path';
import { MANIFEST_FILENAMES, type ResourceTypeName } from '@agentjam/core';

// ---------------------------------------------------------------------------
// Directory Walking
// ---------------------------------------------------------------------------

export interface WalkEntry {
  /** Absolute path. */
  absolutePath: string;
  /** Path relative to the walk root. */
  relativePath: string;
  /** Whether this entry is a directory. */
  isDirectory: boolean;
}

/**
 * Recursively walk a directory, optionally filtering by a predicate.
 * Skips `node_modules`, `.git`, and `dist` by default.
 */
export function walkDirectory(
  rootDir: string,
  filter?: (entry: WalkEntry) => boolean,
): WalkEntry[] {
  const results: WalkEntry[] = [];
  const ignoreDirs = new Set(['node_modules', '.git', 'dist', '.agentjam']);

  function walk(currentDir: string) {
    if (!fs.existsSync(currentDir)) return;
    const entries = fs.readdirSync(currentDir, { withFileTypes: true });

    for (const entry of entries) {
      if (ignoreDirs.has(entry.name)) continue;

      const absolutePath = path.join(currentDir, entry.name);
      const relativePath = path.relative(rootDir, absolutePath).replace(/\\/g, '/');
      const walkEntry: WalkEntry = {
        absolutePath,
        relativePath,
        isDirectory: entry.isDirectory(),
      };

      if (!filter || filter(walkEntry)) {
        results.push(walkEntry);
      }

      if (entry.isDirectory()) {
        walk(absolutePath);
      }
    }
  }

  walk(rootDir);
  return results;
}

// ---------------------------------------------------------------------------
// Manifest Discovery
// ---------------------------------------------------------------------------

export interface DiscoveredResource {
  id: string;
  type: ResourceTypeName;
  path: string;
}

/**
 * Find all manifest files of a specific filename within a directory tree.
 *
 * @example findManifestFiles('/project/agents', 'agent.yaml')
 *   → ['/project/agents/software-engineer/agent.yaml', ...]
 */
export function findManifestFiles(dir: string, targetFileName: string): string[] {
  return walkDirectory(dir, (entry) => !entry.isDirectory && entry.absolutePath.endsWith(targetFileName))
    .map((e) => e.absolutePath);
}

/**
 * Discover all resource directories in a repository or for a specific resource type.
 */
export function discoverResources(
  rootDir: string,
  resourceType?: ResourceTypeName
): DiscoveredResource[] {
  const typesToScan: ResourceTypeName[] = resourceType
    ? [resourceType]
    : ['agent', 'skill', 'tool', 'workflow', 'policy', 'stack', 'language', 'integration', 'prompt', 'template'];

  const results: DiscoveredResource[] = [];

  const pluralMap: Record<string, string> = {
    agent: 'agents',
    skill: 'skills',
    tool: 'tools',
    workflow: 'workflows',
    policy: 'policies',
    stack: 'stacks',
    language: 'languages',
    integration: 'integrations',
    prompt: 'prompts',
    template: 'templates',
  };

  for (const type of typesToScan) {
    const plural = pluralMap[type] ?? `${type}s`;
    const searchDir = path.join(rootDir, plural);
    const manifestFilename = MANIFEST_FILENAMES[type];

    if (!manifestFilename || !fs.existsSync(searchDir)) continue;

    const manifestFiles = findManifestFiles(searchDir, manifestFilename);
    for (const file of manifestFiles) {
      const dirPath = path.dirname(file);
      const id = path.basename(dirPath);
      results.push({
        id,
        type,
        path: dirPath,
      });
    }
  }

  return results;
}
