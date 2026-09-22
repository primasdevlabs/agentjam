/**
 * @agentjam/registry — Registry Types
 */

import type { ResourceType } from '@agentjam/core';

export interface RegistryEntry {
  id: string;
  name: string;
  version: string;
  type: ResourceType | string;
  description: string;
  path: string;
  tags?: string[];
  category?: string;
  metadata?: Record<string, unknown>;
}

export interface RegistryIndex {
  generatedAt: string;
  totalAgents: number;
  totalSkills: number;
  totalTools: number;
  totalWorkflows: number;
  totalStacks: number;
  totalPolicies: number;
  totalIntegrations: number;
  totalLanguages: number;
  entries: RegistryEntry[];
}

export interface RegistrySearchOptions {
  query?: string;
  type?: ResourceType | string;
  category?: string;
  tag?: string;
  limit?: number;
}
