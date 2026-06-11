# Rationale

github-create-pr is a step-by-step workflow for
creating pull requests using the GitHub CLI.

## Why the skill checks for an existing PR first

Agents may re-trigger PR creation on resume. Duplicate
PRs confuse reviewers and clutter the repo. A simple
`gh pr list` check prevents this.

## Why verification runs before pushing

Pushing broken code creates noise — failed CI, reviewer
time wasted, force-push needed to fix. Running the
project's verify step locally catches issues before
they reach the remote.

## Why the skill shows title and description for approval

PR metadata is public and permanent. Agents hallucinate
descriptions when summarizing changes from memory.
Showing the user what will be submitted prevents
embarrassing or inaccurate PR descriptions.

## Why PR descriptions are built from git diff

Grounding the PR description in actual `git diff` and
`git log` output prevents hallucination. The agent
describes what changed based on evidence, not memory.
