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

export function validateRepository(rootDir: string): ValidationResult {
  const errors: ValidationError[] = [];

  // 1. Validate agents
  const agentsDir = path.join(rootDir, 'agents');
  if (fs.existsSync(agentsDir)) {
    const agents = fs.readdirSync(agentsDir, { withFileTypes: true });
    for (const dirent of agents) {
      if (dirent.isDirectory()) {
        const agentPath = path.join(agentsDir, dirent.name);
        const yamlPath = path.join(agentPath, 'agent.yaml');
        if (fs.existsSync(yamlPath)) {
          try {
            const bundle = parseAgent(agentPath);
            if (!bundle.instructions['role.md']) {
              errors.push({
                resourcePath: agentPath,
                field: 'instructions/role.md',
                message: `Agent '${dirent.name}' is missing role.md in instructions/`,
                severity: 'warning',
              });
            }
          } catch (err: any) {
            errors.push({
              resourcePath: agentPath,
              message: `Invalid agent.yaml: ${err.message}`,
              severity: 'error',
            });
          }
        }
      }
    }
  }

  // 2. Validate skills
  const skillsDir = path.join(rootDir, 'skills');
  if (fs.existsSync(skillsDir)) {
    const skills = fs.readdirSync(skillsDir, { recursive: true, withFileTypes: true });
    for (const dirent of skills) {
      if (dirent.isDirectory()) {
        const yamlPath = path.join(dirent.path || skillsDir, dirent.name, 'skill.yaml');
        if (fs.existsSync(yamlPath)) {
          const skillPath = path.dirname(yamlPath);
          try {
            parseSkill(skillPath);
          } catch (err: any) {
            errors.push({
              resourcePath: skillPath,
              message: `Invalid skill.yaml: ${err.message}`,
              severity: 'error',
            });
          }
        }
      }
    }
  }

  // 3. Validate tools
  const toolsDir = path.join(rootDir, 'tools');
  if (fs.existsSync(toolsDir)) {
    const tools = fs.readdirSync(toolsDir, { recursive: true, withFileTypes: true });
    for (const dirent of tools) {
      if (dirent.isDirectory()) {
        const yamlPath = path.join(dirent.path || toolsDir, dirent.name, 'tool.yaml');
        if (fs.existsSync(yamlPath)) {
          const toolPath = path.dirname(yamlPath);
          try {
            parseTool(toolPath);
          } catch (err: any) {
            errors.push({
              resourcePath: toolPath,
              message: `Invalid tool.yaml: ${err.message}`,
              severity: 'error',
            });
          }
        }
      }
    }
  }

  // 4. Validate workflows
  const workflowsDir = path.join(rootDir, 'workflows');
  if (fs.existsSync(workflowsDir)) {
    const workflows = fs.readdirSync(workflowsDir, { recursive: true, withFileTypes: true });
    for (const dirent of workflows) {
      if (dirent.isDirectory()) {
        const yamlPath = path.join(dirent.path || workflowsDir, dirent.name, 'workflow.yaml');
        if (fs.existsSync(yamlPath)) {
          const workflowPath = path.dirname(yamlPath);
          try {
            parseWorkflow(workflowPath);
          } catch (err: any) {
            errors.push({
              resourcePath: workflowPath,
              message: `Invalid workflow.yaml: ${err.message}`,
              severity: 'error',
            });
          }
        }
      }
    }
  }

  // 5. Validate languages
  const languagesDir = path.join(rootDir, 'languages');
  if (fs.existsSync(languagesDir)) {
    const langs = fs.readdirSync(languagesDir, { recursive: true, withFileTypes: true });
    for (const dirent of langs) {
      if (dirent.isDirectory()) {
        const yamlPath = path.join(dirent.path || languagesDir, dirent.name, 'language.yaml');
        if (fs.existsSync(yamlPath)) {
          const langPath = path.dirname(yamlPath);
          try {
            parseLanguage(langPath);
          } catch (err: any) {
            errors.push({
              resourcePath: langPath,
              message: `Invalid language.yaml: ${err.message}`,
              severity: 'error',
            });
          }
        }
      }
    }
  }

  const hasFatalErrors = errors.some((e) => e.severity === 'error');
  return {
    valid: !hasFatalErrors,
    errors,
  };
}
