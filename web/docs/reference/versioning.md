---
title: Versioning and the manifest
description: How artifacts are versioned with git tags and semantic versions, how changelogs are kept, and what the install manifest records.
order: 7
---

# Versioning and the manifest

Every skill, rule and profile is versioned on its own, so each one
changes at its own pace. The site, the docs and the MCP server share
one release number of their own, described at the end of this page.
This page explains how a version is decided and recorded, and what
the `.openkata-manifest.json` file in an installed artifact says.

## Versions

A version has three numbers, `MAJOR.MINOR.PATCH`, as in
[semantic versioning](https://semver.org/). The first release of an
artifact is `1.0.0`.

- **Patch** (`x.y.Z`): documentation, typos and other changes that do
  not change behavior, such as adding a `RATIONALE.md`.
- **Minor** (`x.Y.0`): new features, steps or references, and behavior
  changes that do not break existing use.
- **Major** (`X.0.0`): breaking changes, such as a renamed trigger, a
  removed step, or an output format that others rely on.

### The git tag is the source of truth

A release is a git tag that mirrors the artifact's folder:

```text
skills/<name>/v<version>
rules/<name>/v<version>
profiles/<name>/v<version>
```

For example, `skills/grill-me/v1.0.1`. The newest tag, compared as a
version, is the artifact's current version. Pushing the tag publishes
that version, as described in
[where artifacts are stored](storage.md).

The version also appears in two files, and all three have to agree
at a release:

- the `version` in the frontmatter of the main file (`SKILL.md`,
  `RULE.md` or `<name>.md`)
- the newest section of the artifact's `CHANGELOG.md`

When a tool needs the current version, it uses the tag. A changelog
or a frontmatter can be ahead of the last release, so neither is read
for that.

### Where versions show up

- In the bucket, each version has its own folder. Every version stays
  available.
- In `versions.json`, an artifact's `version` is its newest one. The
  index built from the repository gives `0.0.0` to an artifact that
  has never been tagged, and such an entry is not listed on the site
  or found by the catalog search.
- The MCP list tools return that version, and `install_*` installs it
  unless a version is named. The guide to
  [installing artifacts](../guides/installing-artifacts.md) shows how
  to name one.

## Changelogs

Each artifact has a `CHANGELOG.md` in the
[Keep a Changelog](https://keepachangelog.com/) format:

- an `## [Unreleased]` section at the top, always present
- one `## [x.y.z] - YYYY-MM-DD` section for each release, newest first
- under each release, only the headings Added, Changed, Deprecated,
  Removed, Fixed and Security, and only those that have entries
- entries that say what changed, in a sentence, and not how or why

The changelog is stored with each version, but it is not installed.
The eval runner reads the newest released heading, skipping
`Unreleased`, to find the version it ran against, as described in
[how we use evals](evals.md).

## The manifest

An installed artifact comes with a file named
`.openkata-manifest.json`. It records where the artifact came from
and what its files were when it was installed:

```json
{
  "name": "grill-me",
  "version": "1.0.1",
  "source": "github.com/tovrleaf/openkata",
  "checksums": {
    "SKILL.md": "sha256:<64 hex characters>",
    "references/<file>": "sha256:<64 hex characters>"
  }
}
```

| Field | Meaning |
|-------|---------|
| `name` | The artifact's name |
| `version` | The version that was installed |
| `source` | Where it came from, always `github.com/tovrleaf/openkata` |
| `checksums` | SHA-256 of each installed file, by path in the folder |

The checksums cover exactly the files that were installed. They do
not include the manifest itself or the files an install leaves out
(the changelog, the rationale, evals and tessl metadata). The time of
the install is not recorded. The project that installed the artifact
has it in its git history when it commits the files.

### How you get it

- An **MCP install** returns the manifest next to the files, as the
  `manifest` field of the result. The server does not write it. The
  tool's description presents it as the artifact's
  `.openkata-manifest.json`, for the agent to save in the artifact's
  folder with the other files.
- An **archive download** contains it at
  `<name>/.openkata-manifest.json`.

The content is the same either way. Commit it with the artifact, so
the record travels with the files.

### What it is for

Nothing in OpenKata reads the manifest yet. It is a record, and two
checks are possible with it:

- **Has a file changed?** Hash the file and compare it with its
  entry.
- **Is there a newer version?** Compare `version` with the newest tag,
  or with what `skill_versions`, `rule_versions` or `profile_versions`
  return.

The name starts with `openkata` so that it cannot be mistaken for the
generic `.manifest.json` that other tools use.

## The release of the site

The site, the docs, the MCP server and the tooling that makes them,
such as the eval runner and the scoring, have a release number too. It
is shown in the footer of every page, for example `Release v0.15.0`,
and links to that release on GitHub.

- It is a git tag of the form `vMAJOR.MINOR.PATCH`, with no prefix, so
  it is never mistaken for an artifact tag.
- It starts at `0`. While the major number is 0, a `feat` or a
  breaking change raises the minor number, and anything else raises the
  patch.
- The deploy reads the newest such tag that the deployed commit
  contains and builds it into the binaries. The MCP server reports it in
  its handshake, followed by the commit.
- The release notes are the GitHub release of the tag. There is no
  changelog file in the repository for it.

[Releasing a new version of the site](../contributing/release-site.md)
has the steps.
