# Tasks: Stats Dashboard

## Tasks

### 1. Create stats-fetch CLI skeleton
- **Status**: Done
- **Goal**: New binary that connects to AWS, accepts
  `--since` flag, reads/writes cursor file
- **Boundary**: `cmd/stats-fetch/main.go`, `go.mod`
- **Key files**: `cmd/generate-versions/main.go` (CLI
  pattern), `internal/analytics/analytics.go` (table name)
- **Depends**: None
- **Done when**: `go build -o bin/stats-fetch ./cmd/stats-fetch/`
  compiles; running with no AWS creds prints error and exits
- **Verify**: `go build -o bin/stats-fetch ./cmd/stats-fetch/`

### 2. Fetch DynamoDB download events
- **Status**: Done
- **Goal**: stats-fetch scans `openkata-download-events`
  incrementally, appends to `.local/stats/download-events.json`
- **Boundary**: `cmd/stats-fetch/main.go`
- **Key files**: `internal/analytics/analytics.go` (table
  name, event fields)
- **Depends**: 1
- **Done when**: Running `go run ./cmd/stats-fetch/` creates
  `.local/stats/download-events.json` with events and updates
  cursor; second run fetches only new data
- **Verify**: `go build -o bin/stats-fetch ./cmd/stats-fetch/`

### 3. Fetch CloudWatch Lambda invocation metrics
- **Status**: Done
- **Goal**: stats-fetch queries `AWS/Lambda` Invocations
  metric for `openkata-web`, writes daily data to
  `.local/stats/page-metrics.json`
- **Boundary**: `cmd/stats-fetch/main.go`, `go.mod`
  (add `cloudwatch` SDK)
- **Key files**: None (new SDK usage)
- **Depends**: 1
- **Done when**: `page-metrics.json` contains daily
  invocation counts; incremental on re-run
- **Verify**: `go build -o bin/stats-fetch ./cmd/stats-fetch/`

### 4. Fetch CloudWatch Logs Insights per-path data
- **Status**: Done
- **Goal**: stats-fetch runs Logs Insights query on
  `/aws/lambda/openkata-web` for per-path request counts,
  writes daily entries to `.local/stats/page-paths.json`
- **Boundary**: `cmd/stats-fetch/main.go`, `go.mod`
  (add `cloudwatchlogs` SDK)
- **Key files**: None (new SDK usage)
- **Depends**: 1
- **Done when**: `page-paths.json` contains daily per-path
  entries; incremental on re-run
- **Verify**: `go build -o bin/stats-fetch ./cmd/stats-fetch/`

### 5. Add Makefile target
- **Status**: Done
- **Goal**: `make stats-fetch` runs the CLI
- **Boundary**: `mk/dev.mk`
- **Key files**: `mk/dev.mk` (existing targets), `Makefile`
  (help text)
- **Depends**: 1
- **Done when**: `make stats-fetch` works
- **Verify**: `make -n stats-fetch`

### 6. Stats page template and handler (empty state)
- **Status**: Done
- **Goal**: New `/stats/` route (local-only) that reads
  `.local/stats/` JSON files and renders page. Shows
  "No data" hint when files missing. New types for stats
  data model.
- **Boundary**: `cmd/openkata-web/templates/stats.templ`,
  `cmd/openkata-web/templates/types.go`,
  `cmd/openkata-web/handlers.go`,
  `cmd/openkata-web/main.go`
- **Key files**: `templates/design-system.templ` (local-only
  pattern), `handlers.go` line 166 (handleDesignSystem)
- **Depends**: None
- **Done when**: `/stats/` renders empty state; with sample
  JSON files renders tables (total, per-artifact, per-type,
  per-client, per-country)
- **Verify**: `go build -o bin/openkata-web ./cmd/openkata-web/`

### 7. Page view statistics section
- **Status**: Done
- **Goal**: Add page-view section to stats page — total
  loads, per-path table with type column ("page"/"download")
- **Boundary**: `cmd/openkata-web/templates/stats.templ`,
  `cmd/openkata-web/handlers.go`
- **Key files**: Task 6 output
- **Depends**: 6
- **Done when**: Page shows page-load total and path table
  from `page-metrics.json` and `page-paths.json`
- **Verify**: `go build -o bin/openkata-web ./cmd/openkata-web/`

### 8. Chart.js integration and download charts
- **Status**: Done
- **Goal**: Vendor Chart.js in `web/static/js-local/`,
  render two download chart groups (short-range:
  day/week/month, long-range: month/quarter/year) with
  granularity toggles
- **Boundary**: `web/static/js-local/chart.min.js`,
  `cmd/openkata-web/templates/stats.templ`
- **Key files**: `web/static/js/` (existing JS patterns)
- **Depends**: 6
- **Done when**: Charts render from download-events data;
  toggles switch granularity
- **Verify**: `go build -o bin/openkata-web ./cmd/openkata-web/`

### 9. Page loads chart
- **Status**: Done
- **Goal**: Add page-loads-per-day line chart below the
  page-view section
- **Boundary**: `cmd/openkata-web/templates/stats.templ`
- **Key files**: Task 8 output (chart patterns)
- **Depends**: 7, 8
- **Done when**: Page-loads chart renders from
  `page-metrics.json` data
- **Verify**: `go build -o bin/openkata-web ./cmd/openkata-web/`

### 10. Artifact detail with HTMX and version dropdown
- **Status**: Done
- **Goal**: Clicking an artifact loads detail partial via
  `hx-get="/stats/detail?artifact=X"`. Version dropdown
  refetches with `?version=Y`. Detail shows filtered
  stats and charts for that artifact/version.
- **Boundary**: `cmd/openkata-web/templates/stats.templ`,
  `cmd/openkata-web/handlers.go`,
  `cmd/openkata-web/main.go`
- **Key files**: Task 8 output (chart rendering)
- **Depends**: 8
- **Done when**: Artifact click shows detail; version
  dropdown filters data; charts update
- **Verify**: `go build -o bin/openkata-web ./cmd/openkata-web/`

### 11. Tests
- **Status**: Done
- **Goal**: Tests for stats handler (empty state, with
  data), stats-fetch cursor logic
- **Boundary**: `cmd/openkata-web/handlers_test.go`,
  `cmd/stats-fetch/main_test.go`
- **Key files**: `cmd/openkata-web/handlers_test.go`
  (existing test patterns)
- **Depends**: 6, 10
- **Done when**: `go test ./cmd/openkata-web/... ./cmd/stats-fetch/...`
  passes
- **Verify**: `go test ./cmd/openkata-web/... ./cmd/stats-fetch/...`

## Progress Log

- [2026-06-18] All 11 tasks confirmed complete. stats-fetch
  CLI fetches DynamoDB events, CloudWatch metrics, and Logs
  Insights paths. Stats page renders with Chart.js, HTMX
  detail drilldown, and tests passing.
