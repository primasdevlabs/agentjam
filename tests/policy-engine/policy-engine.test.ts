import { describe, it, expect, beforeEach } from 'vitest';
import { PolicyEngine } from '../../packages/policy-engine/src/index.js';

describe('PolicyEngine', () => {
  let engine: PolicyEngine;

  beforeEach(() => {
    engine = new PolicyEngine();
  });

  it('evaluates dependency policy rules correctly', () => {
    engine.addPolicy({
      id: 'dep-prohibit-legacy',
      name: 'Prohibit Legacy Packages',
      description: 'Prohibits obsolete dependencies',
      category: 'dependencies',
      enforcement: 'strict-block',
      rules: {
        forbidden: ['request', 'moment', 'jquery'],
      },
    });

    const resultAllowed = engine.evaluateDependency('axios');
    expect(resultAllowed.allowed).toBe(true);
    expect(resultAllowed.violations).toHaveLength(0);

    const resultForbidden = engine.evaluateDependency('moment');
    expect(resultForbidden.allowed).toBe(false);
    expect(resultForbidden.violations).toHaveLength(1);
    expect(resultForbidden.violations[0].policyId).toBe('dep-prohibit-legacy');
  });

  it('evaluates design governance rules correctly', () => {
    engine.addPolicy({
      id: 'design-anti-slop',
      name: 'Design Anti-Slop Policy',
      description: 'Prohibits emoji icons and marketing buzzwords',
      category: 'design',
      enforcement: 'strict-block',
    });

    const resultValid = engine.evaluateDesignContent('<div><svg></svg><p>Clean product interface</p></div>');
    expect(resultValid.allowed).toBe(true);

    const resultEmoji = engine.evaluateDesignContent('<div><span>🚀 Feature</span></div>');
    expect(resultEmoji.allowed).toBe(false);

    const resultBuzzword = engine.evaluateDesignContent('<p>Empower your business with AI</p>');
    expect(resultBuzzword.allowed).toBe(false);
  });

  it('evaluates security policy rules correctly', () => {
    engine.addPolicy({
      id: 'security-secrets',
      name: 'Security Secrets Policy',
      description: 'Prohibits hardcoded API keys and secrets',
      category: 'security',
      enforcement: 'strict-block',
    });

    const resultSecret = engine.evaluateSecurity("const api_key = '1234567890abcdef1234';");
    expect(resultSecret.allowed).toBe(false);
  });
});
