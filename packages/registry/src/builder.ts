/**
 * @agentjam/registry — Registry Index Builder
 *
 * Scans repository resources using @agentjam/parser discovery and constructs a full RegistryIndex.
 */

import fs from 'node:fs';
import path from 'node:path';
import { discoverResources, parseResource } from '@agentjam/parser';
import type { RegistryIndex, RegistryEntry } from './types.js';

export function buildRegistryIndex(rootDir: string): RegistryIndex {
  const discovered = discoverResources(rootDir);
  const entries: RegistryEntry[] = [];

  for (const item of discovered) {
    const relPath = path.relative(rootDir, item.path).replace(/\\/g, '/');

    try {
      const parsed = parseResource(item.path, item.type) as any;
      const manifest = parsed.manifest || parsed;

      entries.push({
        id: manifest.id || item.id,
        name: manifest.name || item.id,
        version: manifest.version || '1.0.0',
        type: item.type,
        description: manifest.description || `${item.type} definition for ${manifest.name || item.id}`,
        path: relPath,
        tags: manifest.tags || manifest.categories || [],
        category: manifest.category || undefined,
        metadata: {
          author: manifest.author,
          license: manifest.license,
        },
      });
    } catch {
      // Fallback entry if parsing fails
      entries.push({
        id: item.id,
        name: item.id,
        version: '1.0.0',
        type: item.type,
        description: `Discovered ${item.type} resource at ${relPath}`,
        path: relPath,
      });
    }
  }

  return {
    generatedAt: new Date().toISOString(),
    totalAgents: entries.filter(e => e.type === 'agent').length,
    totalSkills: entries.filter(e => e.type === 'skill').length,
    totalTools: entries.filter(e => e.type === 'tool').length,
    totalWorkflows: entries.filter(e => e.type === 'workflow').length,
    totalStacks: entries.filter(e => e.type === 'stack').length,
    totalPolicies: entries.filter(e => e.type === 'policy').length,
    totalIntegrations: entries.filter(e => e.type === 'integration').length,
    totalLanguages: entries.filter(e => e.type === 'language').length,
    entries,
  };
}

/**
 * Builds and writes registry index to registry.json file.
 */
export function buildAndSaveRegistryIndex(rootDir: string, outputPath?: string): RegistryIndex {
  const index = buildRegistryIndex(rootDir);
  const outFile = outputPath || path.join(rootDir, 'registry.json');
  fs.writeFileSync(outFile, JSON.stringify(index, null, 2), 'utf-8');
  return index;
}
