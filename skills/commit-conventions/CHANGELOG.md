# Changelog

## 1.3.0 — 2026-05-18

### Added

- Branch naming step (step 1) referencing `git-naming` rule
- Guidance for handling "commit all" requests with multiple
  logical changes

### Changed

- Strengthened staging instruction: explicitly prohibit
  `git add .` and `git add -A`, require staging each file
  individually
- Workflow steps renumbered (now 1–8 instead of 1–7)
- Commit format and branch naming now delegate to `git-naming`
  rule with fallback to references
- Removed inline format examples from workflow body (delegated
  to references)

## 1.2.0 — 2026-04-30

### Added

- "Check existing conventions" as step 1 in both workflows
- Example scenario showing end-to-end commit flow
- Common failures section
- references/commit-format.md with full header/body/footer spec
- references/branch-naming.md with format, types, and examples

### Changed

- Restructured as workflow-first skill — format details moved
  to references/
- Description now includes situations alongside actions
- Separate commit and branch workflows

### Removed

- Inline commit format specification (moved to references/)
- Inline branch naming specification (moved to references/)

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
