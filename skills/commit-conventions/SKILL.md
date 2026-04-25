---
name: commit-conventions
version: 1.0.0
description: >
  Enforces Conventional Commits format and branch naming conventions,
  validating commit message structure (type/scope/description header),
  suggesting branch name patterns (feature/*, fix/*, hotfix/*), and
  enforcing breaking change notation. Use when the user asks to commit,
  create a branch, review commit messages, or prepare a pull request.
---

# Git Conventions

This skill enforces consistent commit messages and branch naming across
a project using Conventional Commits and kebab-case branch names.

## Commit Message Format

```
type(scope): description

Body explaining why. Wrap at 72-74 characters.

Footer references.
```

### Header Rules

- **Format:** `type(scope): description`
- **Scope:** Optional. Use the area of the codebase affected.
- **Description:** Imperative mood, lowercase first word, no trailing period
- **Length:** Max 72 characters for the entire header line

### Types

`feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`, `perf`, `ci`, `build`

### Body

Optional for simple changes. Required when:

- The fix isn't immediately obvious
- The change affects behavior in non-obvious ways
- Someone might ask "why did they do it this way?"

Leave a blank line between header and body. Wrap at 72-74 characters.

### Footer

- `Fixes #123` — closes the issue when merged
- `Closes #456` — same as Fixes
- `Relates to #789` — references without closing
- `BREAKING CHANGE: description` — for breaking changes
- `Co-authored-by: Name <email>` — for co-authors

### Examples

Simple commit:

```
feat(dashboard): add loading spinner to dashboard page
```

Bug fix with body and footer:

```
fix(auth): fix race condition in authentication middleware

The middleware was checking token validity before the database
connection was fully established, causing intermittent 401 errors
during server startup.

Adds a connection readiness check before token validation with
retry and exponential backoff (max 3 attempts).

Fixes #142
```

## Branch Naming

### Format

```
type/short-description
type/issue-number-short-description
```

### Types

`feature/`, `fix/`, `refactor/`, `docs/`, `test/`, `chore/`, `hotfix/`

### Rules

- Lowercase, kebab-case (hyphens, not underscores)
- 2-5 words, short but descriptive
- Include issue number when applicable: `fix/123-memory-leak`
- Create from main, delete after merge

### Examples

```
feature/payment-integration
fix/789-session-timeout
refactor/api-error-handling
docs/setup-instructions
chore/upgrade-dependencies
hotfix/security-vulnerability
```

## Commit Workflow

1. Review changes: `git status` and `git diff`
2. Verify changes are atomic — one logical change per commit
3. Stage specific files (prefer `git add <file>` over `git add .`)
4. Write commit message following the format above
5. Commit
6. Validate: run `git log -1` to confirm the message format is correct

## Gotchas

- No "WIP", "temp", or "fix typo" commits — use `git commit --amend`
  or `git rebase -i` to clean up before opening a pull request
- Use `git mv` for file renames to preserve history
- Each commit should represent one logical change — if you need "and"
  in the message, consider splitting into multiple commits
