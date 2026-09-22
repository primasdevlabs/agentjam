# Changelog

All notable changes to AgentJam will be documented in this file.

## [0.3.0] - 2026-09-22

### Added

- **Language Registry & Policy Layer (`languages/`)**: Organized language specifications across 15 ecosystem categories:
  - `systems/` (c, cpp, rust, zig, nim, d)
  - `web/` (javascript, typescript, html, css, scss, less, wasm, webassembly-text)
  - `jvm/` (java, kotlin, scala, groovy, clojure)
  - `dotnet/` (csharp, fsharp, visual-basic)
  - `apple/` (swift, objective-c, objective-cpp)
  - `mobile/` (dart, kotlin, swift, java)
  - `backend/` (php, python, ruby, elixir, erlang, golang, rust, crystal)
  - `data/` (python, r, julia, sql, matlab)
  - `functional/` (haskell, ocaml, elm, clojure, racket, scheme, lisp)
  - `scripting/` (bash, zsh, fish, powershell, perl, lua, tcl)
  - `infrastructure/` (hcl, terraform, cue, dhall, nix)
  - `databases/` (sql, plpgsql, plsql, t-sql, tsql, cypher)
  - `smart-contracts/` (solidity, vyper, move, cairo, clarity)
  - `embedded/` (c, cpp, rust, micropython, arduino)
  - `legacy/` (cobol, fortran, pascal, delphi, ada)
- **Hierarchy Model**: Defined `Language` -> `Framework` -> `Ecosystem` -> `Project Stack` precedence ordering where framework/project conventions override generic language defaults.
- Updated `@agentjam/core`, `@agentjam/parser`, `@agentjam/validator`, and `@agentjam/registry` with `LanguageManifest` support.

## [0.2.0] - 2026-09-22

### Added

- **Design Governance System (`policies/design/`)**: Added 12 non-negotiable design governance policies (`anti-slop.md`, `typography.md`, `color.md`, `icons.md`, `copy.md`, `components.md`, etc.).
- **Existing Project Audit Workflow (`workflows/existing-project-audit/`)**: Mandatory house-cleaning workflow.

## [0.1.0] - 2026-09-22

### Added

- Initial harness-agnostic repository architecture and monorepo TypeScript packages.
