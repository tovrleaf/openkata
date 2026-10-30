---
title: Releasing a new version of the site
description: How a change to the website reaches production, what CI checks, and how to deploy by hand or go back.
order: 2
---

# Releasing a new version of the site

The website, the docs and the MCP server share one release number,
shown in the footer of every page as `Release v0.15.0`. A change
reaches production when it is merged into `main`, and the deploy
workflow ships it. A new release number comes from a git tag, described
below. Skills, rules and profiles are released differently, each with a
tag of its own, as described in
[versioning and the manifest](../reference/versioning.md).

## Before you push

You need Go, and the `templ` command at the version CI uses:

```bash
go install github.com/a-h/templ/cmd/templ@v0.3.1001
```

To see the site while you work, run `make dev`. It starts the site at
`http://localhost:8080` with hot reload, using `air` from your Go `bin`
folder. It regenerates the templates and rebuilds when Go, template or
CSS files change. Outside Lambda the
site reads its assets and artifacts from disk, and serves a few pages
that production does not, such as the design system. `make check`
tests that git, Go, make and tessl are installed. It does not test for
`templ` or `air`.

Before every push, run the checks that CI runs:

```bash
templ generate ./cmd/openkata-web/templates/
go build -o bin/openkata-web ./cmd/openkata-web/
go test ./cmd/openkata-web/...
```

Commit the regenerated `_templ.go` files together with the `.templ`
files they come from. CI regenerates them but does not compare them
with what you committed, so a stale file goes unnoticed until someone
else generates it.

If you added or removed a skill, rule or profile, run `make badges` and
commit the change. The README shows counts of each, and the pull
request build fails if they are out of date.

Branch and commit names follow the
[git-naming rule](https://github.com/tovrleaf/openkata/blob/main/rules/git-naming/RULE.md).
Pull requests are merged with a merge commit or a rebase, never
squashed, so each commit keeps its message and trailers.

## What CI checks on a pull request

The Build workflow runs on every pull request into `main`:

1. install `templ`
2. generate the templates
3. check that the README badge counts are current
4. `go build ./...`
5. `go test ./...`

## The release number

The release is the newest tag of the form `vMAJOR.MINOR.PATCH` that the
deployed commit contains, such as `v0.15.0`. Artifact tags, which look
like `skills/<name>/v1.0.0`, are never mistaken for it. The number is
read when the site is deployed and built into the web server and the MCP
server, so nothing else records it.

To release the site:

1. Decide the number from what changed since the last release tag:

   ```bash
   git tag -l 'v[0-9]*' | sort -V | tail -1
   git log --oneline <that-tag>..HEAD
   ```

   While the major number is 0, a `feat` or a breaking change raises the
   minor number, and anything else raises the patch. Do not release a
   change that no visitor or contributor could notice, such as one that
   only touches specs.

2. Tag the head of `main` once the change is merged. If there is
   something to commit with the release, give the commit the message
   `release(site): v0.15.0`.

   ```bash
   git tag v0.15.0
   ```

3. Push the tag, which deploys again with the new number, and write the
   release notes as a
   [GitHub release](https://github.com/tovrleaf/openkata/releases) for
   the tag. The footer links to it. `make push` pushes tags that were
   not pushed yet, and asks first.

Tag the head of `main`. The deploy of a tag does nothing unless the tag
is on the head of `main`, so a tag on an older commit, such as one made
to record the past, never puts old code in production. The deploy also
only sees a tag that its commit contains, so a pull request merged by
rebasing is released by tagging after the merge, not before.

## What happens on merge

Pushing to `main`, or pushing a release tag, starts the Deploy workflow,
which runs two jobs at the same time. They authenticate to AWS with the `openkata-ci` role, through
GitHub's OIDC, so no keys are stored.

**Web server.** The job checks out the full git history, then runs
`scripts/deploy-web.sh`:

1. generate the templates
2. write the docs last-updated dates from git, with
   `go run ./cmd/docs-dates -require`
3. read the release number with `go run ./cmd/site-version -require`
4. build the `openkata-web` binary for Linux on arm64, with the release
   built in
5. zip it and replace the code of the `openkata-web` Lambda function

**MCP server.** The job also checks out the full history, reads the
release number and builds `cmd/openkata-mcp` the same way, then
replaces the code of the `openkata-mcp` function.

Both jobs run on every merge, whatever changed. A deploy replaces
function code only. It does not touch the bucket, the index or the
tables. Publishing an artifact is a separate workflow, started by a
tag, as described in
[where artifacts are stored](../reference/storage.md).

The full history is needed for the dates and the tags. If the web job
runs on a shallow clone, or a docs page has no date, step 2 fails and
nothing is deployed. If the commit contains no release tag, step 3
fails, so the site never shows `Release dev`. Run locally, the site
shows `Release dev`, without a link.

## Deploying by hand

```bash
make deploy web
make deploy mcp
```

These run the same two scripts, `scripts/deploy-web.sh` and
`scripts/deploy-mcp.sh`. Set `AWS_PROFILE` to a profile that may update
the functions. The web script needs `templ` installed in your Go
`bin` folder, and a clone with its full history.

## Checking a release, and going back

Watch the Deploy run on the Actions page, then open the site and check
the change.

There is no rollback step. A deploy overwrites the function code, and
the scripts create no versions or aliases to switch back to. To go
back, revert the change and merge the revert, which deploys again, or
check out an earlier commit and deploy it by hand.
