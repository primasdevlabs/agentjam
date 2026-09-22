/**
 * @agentjam/parser — Markdown Utilities
 *
 * Frontmatter extraction, instruction file loading, and
 * markdown content processing.
 */

import fs from 'node:fs';
import path from 'node:path';

// ---------------------------------------------------------------------------
// Frontmatter Extraction
// ---------------------------------------------------------------------------

export interface MarkdownFrontmatter {
  /** Parsed YAML frontmatter key-value pairs. */
  attributes: Record<string, unknown>;
  /** The markdown body after frontmatter has been stripped. */
  body: string;
}

/**
 * Extract YAML frontmatter from a markdown string.
 * Frontmatter is delimited by `---` on its own line at the start of the file.
 *
 * @example
 * extractFrontmatter('---\ntitle: Hello\n---\n# Body')
 *   → { attributes: { title: 'Hello' }, body: '# Body' }
 */
export function extractFrontmatter(content: string): MarkdownFrontmatter {
  const trimmed = content.trimStart();

  if (!trimmed.startsWith('---')) {
    return { attributes: {}, body: content };
  }

  const endIndex = trimmed.indexOf('---', 3);
  if (endIndex === -1) {
    return { attributes: {}, body: content };
  }

  const frontmatterBlock = trimmed.slice(3, endIndex).trim();
  const body = trimmed.slice(endIndex + 3).trimStart();

  // Simple key: value parsing (handles single-level YAML without pulling in js-yaml)
  const attributes: Record<string, unknown> = {};
  for (const line of frontmatterBlock.split('\n')) {
    const colonIndex = line.indexOf(':');
    if (colonIndex === -1) continue;

    const key = line.slice(0, colonIndex).trim();
    const rawValue = line.slice(colonIndex + 1).trim();

    // Attempt to parse booleans and numbers
    if (rawValue === 'true') attributes[key] = true;
    else if (rawValue === 'false') attributes[key] = false;
    else if (/^\d+$/.test(rawValue)) attributes[key] = Number(rawValue);
    else attributes[key] = rawValue.replace(/^["']|["']$/g, '');
  }

  return { attributes, body };
}

// ---------------------------------------------------------------------------
// Instruction File Loading
// ---------------------------------------------------------------------------

export interface InstructionFile {
  /** Filename (e.g. 'role.md'). */
  filename: string;
  /** Parsed frontmatter attributes, if any. */
  attributes: Record<string, unknown>;
  /** Markdown body content. */
  body: string;
  /** Full raw content of the file. */
  raw: string;
}

/**
 * Load and parse a single instruction markdown file.
 */
export function parseInstructionFile(filePath: string): InstructionFile {
  const raw = fs.readFileSync(filePath, 'utf-8');
  const { attributes, body } = extractFrontmatter(raw);
  return {
    filename: path.basename(filePath),
    attributes,
    body,
    raw,
  };
}

/**
 * Load all `.md` files from an `instructions/` directory.
 * Returns a record keyed by filename.
 */
export function loadInstructionFiles(instructionsDir: string): Record<string, string> {
  const instructions: Record<string, string> = {};

  if (!fs.existsSync(instructionsDir)) return instructions;

  const files = fs.readdirSync(instructionsDir);
  for (const file of files) {
    if (file.endsWith('.md')) {
      instructions[file] = fs.readFileSync(path.join(instructionsDir, file), 'utf-8');
    }
  }

  return instructions;
}

/**
 * Load all `.md` files from a directory and parse them as InstructionFile objects.
 */
export function loadParsedInstructions(instructionsDir: string): InstructionFile[] {
  if (!fs.existsSync(instructionsDir)) return [];

  const files = fs.readdirSync(instructionsDir).filter((f) => f.endsWith('.md'));
  return files.map((f) => parseInstructionFile(path.join(instructionsDir, f)));
}

/**
 * Concatenate multiple instruction files into a single markdown string,
 * separated by horizontal rules.  Files are concatenated in alphabetical
 * order unless an explicit `order` frontmatter attribute is set.
 */
export function concatenateInstructions(instructions: Record<string, string>): string {
  const entries = Object.entries(instructions);
  if (entries.length === 0) return '';

  // Sort alphabetically by filename
  entries.sort(([a], [b]) => a.localeCompare(b));

  return entries
    .map(([, content]) => content.trim())
    .join('\n\n---\n\n');
}
