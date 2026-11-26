# Validation: Static Asset Delivery and Diagram Rendering

Checked the branch `feature/0025-static-assets-and-diagrams` against
`spec.md` alone. The checks ran on 2026-11-26.

- `go test ./cmd/openkata-web/... ./internal/...` passes
- `go vet ./cmd/openkata-web/... ./internal/...` reports nothing
- `go build -o bin/openkata-web ./cmd/openkata-web/` succeeds
- The committed `_templ.go` files call `assets.URL` where the
  `.templ` sources do

Verdict: **Pass**. Every requirement is met. The deviations below
support the spec and do not conflict with it.

## Requirements Check

### Part 1: Static asset delivery

- [x] `web/static` stays embedded and S3 is not used for assets.
  `main.go:72` builds the server with
  `assets.NewEmbedded(static.FS)`. Nothing in the diff adds an S3
  origin or bucket for assets.
- [x] A `/static/*` CloudFront behavior caches static files at the
  edge. `infra/add-static-cache-behavior.sh` creates the
  `openkata-static` cache policy and a `/static/*` behavior on the
  `lambda` origin. The cache key holds only the `v` query string,
  and the TTLs come from `Cache-Control` (default 300 s, max one
  year). The account owner applied this outside the diff, so the
  live distribution was not checked against the script. A read-only
  `curl -sI https://openkata.dev/static/favicon.svg` returned
  `x-cache: Miss`, then `x-cache: Hit` with `age: 2`, so a cached
  `/static/*` behavior is live. The live response has no
  `Cache-Control` yet, because this branch is not deployed.
- [x] Static responses carry `Cache-Control` and an `ETag`, and a
  browser can revalidate with a `304`. `internal/assets/assets.go:124`
  sets the ETag to the file hash, and lines 125-132 set
  `Cache-Control`. `http.ServeContent` answers `If-None-Match`.
  Tests: `TestServeHTTP/matching_ETag_gets_304` and
  `TestServeHTTP/live_matching_ETag_gets_304`.
- [x] Template URLs carry a per-file content hash (`?v=<hash>`),
  are cached for a year, and a changed file reaches visitors on the
  next page load. `Server.URL` (`assets.go:74`) adds `?v=` with the
  first 10 hex characters of the SHA-256 of the file. A request
  whose `v` matches the current hash gets
  `public, max-age=31536000, immutable`. A changed file gets a new
  hash and so a new URL. An old or wrong hash gets the 5-minute
  value, so a mismatch is never cached for a year. All templates
  use `assets.URL`. `TestStaticLinksCarryHashes` fails on any
  literal `/static/` URL in the templates package. It checks the
  `.templ` files and the hand-written `.go` files. It also fails on
  a named file that does not exist.
  Tests: `TestURL`, `TestPackageURL`,
  `TestServeHTTP/current_hash_is_cached_for_a_year`,
  `TestServeHTTP/old_hash_is_cached_briefly`.
- [x] Files without a hash reach visitors within 5 minutes. They
  get `public, max-age=300` (`assets.go:20`, `cacheShort`). This
  covers the CSS `@import`s in `web/static/css/style.css`,
  `fetch('/static/versions.json')` in `web/static/js/search.js:54`
  and the icons in `web/static/site.webmanifest`. The CloudFront
  policy's default TTL is also 300 s. CloudFront sends an `Age`
  header, so a browser does not add its own 5 minutes on top of the
  edge copy's age. Test: `TestServeHTTP/no_hash_is_cached_briefly`.
- [x] Local development serves `web/static` from disk, and an edit
  shows on the next reload. `main.go:77` uses
  `assets.NewLive(os.DirFS("web/static"))`. Live mode hashes on each
  request and sends `Cache-Control: no-cache`. Tests:
  `TestLiveSeesEdits`,
  `TestServeHTTP/live_files_are_always_checked`.
- [x] Infrastructure changes are scripted under `infra/` like the
  existing scripts. `infra/add-static-cache-behavior.sh` follows the
  pattern of `add-mcp-to-cloudfront.sh`: it looks up the
  distribution, can be re-run safely and uses `jq` on the
  distribution config. `infra/README.md` adds it as step 6, and
  `infra/iam-admin-policy.json` grants the three cache-policy
  actions it needs.

### Part 2: Diagram rendering

- [x] Diagrams are SVG files next to the docs page, linked with
  image syntax and descriptive alt text. The four files are in
  `web/docs/architecture/agent-workflow/`. They are linked from
  `web/docs/architecture/agent-workflow.md` with full-sentence alt
  text. No Makefile, CI or Go code runs
  `generate_diagrams.py`, so nothing in the build generates them.
- [x] The renderer inlines a linked local `.svg` file.
  `docs_render.go` collects images where `isLocalSVG` is true and
  `inlineDiagram` replaces the paragraph with a
  `<figure class="docs-diagram">` that holds the `<svg>` element.
  External, site-absolute and non-SVG images stay `<img>`. Tests:
  `TestRenderDiagrams/local_svg_is_inlined_in_a_figure`,
  `TestRenderDiagrams/unchanged/*`.
