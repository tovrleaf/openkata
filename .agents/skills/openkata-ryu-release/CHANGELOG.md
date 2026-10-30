# Changelog

## [Unreleased]

### Added

- Explicit guidance to ignore plugin.json version during release (ADR 0006)
- Release step to run the evals after the version bump and commit the
  report with the release commit
- Steps for releasing the site, which has a `vMAJOR.MINOR.PATCH` tag of
  its own

## [1.0.0] - 2026-06-11

### Added

- Initial release workflow skill
- Semantic versioning determination rules
- Git-tag-as-source-of-truth principle
- Common failure patterns
