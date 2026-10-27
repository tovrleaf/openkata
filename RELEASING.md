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
5. Push the tag: `git push origin <type>/<name>/v<version>`

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
binary and fails if a page has no date.

### Manual deploy

```bash
AWS_PROFILE=openkata ./scripts/deploy-web.sh
AWS_PROFILE=openkata ./scripts/deploy-mcp.sh
```
