# Commit Message Format

## Header

`type(scope): description`

- **Scope:** Optional. Area of the codebase affected.
- **Description:** Imperative mood, lowercase first word, no
  trailing period.
- **Length:** Max 72 characters for the entire header line.

### Types

`feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`,
`perf`, `ci`, `build`

## Body

Optional for simple, self-evident changes (typo fixes, dependency
bumps, single-line config changes). Required when:

- The change involves a decision or trade-off
- The fix isn't immediately obvious from the diff
- The change affects behavior in non-obvious ways
- Someone might ask "why did they do it this way?"

Leave a blank line between header and body. Wrap at 72-74
characters.

## Footer

- `Fixes #123` — closes the issue when merged
- `Closes #456` — same as Fixes
- `Relates to #789` — references without closing
- `BREAKING CHANGE: description` — for breaking changes
- `Co-authored-by: Name <email>` — for co-authors

## Examples

Simple commit:

```text
feat(dashboard): add loading spinner to dashboard page
```

Bug fix with body and footer:

```text
fix(auth): fix race condition in authentication middleware

The middleware was checking token validity before the database
connection was fully established, causing intermittent 401 errors
during server startup.

Adds a connection readiness check before token validation with
retry and exponential backoff (max 3 attempts).

Fixes #142
```
