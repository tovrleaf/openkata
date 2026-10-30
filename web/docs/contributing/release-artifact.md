---
title: Releasing a skill, rule or profile
description: How to decide what to release, the steps from changelog to tag, what CI does with the tag, and how to rename an artifact.
order: 3
---

# Releasing a skill, rule or profile

A release of an artifact is a commit and a git tag. CI publishes the
tagged files, as described in
[where artifacts are stored](../reference/storage.md), and how the
numbers work is in
[versioning and the manifest](../reference/versioning.md). This page is
the procedure. The project has a local `openkata-ryu-release` skill
that does these steps for an agent, and the steps below are what it
does.

## Decide what to release

These commands list the artifacts with their current version:

```bash
make skills list
make rules
make profiles
```

A `*` after a version means the artifact has changed since its last
tag. For skills the check ignores the changelog, the evals and the
tessl files. `make skills list` also shows, for distributable skills,
whether the evals are up to date and the tessl score. `make skills
status` adds the version on the Tessl registry, and needs the network.

The current version is the newest tag of the artifact:

```bash
git tag -l '<type>/<name>/v*' | sort -V | tail -1
```

Never take it from `CHANGELOG.md` or the frontmatter. Either can be
ahead of the last release. Then pick the next number: patch, minor or
major, as described in the versioning page.

## The steps

1. **Check there is something to release.**

   ```bash
   git diff <latest-tag>..HEAD -- <type>/<name>/
   ```

   If there is no difference, stop.

2. **Update the README counts.** Run `make badges`, and include
   `README.md` in the release commit if it changed.

3. **Update `CHANGELOG.md`.** Add a section after `## [Unreleased]`:
   `## [x.y.z] - YYYY-MM-DD`, with what changed under the usual
   headings.

4. **Bump the version** in the frontmatter of the main file (`SKILL.md`,
   `RULE.md` or `<name>.md`), if the file has a `version` there. Not
   every artifact does, and then the changelog and the tag carry it.

5. **For a skill that has evals, run them.** Do this after the
   version bump and before the commit, with `make eval-local
   skills/<name>`. A result counts only when its `skill_version` is the
   skill's current version, so a report from before the bump is
   ignored and the release would show no effectiveness. The new
   report in `skills/<name>/evals/results/` goes into the release
   commit. How the evals work is in
   [how we use evals](../reference/evals.md).

6. **Commit** with the release message, adding the `Assisted-by`
   trailer if an agent helped:

   ```text
   release(skills): grill-me v1.0.1
   ```

   `release` is one of the commit types in the git-naming rule, for
   the commit that releases an artifact.

7. **Tag the release commit.**

   ```bash
   git tag <type>/<name>/vX.Y.Z
   ```

   Tag the commit you just made, never an earlier one.

8. **Push the commit and the tag.**

   ```bash
   make push
   ```

   It shows the commits and the tags that exist only in your clone,
   asks first, then pushes the branch and each tag with a push of its
   own, so no tag is forgotten. It ends with the check described under
   "Check that releases went out" below. To push by hand, push the
   tag separately: `git push origin <type>/<name>/vX.Y.Z`.

9. **For a skill, optionally publish it to the Tessl registry.**

   ```bash
   cd <type>/<name> && tessl tile publish --bump patch
   ```

   If the skill has not been set up for tessl yet, run
   `tessl skill import` first. Do not edit the `version` in
   `.tessl-plugin/plugin.json`. Tessl manages it, and it is expected to
   differ from the skill's version. See
   [how we use tessl](../reference/tessl.md).

Before you tag, these must all be true:

- the tag's version matches the `version` in the frontmatter, where
  the file has one
- it matches the newest section of `CHANGELOG.md`, and that section is
  not empty
- the tag points at the release commit
- for a skill with evals, the release commit holds a report for this
  version

## What CI does with the tag

Pushing a tag that matches `skills/*/v*`, `rules/*/v*` or
`profiles/*/v*` starts the publish workflow. It uploads the artifact
as it was at the tag to `<type>/<name>/<version>/` in the bucket and
regenerates `versions.json`, reading a skill's eval results from the
uploaded version.

If the tag cannot start the workflow, for example because the AWS
trust policy does not allow tag-triggered runs, start it by hand:

```bash
gh workflow run publish.yaml -f tag=<type>/<name>/vX.Y.Z
```

When a pull request is merged, a second workflow re-pushes the artifact
tags that point at commits in that pull request, which publishes them
again. So a release can go through a pull request too.

## Check that releases went out

Nothing else notices a tag that was never pushed, or a publish run that
failed. `make published` compares the tags in your clone with the
remote, and the pushed tags with the site:

```bash
make published
```

It lists tags that exist only in your clone, and pushed tags whose
files are not on the site, and exits with an error if there are any.
Tags of renamed artifacts are skipped, since their old files are
deleted on purpose. A missing version is published with the `gh
workflow run` command shown above. The check reads the live site and
the remote, so it needs the network. Set `OPENKATA_SITE` or
`OPENKATA_REMOTE` to check something else.

## Renaming an artifact

A rename starts the artifact's versions again at `1.0.0` under the new
name. The old tags stay in git, and the old files are deleted from the
bucket.

1. Rename the folder with `git mv <type>/<old> <type>/<new>`.
2. Change the `name` in the frontmatter.
3. Reset `CHANGELOG.md` to a fresh `1.0.0` entry that says the artifact
   was renamed from the old name.
4. Commit it, for example
   `refactor(skills): rename old-name to new-name`, and tag
   `<type>/<new>/v1.0.0`.
5. Publish the new name by running the workflow by hand, as above, and
   delete the old files:

   ```bash
   aws s3 rm s3://openkata-artifacts/<type>/<old>/ --recursive
   ```

   This needs AWS credentials that may delete from the bucket. The CI
   role cannot.
6. Add a redirect so old links keep working. In
   `cmd/openkata-web/handlers.go`, add the old name to the
   `artifactRedirects` map, under its type, pointing to the new name.
   The site then sends a permanent redirect from the old page.

## Mistakes to avoid

- Reading the current version from the changelog or the frontmatter
  instead of the tag.
- Releasing without a difference since the last tag.
- Running the evals before the version bump, or not committing the
  report. The index is built from what the tag contains.
- Tagging a commit other than the release commit.
- Editing the version in `plugin.json` by hand.
