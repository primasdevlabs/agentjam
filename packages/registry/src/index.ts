import fs from 'node:fs';
import path from 'node:path';
import { parseAgent, parseSkill, parseTool, parseWorkflow, parseLanguage } from '@agentjam/parser';

export interface RegistryEntry {
  name: string;
  version: string;
  type: string;
  description: string;
  path: string;
}

export interface RegistryIndex {
  generatedAt: string;
  totalAgents: number;
  totalSkills: number;
  totalTools: number;
  totalWorkflows: number;
  totalLanguages: number;
  entries: RegistryEntry[];
}

function findManifestFiles(dir: string, targetFileName: string): string[] {
  const results: string[] = [];
  if (!fs.existsSync(dir)) return results;

  function walk(currentDir: string) {
    const files = fs.readdirSync(currentDir, { withFileTypes: true });
    for (const file of files) {
      const fullPath = path.join(currentDir, file.name);
      if (file.isDirectory()) {
        walk(fullPath);
      } else if (file.isFile() && file.name === targetFileName) {
        results.push(fullPath);
      }
    }
  }

  walk(dir);
  return results;
}

export function buildRegistryIndex(rootDir: string): RegistryIndex {
  const entries: RegistryEntry[] = [];

  // Agents
  const agentFiles = findManifestFiles(path.join(rootDir, 'agents'), 'agent.yaml');
  for (const file of agentFiles) {
    const dir = path.dirname(file);
    const relPath = path.relative(rootDir, dir).replace(/\\/g, '/');
    const { manifest } = parseAgent(dir);
    entries.push({
      name: manifest.name,
      version: manifest.version,
      type: 'agent',
      description: manifest.description,
      path: relPath,
    });
  }

  // Skills
  const skillFiles = findManifestFiles(path.join(rootDir, 'skills'), 'skill.yaml');
  for (const file of skillFiles) {
    const dir = path.dirname(file);
    const relPath = path.relative(rootDir, dir).replace(/\\/g, '/');
    const { manifest } = parseSkill(dir);
    entries.push({
      name: manifest.name,
      version: manifest.version,
      type: 'skill',
      description: manifest.description,
      path: relPath,
    });
  }

  // Tools
  const toolFiles = findManifestFiles(path.join(rootDir, 'tools'), 'tool.yaml');
  for (const file of toolFiles) {
    const dir = path.dirname(file);
    const relPath = path.relative(rootDir, dir).replace(/\\/g, '/');
    const { manifest } = parseTool(dir);
    entries.push({
      name: manifest.name,
      version: manifest.version,
      type: 'tool',
      description: manifest.description,
      path: relPath,
    });
  }

  // Workflows
  const workflowFiles = findManifestFiles(path.join(rootDir, 'workflows'), 'workflow.yaml');
  for (const file of workflowFiles) {
    const dir = path.dirname(file);
    const relPath = path.relative(rootDir, dir).replace(/\\/g, '/');
    const { manifest } = parseWorkflow(dir);
    entries.push({
      name: manifest.name,
      version: manifest.version,
      type: 'workflow',
      description: manifest.description,
      path: relPath,
    });
  }

  // Languages
  const languageFiles = findManifestFiles(path.join(rootDir, 'languages'), 'language.yaml');
  for (const file of languageFiles) {
    const dir = path.dirname(file);
    const relPath = path.relative(rootDir, dir).replace(/\\/g, '/');
    const manifest = parseLanguage(dir);
    entries.push({
      name: manifest.name,
      version: '1.0.0',
      type: 'language',
      description: `${manifest.name} language definition (${manifest.ecosystem} ecosystem)`,
      path: relPath,
    });
  }

  return {
    generatedAt: new Date().toISOString(),
    totalAgents: entries.filter((e) => e.type === 'agent').length,
    totalSkills: entries.filter((e) => e.type === 'skill').length,
    totalTools: entries.filter((e) => e.type === 'tool').length,
    totalWorkflows: entries.filter((e) => e.type === 'workflow').length,
    totalLanguages: entries.filter((e) => e.type === 'language').length,
    entries,
  };
}
