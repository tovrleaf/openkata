# Changelog

## 1.0.0 — 2026-05-24

### Added

- Initial release of git-worktree skill
- Create, list, and remove git worktrees for parallel workspaces
- Pre-flight checks to prevent nested worktrees and ensure
  `.worktrees/` is gitignored
- Branch naming conventions for worktree branches
- Parallel execution pattern for batch operations
- Gotchas section covering shared refs, stash, hooks, and
  IDE state
