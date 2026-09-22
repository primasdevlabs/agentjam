import path from 'node:path';
import { buildAndSaveRegistryIndex } from '../packages/registry/src/index.js';

const rootDir = process.cwd();
console.log('📦 Building canonical AgentJam registry.json index...');

const index = buildAndSaveRegistryIndex(rootDir);
console.log(`✅ Successfully generated registry.json at ${path.join(rootDir, 'registry.json')}`);
console.log(`📊 Total entries indexed: ${index.entries.length}`);
