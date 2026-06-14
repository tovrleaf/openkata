# Validation: Stats Dashboard

Date: 2026-06-14

## Requirements Check

### Prerequisite: Remove Existing Stats Page

- [x] Old types removed — no `ArtifactDownloads`, `ClientDownloads`,
  `CountryDownloads` found anywhere in codebase
- [x] New `StatsData`, `ArtifactStats`, `ClientStats`, `CountryStats`
  types in `types.go` — redesigned for the new data model
- [x] Route registration stays, re-pointed to new handler — confirmed
  in `main.go` lines 63–64

### Data Sources

- [x] DynamoDB download events — `cmd/stats-fetch/main.go` scans
  `openkata-download-events` table with timestamp filter
- [x] CloudWatch metrics — fetches `AWS/Lambda` Invocations metric
  for `openkata-web` function
- [x] CloudWatch Logs Insights — runs per-path query on
  `/aws/lambda/openkata-web` log group

### Data Fetching CLI

- [x] Binary at `cmd/stats-fetch/main.go` — confirmed
- [x] `--since YYYY-MM-DD` flag with 30-day default — confirmed
- [x] Uses default AWS credential chain (`config.LoadDefaultConfig`)
  — confirmed
- [x] Incremental via cursor file at `.local/stats/cursor.json` —
  confirmed
- [x] Appends new data to existing JSON files — confirmed for
  `download-events.json`; metrics and paths use overwrite/merge
  (idempotent as spec allows)

### File Layout

- [x] `.local/stats/cursor.json` — confirmed
- [x] `.local/stats/download-events.json` — confirmed
- [x] `.local/stats/page-metrics.json` — confirmed
- [x] `.local/stats/page-paths.json` — confirmed
- [ ] Cursor field names — MINOR DEVIATION: spec shows
  `{ "events_after": "...", "metrics_after": "..." }` but
  implementation uses `{ "downloads": "...", "metrics": "...",
  "paths": "..." }`. Functional behavior is correct; field
  naming differs from spec comment.

### .gitignore

- [x] `.local/` present in `.gitignore` — confirmed

### Stats Page

- [x] Local-only route — registered inside
  `if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == ""` block,
  same pattern as `/design-system/`
- [x] Reads from `.local/stats/` disk files — confirmed in
  `handleStats`
- [x] Empty state shows "No data. Run `make stats-fetch` first."
  — confirmed in template

### Download Statistics

- [x] Total downloads — aggregate count displayed
- [x] Per artifact table — sorted by downloads desc, columns:
  artifact, type, downloads
- [x] Per type breakdown — skills/rules/profiles extracted from
  artifact prefix
- [x] Per client table — sorted by downloads desc
- [x] Per country table — sorted by downloads desc

### Artifact Detail (HTMX partial)

- [x] Clicking artifact row triggers `hx-get="/stats/detail?artifact=..."` —
  confirmed with `hx-target="#stats-detail"` and `hx-swap="innerHTML"`
- [x] `handleStatsDetail` handler returns partial template — confirmed
- [x] Version dropdown with `hx-get` and `hx-include` — filters to
  selected version
- [x] Filtered stats (clients, countries, events) for
  artifact/version — confirmed

### Download Charts (Chart.js)

- [x] Short-range chart group — day / week / month toggles with
  `renderShortChart()` function
- [x] Long-range chart group — month / quarter / year toggles with
  `renderLongChart()` function
- [x] Line charts — `type: 'line'` confirmed
- [x] Detail view includes filtered chart — `chart-detail` canvas in
  `StatsDetail` template filters to artifact/version events

### Page View Statistics

- [x] Total page loads — summed from `PageMetrics` invocations
- [x] Page loads per day chart — `chart-pages` canvas with metrics
  data
- [x] Per-path breakdown table — columns: path, type, count, sorted
  by count desc
- [x] Type column — "page" or "download" based on `/archive` in path
  — confirmed in `handleStats`

### Charting

- [x] Chart.js vendored at `web/static/js-local/chart.min.js` —
  file exists (205KB)
- [x] `js-local/` NOT embedded in production — `embed.go` directive
  lists `all:css all:js all:img` (embeds `js/` not `js-local/`)
- [x] Charts respect theme via CSS custom properties — reads
  `--color-accent` and `--color-surface` via `getComputedStyle`
- [x] Responsive sizing — Chart.js defaults to `responsive: true`;
  canvas elements in block-level containers

### Makefile Target

- [x] `make stats-fetch` target in `mk/dev.mk` — runs
  `go run ./cmd/stats-fetch/`

### AWS Auth

- [x] Default credential chain — `config.LoadDefaultConfig(ctx,
  config.WithRegion("eu-north-1"))` confirmed

## Out of Scope Verified

- Confirmed: no websocket/real-time code
- Confirmed: no IP/session tracking
- Confirmed: no production deployment of stats page (local-only
  route guard)
- Confirmed: no alerting code
- Confirmed: no CSV export
- Confirmed: no authentication on stats page
- Confirmed: no CloudFront log parsing

## Issues Found

1. **Cursor field names differ from spec** — spec comment shows
   `events_after` / `metrics_after`; implementation uses
   `downloads` / `metrics` / `paths`. The three-field design
   (separate cursor per data source) is arguably better than the
   spec's two-field example. Low severity — internal format only.

2. **`chart-toggles` class has no CSS definition** — the template
   uses `class="chart-toggles"` on wrapper divs, but no CSS rule
   exists for this class. The child buttons use `.tab` which is
   styled, so they render correctly, but the container lacks the
   flex layout that `.tabs` provides. The toggle buttons will
   stack vertically instead of appearing in a row. Low severity,
   cosmetic.

3. **No `chart-toggles` active state management in JS** — the
   `renderShortChart` function has a comment-like attempt to
   toggle active class but queries `#chart-short` (the canvas)
   and walks up to parent, which won't correctly find the
   `.chart-toggles` container buttons. `renderLongChart` has no
   active-class management at all. Low severity, cosmetic.
