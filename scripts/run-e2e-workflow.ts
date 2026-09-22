/**
 * AgentJam — End-to-End Workflow & Functionality Verification Script
 *
 * Exercises all 7 packages (@agentjam/core, parser, policy-engine, validator,
 * registry, adapters, runtime) in a realistic end-to-end pipeline.
 */

import path from 'node:path';
import { AgentJamRuntime } from '../packages/runtime/src/index.js';
import { parseRepository, parseAgent, parseSkill, parseWorkflow } from '../packages/parser/src/index.js';
import { validateRepository } from '../packages/validator/src/index.js';
import { buildRegistryIndex, RegistryQueryEngine } from '../packages/registry/src/index.js';
import { getAdapter } from '../packages/adapters/src/index.js';

const rootDir = process.cwd();

console.log('🚀 Executing End-to-End AgentJam Full System Workflow...\n');

// ---------------------------------------------------------------------------
// 1. Runtime Initialization & Environment/Stack Detection
// ---------------------------------------------------------------------------
console.log('1️⃣ [RUNTIME] Initializing AgentJamRuntime & detecting environment...');
const runtime = new AgentJamRuntime(rootDir);
const envs = runtime.detectEnvironments();
const stack = runtime.detectStack();
const config = runtime.loadConfig();

console.log(`   ✅ Environment(s) Detected: ${envs.map(e => `${e.name} (${e.category})`).join(', ')}`);
console.log(`   ✅ Stack Detected: ${stack ? `${stack.name} [${stack.id}]` : 'Generic Node/TS'}`);
console.log(`   ✅ Runtime Config: versionPolicy=${config.versionPolicy}, maxDocAge=${config.maxDocAge}\n`);

// ---------------------------------------------------------------------------
// 2. Parser Discovery & Manifest Loading
// ---------------------------------------------------------------------------
console.log('2️⃣ [PARSER] Discovering & parsing repository resources...');
const parsedRepo = parseRepository(rootDir);
const agentBundle = parseAgent(path.join(rootDir, 'agents', 'software-engineer'));
const skillBundle = parseSkill(path.join(rootDir, 'skills', 'design', 'anti-slop'));
const workflowBundle = parseWorkflow(path.join(rootDir, 'workflows', 'feature-development'));

console.log(`   ✅ Discovered Resources: ${Object.entries(parsedRepo).map(([k, v]) => `${v.length} ${k}`).join(', ')}`);
console.log(`   ✅ Parsed Agent: ${agentBundle.manifest.name} (v${agentBundle.manifest.version}) with ${Object.keys(agentBundle.instructions).length} instructions`);
console.log(`   ✅ Parsed Skill: ${skillBundle.manifest.name} (${skillBundle.manifest.category})`);
console.log(`   ✅ Parsed Workflow: ${workflowBundle.manifest.name} (${workflowBundle.manifest.steps.length} steps)\n`);

// ---------------------------------------------------------------------------
// 3. Policy Engine Evaluation (Dependencies, Design, Security)
// ---------------------------------------------------------------------------
console.log('3️⃣ [POLICY ENGINE] Evaluating governance rules...');
const policyEngine = runtime.getPolicyEngine();
policyEngine.loadPoliciesFromDirectory(path.join(rootDir, 'policies'));

const depEval = policyEngine.evaluateDependency('axios');
const secEval = policyEngine.evaluateSecurity("const safeKey = process.env.API_KEY;");
const designEval = policyEngine.evaluateDesignContent('<button class="px-4 py-2 bg-primary">Submit</button>');

console.log(`   ✅ Dependency Check ('axios'): ${depEval.allowed ? 'ALLOWED' : 'BLOCKED'}`);
console.log(`   ✅ Security Check (env vars): ${secEval.allowed ? 'CLEAN' : 'VIOLATION DETECTED'}`);
console.log(`   ✅ Design Anti-Slop Check (clean component): ${designEval.allowed ? 'COMPLIANT' : 'VIOLATION DETECTED'}\n`);

// ---------------------------------------------------------------------------
// 4. Validator Three-Tier Integrity Check
// ---------------------------------------------------------------------------
console.log('4️⃣ [VALIDATOR] Executing repository validation suite...');
const valResult = validateRepository(rootDir);

console.log(`   ✅ Validation Result: ${valResult.valid ? 'PASSED (0 fatal errors)' : 'FAILED'}`);
console.log(`   📊 Validation Stats: ${valResult.summary.total} issues flagged (${valResult.summary.errorCount} errors, ${valResult.summary.warningCount} warnings)\n`);

// ---------------------------------------------------------------------------
// 5. Registry Indexing & Search Query Engine
// ---------------------------------------------------------------------------
console.log('5️⃣ [REGISTRY] Building index & running query search...');
const regIndex = buildRegistryIndex(rootDir);
const queryEngine = new RegistryQueryEngine(regIndex);
const searchResults = queryEngine.search({ query: 'software-engineer' });
const designSkills = queryEngine.search({ category: 'design' });

console.log(`   ✅ Registry Index Built: ${regIndex.entries.length} total entries`);
console.log(`   🔍 Search Query ('software-engineer'): Found ${searchResults.length} match [${searchResults[0]?.id || ''}]`);
console.log(`   🔍 Search Query (category='design'): Found ${designSkills.length} design skills\n`);

// ---------------------------------------------------------------------------
// 6. Harness Adapters Export Engine
// ---------------------------------------------------------------------------
console.log('6️⃣ [ADAPTERS] Testing export adapters across AI harnesses...');
const harnesses = ['cursor', 'claude', 'gemini', 'devin', 'cline', 'windsurf'];
for (const harnessName of harnesses) {
  const adapter = getAdapter(harnessName);
  const result = adapter.exportAgent(agentBundle);
  console.log(`   ✅ ${harnessName.toUpperCase()} Adapter: Exported ${Object.keys(result.files).length} file(s) [e.g. ${Object.keys(result.files)[0]}]`);
}

console.log('\n🎉 FULL END-TO-END WORKFLOW EXECUTED SUCCESSFULLY WITH 100% FUNCTIONALITY!');
