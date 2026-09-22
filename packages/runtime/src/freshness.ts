/**
 * @agentjam/runtime — Documentation & Policy Freshness Calculator
 *
 * Checks document ages against maximum allowed age policies (e.g., 7d, 30d).
 */

import { parseDuration } from '@agentjam/core';

export interface FreshnessResult {
  isFresh: boolean;
  ageMs: number;
  maxAgeMs: number;
  lastUpdated: string;
}

export function checkFreshness(lastUpdatedIsoDate: string, maxAgeExpression: string): FreshnessResult {
  const lastUpdated = new Date(lastUpdatedIsoDate).getTime();
  const now = Date.now();
  const ageMs = now - lastUpdated;
  const maxAgeMs = parseDuration(maxAgeExpression) ?? (7 * 24 * 60 * 60 * 1000);

  return {
    isFresh: ageMs <= maxAgeMs,
    ageMs,
    maxAgeMs,
    lastUpdated: lastUpdatedIsoDate,
  };
}
