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

export function buildRegistryIndex(rootDir: string): RegistryIndex {
  const entries: RegistryEntry[] = [];

  // Agents
  const agentsDir = path.join(rootDir, 'agents');
  if (fs.existsSync(agentsDir)) {
    const items = fs.readdirSync(agentsDir, { withFileTypes: true });
    for (const item of items) {
      if (item.isDirectory()) {
        const p = path.join(agentsDir, item.name);
        if (fs.existsSync(path.join(p, 'agent.yaml'))) {
          const { manifest } = parseAgent(p);
          entries.push({
            name: manifest.name,
            version: manifest.version,
            type: 'agent',
            description: manifest.description,
            path: `agents/${item.name}`,
          });
        }
      }
    }
  }

  // Skills
  const skillsDir = path.join(rootDir, 'skills');
  if (fs.existsSync(skillsDir)) {
    const items = fs.readdirSync(skillsDir, { recursive: true, withFileTypes: true });
    for (const item of items) {
      if (item.isDirectory()) {
        const yamlPath = path.join(item.path || skillsDir, item.name, 'skill.yaml');
        if (fs.existsSync(yamlPath)) {
          const p = path.dirname(yamlPath);
          const relPath = path.relative(rootDir, p).replace(/\\/g, '/');
          const { manifest } = parseSkill(p);
          entries.push({
            name: manifest.name,
            version: manifest.version,
            type: 'skill',
            description: manifest.description,
            path: relPath,
          });
        }
      }
    }
  }

  // Tools
  const toolsDir = path.join(rootDir, 'tools');
  if (fs.existsSync(toolsDir)) {
    const items = fs.readdirSync(toolsDir, { recursive: true, withFileTypes: true });
    for (const item of items) {
      if (item.isDirectory()) {
        const yamlPath = path.join(item.path || toolsDir, item.name, 'tool.yaml');
        if (fs.existsSync(yamlPath)) {
          const p = path.dirname(yamlPath);
          const relPath = path.relative(rootDir, p).replace(/\\/g, '/');
          const { manifest } = parseTool(p);
          entries.push({
            name: manifest.name,
            version: manifest.version,
            type: 'tool',
            description: manifest.description,
            path: relPath,
          });
        }
      }
    }
  }

  // Workflows
  const workflowsDir = path.join(rootDir, 'workflows');
  if (fs.existsSync(workflowsDir)) {
    const items = fs.readdirSync(workflowsDir, { recursive: true, withFileTypes: true });
    for (const item of items) {
      if (item.isDirectory()) {
        const yamlPath = path.join(item.path || workflowsDir, item.name, 'workflow.yaml');
        if (fs.existsSync(yamlPath)) {
          const p = path.dirname(yamlPath);
          const relPath = path.relative(rootDir, p).replace(/\\/g, '/');
          const { manifest } = parseWorkflow(p);
          entries.push({
            name: manifest.name,
            version: manifest.version,
            type: 'workflow',
            description: manifest.description,
            path: relPath,
          });
        }
      }
    }
  }

  // Languages
  const languagesDir = path.join(rootDir, 'languages');
  if (fs.existsSync(languagesDir)) {
    const items = fs.readdirSync(languagesDir, { recursive: true, withFileTypes: true });
    for (const item of items) {
      if (item.isDirectory()) {
        const yamlPath = path.join(item.path || languagesDir, item.name, 'language.yaml');
        if (fs.existsSync(yamlPath)) {
          const p = path.dirname(yamlPath);
          const relPath = path.relative(rootDir, p).replace(/\\/g, '/');
          const manifest = parseLanguage(p);
          entries.push({
            name: manifest.name,
            version: '1.0.0',
            type: 'language',
            description: `${manifest.name} language definition (${manifest.ecosystem} ecosystem)`,
            path: relPath,
          });
        }
      }
    }
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
