/**
 * @agentjam/parser — Entry Point
 *
 * Exports discovery, markdown utilities, all individual parsers,
 * and unified high-level resource and repository parser functions.
 */

import { ResourceType, ParseError } from '@agentjam/core';
import { parseAgent } from './parsers/agent-parser.js';
import { parseSkill } from './parsers/skill-parser.js';
import { parseTool } from './parsers/tool-parser.js';
import { parseWorkflow } from './parsers/workflow-parser.js';
import { parseStack } from './parsers/stack-parser.js';
import { parsePolicy } from './parsers/policy-parser.js';
import { parseIntegration } from './parsers/integration-parser.js';
import { parseLanguage } from './parsers/language-parser.js';
import { discoverResources } from './discovery.js';

export * from './types.js';
export * from './discovery.js';
export * from './markdown.js';
export * from './parsers/agent-parser.js';
export * from './parsers/skill-parser.js';
export * from './parsers/tool-parser.js';
export * from './parsers/workflow-parser.js';
export * from './parsers/stack-parser.js';
export * from './parsers/policy-parser.js';
export * from './parsers/integration-parser.js';
export * from './parsers/language-parser.js';

/**
 * Dispatcher function to parse any resource given its directory path and resource type.
 */
export function parseResource(dirPath: string, type: ResourceType | string): unknown {
  switch (type) {
    case 'agent':
      return parseAgent(dirPath);
    case 'skill':
      return parseSkill(dirPath);
    case 'tool':
      return parseTool(dirPath);
    case 'workflow':
      return parseWorkflow(dirPath);
    case 'stack':
      return parseStack(dirPath);
    case 'policy':
      return parsePolicy(dirPath);
    case 'integration':
      return parseIntegration(dirPath);
    case 'language':
      return parseLanguage(dirPath);
    default:
      throw new ParseError(`Unknown or unsupported resource type: ${type}`, { resourcePath: dirPath });
  }
}

/**
 * Parse an entire AgentJam repository directory into a collection of resource bundles.
 */
export function parseRepository(rootDir: string): Record<string, unknown[]> {
  const discovered = discoverResources(rootDir);
  const results: Record<string, unknown[]> = {
    agents: [],
    skills: [],
    tools: [],
    workflows: [],
    stacks: [],
    policies: [],
    integrations: [],
    languages: [],
  };

  for (const item of discovered) {
    const parsed = parseResource(item.path, item.type);
    const key = `${item.type}s`;
    if (results[key]) {
      results[key].push(parsed);
    }
  }

  return results;
}
