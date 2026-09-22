# Canonical Stack Profiles

AgentJam stack profiles define project technology defaults, linter/formatter configurations, testing tools, and conventions.

---

## Supported Stack Profiles

### 1. `nextjs-fullstack`
- **Framework**: Next.js (React 19, App Router)
- **Language**: TypeScript (Strict Mode)
- **Styling**: Vanilla CSS / CSS Modules (Tailwind optional upon user confirmation)
- **Tooling**: ESLint, Prettier, Vitest, Playwright

### 2. `react-spa`
- **Framework**: React 19 (Vite)
- **Language**: TypeScript
- **Styling**: Vanilla CSS / Design Tokens
- **Tooling**: ESLint, Prettier, Vitest

### 3. `laravel-monolith`
- **Framework**: Laravel (PHP 8.3+)
- **Architecture**: Modular Monolith / DDD
- **Tooling**: PHPStan, Pint, PHPUnit / Pest

### 4. `node-typescript`
- **Ecosystem**: Node.js 20+
- **Language**: TypeScript
- **Tooling**: ESLint, Prettier, Vitest / Jest

### 5. `python-fastapi`
- **Framework**: FastAPI / Python 3.11+
- **Tooling**: Ruff, Black, Mypy, Pytest, UV

### 6. `go-microservice`
- **Framework**: Go 1.22+
- **Tooling**: golangci-lint, gofmt, go test
