# AGENTS.md

## Project overview

Open Kata is a collection of portable agent skills and rules
following the Agent Skills specification. Skills live in
`skills/`, rules in `rules/`, and local companions in
`.agents/skills/` and `.agents/rules/`.

## Build and test

```bash
# Build binaries (always output to bin/)
go build -o bin/openkata-mcp ./cmd/openkata-mcp/
go build -o bin/openkata-web ./cmd/openkata-web/
go build -o bin/generate-versions ./cmd/generate-versions/

# List skills and rules
make skills
make rules
```

Never run `go build` without `-o bin/`. Binaries at the
project root indicate a build error.

## Before pushing

```bash
templ generate ./cmd/openkata-web/templates/
go build -o bin/openkata-web ./cmd/openkata-web/
go test ./cmd/openkata-web/...
```

Run these before every push so the committed `_templ.go`
files match the `.templ` sources. CI and the deploy
regenerate them, but neither compares them with what is
committed. The pull request build does fail when the README
badge counts are stale: run `make badges` after adding or
removing a skill, rule or profile.

Before opening a pull request for a spec'd feature, run
`/validate-spec`. It has a fresh agent check the diff against
`spec.md` alone and writes
`specs/NNNN-slug/validation-report.md`. See
[`.agents/prompts/validate-spec.md`](.agents/prompts/validate-spec.md).

## Code style

### Go

- Standard `go fmt` formatting
- Build must pass with zero warnings

### Testing

- See [`rules/go-testing/RULE.md`](rules/go-testing/RULE.md)
  for full conventions

### Bash

- `#!/usr/bin/env bash` with `set -euo pipefail`
- 2-space indent, 80 char line limit
- `[[ ]]` not `[ ]`, `"${var}"` not `$var`
- See [`rules/bash-style/RULE.md`](rules/bash-style/RULE.md)
  for full conventions

### Markdown

- Wrap prose at 72–80 characters
- ATX headings, `-` for lists, fenced code blocks
- See [`rules/markdown-style/RULE.md`](rules/markdown-style/RULE.md)
  for full conventions

## Git

Follow [`rules/git-naming/RULE.md`](rules/git-naming/RULE.md)
for branch names, commit messages, and the `Assisted-by`
trailer. See the `commit-conventions` skill for the full
commit workflow.

- Do not push without explicit user confirmation
- Do not force push without explaining the impact first
- Commit freely when asked, but pushing is a separate action
- `make push` pushes the branch and any tags that were not
  pushed yet, after asking. The user runs it, except when the
  agent opens a new pull request: the user has authorized the
  agent to run `make push YES=1` then, before creating the
  PR. Never push at any other time without confirmation

## Skills and rules

- Distributable artifacts must be platform-agnostic
- Skills: SKILL.md under 500 lines, target 95%+ on
  `tessl skill review`
- Rules: RULE.md under 100 lines, every convention literally
  enforceable
- Each artifact gets a CHANGELOG.md (Keep a Changelog format)
- External sources go in `references/ACKNOWLEDGMENTS.md`
- Do not run `tessl` CLI commands without asking the user first

## Symlinks

Distributable artifacts are symlinked into `.agents/` for
local use:

    .agents/skills/create-adr -> ../../skills/create-adr
    .agents/rules/bash-style -> ../../rules/bash-style

Claude Code reads slash commands from `.claude/commands`, which
links to the local prompts, so each file in `.agents/prompts/`
is a command:

    .claude/commands -> ../.agents/prompts

## Releasing

- Distributable artifacts get git tags:
  `skills/<name>/v<version>`
- Local artifacts (in `.agents/`) are not tagged
- Changelogs must not mention dev-only files
  (`.tessl-plugin/`, `.tesslignore`, tessl.json)
