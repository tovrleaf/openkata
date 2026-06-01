---
name: github-create-pr
description: >
  Creates a pull request using GitHub CLI. Checks for
  uncommitted changes, verifies branch state, pushes,
  and opens the PR. Activate when the user says "create
  PR", "open PR", "push and PR", or work is ready for
  review.
metadata:
  version: "1.0.0"
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

4. **Check for existing PR** — Run:
   ```bash
   gh pr list --head $(git branch --show-current)
   ```
   If a PR already exists, show its URL and ask: "A PR
   already exists for this branch. Open it instead?"

5. **Run verification** — Check AGENTS.md or the project's
   build/test configuration for pre-push steps. Run them.
   If any fail, stop and fix before pushing.

6. **Push the branch** — Run:
   ```bash
   git push -u origin $(git branch --show-current)
   ```
   If push fails, report the error and suggest fixes
   (force push not allowed without explicit permission).

7. **Analyze changes** — Run:
   ```bash
   git log main..HEAD --oneline
   git diff main --stat
   ```
   Use this to understand the scope and write an accurate
   PR description. Do not guess — base the description
   on actual changes.

8. **Create the PR** — Run:
   ```bash
   gh pr create --title "<title>" --body "<body>"
   ```
   Show the generated title and description to the user
   for approval before submitting. If `gh` is not
   installed, provide the manual URL:
   `https://github.com/{owner}/{repo}/compare/{branch}`

9. **Ask about browser** — "Open the PR in your browser?"
   If yes, run `gh pr view --web`.

## Conventions

- Title: under 70 characters, plain language summary.
- Description format:
  ```
  Summary paragraph explaining what and why.

  - Concise bullet per logical change
  - No sub-bullets
  - Focus on outcomes, not file paths

  ## Testing (optional, only when changes affect behavior)

  ## Deployment (optional, only when post-merge steps exist)
  ```
- No section headers unless there are multiple distinct
  sections. The default is summary + bullets, flat.
- Do not force push without explicit user permission.
- Do not push to main/master directly.
- If the branch has a single commit, use its message as
  the PR title. If multiple commits, summarize.
