# Validation: Documentation Content and Go-Live

Validated against `spec.md` only (no task plan or progress log read),
the current code on branch `feature/0024-docs-content`, and the
docs-related diff against `main`.

## Method and results of the checks run

- `go test ./cmd/openkata-web/... ./internal/docsdates/
  ./cmd/docs-dates/`: all packages pass (results were cached).
  `go vet` on the same packages: no output.
- `go run ./cmd/docs-dates -dir web/docs -out <scratchpad>/dates.json`
  wrote 23 dates with no warning. The same run with `-require`
  also succeeded. 23 pages exist under `web/docs/` (excluding
  `generated/`), so every page currently has a git date.
- `git check-ignore` confirms `web/docs/generated/dates.json` is
  ignored (`.gitignore:15`, `web/docs/generated/*`), so it is not
  committed. `web/docs/generated/.gitkeep` is the only tracked file
  there.
- Nothing was written into the repository except this report.

## Requirements Check

### Pages: Guides

- [x] Manifesto page pointing to `MANIFESTO.md` — met.
  `web/docs/guides/manifesto.md:13` links to the file on GitHub
  and states it is the only copy (lines 9-11). Content is a
  short summary, not a copy.
- [x] Skills, rules and profiles, with kata vocabulary — met, with
  one inaccurate claim (see finding F1).
  `web/docs/guides/skills-rules-profiles.md:15-20` covers dojo kun and
  sensei profile. The terms the dojo and the ryu match `README.md:27-28`.
  I did not confirm that the page covers "the dojo" and "the ryu"
  beyond the table I read (lines 1-60); the rest of the page was not
  read line by line.
- [x] Installing artifacts through the MCP server — page exists
  (`web/docs/guides/installing-artifacts.md`), loads under the strict
  loader test. Content not fact-checked.
- [x] Using the catalog — page exists
  (`web/docs/guides/using-the-catalog.md`). Content not fact-checked.

### Pages: Reference

- [x] MCP discovery — met.
  `web/docs/reference/mcp-discovery.md` checked against
  `cmd/openkata-mcp/main.go`: nine tools registered
  (`main.go:69-77`) and the table matches; stateless streamable HTTP
  (`main.go:79-81`); tag filter is exact, case-sensitive, trimmed,
  one tag (`main.go:515-522`); versions loaded once at startup, empty
  index on failure (`main.go:58-63`); error strings
  `version %s not found for %s` (`main.go:304`) match.
- [x] Storage and `versions.json` generation, with S3 layout — met.
  `web/docs/reference/storage.md` against
  `.github/workflows/publish.yaml:5-8,52-66` (tag forms, `git archive`,
  `aws s3 sync`, regenerate), `publish-on-merge.yaml` (re-push of tags
  in the merged PR), `cmd/generate-versions/bucket.go:66,79-94` (score
  and published unless private) and `main.go:136-140` (`0.0.0`).
  Table names `openkata-download-events`
  (`internal/analytics/analytics.go:13`) match. The page notes
  infrastructure is public in `infra/`; I did not open `infra/`
  to confirm the bucket name, the four public-access blocks or the IAM
  role names, so those claims are not verified.
- [x] Downloads — met. `web/docs/reference/downloads.md` against
  `cmd/openkata-web/handlers.go:1243-1249` (source `web`, client from
  User-Agent, country from `CloudFront-Viewer-Country`),
  `handlers.go:1371-1385` (atomic `ADD downloads`, error ignored),
  `cmd/openkata-mcp/main.go:326-332` (source `mcp`, client `unknown`,
  empty country), `internal/analytics/analytics.go:47-48` (event write
  error is logged), and the client order in
  `analytics.go:55-63`. `/stats/` is not documented in the page, as
  required.
- [x] Evals and the Docker sandbox — met. `web/docs/reference/evals.md`
  against `cmd/openkata-eval/config.go:81-89` (model
  `claude-sonnet-4.6`, threshold 95, direct timeout 120, concurrency 2,
  image `openkata-eval:latest`, network `host`),
  `main.go:26` (`--threshold`), `Dockerfile.eval:1-5` (Ubuntu 24.04,
  kiro-cli install), `mk/dev.mk:35-42` (`eval-local`, `eval-image`).
  I did not read the sandbox/judge code paths (retry timings, 20
  character rule), so those claims are not verified.
- [x] Scoring (effectiveness and tessl score) — met.
  `web/docs/reference/scoring.md`: registry URL
  `https://tessl.io/registry/openkata/<name>` matches
  `cmd/openkata-web/templates/skill_detail.templ:106-114`; the
  `private` flag handling matches `cmd/generate-versions/main.go:64-68`.
  Bar colour thresholds (95/70) were not verified.
