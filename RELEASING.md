# Releasing

## Skills, rules and profiles

Skills, rules and profiles are released by pushing a git tag.
CI publishes the tagged content to S3 and regenerates the
versions index.

### Steps

1. Update `CHANGELOG.md` in the artifact's directory
2. Bump `version` in the frontmatter of the main file
   (`SKILL.md`, `RULE.md` or `<name>.md`) to match
3. Commit the changelog and the version
4. Create a tag: `git tag <type>/<name>/v<version>`
5. Push the commit and the tag: `make push`. It lists what it
   will push, asks first, and pushes each tag on its own

CI (`.github/workflows/publish.yaml`) handles the rest:
- Extracts the artifact's directory at the tagged commit
- Uploads it to `s3://openkata-artifacts/<type>/<name>/<version>/`
- Regenerates `versions.json` from the bucket and uploads it

The workflow can also be run by hand with a tag as its input.
When a pull request is merged, `publish-on-merge.yaml` re-pushes
the artifact tags that point at commits in that pull request,
which starts the publish again.

### Tag format

```
skills/<name>/v<major>.<minor>.<patch>
rules/<name>/v<major>.<minor>.<patch>
profiles/<name>/v<major>.<minor>.<patch>
```

### What gets published

All files in the artifact's directory at the tagged commit,
including `CHANGELOG.md`, `RATIONALE.md`, `evals/` and the
tessl files (`.tessl-plugin/`, `.tesslignore`). Each version
gets its own folder in the bucket.

### What gets installed by users

Everything published except the files that are left out when
an artifact is installed or downloaded (`internal/exclude`):
- `CHANGELOG.md`
- `RATIONALE.md`
- `references/ACKNOWLEDGMENTS.md`
- `evals/`
- `.tessl-plugin/`, `.tesslignore` and `tile.json`

An install also gets `.openkata-manifest.json` with the name,
version, source, and a SHA-256 checksum for each installed
file.

## Deployments

Deployments happen automatically on merge to main.

| Target | Script | CI job |
|--------|--------|--------|
| Web server | `scripts/deploy-web.sh` | `deploy-web` |
| MCP server | `scripts/deploy-mcp.sh` | `deploy-mcp` |

Both are in `.github/workflows/deploy.yaml`.

The web deploy checks out the full git history. Before it
builds, it runs `go run ./cmd/docs-dates -require`, which
writes each docs page's last-updated date from git into the
binary and fails if a page has no date, and
`go run ./cmd/site-version -require`, which reads the release
number and fails if the commit contains no release tag. Pushing
a release tag (`v*`) starts the same deploy.

### Manual deploy

```bash
AWS_PROFILE=openkata ./scripts/deploy-web.sh
AWS_PROFILE=openkata ./scripts/deploy-mcp.sh
```

### Releasing the site

The site, the docs and the MCP server share a release number,
shown in the footer as `Release v0.x.y`. It is the newest git
tag of the form `vMAJOR.MINOR.PATCH` that the deployed commit
contains. Artifact tags (`skills/<name>/v1.0.0`) are not
counted.

1. Decide the number from the commits since the last release
   tag: while the major is 0, a `feat` or a breaking change
   raises the minor, anything else the patch. Skip a release
   for a change nobody would notice, such as a spec change.
2. Tag the head of `main`: `git tag v0.x.y`
3. Push the tag with `make push`. Pushing a tag deploys again,
   with the new number.
4. Write the release notes as a GitHub release for the tag.
   The footer links to it.
