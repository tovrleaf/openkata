---
status: Done
depth: Standard
---

# A Version for the Site

## Story

As a visitor to openkata.dev, I want to see which release of the site
I am looking at, as asdlc.io does with "Release v0.27.2", so that I
can tell what has changed since I last looked, and as the maintainer,
I want the whole machine, not only the skills, to have a release
history.

## Background

Only skills, rules and profiles have versions today, each with its own
tag. The site, the docs, the processes, the eval runner, the scoring
and the MCP server change on every merge to `main` and have no number.
The MCP server reports `0.3.0+<commit>`, which exists only so that
clients can tell builds apart.

asdlc.io shows the plain text `Release v0.27.2` at the bottom of every
page, after the copyright line. It has no link.

## Decisions

These were confirmed with the user. Where the code differs from the
first wording, the text says so.

- **One version for the project.** The site, the docs, the MCP server
  and the tooling share one number. Artifact versions stay separate.
- **A git tag `v<major>.<minor>.<patch>`**, with no prefix, so it
  cannot be mistaken for an artifact tag (`skills/<name>/v...`).
  Releases start at `v0.1.0`. While the major is 0, a breaking change
  bumps the minor.
- **A release is a tag and a GitHub release.** The notes live in the
  GitHub release, not in a root `CHANGELOG.md`. A root changelog was
  removed on 2026-10-16 (PR 31) because the aggregate went stale and
  nothing read it, and a hand-kept one would risk the same. The tag is
  made on a commit `release(site): v0.x.y` when there is something to
  commit, such as a docs update. The `site` scope is new in the
  git-naming rule.
- **The deploy reads the newest `v*` tag** that the deployed commit
  contains, with `git tag --merged HEAD`, compared by number. This
  replaces the `git describe` of the first draft, which orders by
  ancestry distance, not by number. A consequence is that a mistyped
  high tag, such as `v0.150.0`, wins until it is deleted. The deploy
  checkout has full history. The number is built into the web and MCP
  binaries at deploy time, so it is not committed anywhere else and
  cannot disagree with the tag.
- **Pushing a release tag deploys.** The Deploy workflow also runs on
  `v*` tags, so a release is live as soon as its tag is pushed, even
  when the commit was merged earlier. A tag deploys only when it is on
  the head of `main`, so the tags that record the past, which are on
  old commits, deploy nothing. The deploy fails, on purpose, if the
  commit contains no release tag, so the site never shows
  `Release dev`. That means the past tags must be pushed before the
  first deploy of this change.
- **Commits after the tag**: the footer shows the release only. The
  commits since are not shown (confirmed).
- **Where it shows**: the shared footer of every page, including the
  docs, as `Release v0.x.y`, linking to that release on GitHub,
  `https://github.com/tovrleaf/openkata/releases/tag/v0.x.y`
  (confirmed). With no tag it shows `Release dev`, unlinked.
- **Who picks the number**: the agent, at each push, from the commits
  since the newest `v*` tag: a `feat` or a breaking change bumps the
  minor, anything else the patch (while the major is 0). It states the
  number and why, and the user confirms the push as usual. A push that
  changes nothing a visitor or contributor could notice, such as a
  spec-only change, does not get a release.
- **Pushing tags**: `make push` pushes the current branch and then
  each tag that exists only locally, one at a time, so a tag cannot be
  forgotten and the publish workflow starts for each artifact tag. It
  only pushes tags of the forms `<type>/<name>/vX.Y.Z` and `vX.Y.Z`,
  never every tag, and ends with the check that `make published` runs.
  It is the user's command. The user has authorized the agent to run
  `make push YES=1`, which skips the question, when it opens a new pull
  request, and at no other time.
- **The past**: the history before the first tag was given releases
  from the commit history, listed in [history.md](history.md). The
  tags `v0.1.0` to `v0.14.0` were made on those commits, so they carry
  the commit dates, and were pushed, and a GitHub release was written
  for each. The user did this.
- **The MCP server** reports the same number as its `serverInfo`
  version, replacing `0.3.0+<commit>`, and keeps the commit as build
  metadata, so `0.3.0+<commit>` becomes `<release>+<commit>`.
- **Locally** the footer shows `Release dev` when there is no tag to
  read, the same way the docs show a placeholder for dates.

## Requirements

- Every page of the site shows `Release v<version>` in the footer,
  linked to the GitHub release of that tag
- The version is the newest `v*` tag reachable from the deployed
  commit, and `dev` when there is none
- A tag that is not `v<number>.<number>.<number>` is ignored
- The MCP server's reported version starts with the same number
- `release(site)` is a valid commit scope in the git-naming rule
- The release page for the site describes how to make a release: the
  tag, the GitHub release and what a deploy shows
- The history is accepted or changed before any tag is made (done)

## Out of scope

- A release notes page on the site, and a root `CHANGELOG.md`
- Automatic version bumps from commit messages
- Versioning the infrastructure separately

## Resolved questions

1. The list in [history.md](history.md) was accepted, and the tags
   exist.
2. The history ends at `v0.14.0`, so the first deploy of the new
   footer shows `v0.14.0`. `v0.15.0` is tagged on the head of `main`
   after the merge, and pushing that tag deploys again with it.

## Known limits

- Nothing checks that a new tag follows the previous one, so a typo
  in a tag is not caught.
- A tag made on a branch is not contained in `main` after a rebase
  merge, so the site keeps the old number. Tag after the merge, or
  merge with a merge commit.
- The Deploy workflow has no concurrency group, so a branch push and a
  tag push can deploy at the same time. The code is the same.
- A manual deploy from a shallow clone can report an older release.
- That the tag trigger and the head-of-`main` guard work on GitHub has
  not been seen yet. The first deploy after the merge shows it.
