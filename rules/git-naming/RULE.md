---
name: git-naming
description: >
  Branch naming, commit message format, and trailer conventions
  for git operations.
metadata:
  tags: "category:conventions, tool:git"
---

# Git Naming

Naming conventions for branches, commits, and trailers.

## Branch Names

Format: `type/short-description`

Types: `feature/`, `fix/`, `refactor/`, `docs/`, `test/`,
`chore/`, `hotfix/`

- Lowercase kebab-case, 2–5 words
- Include issue number when applicable:
  `fix/123-memory-leak`
- Create from main, delete after merge

## Commit Messages

Conventional Commits format:

```text
type(scope): description

Body explaining why. Wrap at 72-74 characters.

Footer references.
```

Header rules:
- Max 72 characters
- Imperative mood, lowercase, no trailing period
- Scope is the area of the codebase affected

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`,
`chore`, `perf`, `ci`, `build`, `release`

`release` is for the commit that releases an artifact or the
site. The scope is the artifact type, or `site`, and the
description is the name and version, or just the version for
the site: `release(skills): grill-me v1.0.1`,
`release(site): v0.15.0`.

Body is optional when the diff is under 20 lines and touches
one file. Required when the change involves a decision,
trade-off, or non-obvious behavior.

## Assisted-by Trailer

Add to every AI-assisted commit:

```text
Assisted-by: <agent>:<model>
```

Agent is your name (Kiro, Claude Code, OpenCode). Model is
the model you are running on.

Discovery hints:
- Kiro: `ls -t ~/.kiro/sessions/cli/*.json | head -1 \
    | xargs grep -o '"model_name": "[^"]*"' | cut -d'"' -f4`
- Claude Code: read `ANTHROPIC_MODEL` env var
- Others: use the model name from your session context

## Merging

Never squash-merge a PR. Merge with a merge commit or
rebase so each commit keeps its own message, spec references,
and trailers in the main history.
