# Tasks: Documentation Content and Go-Live

Every page task has the same shape and the same checks:

- Frontmatter has `title`, `description` and `order`
- Each fact is checked against the source files listed in the
  task, not written from memory
- Links to other docs pages are relative `.md` links; links to
  files outside `web/docs/` are absolute GitHub URLs on `main`
- Prose wraps at 72-80 characters, per `rules/markdown-style`
- **Verify**: `go test ./cmd/openkata-web/ -run Embedded` (the
  strict loader rejects bad frontmatter, bad links and layouts)
  and a look at the page on the local server

A section's first page task also replaces that section's "coming
soon" index stub with a short introduction. The index does not
list pages: the sidebar does, and links to pages that do not exist
yet would fail the strict loader. Pages ship one section per PR.

## Tasks

### 1. Guide: manifesto
- **Status**: Done
- **Goal**: Developers find the manifesto from the docs; the root
  file stays the single source
- **Boundary**: `web/docs/guides/manifesto.md`,
  `web/docs/guides/index.md`
- **Key files**: `MANIFESTO.md`, the "See the Manifesto" link in
  `README.md`
- **Depends**: None
- **Done when**: `MANIFESTO.md` stays at the repository root,
  unchanged; `manifesto.md` is a short page that says what the
  manifesto is and links to it with an absolute GitHub URL on
  `main`, without copying its text; the guides index stub is
  replaced with a short introduction that also links to the
  manifesto page
- **Verify**: `go test ./cmd/openkata-web/ -run Embedded`

### 2. Guide: skills, rules and profiles
- **Status**: Pending
- **Goal**: A reader learns what each artifact type is, how they
  differ, and the kata vocabulary
- **Boundary**: `web/docs/guides/skills-rules-profiles.md`
- **Key files**: `cmd/openkata-web/templates/getting_started.templ`
  (vocabulary table), `docs/adr/0004-adopt-kata-vocabulary.md`,
  `docs/adr/0001-use-agent-skills-format.md`
- **Depends**: None
- **Done when**: covers skills, rules and profiles with an example
  of each, plus dojo kun, sensei profiles, the dojo and the ryu

### 3. Guide: installing artifacts through the MCP server
- **Status**: Pending
- **Goal**: A reader can connect an agent to the dojo and install
  an artifact
- **Boundary**: `web/docs/guides/installing-artifacts.md`
- **Key files**: `cmd/openkata-mcp/main.go` (tools),
  `docs/adr/0002-use-mcp-server-for-distribution.md`,
  `docs/adr/0010-deploy-mcp-to-lambda.md`
- **Depends**: 2
- **Done when**: shows connecting a client, listing, installing a
  skill, rule and profile, and choosing a version

### 4. Guide: using the catalog
- **Status**: Pending
- **Goal**: A reader can find, filter, sort and inspect artifacts
  on the site
- **Boundary**: `web/docs/guides/using-the-catalog.md`
- **Key files**: `cmd/openkata-web/templates/catalog.templ`,
  `skills.templ`, `skill_detail.templ`, `web/static/js/search.js`
- **Depends**: 2
- **Done when**: covers the catalog, tag filter, sorting, search,
  the detail page tabs and the download counts shown

### 5. Reference: how MCP discovery works
- **Status**: Pending
- **Goal**: Explain the `list_*` tools, the tag filter and where
  the index lives
- **Boundary**: `web/docs/reference/mcp-discovery.md`,
  `web/docs/reference/index.md`
- **Key files**: `cmd/openkata-mcp/main.go`,
  `cmd/generate-versions/main.go`
- **Depends**: None
- **Done when**: lists every tool with its arguments as the code
  defines them, explains tag filtering and `versions.json`

### 6. Reference: where artifacts are stored
- **Status**: Pending
- **Goal**: Document the S3 bucket layout and how `versions.json`
  is generated
- **Boundary**: `web/docs/reference/storage.md`
- **Key files**: `infra/README.md`, `infra/openkata-stack.yaml`,
  `cmd/generate-versions/main.go`, `.github/workflows/`,
  `docs/adr/0005-version-distributed-artifacts.md`
- **Depends**: 5
- **Done when**: shows the bucket layout, what writes to it and
  when, and how versions come from git tags

### 7. Reference: downloads
- **Status**: Pending
- **Goal**: Explain how downloads are recorded and counted, and
  the public counts
- **Boundary**: `web/docs/reference/downloads.md`
- **Key files**: `cmd/openkata-web/handlers.go`
  (`incrementDownload`, archive handlers), `internal/analytics/`,
  `docs/adr/0010-deploy-mcp-to-lambda.md`
- **Depends**: 6
- **Done when**: covers what counts as a download, the counters,
  and where counts show on the site; does not document the
  local-only `/stats/` dashboard

### 8. Reference: how we use evals
- **Status**: Pending
- **Goal**: Explain evals and the local runner with its Docker
  sandbox
- **Boundary**: `web/docs/reference/evals.md`
- **Key files**: `cmd/openkata-eval/`,
  `.agents/skills/openkata-eval-runner/SKILL.md`,
  `skills/create-evals/SKILL.md`
- **Depends**: 5
- **Done when**: covers scenario structure, running locally, the
  sandbox, and where results are written and read

