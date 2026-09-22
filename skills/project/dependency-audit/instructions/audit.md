# Dependency Audit Instructions

## Protocol
1. Read lockfiles (`package-lock.json`, `pnpm-lock.yaml`, `composer.lock`, `Cargo.lock`, `poetry.lock`).
2. Verify pinned versions against declared package manifests.
3. Check for security advisories and outdated major versions.
4. Ensure no unvetted or unmaintained third-party dependencies are added without user review.