- [x] Tessl — page exists (`web/docs/reference/tessl.md`); not
  fact-checked.
- [x] Versioning and `.openkata-manifest.json` — met, with one
  inaccurate claim (F1 does not apply here).
  `web/docs/reference/versioning.md`: manifest fields and checksums match
  `cmd/openkata-web/handlers.go:1196-1209` and
  `cmd/openkata-mcp/main.go:309-320`; excluded files match
  `internal/exclude/exclude.go:7-12` and its use
  (`handlers.go:1172`, `main.go:430`).

### Pages: Architecture

- [x] Agent workflow — `web/docs/architecture/agent-workflow.md` exists.
- [x] Infrastructure diagram —
  `web/docs/architecture/infrastructure-diagram.md` exists. Content not
  checked against the real infrastructure.
- Moves from spec 0023 are described by the spec as already done; not
  re-verified.

### Pages: Contributing

- [x] Authoring skills, rules and profiles —
  `web/docs/contributing/authoring.md` exists.
- [x] Release workflow, one page each for site, artifact, MCP server —
  `release-site.md`, `release-artifact.md`, `release-mcp.md` exist.
  `release-site.md` checked: templ version `v0.3.1001`
  (`.github/workflows/build.yaml:18`, `deploy.yaml:51`); Build steps
  (`build.yaml:17-32`); tag-on-head-of-main rule
  (`deploy.yaml:26-35`); deploy-web script order
  (`scripts/deploy-web.sh:7-26`); full history
  (`deploy.yaml:22`); `make deploy web|mcp`, `make push`,
  `make badges`, `make versions` (`mk/dev.mk`, `mk/deploy.mk:7`);
  `make check` covers git, go, make, tessl (`scripts/check-prereqs.sh`)
  and not templ or air. All match. `release-mcp.md` spot-checked
  against `scripts/deploy-mcp.sh:16-25`: matches. `release-artifact.md`
  was not fact-checked.
- [x] Repository conventions — `web/docs/contributing/conventions.md`
  exists; not fact-checked.
- [x] Docs home and section indexes exist (`web/docs/index.md` and
  `index.md` in each of the four sections). The strict-loader test
  `go test ./cmd/openkata-web/...` passes, so front matter, structure
  and relative links in the embedded tree load.
- Section-per-PR note: informational, no check needed.

### Last-Updated Date

- [x] Shown on every docs page, including home and section indexes,
  below the previous/next links — met.
  `cmd/openkata-web/templates/docs.templ:19` (home) and `:41-42`
  (pages: `docsPager` then `docsUpdated`). Section indexes use the
  same `DocsPage` template (not separately confirmed, but the handler
  uses the same `updatedFor` at `docs_handler.go:99`). Format
  `2 January 2006` at `docs_handler.go:115`
  (`Last updated <time>25 October 2026</time>`).
- [x] Date is author date of the last commit touching the file, not
  written in the page — met. `internal/docsdates/docsdates.go:62-66`
  uses `git log --format=%as --name-only`, newest first, first
  sighting wins; `--relative`, `-M` handles renames. Front matter in the
  pages I read contains no date field.
- [x] Local: read from git; placeholder otherwise, same place and
  style — met. `docs_handler.go:23` calls `docsdates.Lookup`;
  `docs.templ:51-58` renders `Last updated: <span ...>not committed
  yet</span>` in the same `<footer class="docs-updated">`. Lookup drops
  untracked files and reused paths (`docsdates.go:90-107`). A shallow
  clone, no git or no repository yields an empty map plus an error that
  the local handler ignores (`docs_handler.go:21-23`), so locally the
  placeholder shows.
- [x] Production: dates stamped at deploy time from full history,
  embedded, not committed, no CI push — met.
  `scripts/deploy-web.sh:10-13` runs `go run ./cmd/docs-dates -require`
  before the build; `cmd/docs-dates/main.go:61-68` returns an error on
  any page without a date when `-require` is set;
  `main.go:41-43` fails on a shallow clone with a hint about
  `fetch-depth: 0`; `.github/workflows/deploy.yaml:22` sets
  `fetch-depth: 0`. `web/docs/embed.go:10` embeds `all:generated`;
  `.gitignore:15` ignores the output. No workflow commits or pushes
  dates (grep of `.github/workflows` shows only the existing
  tag-retrigger push in `publish-on-merge.yaml`).
  Caveat: `siteFromFS` (`docs.go:211-234`) treats a missing
  `dates.json` as "no dates, not an error", so the guarantee that
  production never shows the placeholder rests entirely on the deploy
  script. It is also moot today because `/docs/` is not served in
  production (see Go-Live).
