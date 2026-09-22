/**
 * @agentjam/runtime — Precedence Order Pipeline
 *
 * Defines the strict 7-level precedence hierarchy for policy and instruction resolution.
 */

import { PRECEDENCE_ORDER, type PrecedenceLevel } from '@agentjam/core';

export interface PrecedenceItem<T> {
  level: PrecedenceLevel;
  sourceName: string;
  value: T;
}

export class PrecedenceResolver {
  /**
   * Resolve highest-precedence item from candidates according to canonical PRECEDENCE_ORDER.
   */
  static resolveHighest<T>(candidates: PrecedenceItem<T>[]): PrecedenceItem<T> | undefined {
    if (candidates.length === 0) return undefined;

    return [...candidates].sort((a, b) => {
      const idxA = PRECEDENCE_ORDER.indexOf(a.level);
      const idxB = PRECEDENCE_ORDER.indexOf(b.level);
      return (idxA === -1 ? 999 : idxA) - (idxB === -1 ? 999 : idxB);
    })[0];
  }
}
