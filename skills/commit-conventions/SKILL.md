---
name: commit-conventions
version: 1.0.0
description: >
  Enforces Conventional Commits format and branch naming conventions.
  Activate when the user asks to commit, create a branch, or review
  commit messages. Also activate when preparing a pull request.
---

# Git Conventions

This skill enforces consistent commit messages and branch naming across
a project using Conventional Commits and kebab-case branch names.

## Commit Message Format

```
type(scope): description

Body explaining why, not what. Wrap at 72-74 characters.

Footer references.
```

### Header Rules

- **Format:** `type(scope): description`
- **Scope:** Optional. Use the area of the codebase affected.
- **Description:** Imperative mood ("Add feature" not "Added feature"),
  capitalize the first word, no period at the end
- **Length:** Max 72 characters for the entire header line

### Types

| Type | When to use |
|------|-------------|
| `feat` | New feature or capability |
| `fix` | Bug fix |
| `docs` | Documentation only |
| `style` | Formatting, no code change |
| `refactor` | Code change that neither fixes a bug nor adds a feature |
| `test` | Adding or updating tests |
| `chore` | Maintenance, dependencies, CI |
| `perf` | Performance improvement |
| `ci` | CI/CD changes |
| `build` | Build system or external dependencies |

### Body

Optional for simple, self-explanatory changes. Required when:

- The fix isn't immediately obvious
- The change affects behavior in non-obvious ways
- Someone might ask "why did they do it this way?"

**Include:** Why the change is needed, how the solution works, side
effects, alternatives considered.

**Do not:** Describe what changed (the diff shows that), be vague
("fixed some stuff"), or apologize.

Leave a blank line between header and body. Wrap at 72-74 characters.

### Footer

- `Fixes #123` — closes the issue when merged
- `Closes #456` — same as Fixes
- `Relates to #789` — references without closing
- `BREAKING CHANGE: description` — for breaking changes
- `Co-authored-by: Name <email>` — for co-authors

### Examples

Simple commit (no body needed):

```
feat(dashboard): Add loading spinner to dashboard page
```

Bug fix with explanation:

```
fix(auth): Fix race condition in authentication middleware

The middleware was checking token validity before the database
connection was fully established, causing intermittent 401 errors
during server startup.

Adds a connection readiness check before token validation with
retry and exponential backoff (max 3 attempts).

Fixes #142
```

Feature with detailed body:

```
feat(export): Add user profile export functionality

Users can now export their profile data in JSON or CSV format
from account settings. This addresses GDPR data portability
requirements.

Large exports are processed asynchronously and delivered via
email with a secure download link valid for 7 days.

Closes #234
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

## Gotchas

- No "WIP", "temp", or "fix typo" commits — use `git commit --amend`
  or `git rebase -i` to clean up before opening a pull request
- Use `git mv` for file renames to preserve history
- Each commit should represent one logical change — if you need "and"
  in the message, consider splitting into multiple commits
