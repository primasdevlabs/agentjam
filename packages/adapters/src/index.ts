/**
 * @agentjam/adapters — Entry Point
 *
 * Exports base adapter interfaces, individual harness adapters, and factory lookup.
 */

import { HarnessAdapter } from './base-adapter.js';
import { CursorAdapter } from './adapters/cursor-adapter.js';
import { ClaudeCodeAdapter } from './adapters/claude-code-adapter.js';
import { GeminiAdapter } from './adapters/gemini-adapter.js';
import { ClineAdapter } from './adapters/cline-adapter.js';
import { WindsurfAdapter } from './adapters/windsurf-adapter.js';
import { DevinAdapter } from './adapters/devin-adapter.js';
import { GenericAdapter } from './adapters/generic-adapter.js';

export * from './base-adapter.js';
export * from './adapters/cursor-adapter.js';
export * from './adapters/claude-code-adapter.js';
export * from './adapters/gemini-adapter.js';
export * from './adapters/cline-adapter.js';
export * from './adapters/windsurf-adapter.js';
export * from './adapters/devin-adapter.js';
export * from './adapters/generic-adapter.js';

/**
 * Factory function to instantiate a HarnessAdapter by harness name.
 */
export function getAdapter(harnessName: string): HarnessAdapter {
  switch (harnessName.toLowerCase()) {
    case 'cursor':
      return new CursorAdapter();
    case 'claude':
    case 'claude-code':
      return new ClaudeCodeAdapter();
    case 'gemini':
      return new GeminiAdapter();
    case 'cline':
      return new ClineAdapter();
    case 'windsurf':
      return new WindsurfAdapter();
    case 'devin':
      return new DevinAdapter();
    case 'generic':
    default:
      return new GenericAdapter();
  }
}
