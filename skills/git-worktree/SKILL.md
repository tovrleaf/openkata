---
name: git-worktree
description: >
  Creates, lists, and removes git worktrees to manage parallel
  workspaces. Activate when the user needs to work on multiple
  branches simultaneously, run parallel tasks (evals, builds,
  reviews) in isolation, or when an agent needs its own
  workspace. Also activate when the user says "worktree",
  "parallel branch", "work on two things at once", or "run
  these in parallel".
metadata:
  version: "1.0.0"
  tags: "category:workflow, tool:git"
---

# Git Worktrees

Worktrees give each parallel task its own filesystem while
sharing one repository. Every worktree has a branch; every
parallel agent gets its own worktree.

## Steps

**Pre-flight (mandatory before any worktree creation):**

1. Run these two commands and compare output:
   ```bash
   git rev-parse --git-dir
   git rev-parse --git-common-dir
   ```
   If the outputs differ, you are inside a worktree.
   **Stop immediately** and warn: "You are already inside
   a worktree. Switch to the main checkout first."
   Do not proceed.

2. Run:
   ```bash
   git check-ignore -q .worktrees
   ```
   If exit code is non-zero, add `.worktrees/` to
   `.gitignore` before proceeding.

Both checks are mandatory. Never skip them — nested
worktrees corrupt repository state.

1. **Create the worktree** — Run:
   ```bash
   git worktree add .worktrees/<name> -b <branch-name>
   ```
   Convention: name the directory after the task, name the
   branch by type (`eval/`, `release/`, `feature/`):
   ```bash
   git worktree add .worktrees/eval-create-skill -b eval/create-skill
   git worktree add .worktrees/release-commit-conventions -b release/commit-conventions
   ```
   If the branch already exists:
   ```bash
   git worktree add .worktrees/<name> <existing-branch>
   ```

2. **Work in the worktree** — `cd .worktrees/<name>` and
   operate normally. Run the project's tests after entering
   to confirm a clean baseline.

3. **List active worktrees** — After creating worktrees,
   always run:
   ```bash
   git worktree list
   ```
   This confirms the worktrees are set up correctly and
   shows the full set of active workspaces.

4. **Merge results** — When work is complete:
   - Push the branch and create a PR, or
   - From the main checkout: `git merge <branch-name>`

5. **Remove the worktree** — After merging, always run
   both commands:
   ```bash
   git worktree remove .worktrees/<name>
   git branch -d <branch-name>
   ```
   Both are required — `worktree remove` detaches the
   directory, `branch -d` deletes the branch. If the
   directory was deleted manually, run `git worktree prune`
   instead.

## Parallel Execution Pattern

For batch operations across multiple tasks, see
[parallel-execution.md](references/parallel-execution.md).

## Conventions

- Store worktrees in `.worktrees/` at the project root
- Name directories after the task (not the branch)
- Branch naming: `eval/<skill>`, `release/<skill>`,
  `feature/<name>`
- Add `.worktrees/` to `.gitignore` (one-time setup)
- Remove worktrees after merging — they are temporary

## Gotchas

- **Same branch twice** — Git refuses to check out a
  branch in two worktrees simultaneously. Use a different
  branch name or detach HEAD.
- **Shared refs** — All worktrees share `.git/refs`. A
  tag or remote update in one is visible in all.
- **Shared stash** — `git stash` is global. Name stashes:
  `git stash push -m "worktree: description"`
- **Hooks** — Git hooks are shared (`.git/hooks`). A
  pre-commit hook runs in whichever worktree triggers it.
- **IDE state** — Open each worktree as a separate
  project window to avoid conflicts.
- **Skipping pre-flight** — Agents consistently skip
  pre-flight checks when eager to create worktrees. The
  guards exist because nested worktrees corrupt state and
  unignored `.worktrees/` pollutes git status.

## Boundaries

- DOES create, list, and remove worktrees
- DOES establish branch naming conventions for worktrees
- DOES demonstrate parallel execution patterns
- Does NOT manage PRs or merges (use create-pr skill)
- Does NOT handle CI/CD or deployment
- Does NOT replace branch workflow — extends it with
  parallelism
