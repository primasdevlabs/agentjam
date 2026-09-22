/**
 * @agentjam/core — Utility Functions
 *
 * Shared helpers used across AgentJam packages for path resolution,
 * string normalisation, deep merging, and version comparison.
 */

// ---------------------------------------------------------------------------
// String Utilities
// ---------------------------------------------------------------------------

/**
 * Convert a resource name into a URL-safe, filesystem-safe slug.
 *
 * @example slugify('Frontend Engineer') → 'frontend-engineer'
 * @example slugify('code_quality')      → 'code-quality'
 */
export function slugify(input: string): string {
  return input
    .toLowerCase()
    .trim()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '');
}

/**
 * Normalise a resource path to forward slashes and strip trailing separators.
 *
 * @example normalisePath('agents\\software-engineer\\') → 'agents/software-engineer'
 */
export function normalisePath(p: string): string {
  return p.replace(/\\/g, '/').replace(/\/+$/, '');
}

// ---------------------------------------------------------------------------
// Version Utilities
// ---------------------------------------------------------------------------

interface SemVer {
  major: number;
  minor: number;
  patch: number;
  prerelease?: string;
}

/**
 * Parse a SemVer string (with optional leading `^` or `~`) into components.
 * Returns `null` if the string is not a valid SemVer.
 */
export function parseSemVer(version: string): SemVer | null {
  const match = version.match(/^[~^]?(\d+)\.(\d+)\.(\d+)(?:-(.+))?$/);
  if (!match) return null;
  return {
    major: Number(match[1]),
    minor: Number(match[2]),
    patch: Number(match[3]),
    prerelease: match[4],
  };
}

/**
 * Compare two SemVer strings.
 * Returns negative if a < b, zero if equal, positive if a > b.
 * Pre-release versions sort lower than the same version without pre-release.
 */
export function compareSemVer(a: string, b: string): number {
  const pa = parseSemVer(a);
  const pb = parseSemVer(b);
  if (!pa || !pb) return 0;

  const majorDiff = pa.major - pb.major;
  if (majorDiff !== 0) return majorDiff;

  const minorDiff = pa.minor - pb.minor;
  if (minorDiff !== 0) return minorDiff;

  const patchDiff = pa.patch - pb.patch;
  if (patchDiff !== 0) return patchDiff;

  // Both have prerelease → lexicographic; one has, one doesn't → pre < release
  if (pa.prerelease && pb.prerelease) return pa.prerelease.localeCompare(pb.prerelease);
  if (pa.prerelease) return -1;
  if (pb.prerelease) return 1;
  return 0;
}

/**
 * Format a SemVer object back to a string.
 */
export function formatSemVer(ver: SemVer): string {
  const base = `${ver.major}.${ver.minor}.${ver.patch}`;
  return ver.prerelease ? `${base}-${ver.prerelease}` : base;
}

// ---------------------------------------------------------------------------
// Deep Merge
// ---------------------------------------------------------------------------

/**
 * Type guard: value is a non-null plain object (not an array, Date, etc.).
 */
function isPlainObject(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null && !Array.isArray(value) && !(value instanceof Date);
}

/**
 * Deep-merge two plain objects.  Arrays are replaced, not concatenated.
 * Returns a new object without mutating either input.
 *
 * Used for config layering: project config ← stack config ← global defaults.
 */
export function mergeDeep<T extends Record<string, unknown>>(
  base: T,
  override: Partial<T>,
): T {
  const result = { ...base } as Record<string, unknown>;

  for (const key of Object.keys(override)) {
    const baseVal = result[key];
    const overVal = (override as Record<string, unknown>)[key];

    if (isPlainObject(baseVal) && isPlainObject(overVal)) {
      result[key] = mergeDeep(
        baseVal as Record<string, unknown>,
        overVal as Record<string, unknown>,
      );
    } else if (overVal !== undefined) {
      result[key] = overVal;
    }
  }

  return result as T;
}

// ---------------------------------------------------------------------------
// Resource Path Resolution
// ---------------------------------------------------------------------------

/**
 * Given a resource type and name, return the canonical relative path
 * within the AgentJam repository.
 *
 * @example resolveResourcePath('agent', 'software-engineer')
 *   → 'agents/software-engineer'
 * @example resolveResourcePath('skill', 'development/architecture')
 *   → 'skills/development/architecture'
 */
export function resolveResourcePath(resourceType: string, resourceName: string): string {
  const pluralMap: Record<string, string> = {
    agent: 'agents',
    skill: 'skills',
    tool: 'tools',
    workflow: 'workflows',
    policy: 'policies',
    stack: 'stacks',
    language: 'languages',
    integration: 'integrations',
    prompt: 'prompts',
    template: 'templates',
  };

  const plural = pluralMap[resourceType] ?? `${resourceType}s`;
  return normalisePath(`${plural}/${resourceName}`);
}

// ---------------------------------------------------------------------------
// Misc Helpers
// ---------------------------------------------------------------------------

/**
 * Create a unique set of strings from an array, preserving insertion order.
 */
export function uniqueStrings(items: string[]): string[] {
  return [...new Set(items)];
}

/**
 * Truncate a string to a maximum length, appending '…' if truncated.
 */
export function truncate(str: string, maxLength: number): string {
  if (str.length <= maxLength) return str;
  return str.slice(0, maxLength - 1) + '…';
}

/**
 * Parse a duration string like "7d", "24h", "30m" into milliseconds.
 * Supports d(ays), h(ours), m(inutes), s(econds).
 * Returns `null` if the format is unrecognised.
 */
export function parseDuration(duration: string): number | null {
  const match = duration.match(/^(\d+)\s*(d|h|m|s)$/i);
  if (!match) return null;

  const value = Number(match[1]);
  const unit = match[2].toLowerCase();

  const multipliers: Record<string, number> = {
    s: 1_000,
    m: 60_000,
    h: 3_600_000,
    d: 86_400_000,
  };

  return value * (multipliers[unit] ?? 0);
}
