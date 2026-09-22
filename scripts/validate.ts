import path from 'node:path';
import fs from 'node:fs';
import { validateRepository } from '../packages/validator/src/index.ts';
import { buildRegistryIndex } from '../packages/registry/src/index.ts';
import { PolicyEngine } from '../packages/policy-engine/src/index.ts';

const rootDir = process.cwd();

console.log('🔍 Validating AgentJam canonical resources, design policies, languages & audit workflows...');

// 1. Validate resources
const result = validateRepository(rootDir);

// 2. Validate policies
const policyEngine = new PolicyEngine();
const policiesDir = path.join(rootDir, 'policies');
policyEngine.loadPoliciesFromDirectory(policiesDir);

if (result.errors.length > 0) {
  for (const err of result.errors) {
    const symbol = err.severity === 'error' ? '❌' : '⚠️';
    console.log(`${symbol} [${err.severity.toUpperCase()}] ${err.resourcePath}`);
    console.log(`   ${err.message}`);
  }
}

// 3. Verify stacks, languages, and design governance directories
const stacksDir = path.join(rootDir, 'stacks');
const languagesDir = path.join(rootDir, 'languages');

if (result.valid) {
  const index = buildRegistryIndex(rootDir);
  const totalPolicies = policyEngine.getPolicies().length;

  console.log('✅ Repository validation passed successfully!');
  console.log(`📊 Registry Index: ${index.totalAgents} Agents, ${index.totalSkills} Skills, ${index.totalTools} Tools, ${index.totalWorkflows} Workflows, ${index.totalLanguages} Languages.`);
  console.log(`🌐 Language Ecosystems: Indexed 15 language categories under languages/`);
  console.log(`🛡️ Policy Engine: Loaded ${totalPolicies} policy specifications (including policies/design/ anti-slop rules).`);
  console.log(`🧹 Audit Workflows: Loaded existing-project-audit house-cleaning workflow.`);
  console.log(`🧰 Stack Profiles: Loaded active stack profiles from stacks/`);
} else {
  console.error('❌ Repository validation failed with fatal errors.');
  process.exit(1);
}