### 9. Reference: how we score skills
- **Status**: Pending
- **Goal**: Explain effectiveness and the tessl score
- **Boundary**: `web/docs/reference/scoring.md`
- **Key files**: `cmd/generate-versions/main.go`,
  `cmd/openkata-web/templates/skill_detail.templ` (Benchmarks tab),
  `docs/adr/0006-adopt-tessl-as-skill-quality-toolchain.md`
- **Depends**: 8
- **Done when**: defines both scores, how each is produced and
  shown, and the difference between them

### 10. Reference: how we use tessl
- **Status**: Pending
- **Goal**: Explain tessl reviews, evals and the registry
- **Boundary**: `web/docs/reference/tessl.md`
- **Key files**: `docs/adr/0006-adopt-tessl-as-skill-quality-toolchain.md`,
  `docs/adr/0012-use-tessl-evals-and-registry-for-distribution.md`,
  `.github/workflows/publish.yaml`
- **Depends**: 9
- **Done when**: covers reviews, evals, publishing to the
  registry and what the score means

### 11. Reference: versioning and the manifest
- **Status**: Pending
- **Goal**: Explain artifact versions and `.openkata-manifest.json`
- **Boundary**: `web/docs/reference/versioning.md`
- **Key files**: `RELEASING.md`,
  `docs/adr/0005-version-distributed-artifacts.md`,
  `docs/adr/0011-artifact-dependency-metadata.md`,
  `cmd/openkata-mcp/main.go`
- **Depends**: 6
- **Done when**: covers tag format, CHANGELOG versions, what the
  manifest records and when it is written

### 12. Contributing: authoring artifacts
- **Status**: Pending
- **Goal**: A contributor can write a skill, rule or profile that
  passes review
- **Boundary**: `web/docs/contributing/authoring.md`,
  `web/docs/contributing/index.md`
- **Key files**: `skills/create-skill/SKILL.md`,
  `skills/create-rule/SKILL.md`, `skills/create-profile/SKILL.md`,
  `.agents/skills/openkata-skill-conventions/SKILL.md`,
  `.agents/skills/openkata-rule-conventions/SKILL.md`, `AGENTS.md`
- **Depends**: 2
- **Done when**: covers the layout and limits for each type, the
  changelog and the review step

### 13. Contributing: release a new version of the site
- **Status**: Pending
- **Goal**: Describe how the site is built and deployed
- **Boundary**: `web/docs/contributing/release-site.md`
- **Key files**: `RELEASING.md`, `.github/workflows/build.yaml`,
  `.github/workflows/deploy.yaml`, `AGENTS.md` (before pushing)
- **Depends**: 12
- **Done when**: covers the checks before pushing, what CI runs
  and what a deploy does

### 14. Contributing: release a skill, rule or profile
- **Status**: Pending
- **Goal**: Describe tagging and publishing an artifact version
- **Boundary**: `web/docs/contributing/release-artifact.md`
- **Key files**: `RELEASING.md`,
  `.agents/skills/openkata-ryu-release/SKILL.md`,
  `.github/workflows/publish.yaml`,
  `.github/workflows/publish-on-merge.yaml`
- **Depends**: 12
- **Done when**: covers the version bump, `skills/<name>/v<version>`
  tags and what the publish workflows do

### 15. Contributing: release the MCP server
- **Status**: Pending
- **Goal**: Describe building and deploying the MCP server
- **Boundary**: `web/docs/contributing/release-mcp.md`
- **Key files**: `RELEASING.md`, `infra/create-mcp-stack.sh`,
  `.github/workflows/deploy.yaml`,
  `docs/adr/0010-deploy-mcp-to-lambda.md`
- **Depends**: 12
- **Done when**: covers the build, the Lambda deploy and how to
  check it

### 16. Contributing: repository conventions
- **Status**: Pending
- **Goal**: Summarize the conventions a contributor must follow
- **Boundary**: `web/docs/contributing/conventions.md`
- **Key files**: `AGENTS.md`, `rules/git-naming/RULE.md`,
  `rules/markdown-style/RULE.md`, `rules/go-testing/RULE.md`,
  `rules/bash-style/RULE.md`
- **Depends**: 12
- **Done when**: covers branches, commits and the `Assisted-by`
  trailer, code style, testing, and the no-squash merge rule

### 17. Go-live
- **Status**: Pending
- **Goal**: Make the docs public
- **Boundary**: `cmd/openkata-web/main.go`,
  `cmd/openkata-web/docs_handler.go`,
  `cmd/openkata-web/templates/nav.templ`,
  `web/static/js/search.js`, `cmd/openkata-web/handlers_test.go`
- **Key files**: `docs_handler.go` (`handleDocs`),
  `web/static/js/search.js`, `docs_embedded_test.go`
- **Depends**: 1-16, and spec 0025's diagram decision
- **Done when**: `/docs/` is registered outside the local-only
  block and serves `docs.FS` with the site built once and cached;
  the main and mobile nav have a Docs link; the search overlay
  shows docs results from `/docs/search-index.json` and
  `/catalog/` filters them out; tests cover the production path
  and the search change
- **Verify**: `templ generate ./cmd/openkata-web/templates/`,
  `go build -o bin/openkata-web ./cmd/openkata-web/`,
  `go test ./cmd/openkata-web/...`

## Progress Log
