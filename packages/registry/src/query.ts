/**
 * @agentjam/registry — Query Engine
 *
 * Provides search, filtering, and lookup utilities over a RegistryIndex.
 */

import type { RegistryIndex, RegistryEntry, RegistrySearchOptions } from './types.js';

export class RegistryQueryEngine {
  constructor(private index: RegistryIndex) {}

  /**
   * Search entries using query string, resource type, category, or tags.
   */
  search(options: RegistrySearchOptions): RegistryEntry[] {
    let results = [...this.index.entries];

    if (options.type) {
      results = results.filter(e => e.type.toLowerCase() === options.type!.toLowerCase());
    }

    if (options.category) {
      results = results.filter(e => e.category?.toLowerCase() === options.category!.toLowerCase());
    }

    if (options.tag) {
      const targetTag = options.tag.toLowerCase();
      results = results.filter(e => e.tags?.some(t => t.toLowerCase() === targetTag));
    }

    if (options.query) {
      const q = options.query.toLowerCase();
      results = results.filter(
        e =>
          e.name.toLowerCase().includes(q) ||
          e.id.toLowerCase().includes(q) ||
          e.description.toLowerCase().includes(q) ||
          e.tags?.some(t => t.toLowerCase().includes(q))
      );
    }

    if (options.limit && options.limit > 0) {
      results = results.slice(0, options.limit);
    }

    return results;
  }

  /**
   * Find entry by exact ID.
   */
  findById(id: string): RegistryEntry | undefined {
    return this.index.entries.find(e => e.id.toLowerCase() === id.toLowerCase());
  }

  /**
   * Get all entries of a specific resource type.
   */
  getByType(type: string): RegistryEntry[] {
    return this.search({ type });
  }
}
