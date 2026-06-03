# Changelog

## [1.0.1] - 2026-06-03

### Added

- Regenerate versions.json after tagging to ensure local visibility

## 1.0.0 — 2026-05-15

### Added

- 10-step release workflow: identify, diff, recommend bump,
  update changelog, generate tags, regenerate root changelog,
  commit, tag, confirm
- Auto-generates namespaced tags (`category:`, `language:`,
  `tool:`) from artifact content during release
- Tags presented to user for review before committing
- Bump examples table for semver calibration
- Tag format mirroring directory paths (ADR 0005)
- Support for skills (SKILL.md), rules (RULE.md), and local
  artifacts (.agents/)
- Local artifacts get version bumps and changelogs without
  git tags
