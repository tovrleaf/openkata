# AGENTS.md

## Project overview

Open Kata is a collection of portable agent skills and rules
following the Agent Skills specification. Skills live in
`skills/`, rules in `rules/`, and local companions in
`.agents/skills/` and `.agents/rules/`.

## Build and test

```bash
# Build the MCP server
cd cmd/openkata-mcp && go build .

# List skills and rules
make skills
make rules
```

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
- See [`rules/markdown-consistency/RULE.md`](rules/markdown-consistency/RULE.md)
  for full conventions

## Commits

Conventional Commits, lowercase descriptions:

    type(scope): description

Scope is the artifact name when applicable. See the
`commit-conventions` skill for full details.

Add an `Assisted-by` trailer to every AI-assisted commit.
Format: `Assisted-by: Agent:model_name`

This follows the Linux kernel convention from
`Documentation/process/coding-assistants.rst`. We use
`Assisted-by` rather than `Co-authored-by` because the
human is the author; the AI assisted.

For Kiro, read the model from the current session:

    ls -t ~/.kiro/sessions/cli/*.json | head -1 \
      | xargs grep -o '"model_name": "[^"]*"' \
      | cut -d'"' -f4

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
