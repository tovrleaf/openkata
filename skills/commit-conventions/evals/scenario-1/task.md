# Setting Up a Development Branch for a Tracked Bug

## Problem/Feature Description

A development team uses GitHub Issues to track their work. Issue #312 has just been assigned to you — it's a bug where sessions time out unexpectedly after users switch networks mid-session. As part of fixing it, the team also discovered that a module file needs to be renamed: `session.js` should be renamed to `session-manager.js` to reflect the fact that it now handles the full lifecycle of a user session, not just token creation. Keeping the rename in git's history correctly is important to the team because they frequently use `git log --follow` to trace changes.

Your job is to get the development environment ready: set up the right branch and handle the file rename in a way that preserves git history. The team works from `main` and cleans up branches after they're merged.

## Output Specification

Produce a shell script named `branch-setup.sh` that contains all the git commands you would run to accomplish this, in order — from creating the branch through completing the rename. Include comments in the script explaining each step.

Also create a `branch-decision.md` file documenting:
- The exact branch name you chose and why
- How you determined the branch type
- How the file rename was handled and why that approach was used
