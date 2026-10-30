---
status: Draft
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

## Decisions to confirm

These are proposals. The open questions are marked.

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
- **The deploy reads the newest `v*` tag** with
  `git describe --tags --match 'v[0-9]*'`. The deploy checkout already
  has full history. The number is built into the web binary at deploy
  time, so it is not committed anywhere else and cannot disagree with
  the tag.
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
  never every tag, and ends with `make published`. It is the
  user's command: agents still do not push on their own.
- **The past**: the history before the first tag is given releases
  from the commit history, proposed in [history.md](history.md). The
  tags are made on those commits, so they carry the commit dates.
  Making and pushing them, and writing the GitHub releases, is yours.
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
- The proposed history is accepted or changed before any tag is made

## Out of scope

- A release notes page on the site, and a root `CHANGELOG.md`
- Automatic version bumps from commit messages
- Versioning the infrastructure separately

## Open questions

1. Is the list in [history.md](history.md) right? It decides the
   numbers, so it has to be settled before the tags are made.
2. The history ends at `v0.14.0`, so the first deploy of the new
   footer would show `v0.14.0` until `v0.15.0` is tagged. Should
   `v0.15.0` be tagged just before the go-live merge?
