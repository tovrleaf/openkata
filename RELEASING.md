# Releasing

## Skills and rules

Skills and rules are released by pushing a git tag.
CI publishes the tagged content to S3 and regenerates
the versions index.

### Steps

1. Update `CHANGELOG.md` in the skill/rule directory
2. Commit the changelog
3. Create a tag: `git tag skills/<name>/v<version>`
4. Push the tag: `git push --tags`

CI (`.github/workflows/publish.yaml`) handles the rest:
- Extracts files at the tagged commit
- Uploads to `s3://openkata-artifacts/<type>/<name>/<version>/`
- Regenerates `versions.json`

### Tag format

```
skills/<name>/v<major>.<minor>.<patch>
rules/<name>/v<major>.<minor>.<patch>
```

### What gets published

All files in the skill/rule directory at the tagged
commit, except:
- `tile.json` — tessl metadata, not distributed
- `CHANGELOG.md` — stays in the repo, not installed

### What gets installed by users

- All skill/rule files (SKILL.md, references/, assets/)
- `.manifest.json` with version, source, and checksums
- Excludes `tile.json` and `references/ACKNOWLEDGMENTS.md`

## Deployments

Deployments happen automatically on merge to main.

| Target | Script | CI job |
|--------|--------|--------|
| Web server | `scripts/deploy-web.sh` | `deploy-web` |
| MCP server | `scripts/deploy-mcp.sh` | `deploy-mcp` |

Both are in `.github/workflows/deploy.yaml`.

### Manual deploy

```bash
AWS_PROFILE=openkata ./scripts/deploy-web.sh
AWS_PROFILE=openkata ./scripts/deploy-mcp.sh
```
