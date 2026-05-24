---
name: git-worktree
description: >
  Manages parallel workspaces using git worktrees. Activate
  when the user needs to work on multiple branches
  simultaneously, run parallel tasks (evals, builds, reviews)
  in isolation, or when an agent needs its own workspace.
  Also activate when the user says "worktree", "parallel
  branch", "work on two things at once", or "run these in
  parallel".
metadata:
  tags: "category:workflow, tool:git"
---

# Git Worktrees

Worktrees give each parallel task its own filesystem while
sharing one repository. Every worktree has a branch; every
parallel agent gets its own worktree.

## Mental Model

- **Branch** = the logical unit of work (a pointer to commits)
- **Worktree** = the physical workspace (a directory checked
  out to a branch)
- You cannot have a worktree without a branch
- You can have a branch without a worktree (just checkout)
- Multiple worktrees = multiple branches active simultaneously

Use worktrees whenever work should happen in parallel without
blocking the current workspace.

## When to Use

- Running long tasks (evals, deploys) while continuing
  development on another branch
- Multiple agents working on different tasks simultaneously
- Reviewing a PR while keeping current work intact
- Working on a hotfix without disrupting feature work
- Batch operations across multiple skills or modules

## When to Skip

- Quick fix on the current branch — just commit in place
- Already inside a worktree — do not nest, continue working
- Single sequential task — regular checkout suffices

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

4. **Work in the worktree** — `cd .worktrees/<name>` and
   operate normally. Commits, pushes, and pulls work as
   expected. The worktree shares the object store — no
   extra disk for git history.

   Run the project's tests after entering to confirm a
   clean baseline. If tests fail before you change
   anything, the issue is pre-existing.

5. **Merge results** — When work is complete:
   - Push the branch and create a PR, or
   - From the main checkout: `git merge <branch-name>`

6. **Remove the worktree** — After merging:
   ```bash
   git worktree remove .worktrees/<name>
   git branch -d <branch-name>
   ```
   Or if the directory was deleted manually:
   ```bash
   git worktree prune
   ```

## Parallel Execution Pattern

For batch operations (e.g., running evals for all skills):

```bash
# Create one worktree per task
git worktree add .worktrees/eval-skill-a -b eval/skill-a
git worktree add .worktrees/eval-skill-b -b eval/skill-b
git worktree add .worktrees/eval-skill-c -b eval/skill-c

# Run tasks in parallel (each in its own directory)
cd .worktrees/eval-skill-a && tessl eval run skills/skill-a/ &
cd .worktrees/eval-skill-b && tessl eval run skills/skill-b/ &
cd .worktrees/eval-skill-c && tessl eval run skills/skill-c/ &
wait

# Merge results back
git merge eval/skill-a eval/skill-b eval/skill-c

# Cleanup
git worktree remove .worktrees/eval-skill-a
git worktree remove .worktrees/eval-skill-b
git worktree remove .worktrees/eval-skill-c
git branch -d eval/skill-a eval/skill-b eval/skill-c
```

## Listing and Status

```bash
git worktree list          # Show all worktrees and branches
git worktree list --porcelain  # Machine-readable output
```

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
- **Shared stash** — `git stash` is global. Name stashes
  to avoid confusion:
  `git stash push -m "worktree: description"`
- **Hooks** — Git hooks are shared (`.git/hooks`). A
  pre-commit hook runs in whichever worktree triggers it.
- **IDE state** — Open each worktree as a separate
  project window. IDEs may conflict if pointed at the
  same directory.

## Boundaries

- DOES create, list, and remove worktrees
- DOES establish branch naming conventions for worktrees
- DOES demonstrate parallel execution patterns
- Does NOT manage PRs or merges (use create-pr skill)
- Does NOT handle CI/CD or deployment
- Does NOT replace branch workflow — extends it with
  parallelism
