# Tasks: A Version for the Site

## Tasks

### 1. Read the release version from git
- **Status**: Pending
- **Goal**: One function turns the tags of a clone into the site's version
- **Boundary**: `internal/siteversion/` (new)
- **Key files**: `cmd/openkata-mcp/version.go` for the existing stamp
- **Depends**: None
- **Done when**: a pure function picks the newest tag of the form
  `vMAJOR.MINOR.PATCH` from a list of tags, by number, so `v0.10.0`
  beats `v0.9.0`, and ignores every other tag including artifact tags;
  a second function runs `git describe --tags --match 'v[0-9]*'` and
  returns `dev` when there is no tag or no git; tests cover several
  tags, none, a pre-release suffix, and artifact tags only
- **Verify**: `go test ./internal/siteversion/`

### 2. A shared footer with the release
- **Status**: Pending
- **Goal**: Every page shows `Release v0.x.y`, linked to GitHub
- **Boundary**: `cmd/openkata-web/templates/*.templ` and the generated
  files, `cmd/openkata-web/templates/docs.templ`,
  `web/static/css/components/` for the footer style
- **Key files**: the 11 page templates that repeat the footer, and
  `docs.templ`
- **Depends**: 1
- **Done when**: one footer component replaces the copies; it shows
  the copyright line and the licence as today, then `Release v<version>`
  linking to `https://github.com/tovrleaf/openkata/releases/tag/v<version>`;
  for `dev` it shows `Release dev` without a link; the style uses
  design tokens only; a test renders the component for a version and
  for `dev`
- **Verify**: `templ generate ./cmd/openkata-web/templates/` then
  `go test ./cmd/openkata-web/...`, and look at a page with `make dev`

### 3. Build the version into the binaries
- **Status**: Pending
- **Goal**: The deployed site and MCP server know their release
- **Boundary**: `cmd/openkata-web/`, `cmd/openkata-mcp/`,
  `scripts/deploy-web.sh`, `scripts/deploy-mcp.sh`
- **Key files**: `serverVersion` in `cmd/openkata-mcp/version.go`
- **Depends**: 1, 2
- **Done when**: both deploy scripts read the version with the function
  from task 1 and pass it with `-ldflags -X`; the web footer shows it;
  the MCP server reports `<release>+<commit>` in place of
  `0.3.0+<commit>`; both fall back to `dev` when run from source; tests
  cover the MCP version with and without a stamp
- **Verify**: `go test ./cmd/openkata-mcp/ ./cmd/openkata-web/`, and
  build with the flag and read the version

### 4. Describe the site release in the rules and docs
- **Status**: Pending
- **Goal**: A release of the site is a known procedure
- **Boundary**: `rules/git-naming/RULE.md` and its changelog,
  `web/docs/contributing/release-site.md`,
  `web/docs/reference/versioning.md`,
  `.agents/skills/openkata-ryu-release/SKILL.md` and its changelog,
  `RELEASING.md`
- **Key files**: the release steps in the release skill
- **Depends**: 3
- **Done when**: `site` is a valid commit scope for `release`; the site
  release page says the site has a version, how the agent picks the
  bump at each push (`feat` or a breaking change bumps the minor,
  anything else the patch, no release for a change nobody would
  notice), the tag, the GitHub release and what a deploy shows; the
  release skill has the site steps; the docs that say the site has no
  version are corrected
- **Verify**: `go test ./cmd/openkata-web/...` for the docs tree

### 5. `make push`
- **Status**: Pending
- **Goal**: A tag is never forgotten when pushing
- **Boundary**: `scripts/push.sh` (new), `mk/dev.mk`
- **Key files**: `scripts/check-published.sh`
- **Depends**: None
- **Done when**: the script prints the commits and the tags that
  exist only locally and waits for a yes; it then pushes the current
  branch and each of those tags with its own `git push`; it only
  considers `<type>/<name>/vX.Y.Z` and `vX.Y.Z`; with nothing to push
  it says so; it ends with `make published`; it follows the bash-style
  rule; `make push` runs it
- **Verify**: run it with the default answer on a clone that has an
  unpushed tag, and check it lists the 5 known tags without pushing

### 6. The release docs use `make push`
- **Status**: Pending
- **Goal**: One way to push, in every description of a release
- **Boundary**: `.agents/skills/openkata-ryu-release/SKILL.md`,
  `web/docs/contributing/release-artifact.md`,
  `web/docs/contributing/release-site.md`, `RELEASING.md`,
  `AGENTS.md`
- **Key files**: the push step in each
- **Depends**: 4, 5
- **Done when**: the push step of each release procedure says to run
  `make push`, which asks first, and the agent rule that pushing needs
  the user's confirmation is unchanged
- **Verify**: search the docs for `git push origin` and check each hit
  is deliberate

### 7. The commands for the proposed history
- **Status**: Pending
- **Goal**: The user can create the past releases in one go
- **Boundary**: `specs/0027-site-version/history.md`, a script under
  `scripts/` that prints commands
- **Key files**: `specs/0027-site-version/history.md`
- **Depends**: None
- **Done when**: a script prints, and does not run, a `git tag` command
  for each release in `history.md` on its commit, and a
  `gh release create` command with its summary as notes, in order;
  it checks that each commit is on `main`; the user runs them
- **Verify**: run it and compare the output with the table

## Progress Log
