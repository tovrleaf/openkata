# Rationale

git-worktree manages parallel git workspaces for
simultaneous multi-branch development.

## Why the skill runs pre-flight checks for nested worktrees

Creating a worktree inside another worktree corrupts
git state silently. Recovery requires manual cleanup.
The pre-flight check prevents an irreversible mistake.

## Why directory names and branch names are separate

Directory is named for the task (what you're doing).
Branch is named by type convention (why it exists).
This separation makes worktree lists readable while
keeping branch conventions consistent.

## Why cleanup requires both worktree remove and branch delete

`git worktree remove` leaves the branch behind.
Orphan branches accumulate and pollute `git branch`
output. Explicit cleanup of both prevents drift.

## Why the skill warns about shared stash

Git stash is global across all worktrees. Stashing
in one worktree then popping in another causes silent
data mixing. Named stash convention makes this
footgun visible.
