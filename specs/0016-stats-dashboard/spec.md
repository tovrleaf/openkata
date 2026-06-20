---
status: Done
depth: Standard
---

# Stats Dashboard

## Story

As the site maintainer, I want a local-only stats page
that visualizes download analytics and page-view data so
I can understand usage patterns without logging into AWS
consoles.

## Requirements

### Prerequisite: Remove Existing Stats Page

Delete the prematurely built stats implementation:

- `handleStats` handler in `handlers.go`
- `templates/stats.templ` and `templates/stats_templ.go`
- `StatsData`, `ArtifactStats`, `ClientStats`,
  `CountryStats` types from `types.go`
- Route registration in `main.go` stays (re-pointed to
  new handler)

### Data Sources

Two data sources, both fetched to local disk on command:

1. **DynamoDB download events**
   (`openkata-download-events` table) — artifact,
   version, source, client, country, timestamp
2. **CloudWatch metrics/logs** — Lambda invocation
   counts (page loads) from `AWS/Lambda` metrics, and
   per-path request breakdown from CloudWatch Logs
   Insights on `/aws/lambda/openkata-web`

### Data Fetching CLI (`cmd/stats-fetch/`)

New binary at `cmd/stats-fetch/main.go`. Makefile target:

```makefile
.PHONY: stats-fetch
stats-fetch:
	@go run ./cmd/stats-fetch/
```

Behavior:

- Uses default AWS credential chain
  (`config.LoadDefaultConfig`)
- Accepts `--since YYYY-MM-DD` flag (default: 30 days
  ago on first run)
- Incremental: stores last-fetched timestamp in
  `.local/stats/cursor.json`; subsequent runs query
  only data after that timestamp
- Appends new data to existing JSON files on disk
- DynamoDB: filters by `timestamp > cursor`
- CloudWatch metrics: queries with `StartTime` = cursor
- Logs Insights: queries with start time = cursor,
  stores daily-granularity entries (idempotent — if a
  day exists, overwrite it)

### File Layout

```
.local/stats/
  cursor.json          # { "events_after": "...", "metrics_after": "..." }
  download-events.json # array of event objects
  page-metrics.json    # [{ "date": "...", "invocations": N }, ...]
  page-paths.json      # [{ "date": "...", "path": "...", "count": N }, ...]
```

Add `.local/` to `.gitignore`.

### Stats Page (`/stats/`)

Local-only route (same pattern as `/design-system/`).
Reads data from `.local/stats/` JSON files.

**Empty state**: if files don't exist, show
"No data. Run `make stats-fetch` first."

#### Download Statistics

- **Total downloads** — aggregate count
- **Per artifact** — table sorted by downloads desc
  (columns: artifact, type, downloads)
- **Per type** — breakdown by skills/rules/profiles
- **Per client** — table (Kiro, Claude-Desktop,
  Cursor, curl, browser, unknown)
- **Per country** — table sorted by downloads desc

#### Artifact Detail (HTMX partial)

Clicking an artifact loads a detail section via
`hx-get="/stats/detail?artifact=skills/create-adr"`.
Shows:

- Version dropdown — selecting a version refetches
  with `?version=1.2.0`, filtering all stats to that
  version only
- Filtered charts for that artifact/version

#### Download Charts (Chart.js)

Two chart groups on the main page:

1. **Short-range** — toggle between day / week / month
2. **Long-range** — toggle between month / quarter / year

Line charts. When viewing a specific artifact or
version in the detail view, charts filter to that
selection.

#### Page View Statistics

- **Total page loads** — from CloudWatch Lambda
  invocation metric
- **Page loads per day** — line chart
- **Per-path breakdown** — table with columns: path,
  type ("page" or "download" based on `/archive` in
  path), count — sorted by count desc

### Charting

- Chart.js vendored at `web/static/js-local/chart.min.js`
- `web/static/js-local/` is NOT embedded in production
  (not listed in `embed.go` directive)
- Charts respect active theme via CSS custom properties
- Responsive sizing

### Infrastructure

No infrastructure changes required. CloudWatch metrics
and logs already exist for the Lambda. DynamoDB events
table exists from spec 0015.

New AWS SDK dependency: `cloudwatchlogs` and
`cloudwatch` packages added to `go.mod`.

## Out of Scope

- Real-time streaming/websocket updates
- User session tracking or IP storage
- Deploying the stats page to production
- Alerting or notifications
- Data export (CSV, etc.)
- Authentication on the stats page
- CloudFront access log parsing

## Rejected Approaches

- **CloudFront access logs to S3** — requires enabling
  logging, an S3 bucket, and log parsing. CloudWatch
  metrics give equivalent data without infra changes.
- **DynamoDB scan on every page load** — too many
  requests for a local dev tool. Fetch-to-disk is
  faster and works offline.
- **Chart.js in `web/static/js/`** — would embed 200KB
  in production binary for a page that never deploys.
  Separate `js-local/` directory avoids this.
- **5 simultaneous charts** — too visually dense. Two
  chart groups with granularity toggles give all views
  in less space.
- **Separate URL for artifact detail** — navigation
  overhead. HTMX partial keeps context on one page.
- **Keeping old stats types** — fundamentally different
  data model (time-series, versions, paths). Clean
  redesign avoids confusing shape changes.

## Open Questions

None.

Date: 2026-06-14
