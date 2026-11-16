---
status: Draft
depth: Standard
---

# Skill Frontmatter

## Story

As the maintainer, I want to know which Agent Skills fields our
skills should use and how to check them. Then every published skill
follows the specification and tells agents and people all it can.

## Background

The [Agent Skills specification](https://agentskills.io/specification)
defines these frontmatter fields:

| Field | Required | Constraints |
|---|---|---|
| `name` | Yes | 1–64 characters: lowercase letters, numbers and hyphens. No hyphen at the start or end, no `--`. Must match the folder name. |
| `description` | Yes | 1–1024 characters. Says what the skill does and when to use it. |
| `license` | No | A license name or a bundled license file. Keep it short. |
| `compatibility` | No | 1–500 characters. Environment needs: product, system packages, network access. Most skills do not need it. |
| `metadata` | No | A map from string keys to string values. Keys should be unique enough not to clash. |
| `allowed-tools` | No | Space-separated pre-approved tools, such as `Bash(git:*) Read`. Experimental: support varies between agents. |

The specification points to `skills-ref validate ./my-skill` as the
reference validator.

What the repo has today:

- All 14 skills in `skills/` use only `name`, `description`,
  `metadata.version` and `metadata.tags`.
- The repo license is MIT.
- `cmd/generate-versions` reads frontmatter line by line, not as
  YAML, and knows only `description` and `tags`.
- The templates in `create-skill` and `openkata-skill-conventions`
  decide which fields new skills get.
- Spec `0021-quality-lint` (Draft) plans a static lint that also
  checks frontmatter.

## Requirements

This spec starts as an investigation. Its first result is a set of
decisions, written back into this spec, before any skill changes.

### Investigation

- Check every skill against each rule of the specification, and
  list what fails.
- For each optional field, decide whether OpenKata uses it, and
  for which skills:
  - `license`
  - `compatibility`: which skills have real environment needs
  - `metadata` keys beyond `version` and `tags`, such as `author`
  - `allowed-tools`: whether an experimental, agent-specific field
    fits skills that must be platform-agnostic
- Decide how skills are checked against the specification: a Go
  test, `skills-ref validate` in CI, or both. Decide how this
  relates to spec `0021-quality-lint`.
- Decide whether the site, the MCP server or `versions.json` show
  any new field.
- Decide how the change is released: a patch version and tag for
  each changed skill, or none until each skill's next release.
- Decide whether rules (`RULE.md`) and profiles are in scope.

### Changes (to be fixed by the investigation)

- Every skill passes the chosen check.
- The skill templates and conventions describe the chosen fields,
  so new skills get them.
- The check runs in CI and fails the build when a skill breaks the
  specification.

## Out of Scope

- Rewriting skill bodies or descriptions, except where a check
  fails.
- Fields that the Agent Skills specification does not define.

## Rejected Approaches

- None yet.

## Open Questions

- Which optional fields to adopt for which skills.
- Go test, `skills-ref validate`, or both, and whether this spec
  absorbs `0021-quality-lint` or builds on it.
- Whether new fields appear on the site, in the MCP server or in
  `versions.json`.
- Whether every changed skill gets a version bump and tag.
- Whether rules and profiles are in scope.
