---
status: ACCEPTED
date: 2026-04-12
authors: [niko.kivela]
---

# 0005. Version distributed artifacts with semver and git tags

## Context

OpenKata distributes multiple artifact types (kata, sensei profiles, kata
forms) and the dojo (MCP server binary). Users need to know what version
they have installed, what changed between versions, and whether updates
are available.

The repo itself is a web application and a collection — not something
consumers pin to a version.

## Decision Drivers

- Users need to pin and track versions of individual artifacts
- The dojo needs to compare installed vs available versions
- Different artifacts evolve at different rates
- Git tags should be unambiguous across artifact types
- Changelog history must be human-readable

## Decision

We will version all distributed artifacts using semver and git tags that
mirror directory paths. The repo itself has no version.

### Versioning scheme

| What | Versioning | Example tag |
|------|-----------|-------------|
| Skills, roles, prompts | Semver | `skills/create-adr/v1.2.0` |
| Dojo (MCP server binary) | Semver | `dojo/v0.1.0` |
| Repo | None | — |

### Implementation

- Git tags are the single source of truth for artifact versions
- Each artifact (distributable and local) has a `CHANGELOG.md`
- Git tags follow the pattern `<directory-path>/v<major>.<minor>.<patch>`
- The dojo binary version is set in the Go source
- The dojo resolves the version from the git tag at install time and
  writes it to `.openkata-manifest.json` — frontmatter does not carry
  version
- `CHANGELOG.md` is not copied to the target project on install — it
  stays in the source repository only

### Version resolution chain

Tools that need an artifact's version resolve it in this order:

1. **Git tag** — authoritative for released distributable artifacts
2. **CHANGELOG.md** — first `## ` heading, with or without brackets
   (`## 1.0.0` and `## [1.0.0]` both work)
3. **`versions.json`** — embedded build artifact for environments
   without git (e.g. Lambda deployments). Generated at build time
   from git tags.

Local artifacts (`.agents/`) are not tagged. Their version comes
from CHANGELOG.md only.

### Install manifest

When the dojo installs an artifact, it writes a hidden
`.openkata-manifest.json` in the installed directory as a mark of
provenance:

```json
{
  "name": "create-adr",
  "version": "1.0.0",
  "source": "github.com/tovrleaf/openkata"
}
```

This allows the dojo to compare installed versions against available
versions and trace any installed artifact back to its origin. The
file is a record next to the artifact, not an operational file. The
filename is namespaced with `openkata` to avoid collisions with other
tools that use a generic `.manifest.json`. Install time is not
stored; it is visible in the git history of the project that
committed the install.

### CHANGELOG format

Follows [Keep a Changelog](https://keepachangelog.com/) categories:

```markdown
# Changelog

## 1.1.0 — 2026-04-15
### Added
- Reversibility section to template
### Changed
- Status values now uppercase

## 1.0.0 — 2026-04-12
### Added
- Initial release
```

Categories: **Added**, **Changed**, **Deprecated**, **Removed**, **Fixed**,
**Security**. Omit categories with no entries.

## Alternatives Considered

### Repo-level semver for everything

- **Pros:** Simple, one version to track
- **Cons:** A typo fix in one prompt bumps the version for all artifacts,
  version number is meaningless for a collection
- **Rejected because:** The repo is not a distributable unit — individual
  artifacts are

### Git commit hashes only (skills.sh approach)

- **Pros:** Zero maintenance, automatic
- **Cons:** No human-readable version history, no way to communicate
  breaking changes, no changelogs
- **Rejected because:** Users can't tell what changed or whether an
  update is safe

### Date-based versioning for the repo

- **Pros:** Simple, no semver ambiguity for collections
- **Cons:** Doesn't apply to individual artifacts, no breaking change
  signal
- **Rejected because:** Artifacts need semver; the repo doesn't need
  versioning at all

## Consequences

### Positive

- Each artifact versions independently at its own pace
- Git tags are unambiguous — `skills/create-adr/v1.2.0` can't collide
  with a prompt of the same name
- Changelogs give users human-readable release history
- The dojo can compare manifest versions against git tags to detect
  updates

### Negative

- Maintainers must remember to update changelog and create a git tag
  for each release
- Version is not visible by reading SKILL.md alone — check the tag
  or changelog
- More tags in the repo as the number of artifacts grows

### Neutral

- Automation can be added later to enforce version/tag/changelog
  consistency

## Update — 2026-10-28

The decision stands. The implementation differs from the text in
these ways:

- The install manifest also records a SHA-256 checksum for each
  installed file (`checksums`), next to `name`, `version` and
  `source`.
- Frontmatter does carry the version, as `metadata.version`. The
  release workflow bumps it and checks it against the tag and the
  newest changelog section.
- The MCP server does not read git tags at install time. It takes
  the version from `versions.json`, or the version the caller names,
  and reads that version's folder from the bucket.
- `versions.json` is not an embedded build artifact. The publish
  workflow builds it with `cmd/generate-versions` from the bucket
  and uploads it to the bucket root. With `--local` the command
  builds it from the repository into `web/static/versions.json`.
- `CHANGELOG.md` is not the only file left out of an install. The
  full list is in `internal/exclude`.

## References

- [Semantic Versioning](https://semver.org/)
- ADR 0002 — MCP server for distribution (the dojo needs version info)
- ADR 0004 — Kata vocabulary (artifact type names)

## Update — 2026-10-30

The repo itself now has a release number: the site, docs and MCP server
share a `vMAJOR.MINOR.PATCH` git tag, shown in the footer of the site.
Artifact versioning is unchanged. See `web/docs/reference/versioning.md`.
