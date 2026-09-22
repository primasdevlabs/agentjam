/**
 * @agentjam/core — Error Classes
 *
 * Structured error hierarchy for all AgentJam packages.
 * Every error carries a machine-readable `code`, optional `resourcePath`,
 * and optional `context` bag for structured logging.
 */

// ---------------------------------------------------------------------------
// Base Error
// ---------------------------------------------------------------------------

export interface AgentJamErrorContext {
  /** The resource path (file or directory) related to this error. */
  resourcePath?: string;
  /** The field or property that caused the error. */
  field?: string;
  /** Additional structured data for logging / debugging. */
  [key: string]: unknown;
}

/**
 * Base error class for all AgentJam errors.
 * Extends `Error` with a machine-readable `code` and optional `context`.
 */
export class AgentJamError extends Error {
  /** Machine-readable error code, e.g. `PARSE_YAML_INVALID`. */
  readonly code: string;

  /** Optional resource path this error relates to. */
  readonly resourcePath?: string;

  /** Structured context for logging / debugging. */
  readonly context: AgentJamErrorContext;

  constructor(message: string, code: string, context: AgentJamErrorContext = {}) {
    super(message);
    this.name = 'AgentJamError';
    this.code = code;
    this.resourcePath = context.resourcePath;
    this.context = context;

    // Maintain proper prototype chain for instanceof checks
    Object.setPrototypeOf(this, new.target.prototype);
  }
}

// ---------------------------------------------------------------------------
// Parse Errors
// ---------------------------------------------------------------------------

/**
 * Thrown when a YAML manifest or markdown instruction file fails to parse
 * or fails Zod schema validation.
 */
export class ParseError extends AgentJamError {
  /** Zod or YAML-level validation issues, if available. */
  readonly issues?: Array<{ path: string; message: string }>;

  constructor(
    message: string,
    context: AgentJamErrorContext = {},
    issues?: Array<{ path: string; message: string }>,
  ) {
    super(message, 'PARSE_ERROR', context);
    this.name = 'ParseError';
    this.issues = issues;
  }
}

// ---------------------------------------------------------------------------
// Policy Violation Errors
// ---------------------------------------------------------------------------

export interface PolicyViolationDetail {
  policyId: string;
  policyName: string;
  category: string;
  message: string;
  severity: 'strict-block' | 'warning' | 'info';
  suggestion?: string;
  affectedFiles?: string[];
}

/**
 * Thrown when a policy evaluation finds one or more violations
 * at the `strict-block` enforcement level.
 */
export class PolicyViolationError extends AgentJamError {
  readonly violations: PolicyViolationDetail[];

  constructor(
    message: string,
    violations: PolicyViolationDetail[],
    context: AgentJamErrorContext = {},
  ) {
    super(message, 'POLICY_VIOLATION', context);
    this.name = 'PolicyViolationError';
    this.violations = violations;
  }
}

// ---------------------------------------------------------------------------
// Validation Errors
// ---------------------------------------------------------------------------

export interface ValidationIssue {
  resourcePath: string;
  field?: string;
  message: string;
  severity: 'error' | 'warning';
  suggestion?: string;
}

/**
 * Thrown when repository structural or semantic validation fails.
 */
export class RepositoryValidationError extends AgentJamError {
  readonly issues: ValidationIssue[];

  constructor(
    message: string,
    issues: ValidationIssue[],
    context: AgentJamErrorContext = {},
  ) {
    super(message, 'VALIDATION_ERROR', context);
    this.name = 'RepositoryValidationError';
    this.issues = issues;
  }
}

// ---------------------------------------------------------------------------
// Adapter Errors
// ---------------------------------------------------------------------------

/**
 * Thrown when exporting canonical resources to a target harness
 * format fails (e.g. incompatible resource, missing capabilities).
 */
export class AdapterError extends AgentJamError {
  /** The target harness that failed. */
  readonly harness: string;

  constructor(message: string, harness: string, context: AgentJamErrorContext = {}) {
    super(message, 'ADAPTER_ERROR', context);
    this.name = 'AdapterError';
    this.harness = harness;
  }
}

// ---------------------------------------------------------------------------
// Configuration Errors
// ---------------------------------------------------------------------------

/**
 * Thrown when `.agentjam/config.yaml` or project configuration is
 * missing, malformed, or contains invalid values.
 */
export class ConfigError extends AgentJamError {
  constructor(message: string, context: AgentJamErrorContext = {}) {
    super(message, 'CONFIG_ERROR', context);
    this.name = 'ConfigError';
  }
}

// ---------------------------------------------------------------------------
// Discovery Errors
// ---------------------------------------------------------------------------

/**
 * Thrown when a required resource, manifest, or directory cannot be found
 * during discovery or resolution.
 */
export class ResourceNotFoundError extends AgentJamError {
  readonly resourceType: string;
  readonly resourceName: string;

  constructor(
    resourceType: string,
    resourceName: string,
    context: AgentJamErrorContext = {},
  ) {
    super(
      `${resourceType} '${resourceName}' not found`,
      'RESOURCE_NOT_FOUND',
      context,
    );
    this.name = 'ResourceNotFoundError';
    this.resourceType = resourceType;
    this.resourceName = resourceName;
  }
}
