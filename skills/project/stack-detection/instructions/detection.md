# Stack Detection Instructions

## Overview
Analyzes codebase manifests and directory signatures to accurately identify the active tech stack without reliance on static model defaults.

## Detection Rules
1. **JavaScript/TypeScript**: Inspect `package.json` for `next`, `react`, `vue`, `laravel-vite-plugin`, `@nestjs/core`, `vite`, `express`.
2. **PHP**: Inspect `composer.json` for `laravel/framework`, `symfony/framework-bundle`, `livewire/livewire`.
3. **Python**: Inspect `pyproject.toml`, `requirements.txt`, or `Pipfile` for `django`, `fastapi`, `flask`.
4. **Rust**: Inspect `Cargo.toml` for `actix-web`, `axum`, `tokio`.
5. **Output**: Record detected stack, exact version constraints, runtime requirement, and package manager.
