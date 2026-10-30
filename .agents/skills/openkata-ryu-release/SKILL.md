---
name: openkata-ryu-release
description: >
  Releases distributable artifacts (skills, rules, profiles)
  with proper versioning, changelog, git tags, and optional
  tessl registry publish.
---

# openkata-ryu-release

Release a distributable artifact following Semantic
Versioning. The git tag is the single source of truth
for the current published version.

## Trigger

Activate when the user says "release", "bump", "version",
"tag", or "publish" an artifact, or "release the site". The
site steps are under Releasing the Site.

## Determine Current Version

```bash
git tag -l '<type>/<name>/v*' | sort -V | tail -1
```

NEVER read the version from CHANGELOG.md or SKILL.md to
determine the current version. Those files may be ahead
of (or behind) the actual tagged release.

## Determine Next Version

Follow [Semantic Versioning](https://semver.org/):

- **PATCH** (x.y.Z): documentation additions, typo fixes,
  non-functional changes (added RATIONALE.md, reformatted
  ACKNOWLEDGMENTS.md, fixed examples)
- **MINOR** (x.Y.0): new features, new steps, new
  references, behavioral changes that don't break
  existing usage
- **MAJOR** (X.0.0): breaking changes to the skill's
  interface (renamed triggers, removed steps, changed
  output format that downstream consumers rely on)

When in doubt, ask the user.

## Release Steps

1. Confirm the artifact has uncommitted changes vs its
   latest tag:
   ```bash
   git diff <latest-tag>..HEAD -- <type>/<name>/
   ```
   If no diff, abort — nothing to release.

2. Run the badge update script to sync counts:
   ```bash
   ./scripts/update-readme-badges.sh
   ```
   If README.md changed, stage it for the release commit.

3. Update `CHANGELOG.md` — add a version section after
   `[Unreleased]` with today's date and what changed.

4. Bump version in the frontmatter of the main doc
   (`SKILL.md`, `RULE.md`, or `<name>.md`).

5. For a skill that has evals, run them now, after the version
   bump and before the commit. A result counts only when its
   `skill_version` equals the skill's version, so results from
   before the bump are ignored. The new report in
   `evals/results/` goes into the release commit.

6. Commit:
   ```text
   release(<type>): <name> vX.Y.Z

   Assisted-by: <agent>:<model>
   ```

7. Tag:
   ```bash
   git tag <type>/<name>/vX.Y.Z
   ```

8. Push by telling the user to run `make push`. It lists
   the commits and the tags that exist only locally, asks,
   then pushes the branch and each tag separately. Run it
   yourself (`make push YES=1`) only when opening a new pull
   request, which the user has authorized.

9. Ask the user: "Release to tessl registry as well? (y/n)"
   If yes:
   ```bash
   cd <type>/<name> && tessl tile publish --bump patch
   ```
   If the tile doesn't exist yet, run `tessl skill import`
   first, then `tessl tile publish`.

## Validation

- Tag version must match SKILL.md frontmatter version
- Tag version must match latest CHANGELOG.md section
- Tag must point to the release commit (not an
  earlier commit)
- For a skill with evals, the release commit holds a report
  whose `skill_version` is the tag's version
- CHANGELOG entry must not be empty

## Renaming an Artifact

Renaming starts versioning from v1.0.0 under the new name.
Old git tags are preserved — old S3 artifacts are deleted.

1. Rename the directory:
   ```bash
   git mv <type>/<old-name> <type>/<new-name>
   ```

2. Update frontmatter `name:` field.

3. Reset CHANGELOG.md to a fresh v1.0.0 entry. Mention
   the old name:
   ```
   ## [1.0.0] - YYYY-MM-DD

   ### Changed

   - Renamed from <old-name>
   ```

4. Commit, tag, push:
   ```bash
   git commit -m 'refactor(<type>): rename <old-name> to <new-name>'
   git tag <type>/<new-name>/v1.0.0
   ```

5. Publish new name, delete old S3 prefix:
   ```bash
   gh workflow run publish.yaml -f tag=<type>/<new-name>/v1.0.0
   aws s3 rm s3://openkata-artifacts/<type>/<old-name>/ --recursive
   ```

6. Add redirect in `artifactRedirects` map in handlers.go:
   ```go
   "<old-name>": "<new-name>",
   ```

## Releasing the Site

The site, docs and MCP server share one release number, the
newest tag of the form `vMAJOR.MINOR.PATCH` (no prefix). The
footer shows it and links to its GitHub release. Decide it at
each push, and state the number and why before the user
confirms the push.

1. Find the current release and what changed since:
   ```bash
   git tag -l 'v[0-9]*' | sort -V | tail -1
   git log --oneline <that-tag>..HEAD
   ```

2. Pick the next number. While the major is 0, a `feat` or
   a breaking change raises the minor; anything else raises
   the patch. If nothing a visitor or contributor could
   notice changed (for example only specs), do not release.

3. If there is something to commit with the release, use
   `release(site): vX.Y.Z` as its message.

4. Tag the head of `main`, after the change is merged:
   ```bash
   git tag vX.Y.Z
   ```
   A tag that is not on the head of `main` deploys nothing.

5. Push (only with user confirmation). Pushing the tag
   deploys again with the new number. Then write the release
   notes as a GitHub release for the tag.

## Common Failures

- NEVER determine current version from CHANGELOG.md —
  it may contain unreleased or untagged entries
- NEVER run the evals before the version bump — the report
  would be for the previous version and the site would show
  no effectiveness for the release
- NEVER skip the diff check — releasing unchanged
  artifacts creates noise
- NEVER use `git tag` on a commit other than HEAD
- MUST push the tag separately from the commit if the
  IAM trust policy doesn't allow tag-triggered workflows
  (use `gh workflow run publish.yaml -f tag=...` as
  workaround)
- NEVER manually edit `.tessl-plugin/plugin.json` version —
  tessl manages it via `tessl tile publish --bump`. A
  difference between SKILL.md version and plugin.json
  version is expected (ADR 0006).
