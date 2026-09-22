import { describe, it, expect } from 'vitest';
import { EnvironmentDetector, StackDetector, checkFreshness, PrecedenceResolver } from '../../packages/runtime/src/index.js';

describe('Runtime Package', () => {
  it('detects environments in current workspace', () => {
    const environments = EnvironmentDetector.detectEnvironments(process.cwd());
    expect(environments.length).toBeGreaterThan(0);
  });

  it('detects stack in current workspace', () => {
    const stack = StackDetector.detectStack(process.cwd());
    expect(stack).toBeDefined();
    expect(stack?.id).toBe('node-typescript');
  });

  it('calculates freshness correctly', () => {
    const nowIso = new Date().toISOString();
    const resultFresh = checkFreshness(nowIso, '7d');
    expect(resultFresh.isFresh).toBe(true);

    const oldIso = new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString();
    const resultStale = checkFreshness(oldIso, '7d');
    expect(resultStale.isFresh).toBe(false);
  });

  it('resolves precedence order correctly', () => {
    const candidates = [
      { level: 'Environment defaults' as const, sourceName: 'defaults', value: 'v1' },
      { level: 'User explicit instruction' as const, sourceName: 'user', value: 'v2' },
      { level: 'Project configuration (.agentjam/config.yaml)' as const, sourceName: 'config', value: 'v3' },
    ];

    const highest = PrecedenceResolver.resolveHighest(candidates);
    expect(highest?.value).toBe('v2');
  });
});
