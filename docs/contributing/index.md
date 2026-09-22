# Contributing to AgentJam

Thank you for contributing to AgentJam!

## Guidelines
1. **Canonical Neutrality**: All agents, skills, and policies must remain harness-neutral.
2. **Quality & Validation**: All proposed changes must pass monorepo build, tests, and validator check:
   ```bash
   npm run build
   npm test
   npx tsx scripts/validate.ts
   ```
3. **No Marketing Buzzwords**: Keep docs objective, technical, and accurate.
