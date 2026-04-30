---
name: commit-conventions
version: 2.0.0
description: >
  Enforces Conventional Commits format and branch naming conventions,
  validating commit message structure (type/scope/description header),
  suggesting branch name patterns (feature/*, fix/*, hotfix/*), and
  enforcing breaking change notation. Use when the user is about to
  commit, creating a branch, reviewing commit history, preparing a
  pull request, setting up commit linting, or asking about commit
  message format, branch naming, or Conventional Commits.
---

# Git Conventions

Enforce consistent commit messages and branch naming using
Conventional Commits and kebab-case branch names.

## Commit Workflow

1. **Check existing conventions** — Before applying defaults,
   look for `.commitlintrc`, `commitlint.config.js`, `.czrc`,
   `CONTRIBUTING.md`, or recent commit history. If the repo has
   established conventions, follow those instead.

2. **Review changes** — Run `git status` and `git diff`.

3. **Verify atomicity** — One logical change per commit. If you
   need "and" in the message, split into multiple commits.

4. **Stage specific files** — Prefer `git add <file>` over
   `git add .`.

5. **Write the commit message** — Follow the
   [commit format](references/commit-format.md):

   ```text
   type(scope): description

   Body explaining why. Wrap at 72-74 characters.

   Footer references.
   ```

   Key rules:
   - Imperative mood, lowercase, no trailing period
   - Max 72 characters for the header
   - Body required when the reasoning matters — decisions,
     refactors, non-trivial fixes
   - Header-only is fine when the diff tells the full story

6. **Commit.**

7. **Validate** — Run `git log -1` to confirm the format. If
   incorrect, run `git commit --amend` to fix.

## Branch Workflow

1. **Check existing conventions** — Look for branch protection
   rules, CI config, or existing branch patterns.

2. **Create the branch** — Follow the
   [branch naming format](references/branch-naming.md):

   ```text
   type/short-description
   type/issue-number-short-description
   ```

   Types: `feature/`, `fix/`, `refactor/`, `docs/`, `test/`,
   `chore/`, `hotfix/`

3. **Create from main, delete after merge.**

## Example Scenario

User: "commit these changes"

1. Skill checks for `.commitlintrc` — none found, uses defaults
2. Runs `git status` and `git diff` to understand the changes
3. Changes touch auth middleware — suggests scope `auth`
4. Writes: `fix(auth): handle expired tokens during reconnect`
5. Stages specific files, commits, validates with `git log -1`

## Common Failures

- **Repo already has conventions** — applying Conventional
  Commits over an existing format creates inconsistent history.
  Always check first (step 1).
- **Scope too broad or missing** — `fix: stuff` tells nothing.
  `fix(auth): resolve token expiry race condition` tells the
  story.
- **Body on trivial commits** — a one-line typo fix doesn't
  need three paragraphs of reasoning.
- **Non-atomic commits** — bundling unrelated changes makes
  history useless. If you need "and", split.

## Gotchas

- No "WIP" or "temp" commits — use `git commit --amend` or
  `git rebase -i` to clean up before opening a pull request
- Use `git mv` for file renames to preserve history
