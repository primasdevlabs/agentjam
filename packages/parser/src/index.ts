import fs from 'node:fs';
import path from 'node:path';
import yaml from 'js-yaml';
import {
  AgentManifestSchema,
  SkillManifestSchema,
  ToolManifestSchema,
  WorkflowManifestSchema,
  IntegrationManifestSchema,
  LanguageManifestSchema,
  type AgentManifest,
  type SkillManifest,
  type ToolManifest,
  type WorkflowManifest,
  type IntegrationManifest,
  type LanguageManifest
} from '@agentjam/core';

export interface ResourceBundle<T> {
  manifest: T;
  instructions: Record<string, string>;
  basePath: string;
}

export function parseYamlFile<T>(filePath: string): T {
  const content = fs.readFileSync(filePath, 'utf-8');
  return yaml.load(content) as T;
}

export function parseAgent(dirPath: string): ResourceBundle<AgentManifest> {
  const manifestPath = path.join(dirPath, 'agent.yaml');
  const raw = parseYamlFile(manifestPath);
  const manifest = AgentManifestSchema.parse(raw);

  const instructions: Record<string, string> = {};
  const instructionsDir = path.join(dirPath, 'instructions');
  if (fs.existsSync(instructionsDir)) {
    const files = fs.readdirSync(instructionsDir);
    for (const file of files) {
      if (file.endsWith('.md')) {
        instructions[file] = fs.readFileSync(path.join(instructionsDir, file), 'utf-8');
      }
    }
  }

  return { manifest, instructions, basePath: dirPath };
}

export function parseSkill(dirPath: string): ResourceBundle<SkillManifest> {
  const manifestPath = path.join(dirPath, 'skill.yaml');
  const raw = parseYamlFile(manifestPath);
  const manifest = SkillManifestSchema.parse(raw);

  const instructions: Record<string, string> = {};
  const instructionsDir = path.join(dirPath, 'instructions');
  if (fs.existsSync(instructionsDir)) {
    const files = fs.readdirSync(instructionsDir);
    for (const file of files) {
      if (file.endsWith('.md')) {
        instructions[file] = fs.readFileSync(path.join(instructionsDir, file), 'utf-8');
      }
    }
  }

  return { manifest, instructions, basePath: dirPath };
}

export function parseTool(dirPath: string): ResourceBundle<ToolManifest> {
  const manifestPath = path.join(dirPath, 'tool.yaml');
  const raw = parseYamlFile(manifestPath);
  const manifest = ToolManifestSchema.parse(raw);
  return { manifest, instructions: {}, basePath: dirPath };
}

export function parseWorkflow(dirPath: string): ResourceBundle<WorkflowManifest> {
  const manifestPath = path.join(dirPath, 'workflow.yaml');
  const raw = parseYamlFile(manifestPath);
  const manifest = WorkflowManifestSchema.parse(raw);

  const instructions: Record<string, string> = {};
  const instructionsDir = path.join(dirPath, 'instructions');
  if (fs.existsSync(instructionsDir)) {
    const files = fs.readdirSync(instructionsDir);
    for (const file of files) {
      if (file.endsWith('.md')) {
        instructions[file] = fs.readFileSync(path.join(instructionsDir, file), 'utf-8');
      }
    }
  }

  return { manifest, instructions, basePath: dirPath };
}

export function parseIntegration(dirPath: string): IntegrationManifest {
  const manifestPath = path.join(dirPath, 'manifest.yaml');
  const raw = parseYamlFile(manifestPath);
  return IntegrationManifestSchema.parse(raw);
}

export function parseLanguage(dirPath: string): LanguageManifest {
  const manifestPath = path.join(dirPath, 'language.yaml');
  const raw = parseYamlFile(manifestPath);
  return LanguageManifestSchema.parse(raw);
}

