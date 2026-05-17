---
name: openkata-ryu-release
description: >
  Release a new version of an OpenKata artifact (skill or rule). Detects
  changes since the last release, recommends a semver bump, updates
  frontmatter and changelog, commits, and creates a git tag. Use when the
  maintainer says "release", "bump", "version", or "tag" followed by an
  artifact name.
---

# Ryu Release

Internal skill for releasing new versions of OpenKata artifacts. Supports
skills (SKILL.md) and rules (RULE.md). This skill is not distributed — it's
for maintainers of this repository only.

## Workflow

1. **Identify the artifact** — Determine which artifact to release from
   the user's request (e.g., "release create-adr"). Resolve the directory
   path by checking both `skills/` and `rules/` for distributable artifacts,
   and `.agents/skills/` and `.agents/rules/` for local ones.

2. **Find the last release** — Look for the most recent git tag matching
   `<directory-path>/v*` (e.g., `skills/create-adr/v1.0.0`). If no tag
   exists, this is the first release.

3. **Diff since last release** — Run `git diff <last-tag> -- <directory-path>`
   to see what changed. If no previous tag, diff against the initial commit.
   Read the actual changes to understand what was modified.

4. **Recommend a bump level** — Based on the diff, recommend one of:
   - **major** — Breaking changes: removed sections, renamed fields,
     changed behavior that would break existing users' workflows
   - **minor** — New features: added sections, new workflow steps,
     new optional fields
   - **patch** — Fixes: typo corrections, wording improvements,
     clarifications that don't change behavior

   Present the recommendation with justification. Let the user confirm
   or override.

5. **Update changelog** — Prepend a new entry to the artifact's
   CHANGELOG.md using [Keep a Changelog](https://keepachangelog.com/)
   categories (Added, Changed, Deprecated, Removed, Fixed, Security):
   - The new version number and today's date
   - Changes grouped by category, omitting empty categories

   If CHANGELOG.md doesn't exist, create it.

6. **Generate tags** — Read the artifact's SKILL.md (or RULE.md)
   content — title, description, workflow steps, and section
   headings — and generate 2–5 namespaced tags describing what
   the artifact does.

   Tags use `namespace:value` format:
   - `category:` — what the skill does (e.g., `category:documentation`,
     `category:testing`, `category:workflow`)
   - `language:` — programming language if applicable (e.g.,
     `language:go`, `language:bash`)
   - `tool:` — external tool the skill directly invokes or
     operates on (e.g., `tool:git`, `tool:makefile`). Do not
     add a tool tag just because the domain is related — only
     if the skill actually calls or configures the tool.

   Tags must be:
   - Lowercase, hyphenated values (e.g., `category:code-review`)
   - Derived from the skill's domain, actions, and outputs
   - 2–5 tags per artifact
   - Not redundant — avoid tagging the same concept in multiple
     namespaces (e.g., don't use both `category:git` and
     `tool:git`)

   Write them into the frontmatter as `metadata.tags`:

   ```yaml
   ---
   name: create-adr
   description: Creates Architecture Decision Records
   metadata:
     tags: "category:documentation, category:architecture, tool:git"
   ---
   ```

   Present the proposed tags to the user for review. Accept
   edits before continuing.

7. **Regenerate root changelog** — Run `make changelog` to update
   the aggregate CHANGELOG.md at the repo root.

8. **Commit** — Stage the changed files (including root CHANGELOG.md)
   and commit with message:
   `release(<artifact-name>): v<new-version>`

9. **Tag** — For distributable artifacts (`skills/`, `rules/`) only:
   create a git tag `<directory-path>/v<new-version>`. Skip tagging for
   local artifacts in `.agents/` — they are versioned but not distributed.

10. **Publish to registry** — If the artifact is distributable
    (`skills/` or `rules/`) and has a `tile.json`, ask: "Want me
    to publish this to the Tessl registry?" If yes, run
    `tessl tile publish <directory-path>`.

11. **Confirm** — Show the user what was done: version bumped, changelog
    entry, commit hash, and tag name (if tagged). Ask if they want to
    push.

## Bump Examples

| Change | Bump | Why |
|--------|------|-----|
| Fixed typo in Context section guidance | patch | No behavior change |
| Added optional Non-goals section to template | minor | New feature, non-breaking |
| Reworded existing guidance for clarity | patch | Clarification only |
| Changed status values from lowercase to uppercase | major | Breaks existing ADRs using old format |
| Added new workflow step (explore codebase) | minor | New capability |
| Removed Alternatives Considered section | major | Breaking removal |

## Tag Format

Tags mirror directory paths per ADR 0005:

```text
skills/create-adr/v1.0.0
rules/markdown-consistency/v1.0.0
dojo/v0.1.0
```

List all tags for an artifact: `git tag -l "skills/create-adr/v*"`

## Gotchas

- Always check for uncommitted changes in the artifact directory before
  starting. Refuse to release with a dirty working tree.
- Never skip the changelog — it's the only human-readable release history.
- If the user disagrees with the recommended bump level, use theirs.
- Local artifacts (`.agents/skills/`, `.agents/rules/`) get version bumps
  and changelogs but no git tags — tags are reserved for distributable
  artifacts (ADR 0005).
- Changelogs document artifact-facing changes only. Dev-only files
  (tile.json, tessl.json) are not changelog-worthy.
