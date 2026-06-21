---
name: openkata-ryu-release
description: >
  Releases distributable artifacts (skills, rules, profiles)
  with proper versioning, changelog, and git tags.
---

# openkata-ryu-release

Release a distributable artifact following Semantic
Versioning. The git tag is the single source of truth
for the current published version.

## Trigger

Activate when the user says "release", "bump", "version",
"tag", or "publish" an artifact.

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

2. Update `CHANGELOG.md` — add a version section after
   `[Unreleased]` with today's date and what changed.

3. Bump version in the frontmatter of the main doc
   (`SKILL.md`, `RULE.md`, or `<name>.md`).

4. Commit:
   ```text
   release(<type>): <name> vX.Y.Z

   Assisted-by: <agent>:<model>
   ```

5. Tag:
   ```bash
   git tag <type>/<name>/vX.Y.Z
   ```

6. Push (only with user confirmation):
   ```bash
   git push origin main
   git push origin <type>/<name>/vX.Y.Z
   ```

## Validation

- Tag version must match SKILL.md frontmatter version
- Tag version must match latest CHANGELOG.md section
- Tag must point to the release commit (not an
  earlier commit)
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

## Common Failures

- NEVER determine current version from CHANGELOG.md —
  it may contain unreleased or untagged entries
- NEVER skip the diff check — releasing unchanged
  artifacts creates noise
- NEVER use `git tag` on a commit other than HEAD
- MUST push the tag separately from the commit if the
  IAM trust policy doesn't allow tag-triggered workflows
  (use `gh workflow run publish.yaml -f tag=...` as
  workaround)
- NEVER touch `.tessl-plugin/plugin.json` version during
  release — it is Tessl-internal and independent of our
  versioning (ADR 0006). A difference between SKILL.md
  version and plugin.json version is expected, not a bug.
