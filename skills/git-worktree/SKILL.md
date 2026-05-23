---
name: git-worktree
description: >
  Manages parallel workspaces using git worktrees. Activate
  when the user needs to work on multiple branches
  simultaneously, run long tasks in isolation, or review
  code while keeping current work untouched. Also activate
  when the user says "worktree", "parallel branch", or
  "work on two things at once".
metadata:
  tags: "category:workflow, tool:git"
---

# Git Worktrees

Work on multiple branches simultaneously without stashing
or switching. Each worktree is a full working directory
linked to the same repository.

## When to Use

- Running long tasks (evals, deploys) while continuing
  development on another branch
- Reviewing a PR while keeping current work intact
- Working on a hotfix without disrupting feature work
- Comparing behavior across branches side by side

## When to Skip

- Quick fix on the current branch — just commit in place
- Already inside a worktree — do not nest, continue working
- Single task with no parallel need — regular checkout suffices

## Steps

1. **Check if already in a worktree** — Run:
   ```bash
   git rev-parse --git-dir
   git rev-parse --git-common-dir
   ```
   If they differ, you are already in a linked worktree.
   Do not create another one inside it.

2. **Ensure .worktrees/ is ignored** — Check:
   ```bash
   git check-ignore -q .worktrees
   ```
   If not ignored, add `.worktrees/` to `.gitignore` and
   commit before proceeding. Worktree contents must never
   be tracked.

3. **Create the worktree** — Run:
   ```bash
   git worktree add .worktrees/<branch-name> <branch-name>
   ```
   For a new branch:
   ```bash
   git worktree add .worktrees/<branch-name> -b <branch-name>
   ```
   This creates a full working directory at
   `.worktrees/<branch-name>/` checked out to that branch.

4. **Work in the worktree** — `cd .worktrees/<branch-name>`
   and operate normally. Commits, pushes, and pulls work
   as expected. The worktree shares the object store with
   the main repo — no extra disk for git history.

   Run the project's tests after entering to confirm a clean
   baseline. If tests fail before you change anything, the
   issue is pre-existing — not caused by your work.

5. **List worktrees** — Run:
   ```bash
   git worktree list
   ```
   Shows all linked worktrees and their checked-out
   branches.

6. **Remove when done** — Run:
   ```bash
   git worktree remove .worktrees/<branch-name>
   ```
   Or delete the directory and prune:
   ```bash
   rm -rf .worktrees/<branch-name>
   git worktree prune
   ```

## Conventions

- Store worktrees in `.worktrees/` at the project root
- Name the directory after the branch
- Add `.worktrees/` to `.gitignore` (one-time setup)
- Remove worktrees after merging the branch

## Gotchas

- **Same branch twice** — Git refuses to check out a
  branch that is already checked out in another worktree.
  Use `git checkout --detach` or create a new branch.
- **Shared refs** — All worktrees share the same
  `.git/refs`. A tag or remote update in one is visible
  in all.
- **Shared stash** — `git stash` is global across
  worktrees. Name stashes to avoid confusion:
  `git stash push -m "worktree: description"`
- **Hooks** — Git hooks are shared (they live in
  `.git/hooks`). A pre-commit hook runs in whichever
  worktree triggers the commit.
- **IDE state** — Open each worktree as a separate
  project/window. IDEs that lock files may conflict if
  pointed at the same worktree.
