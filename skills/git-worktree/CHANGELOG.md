# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.1] - 2026-06-11

### Added

- RATIONALE.md explaining design decisions

## [1.0.0] - 2026-05-24

### Added

- Initial release of git-worktree skill
- Create, list, and remove git worktrees for parallel workspaces
- Pre-flight checks to prevent nested worktrees and ensure
  `.worktrees/` is gitignored
- Branch naming conventions for worktree branches
- Parallel execution pattern for batch operations
- Gotchas section covering shared refs, stash, hooks, and
  IDE state
