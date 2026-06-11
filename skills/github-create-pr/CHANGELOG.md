# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.1.0] - 2026-06-11

### Added

- Boundaries section (DOES / Does NOT pattern)
- Proactive rebase offer when branch is behind remote

### Changed

- Switched to `--force-with-lease` push for safety after rebase

## [1.0.0] - 2026-05-24

### Added

- Initial release of create-pr skill
- PR creation workflow using GitHub CLI (`gh`)
- Uncommitted changes check with user prompt
- Branch validation (prevents PR from main/master)
- Remote state check with behind-remote detection
- Pre-push verification step
- Browser open offer after PR creation
