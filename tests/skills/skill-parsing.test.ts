import { describe, it, expect } from 'vitest';
import path from 'node:path';
import { parseSkill } from '../../packages/parser/src/index.js';

describe('Skill Parsing', () => {
  it('should parse code-review skill successfully', () => {
    const skillDir = path.resolve(process.cwd(), 'skills/code-review');
    const bundle = parseSkill(skillDir);
    expect(bundle.manifest.name).toBe('code-review');
  });
});
