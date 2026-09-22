import fs from 'node:fs';
import path from 'node:path';
import { parseAgent, parseSkill, parseTool, parseWorkflow, parseLanguage } from '@agentjam/parser';

export interface ValidationError {
  resourcePath: string;
  field?: string;
  message: string;
  severity: 'error' | 'warning';
}

export interface ValidationResult {
  valid: boolean;
  errors: ValidationError[];
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

export function validateRepository(rootDir: string): ValidationResult {
  const errors: ValidationError[] = [];

  // 1. Validate agents
  const agentFiles = findManifestFiles(path.join(rootDir, 'agents'), 'agent.yaml');
  for (const file of agentFiles) {
    const agentPath = path.dirname(file);
    try {
      const bundle = parseAgent(agentPath);
      if (!bundle.instructions['role.md']) {
        errors.push({
          resourcePath: agentPath,
          field: 'instructions/role.md',
          message: `Agent '${bundle.manifest.name}' is missing role.md in instructions/`,
          severity: 'warning',
        });
      }
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      errors.push({
        resourcePath: agentPath,
        message: `Invalid agent.yaml: ${msg}`,
        severity: 'error',
      });
    }
  }

  // 2. Validate skills
  const skillFiles = findManifestFiles(path.join(rootDir, 'skills'), 'skill.yaml');
  for (const file of skillFiles) {
    const skillPath = path.dirname(file);
    try {
      parseSkill(skillPath);
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      errors.push({
        resourcePath: skillPath,
        message: `Invalid skill.yaml: ${msg}`,
        severity: 'error',
      });
    }
  }

  // 3. Validate tools
  const toolFiles = findManifestFiles(path.join(rootDir, 'tools'), 'tool.yaml');
  for (const file of toolFiles) {
    const toolPath = path.dirname(file);
    try {
      parseTool(toolPath);
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      errors.push({
        resourcePath: toolPath,
        message: `Invalid tool.yaml: ${msg}`,
        severity: 'error',
      });
    }
  }

  // 4. Validate workflows
  const workflowFiles = findManifestFiles(path.join(rootDir, 'workflows'), 'workflow.yaml');
  for (const file of workflowFiles) {
    const workflowPath = path.dirname(file);
    try {
      parseWorkflow(workflowPath);
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      errors.push({
        resourcePath: workflowPath,
        message: `Invalid workflow.yaml: ${msg}`,
        severity: 'error',
      });
    }
  }

  // 5. Validate languages
  const languageFiles = findManifestFiles(path.join(rootDir, 'languages'), 'language.yaml');
  for (const file of languageFiles) {
    const langPath = path.dirname(file);
    try {
      parseLanguage(langPath);
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : String(err);
      errors.push({
        resourcePath: langPath,
        message: `Invalid language.yaml: ${msg}`,
        severity: 'error',
      });
    }
  }

  const hasFatalErrors = errors.some((e) => e.severity === 'error');
  return {
    valid: !hasFatalErrors,
    errors,
  };
}
