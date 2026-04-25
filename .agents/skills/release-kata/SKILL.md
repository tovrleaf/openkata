---
name: release-kata
version: 1.0.0
description: >
  Release a new version of an OpenKata artifact. Detects changes since the
  last release, recommends a semver bump, updates frontmatter and changelog,
  commits, and creates a git tag. Use when the maintainer says "release",
  "bump", "version", or "tag" followed by a skill name.
---

# Release Kata

Internal skill for releasing new versions of OpenKata artifacts. This skill
is not distributed — it's for maintainers of this repository only.

## Workflow

1. **Identify the artifact** — Determine which artifact to release from
   the user's request (e.g., "release create-adr"). Resolve the directory
   path (e.g., `skills/create-adr`). Determine if it is distributable
   (`skills/`) or local (`.agents/skills/`).

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

5. **Update version** — Bump the `version` field in the artifact's
   frontmatter (SKILL.md or equivalent).

6. **Update changelog** — Prepend a new entry to the artifact's
   CHANGELOG.md using [Keep a Changelog](https://keepachangelog.com/)
   categories (Added, Changed, Deprecated, Removed, Fixed, Security):
   - The new version number and today's date
   - Changes grouped by category, omitting empty categories
   
   If CHANGELOG.md doesn't exist, create it.

7. **Commit** — Stage the changed files and commit with message:
   `release(<artifact-name>): v<new-version>`

8. **Tag** — For distributable artifacts in `skills/` only: create a git
   tag `<directory-path>/v<new-version>` (e.g., `skills/create-adr/v1.1.0`).
   Skip tagging for local skills in `.agents/skills/` — they are versioned
   but not distributed.

9. **Confirm** — Show the user what was done: version bumped, changelog
   entry, commit hash, and tag name (if tagged). Ask if they want to push.

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

```
skills/create-adr/v1.0.0
skills/create-adr/v1.1.0
dojo/v0.1.0
```

List all tags for an artifact: `git tag -l "skills/create-adr/v*"`

## Gotchas

- Always check for uncommitted changes in the artifact directory before
  starting. Refuse to release with a dirty working tree.
- The version in frontmatter and the git tag must match.
- Never skip the changelog — it's the only human-readable release history.
- If the user disagrees with the recommended bump level, use theirs.
- Local skills (`.agents/skills/`) get version bumps and changelogs but
  no git tags — tags are reserved for distributable artifacts (ADR 0005).
