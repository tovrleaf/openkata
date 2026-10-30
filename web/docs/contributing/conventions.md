---
title: Repository conventions
description: How branches, commits, pull requests, code, tests, scripts, markdown, changelogs and styles are written in this repository.
order: 5
---

# Repository conventions

These are the conventions that changes to this repository follow. Most
of them are written down as rules in `rules/`, the same rules that are
published in the catalog, and `AGENTS.md` summarizes them for agents.
Reviewers expect them, and agents are told to follow them.

## Branches

A branch name is `type/short-description`.

- The types are `feature/`, `fix/`, `refactor/`, `docs/`, `test/`,
  `chore/` and `hotfix/`.
- The description is lowercase kebab-case, two to five words.
- Include the issue number when there is one, as in
  `fix/123-memory-leak`.
- Create the branch from `main`, and delete it after it is merged.

## Commits

Commit messages follow Conventional Commits:

```text
type(scope): description

Body explaining why. Wrap at 72-74 characters.

Footer references.
```

- **The header** is at most 72 characters, in the imperative mood,
  lowercase, with no period at the end. The scope is the area of the
  codebase that changed.
- **The types** are `feat`, `fix`, `docs`, `style`, `refactor`, `test`,
  `chore`, `perf`, `ci`, `build` and `release`. `release` is for the
  commit that releases an artifact, as in
  `release(skills): grill-me v1.0.1`.
- **The body** is optional when the diff is under 20 lines and touches
  one file. It is required when the change involves a decision, a
  trade-off, or behavior that is not obvious. If someone would ask
  "why?", write the body.

### The Assisted-by trailer

Every commit that an AI agent helped with ends with a trailer naming
the agent and the model:

```text
Assisted-by: Claude Code:claude-sonnet-5-5
```

The agent is the tool you used, such as Kiro, Claude Code or OpenCode,
and the model is the one it ran on. The trailer makes AI contributions
visible in the history, so they can be filtered and reviewed.

## Pull requests and merging

Pull requests are never squash-merged. Merge with a merge commit, or
rebase. Squashing replaces many described commits with one, and loses
the reasoning of each step, its spec reference and its trailers.

Agents are asked to follow three rules about git:

- Do not push without explicit confirmation.
- Do not force push without first explaining what it will do.
- Committing is fine when asked, but pushing is a separate action.

[Releasing a new version of the site](release-site.md) lists the checks
to run before you push.

## Go

- Format with `go fmt`, and build with zero warnings.
- Always build to `bin/`, as in
  `go build -o bin/openkata-web ./cmd/openkata-web/`. A binary at the
  project root is a sign of a build error.

### Tests

- Use the standard library only. No testify or other frameworks.
- Run `go test ./...`, and `go test -race ./...` for races.
- Keep a test next to the code it tests, so `handlers.go` has
  `handlers_test.go`. Fixtures go in `testdata/`.
- Test HTTP handlers with `net/http/httptest`.
- Write table-driven tests with named subtests.
- A failure message names the function, the inputs, what it got and
  what it wanted. A message such as "unexpected result" is not
  allowed.
- Compare errors with `errors.Is` and `errors.As`, never by their
  text.
- A test helper calls `t.Helper()` first.
- Test behavior, not implementation. Skip trivial logic, and do not
  aim for 100% coverage.

## Bash

- Start with `#!/usr/bin/env bash` and `set -euo pipefail`, then a
  `#` line saying what the script does.
- Indent with 2 spaces, and keep lines to 80 characters.
- Quote variables with braces, as in `"${var}"`.
- Use `[[ ]]` and `(( ))`, not `[ ]`, `test` or `let`.
- Name functions and variables in `lower_snake_case`, and constants in
  `UPPER_SNAKE_CASE` declared with `readonly`.
- Use `local` for function variables. A script with functions has a
  `main` function, called last with `main "$@"`.
- Check return values, use arrays for argument lists, and never use
  `eval`. Errors go to stderr.
- Keep scripts under 100 lines. Shell is for small utilities.

## Markdown

These apply to every markdown file, including these docs.

- Start with one `#` heading, use ATX headings only, and never skip a
  level.
- Wrap prose at 72 to 80 characters. Links, tables, headings and code
  blocks may be longer.
- Use `-` for unordered lists, and indent nested lists by 2 spaces.
- Use fenced code blocks, always with a language. Use `text` when none
  applies.
- Use descriptive link text, and relative links within the repository.
- Keep table cells under 60 characters, and use a list when the data
  has no uniform structure.
- One blank line around headings, lists, code blocks and tables. No
  trailing spaces, and a single newline at the end of the file.

## Changelogs

Every changelog follows Keep a Changelog, with an `## [Unreleased]`
section at the top and versions below it, newest first, as described in
[versioning and the manifest](../reference/versioning.md). Entries say
what changed in a sentence. "Fixed bug" is not enough, and "Fixed crash
when uploading files larger than 10MB" is. A changelog does not mention
the tessl files.

## Templates and CSS

The website's styles follow the design system rule.

- Use `var(--token)` for every color, spacing and font. Never hard-code
  them.
- Use the semantic color tokens, such as `--color-ink` and
  `--color-surface`, for anything themed. New components must work in
  every theme.
- Put each component in its own file in `components/`, each layout in
  `layouts/`, and import them from `style.css`.
- Use solid borders, and no inline styles.
- Name components with BEM: `.block`, `.block__element` and
  `.block__element--modifier`.
- Define fonts as tokens, and vendor them locally. Do not load fonts
  from an external CDN.
- After a template changes, run `templ generate` and commit the
  regenerated files.
