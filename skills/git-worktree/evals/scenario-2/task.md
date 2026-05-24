# Worktree Lifecycle Management Script

## Problem/Feature Description

Your team has been using git worktrees informally for a few months, but has run into several painful incidents: stash entries from one context showing up unexpectedly in another, lingering stale worktree references after developers deleted directories manually instead of using the proper git command, and confusion when two people tried to use the same branch name in different worktrees.

You've been asked to write a comprehensive worktree management script that demonstrates correct lifecycle handling — from creation through teardown — with safeguards against these known failure modes. The script will serve as both a working tool and a reference for new team members onboarding to the worktree-based workflow.

The scenario: a developer needs to work on a release task and a hotfix simultaneously. They may have uncommitted changes they need to stash while switching between contexts. After finishing, they should clean up properly — even if they accidentally deleted one of the worktree directories by hand rather than using git commands.

## Output Specification

Produce a shell script named `worktree_lifecycle.sh` that demonstrates the full lifecycle:

1. Setting up two worktrees for a release task and a hotfix
2. Simulating stashing work-in-progress changes with descriptive stash names
3. Simulating completion and teardown — including handling the case where a worktree directory has been manually deleted rather than removed with git
4. Final verification that no stale worktree references remain

The script should be runnable bash and include inline comments explaining each key decision. Also produce a brief `GOTCHAS.md` file documenting at least three specific gotchas a developer should know about when working with git worktrees in a team environment.
