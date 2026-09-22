import { describe, it, expect } from 'vitest';
import path from 'node:path';
import { parseAgent } from '../../packages/parser/src/index.js';

describe('Agent Parsing', () => {
  it('should parse code-reviewer agent successfully', () => {
    const agentDir = path.resolve(process.cwd(), 'agents/code-reviewer');
    const bundle = parseAgent(agentDir);
    expect(bundle.manifest.name).toBe('code-reviewer');
    expect(bundle.manifest.skills).toContain('code-review');
    expect(bundle.instructions['role.md']).toBeDefined();
  });
});
