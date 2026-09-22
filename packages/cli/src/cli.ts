#!/usr/bin/env node

/**
 * @agentjam/cli — Command Line Interface
 *
 * Provides workspace initialization, validation, registry index building,
 * and harness rules export functionality.
 */

import fs from 'node:fs';
import path from 'node:path';
import { parseAgent } from '@agentjam/parser';
import { validateRepository } from '@agentjam/validator';
import { buildAndSaveRegistryIndex } from '@agentjam/registry';
import { getAdapter } from '@agentjam/adapters';

const rootDir = process.cwd();
const args = process.argv.slice(2);
const command = args[0];

function printHelp() {
  console.log(`
AgentJam CLI v0.1.0 — Governance & Harness Rules Management

Usage:
  npx agentjam <command> [options]

Commands:
  init                           Initialize .agentjam/ workspace configuration and export default rules
  validate                       Run repository schema, cross-ref, and instruction validation
  export --agent <id> --harness  Export canonical agent instructions to target AI harness format
  build-registry                 Scan repository and generate root registry.json index
  help                           Show command usage instructions
`);
}

function handleInit() {
  console.log('🚀 Initializing AgentJam workspace configuration...');
  const agentjamDir = path.join(rootDir, '.agentjam');
  const policiesDir = path.join(agentjamDir, 'policies');

  if (!fs.existsSync(agentjamDir)) {
    fs.mkdirSync(agentjamDir, { recursive: true });
  }
  if (!fs.existsSync(policiesDir)) {
    fs.mkdirSync(policiesDir, { recursive: true });
  }

  const configPath = path.join(agentjamDir, 'config.yaml');
  if (!fs.existsSync(configPath)) {
    const defaultConfig = `# AgentJam Project Configuration
versionPolicy: current-stable
freshnessRequired: true
maxDocAge: 7d
`;
    fs.writeFileSync(configPath, defaultConfig, 'utf-8');
    console.log(`   ✅ Created .agentjam/config.yaml`);
  } else {
    console.log(`   ℹ️ .agentjam/config.yaml already exists`);
  }

  // Export default harness rules (AGENTS.md)
  const agentsMdPath = path.join(rootDir, 'AGENTS.md');
  if (!fs.existsSync(agentsMdPath)) {
    const agentsMdContent = `# AgentJam Workspace Guidance (AGENTS.md)

## Operational Rules & Guardrails
1. **Toolchain Preflight**: Run lint, format, typecheck, and unit tests before shipping changes.
2. **Modular Boundaries**: Use public package APIs only. Never import private internal paths.
3. **Security Standards**: Parameterize SQL queries, sanitize inputs, and use secret managers.
4. **Design Governance**: Follow HSL color tokens, use SVG icons (no emojis), and write human copy.
5. **Empirical Verification**: Verify changes with test and build commands before declaring task completion.
`;
    fs.writeFileSync(agentsMdPath, agentsMdContent, 'utf-8');
    console.log(`   ✅ Created root AGENTS.md workspace guidance file`);
  }

  console.log('\n🎉 AgentJam workspace initialized successfully!');
}

function handleValidate() {
  console.log('🔍 Validating AgentJam repository...');
  const result = validateRepository(rootDir);

  if (result.errors.length > 0) {
    for (const err of result.errors) {
      const symbol = err.severity === 'error' ? '❌' : '⚠️';
      console.log(`${symbol} [${err.severity.toUpperCase()}] ${err.resourcePath}`);
      console.log(`   ${err.message}`);
    }
  }

  if (result.valid) {
    console.log(`✅ Repository validation passed successfully! (${result.summary.total} total checks run)`);
  } else {
    console.error(`❌ Validation failed with ${result.summary.errorCount} error(s).`);
    process.exit(1);
  }
}

function handleExport() {
  const agentIdx = args.indexOf('--agent');
  const harnessIdx = args.indexOf('--harness');

  const agentId = agentIdx !== -1 ? args[agentIdx + 1] : 'software-engineer';
  const harnessName = harnessIdx !== -1 ? args[harnessIdx + 1] : 'cursor';

  console.log(`📦 Exporting agent '${agentId}' for harness '${harnessName}'...`);
  const agentDir = path.join(rootDir, 'agents', agentId);

  if (!fs.existsSync(agentDir)) {
    console.error(`❌ Agent directory not found: ${agentDir}`);
    process.exit(1);
  }

  try {
    const bundle = parseAgent(agentDir);
    const adapter = getAdapter(harnessName);
    const exportResult = adapter.exportAgent(bundle);

    for (const [relPath, content] of Object.entries(exportResult.files)) {
      const outPath = path.join(rootDir, relPath);
      const outDir = path.dirname(outPath);
      if (!fs.existsSync(outDir)) {
        fs.mkdirSync(outDir, { recursive: true });
      }
      fs.writeFileSync(outPath, content, 'utf-8');
      console.log(`   ✅ Wrote ${relPath}`);
    }

    console.log(`\n🎉 Successfully exported rules for ${harnessName}!`);
  } catch (err: unknown) {
    const msg = err instanceof Error ? err.message : String(err);
    console.error(`❌ Export failed: ${msg}`);
    process.exit(1);
  }
}

function handleBuildRegistry() {
  console.log('📦 Building registry.json index...');
  const index = buildAndSaveRegistryIndex(rootDir);
  console.log(`✅ Successfully generated registry.json (${index.entries.length} entries indexed)`);
}

switch (command) {
  case 'init':
    handleInit();
    break;
  case 'validate':
    handleValidate();
    break;
  case 'export':
    handleExport();
    break;
  case 'build-registry':
    handleBuildRegistry();
    break;
  case 'help':
  case '-h':
  case '--help':
  default:
    printHelp();
    break;
}
