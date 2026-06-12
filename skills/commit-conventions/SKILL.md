---
name: commit-conventions
description: >
  Enforces Conventional Commits format and branch naming conventions,
  validating commit message structure (type/scope/description header),
  suggesting branch name patterns (feature/*, fix/*, hotfix/*), and
  enforcing breaking change notation. Use when the user is about to
  commit, creating a branch, reviewing commit history, preparing a
  pull request, setting up commit linting, or asking about commit
  message format, branch naming, or Conventional Commits.
metadata:
  version: "1.4.0"
  tags: "category:conventions, category:version-control, tool:git"
---

# Git Conventions

Enforce consistent commit messages and branch naming using
Conventional Commits and kebab-case branch names.

## Commit Workflow

1. **Name branches by type** — Follow the branch naming format
   in the `git-naming` rule. If no rule is installed, see
   [branch-naming.md](references/branch-naming.md).

2. **Check existing conventions** — Before applying defaults,
   look for `.commitlintrc`, `commitlint.config.js`, `.czrc`,
   `CONTRIBUTING.md`, or recent commit history. If the repo has
   established conventions, follow those instead.

3. **Review changes** — Run `git status` and `git diff`.

4. **Verify atomicity** — One logical change per commit. If you
   need "and" in the message, split into multiple commits.
   When asked to "commit all" or "commit everything", first
   review the diff for unrelated changes. If multiple logical
   changes are present, propose a split with one commit per
   concern and confirm the grouping before proceeding.

5. **Stage specific files** — Do not use `git add .` or
   `git add -A`. Stage each file individually to avoid
   committing unrelated changes.

6. **Include generated files** — If the project uses code
   generation (e.g., `templ generate`), ensure generated
   output files are staged alongside their source files.
   Run the generator before staging if unsure.

7. **Write the commit message** — Follow the `git-naming` rule
   for format. If no rule is installed, see
   [commit-format.md](references/commit-format.md).

   Key rules:
   - Imperative mood, lowercase, no trailing period
   - Max 72 characters for the header
   - Body required when the reasoning matters — decisions,
     refactors, non-trivial fixes
   - Header-only is fine when the diff tells the full story
   - Breaking changes always require a body explaining the
     impact and migration path
   - Separate header, body, and footer with blank lines

8. **Commit.**

9. **Validate** — Run `git log -1` to confirm the format. If
   incorrect, run `git commit --amend` to fix.

## Branch Workflow

1. **Check existing conventions** — Look for branch protection
   rules, CI config, or existing branch patterns.

2. **Create the branch** — Follow the `git-naming` rule for
   format. If no rule is installed, see
   [branch-naming.md](references/branch-naming.md).

3. **Create from main, delete after merge.**

## Checkpoint Cleanup

Before pushing, check for WIP or checkpoint commits:

1. Run `git log --oneline origin/main..HEAD`
2. If any commits are checkpoints (prefixed "checkpoint",
   "wip", "temp"):
   - `git reset --soft origin/main` to unstage all changes
   - Review with `git diff --cached`
   - Re-stage selectively and create logical, atomic commits
   - Do NOT squash everything into one commit
3. If history is clean, skip this step

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

## Boundaries

- DOES stage and commit changes
- DOES validate commit message format
- DOES restructure checkpoint commits
- DOES suggest branch names
- Does NOT push to remote
- Does NOT modify source code
- Does NOT create pull/merge requests
