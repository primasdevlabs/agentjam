import { defineConfig } from 'vitest/config';
import path from 'node:path';

export default defineConfig({
  test: {
    globals: true,
    environment: 'node',
  },
  resolve: {
    alias: {
      '@agentjam/core': path.resolve(__dirname, 'packages/core/src/index.ts'),
      '@agentjam/parser': path.resolve(__dirname, 'packages/parser/src/index.ts'),
      '@agentjam/validator': path.resolve(__dirname, 'packages/validator/src/index.ts'),
      '@agentjam/registry': path.resolve(__dirname, 'packages/registry/src/index.ts'),
      '@agentjam/adapters': path.resolve(__dirname, 'packages/adapters/src/index.ts'),
      '@agentjam/policy-engine': path.resolve(__dirname, 'packages/policy-engine/src/index.ts'),
      '@agentjam/runtime': path.resolve(__dirname, 'packages/runtime/src/index.ts'),
    },
  },
});
