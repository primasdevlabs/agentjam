// Package agentjam exposes the embedded canonical resource tree so installed
// binaries can self-bootstrap a workspace without a repository checkout.
package agentjam

import "embed"

// CanonicalFS contains the canonical resource tree (agents, skills, tools,
// workflows, policies, stacks, languages, integrations, prompts, templates),
// embedded at build time.
//
//go:embed agents skills tools workflows policies stacks languages integrations prompts templates
var CanonicalFS embed.FS
