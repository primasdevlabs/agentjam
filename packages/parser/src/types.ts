/**
 * @agentjam/parser — Types
 */

export interface ResourceBundle<T = unknown> {
  manifest: T;
  instructions: Record<string, string>;
  basePath: string;
  metadata?: Record<string, unknown>;
}
