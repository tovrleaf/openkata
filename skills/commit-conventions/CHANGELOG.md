# Changelog

## 1.1.1 — 2026-04-26

### Fixed

- Applied markdown-consistency rule: added language specifiers to
  code blocks, blank lines around headings in changelog

## 1.1.0 — 2026-04-25

### Added

- Validation step (`git log -1`) and error recovery
  (`git commit --amend`) to workflow

### Changed

- Adopted lowercase convention for commit descriptions
  (Angular/commitlint standard)
- Compacted types table into inline list
- Improved description with specific actions and trigger terms
- Expanded body guidance with clearer criteria for when reasoning
  matters
- Trimmed redundant explanations and removed third example

### Fixed

- CHANGELOG.md formatting (added missing category header)

## 1.0.0 — 2026-04-16

### Added

- Conventional Commits format with optional scope
- Commit types: feat, fix, docs, style, refactor, test, chore,
  perf, ci, build
- Header, body, and footer guidelines with examples
- Branch naming: type/short-description in kebab-case
- Commit workflow for agents
