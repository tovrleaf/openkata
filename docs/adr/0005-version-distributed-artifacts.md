---
status: PROPOSED
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

- Each distributed artifact has a `version` field in its frontmatter
- Each distributed artifact has a `CHANGELOG.md` in its directory
- Git tags follow the pattern `<directory-path>/v<major>.<minor>.<patch>`
- The dojo binary version is set in the Go source
- `CHANGELOG.md` is not copied to the target project on install — it
  stays in the source repository only

### Install manifest

When the dojo installs an artifact, it writes a `.manifest.json` in the
installed directory to track provenance:

```json
{
  "name": "create-adr",
  "version": "1.0.0",
  "source": "github.com/tovrleaf/openkata",
  "installedAt": "2026-04-16T20:26:45Z"
}
```

This allows the dojo to compare installed versions against available
versions and trace any installed artifact back to its origin.

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
- The dojo can compare frontmatter versions to detect updates

### Negative

- Maintainers must remember to bump version, update changelog, and
  create a git tag for each release
- More tags in the repo as the number of artifacts grows

### Neutral

- Automation can be added later to enforce version/tag/changelog
  consistency

## References

- [Semantic Versioning](https://semver.org/)
- ADR 0002 — MCP server for distribution (the dojo needs version info)
- ADR 0004 — Kata vocabulary (artifact type names)
