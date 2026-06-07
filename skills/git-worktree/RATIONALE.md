# Rationale

## Why pre-flight checks for nested worktrees

Creating a worktree inside another worktree corrupts
git state silently. Recovery requires manual cleanup.
The pre-flight check prevents an irreversible mistake.

## Why separate directory name from branch name

Directory is named for the task (what you're doing).
Branch is named by type convention (why it exists).
This separation makes worktree lists readable while
keeping branch conventions consistent.

## Why require both remove and branch delete

`git worktree remove` leaves the branch behind.
Orphan branches accumulate and pollute `git branch`
output. Explicit cleanup of both prevents drift.

## Why the shared stash warning

Git stash is global across all worktrees. Stashing
in one worktree then popping in another causes silent
data mixing. Named stash convention makes this
footgun visible.
