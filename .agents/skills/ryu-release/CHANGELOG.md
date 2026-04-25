# Changelog

## 2.0.0 — 2026-04-25

### Changed

- Renamed from `release-kata` to `ryu-release`
- Extended to support both skills (SKILL.md) and rules (RULE.md)
- Artifact resolution now checks skills/, rules/,
  .agents/skills/, and .agents/rules/
- Tagging applies to all distributable directories
  (skills/, rules/), not just skills/

## 1.1.0 — 2026-04-25

### Added

- Local skill support: version bumps and changelogs without git
  tags
- Gotcha: changelogs document skill-facing changes only, not dev
  artifacts
- Version field in frontmatter for consistency with other skills

## 1.0.0 — 2026-04-25

### Added

- 9-step release workflow: identify, diff, recommend bump, update
  version, changelog, commit, tag, confirm
- Bump examples table for semver calibration
- Tag format mirroring directory paths (ADR 0005)