- Not verified: the content of the `docs-dates` and `docsdates` tests
  beyond that they pass.

### Go-Live

- [ ] Go-live waits for spec 0026 (benchmarks and version list live) —
  not built and not decided by this feature. The spec itself says
  go-live must wait for this. I did not inspect the status of spec
  0026 (not an allowed input), so I cannot say whether it is done.
- [ ] Go-live waits for spec 0025 (diagram rendering and library
  hosting) — same: the spec says it must wait; not verified here.
  The pages themselves are allowed not to wait, per the spec.
- [ ] `/docs/` moved out of the local-only block and served from the
  embedded tree, with the `Site` built once and cached — NOT BUILT.
  `cmd/openkata-web/main.go:86-92` still registers `/docs/` and sets
  `templates.DocsEnabled = true` inside
  `if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == ""`. The function
  `loadEmbeddedSite` (`docs_handler.go:30`) exists but nothing outside
  tests calls it. `docsHandler` calls `load()` on every request
  (`docs_handler.go:49`), so no once-only caching exists yet, and the
  local handler rebuilds on each request by design (`:17-25`).
- [ ] "Docs" link in the main nav and the mobile nav menu — NOT BUILT.
  `cmd/openkata-web/templates/nav.templ:11-14` and `:26-29` list only
  Getting Started, Skills, Rules and Profiles. A "Documentation" link
  exists only in the footer, and only when `DocsEnabled`
  (`templates/footer.templ:25-27`), which is not required by the
  spec and is local-only.
- [ ] `search.js` loads `/docs/search-index.json`, shows docs results,
  and `/catalog/` filters them out — NOT BUILT. `grep` for `docs` and
  `search-index` in `web/static/js/search.js` finds nothing. The
  server side exists: `docs_handler.go:15,52-53` serves the index at
  that path (local only).

## Out of Scope Verified

- Confirmed: the local-only `/stats/` dashboard is not documented in
  `web/docs/reference/downloads.md`.
- Spec 0023 engine work and specs 0025/0026 work are not part of this
  feature's page content; I did not look for them.

## Issues Found

- F1 (low): Inaccurate claim in
  `web/docs/guides/skills-rules-profiles.md:23-26`: "Every artifact
  folder also holds a `CHANGELOG.md` and a `RATIONALE.md`". Three
  artifact folders have no `RATIONALE.md`: `skills/create-evals`,
  `skills/critical-thinking`, `rules/go-testing`. All have a changelog.
- F2 (high, Go-Live): `/docs/` is still local-only, no main-nav Docs
  link, and `search.js` does not load the docs index. See Go-Live
  above. These items are unbuilt, and the spec says go-live must wait
  for specs 0025 and 0026, so they may be intentionally open.
- F3 (medium): Production caching of the `Site` ("built once and
  cached") is not implemented, because `docsHandler` rebuilds via the
  load function on every request. `loadEmbeddedSite` returns a fresh
  site each call, with no `sync.Once`.
- F4 (low): `siteFromFS` (`cmd/openkata-web/docs.go:211-234`) accepts a
  missing `generated/dates.json`, so a binary built without the deploy
  script would show the placeholder. The spec requires production never
  to show it; this is enforced only by `scripts/deploy-web.sh`. The
  build workflow (`go build ./...`) does not stamp dates, which is
  fine since it does not deploy.
- F5 (info): The spec status is still `Draft` (`spec.md:2`).

## Not verified

- Page contents not fact-checked: `guides/installing-artifacts.md`,
  `guides/using-the-catalog.md`, `reference/tessl.md`,
  `contributing/authoring.md`, `contributing/conventions.md`,
  `contributing/release-artifact.md`, the architecture pages.
- Infrastructure facts (bucket name, public-access blocks, IAM role
  names, DynamoDB on-demand) in `storage.md` and `downloads.md`: `infra/`
  not read.
- Eval runner sandbox and judge details (retry waits, 20 character
  suspicious-result rule, snapshot diff).
- Real rendering in a browser; section index pages' footer was
  confirmed through templates and the shared handler only.
- Status of specs 0025 and 0026.
- The tests ran from cache; I did not force a re-run.

## Summary counts

Requirements counted by spec bullet group: pages met 16 of 16 (two
with accuracy caveats, F1, and unverified content in several); date
behavior met 4 of 4 (with caveat F4); Go-Live: 0 of 5 built (2 depend
on specs 0025 and 0026 per the spec itself, 3 are unbuilt feature work).
