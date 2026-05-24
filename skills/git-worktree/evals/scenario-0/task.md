# Setting Up Parallel Development Workspaces

## Problem/Feature Description

Your team has a small open-source library that two developers need to work on simultaneously — one fixing a critical authentication bug while the other adds a new caching feature. Rather than blocking each other on the same working tree, you want to demonstrate how to set up completely isolated development environments within a single repository so both workstreams can proceed without interference.

You've been asked to set up the repository scaffolding for this parallel workflow from scratch. The repository is fresh and doesn't yet have any worktree-related configuration. Your task is to initialize a git repository (or use one that already exists), then prepare it for parallel branch development by creating two separate isolated workspaces — one for the bug fix and one for the feature — following your team's standard workflow conventions.

## Output Specification

Produce a shell script named `setup_worktrees.sh` that captures the full sequence of commands needed to:

1. Initialize or prepare the repository for parallel worktree development
2. Create two isolated workspaces: one for a bug fix task and one for a caching feature
3. List all worktrees at the end to confirm the setup

The script should be executable and correct bash. Also produce a file `worktree_status.txt` that contains the output you'd expect from listing all worktrees once setup is complete (you can generate this by actually running the commands on a test repo, or write it out manually based on expected output).

Do not include any teardown or merge steps in this task — focus on initial setup only.
