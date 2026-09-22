import { describe, it, expect } from 'vitest';
import { getAdapter, CursorAdapter, ClaudeCodeAdapter, GeminiAdapter, DevinAdapter } from '../../packages/adapters/src/index.js';

describe('Adapters Package', () => {
  const dummyBundle = {
    manifest: {
      name: 'software-engineer',
      version: '1.0.0',
      type: 'agent' as const,
      description: 'Generalist software engineer',
      skills: ['architecture', 'testing'],
      tools: ['filesystem', 'terminal'],
      inputs: ['requirement-spec'],
      outputs: ['source-code'],
    },
    instructions: {
      'role.md': '# Software Engineer\n\nImplement features cleanly.',
    },
    basePath: '/fake/path',
  };

  it('instantiates harness adapters via getAdapter factory', () => {
    const cursor = getAdapter('cursor');
    expect(cursor).toBeInstanceOf(CursorAdapter);

    const claude = getAdapter('claude');
    expect(claude).toBeInstanceOf(ClaudeCodeAdapter);

    const gemini = getAdapter('gemini');
    expect(gemini).toBeInstanceOf(GeminiAdapter);

    const devin = getAdapter('devin');
    expect(devin).toBeInstanceOf(DevinAdapter);
  });

  it('exports agent configuration for Cursor', () => {
    const adapter = new CursorAdapter();
    const result = adapter.exportAgent(dummyBundle);

    expect(result.harness).toBe('cursor');
    expect(result.files['.cursorrules']).toContain('Software Engineer');
    expect(result.files['.cursor/rules/software-engineer.mdc']).toContain('description: Generalist software engineer');
  });

  it('exports agent configuration for Claude Code', () => {
    const adapter = new ClaudeCodeAdapter();
    const result = adapter.exportAgent(dummyBundle);

    expect(result.harness).toBe('claude-code');
    expect(result.files['CLAUDE.md']).toContain('# software-engineer');
  });

  it('exports agent configuration for Gemini', () => {
    const adapter = new GeminiAdapter();
    const result = adapter.exportAgent(dummyBundle);

    expect(result.harness).toBe('gemini');
    expect(result.files['GEMINI.md']).toContain('# software-engineer');
  });

  it('exports agent configuration for Devin', () => {
    const adapter = new DevinAdapter();
    const result = adapter.exportAgent(dummyBundle);

    expect(result.harness).toBe('devin');
    expect(result.files['DEVIN.md']).toContain('# software-engineer');
    expect(result.files['.devin/rules/software-engineer.md']).toBeDefined();
  });
});
