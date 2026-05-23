---
name: create-pr
description: >
  Creates a pull request using GitHub CLI. Checks for
  uncommitted changes, verifies branch state, pushes,
  and opens the PR. Activate when the user says "create
  PR", "open PR", "push and PR", or work is ready for
  review.
metadata:
  tags: "category:workflow, tool:git, tool:github"
---

# Create Pull Request

Open a pull request for the current branch using `gh`.

## Steps

1. **Check for uncommitted changes** — Run `git status
   --short`. If output is non-empty, show the list and
   ask: "There are uncommitted changes. Commit and
   include them, or create the PR without them?" Do not
   proceed without an answer.

2. **Verify branch** — Run `git branch --show-current`.
   If on `main` or `master`, stop and warn: "You are on
   the default branch. Create a feature branch first."

3. **Check remote state** — Run `git status -sb` to
   detect if the branch is behind the remote. If behind,
   warn and ask whether to pull first.

4. **Push the branch** — Run:
   ```bash
   git push -u origin $(git branch --show-current)
   ```
   If push fails, report the error and suggest fixes
   (force push not allowed without explicit permission).

5. **Create the PR** — Run:
   ```bash
   gh pr create --fill
   ```
   Show the generated title and URL to the user. If `gh`
   is not installed, provide the manual URL:
   `https://github.com/{owner}/{repo}/compare/{branch}`

6. **Ask about browser** — "Open the PR in your browser?"
   If yes, run `gh pr view --web`.

## Conventions

- Title comes from the first commit on the branch (via
  `--fill`). If the user wants a custom title, use
  `--title` flag instead.
- Do not force push without explicit user permission.
- Do not push to main/master directly.
- If the branch has a single commit, use its message as
  the PR title. If multiple commits, summarize.
