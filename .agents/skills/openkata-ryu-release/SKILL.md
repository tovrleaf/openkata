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

   Assisted-by: Kiro:claude-opus-4.6
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
