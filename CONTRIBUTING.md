# Contributing to Open Kata

## Types of contributions

| Artifact | Location | Scope |
|----------|----------|-------|
| Distributable skill | `skills/<name>/` | Portable, platform-agnostic |
| Distributable rule | `rules/<name>/` | Portable, platform-agnostic |
| Local skill | `.agents/skills/<name>/` | Repo-specific, not distributed |
| Local rule | `.agents/rules/<name>/` | Repo-specific, not distributed |

## Prerequisites

- [Go 1.26+](https://go.dev/dl/) — for building the MCP server
- [Tessl CLI](https://docs.tessl.io/introduction-to-tessl/installation) —
  for reviewing and evaluating kata

## Adding a skill

1. Use the `create-skill` skill or follow its workflow manually
2. Place in `skills/<name>/` for distributable,
   `.agents/skills/<name>/` for local
3. Symlink distributable skills into `.agents/skills/` for
   local use
4. Include CHANGELOG.md (Keep a Changelog format)
5. Include `references/ACKNOWLEDGMENTS.md` if drawing on
   external sources
6. Run `tessl skill review` — target 95%+

## Adding a rule

1. Use the `create-rule` skill or follow its workflow manually
2. Place in `rules/<name>/` for distributable,
   `.agents/rules/<name>/` for local
3. Symlink distributable rules into `.agents/rules/`
4. Validate against the rule design checklist in
   `skills/create-rule/references/`
5. Include CHANGELOG.md and ACKNOWLEDGMENTS.md as needed

## Quality bar

- Skills must score 95%+ on `tessl skill review`
- Rules must pass the rule design checklist
- All conventions in SKILL.md/RULE.md must be specific and
  literally enforceable
- Token cost matters — keep SKILL.md under 500 lines,
  RULE.md under 100 lines

## Commits

Follow Conventional Commits with lowercase descriptions:

    type(scope): description

See the `commit-conventions` skill for full details.

## Changelogs

- Follow Keep a Changelog format
- Blank lines around `###` headings
- Do not mention dev-only artifacts (tile.json, tessl.json)

## Design principles

- **Platform-agnostic** — no coupling to specific agents,
  MCP servers, or tooling
- **Progressive disclosure** — keep main files lean, move
  detail to `references/`
- **Explain the why** — reasoning over rigid ALWAYS/NEVER rules
- **Self-contained** — each artifact must work when installed
  standalone
