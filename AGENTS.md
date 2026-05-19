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

## Code style

### Go

- Standard `go fmt` formatting
- Build must pass with zero warnings

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

## Releasing

- Distributable artifacts get git tags:
  `skills/<name>/v<version>`
- Local artifacts (in `.agents/`) are not tagged
- Changelogs must not mention dev-only files (tile.json,
  tessl.json)
