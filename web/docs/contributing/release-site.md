---
title: Releasing a new version of the site
description: How a change to the website reaches production, what CI checks, and how to deploy by hand or go back.
order: 2
---

# Releasing a new version of the site

The website and the MCP server have no version number. A change
reaches production when it is merged into `main`, and the deploy
workflow ships it. This page covers that path. Skills, rules and
profiles are released differently, with a git tag, as described in
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

## What happens on merge

Pushing to `main` starts the Deploy workflow, which runs two jobs at the
same time. They authenticate to AWS with the `openkata-ci` role, through
GitHub's OIDC, so no keys are stored.

**Web server.** The job checks out the full git history, then runs
`scripts/deploy-web.sh`:

1. generate the templates
2. write the docs last-updated dates from git, with
   `go run ./cmd/docs-dates -require`
3. build the `openkata-web` binary for Linux on arm64
4. zip it and replace the code of the `openkata-web` Lambda function

**MCP server.** The job builds `cmd/openkata-mcp` the same way and
replaces the code of the `openkata-mcp` function.

Both jobs run on every merge, whatever changed. A deploy replaces
function code only. It does not touch the bucket, the index or the
tables. Publishing an artifact is a separate workflow, started by a
tag, as described in
[where artifacts are stored](../reference/storage.md).

The full history is needed for the dates. If the web job runs on a
shallow clone, or a docs page has no date, step 2 fails and nothing is
deployed.

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