- [x] SVG colors use design tokens with fallbacks for both color
  schemes. Each file uses `var(--color-surface-raised, ...)`,
  `--color-surface`, `--color-ink`, `--color-ink-muted`,
  `--color-accent-muted` and `--font-body`. All of these are defined
  in `web/static/css/themes/japandi.css`, `themes/indigo.css` and
  `tokens.css`. Each file repeats the color rules in a
  `@media (prefers-color-scheme: dark)` block with the indigo
  values as fallbacks. No color is set outside these rules.
- [x] The four `mermaid` diagrams in `agent-workflow.md` are
  redrawn. The diff removes all four `mermaid` blocks and adds four
  image links: `feature-development.svg`,
  `parallel-implementation.svg`, `agent-responsibilities.svg` and
  `routing.svg`.
- [x] The strict loader rejects a page that links a missing SVG.
  `inlineDiagram` returns `errDocsLink` when the file cannot be
  read, and `buildSite` returns that error. Test:
  `TestRenderDiagrams/error/missing_file`.
- [x] No Mermaid remains. `grep -rni mermaid cmd web internal`
  finds only the sentence "Keep no Mermaid or other source in the
  file" in `web/docs/contributing/conventions.md:178`. The
  `mermaid-source` rewrite in `renderPage` and its tests are
  removed. A `mermaid` block now renders as any code block. The
  SVG files hold no Mermaid source.
- [x] visual-explainer is not a dependency or a required step. It
  does not appear in the code, the build or the contributor docs.

## Out of Scope Verified

- HTML pages and docs Markdown are not moved to S3
- Static files are not served from S3
- The only CloudFront change is one cache behavior and its cache
  policy. The Lambda architecture is unchanged
- visual-explainer is not adopted
- The ASCII diagram in the "Task Parallelization Example" `text`
  block is unchanged
- Nothing generates diagrams from Mermaid source. Spec 0030 is
  added as a Draft for that work, with a start condition

## Deviations

None of these conflict with the spec.

- The docs loader now allows a directory in a section that is
  named after one of its pages (`checkPageFileDirs` in `docs.go`).
  The directory may hold no pages or directories. Before, every
  directory in a section was an error. This is needed to put
  diagrams next to their page. Tests: the four new
  `TestBuildSite` error cases and
  `TestBuildSite/file_directory_named_after_a_page`.
- The renderer adds rules the spec does not state. A diagram must
  be alone in its paragraph, must stay inside `web/docs` and must
  start with `<svg`. An XML declaration is removed. Each breaks
  the build with `errDocsLink`.
- `TestDiagramTitlesMatchAltText` requires each SVG `<title>` to
  match the image's alt text word for word.
- Local development no longer lists directories or serves
  `index.html` for a directory, as `http.FileServer` did. It returns
  404. It now also sends `Cache-Control: no-cache` and an `ETag`.
  Files are still served from disk, as the spec requires.
- `web/docs/contributing/conventions.md` gains a "Diagrams" section
  with SVG conventions and a template. It links to
  `generate_diagrams.py` as an optional starting point for new
  diagrams.
- `specs/0025-static-assets-and-diagrams/generate_diagrams.py` is
  committed. It made the first versions of the SVGs. It is not part
  of the build, so the requirement that nothing in the build
  generates the SVGs still holds.
- `specs/0030-mermaid-diagrams/spec.md` is a new Draft spec.
- `spec.md` front matter changes from `Draft`/`Standard` to
  `Implementing`/`Deep`.
- `infra/iam-admin-policy.json` adds `ListCachePolicies`,
  `GetCachePolicy` and `CreateCachePolicy`.
- `web/static/css/components/docs.css` adds `.docs-diagram` styles.

## Issues Found

No blocking issues. Notes for the implementer:

- The live CloudFront distribution was not checked against the
  script, because the account owner applied it outside the diff.
  Confirm that the live `/static/*` behavior uses the
  `openkata-static` policy and not a managed policy. If it uses a
  policy without `v` in the cache key, hashed URLs would not bust
  the edge cache.

  Added after validation: checked on the live site with `curl -sI`
  on `/static/favicon.svg`. A new `v` value was a `Miss`, and the
  same `v` with an extra `x` parameter was a `Hit`. So `v` is in
  the cache key and nothing else is, which is the `openkata-static`
  policy.
- "The Lambda only sees the first request for each URL" holds for
  hashed URLs. Unhashed URLs reach the Lambda once per edge every
  5 minutes, as the 5-minute requirement implies. Those requests
  are revalidations, so most get a `304`.
- `stats.templ` links `js-local/chart.min.js`, which is not in the
  `//go:embed` list in `web/static/embed.go`. The stats page only
  runs locally, where the file is served from disk, so this has no
  effect today. `TestStaticLinksCarryHashes` checks the file exists
  on disk, not in the embedded set.
