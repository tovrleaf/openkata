---
status: Draft
depth: Standard
---

# Download Analytics

## Story

As the site maintainer, I want each download recorded
as an event with version, source, client, timestamp,
and country so I can analyze usage patterns later.

## Requirements

### Data Model

- Each download recorded as an individual event (not
  a counter increment)
- Event fields:
  - `artifact` — e.g., `skills/commit-conventions`
  - `version` — the version downloaded
  - `source` — `web` or `mcp`
  - `client` — parsed from User-Agent header (e.g.,
    `Claude-Desktop`, `Cursor`, `Kiro`, `browser`,
    `curl`, `unknown`); best-effort for MCP (depends
    on library exposing HTTP headers via context)
  - `timestamp` — ISO 8601, raw (no bucketing)
  - `country` — from CloudFront-Viewer-Country header
    (available for web downloads via CloudFront; MCP
    uses direct Lambda Function URL without CloudFront,
    so country will be empty for MCP events)
- Events kept forever (no TTL)
- Existing total download count on listing/detail
  pages continues to work (dual-write)

### DynamoDB Table: `openkata-download-events`

- Partition key: `artifact` (String)
- Sort key: `timestamp` (String, ISO 8601)
- Attributes: `version`, `source`, `client`, `country`
- Billing: PAY_PER_REQUEST (same as existing table)

### Infrastructure Changes

- Add table creation to `infra/create-mcp-stack.sh`
- Update IAM policies (`iam-mcp-role-policy.json`,
  `iam-web-role-mcp-policy.json`) to grant PutItem
  on `openkata-download-events`
- Update `infra/README.md` to document the new table
- **Gate**: infra changes must be applied manually
  before deploying code that writes to the new table

### Recording Events

- Web archive downloads: record on
  `/skills/:name/archive`, `/rules/:name/archive`,
  `/profiles/:name/archive`
- MCP installs: record in `installArtifact` function
- Include version (from URL or resolved latest)
- Source: `web` for archive handler, `mcp` for MCP
  install handler
- Client: parse User-Agent header into a known client
  name or `unknown`
- Country: read `CloudFront-Viewer-Country` header,
  empty string if not present

### Shared Package

- Extract event recording into `internal/analytics`
- Both `cmd/openkata-web` and `cmd/openkata-mcp` import
  the shared package
- Single `RecordDownload(ctx, event)` function handles
  the PutItem to events table
- Both binaries continue to increment the old counter
  table as before (dual-write)

### Client Parsing

- `Claude-Desktop` — User-Agent contains "Claude"
- `Cursor` — User-Agent contains "Cursor"
- `Kiro` — User-Agent contains "Kiro"
- `curl` — User-Agent contains "curl"
- `browser` — User-Agent contains "Mozilla"
- `unknown` — anything else

## Out of Scope

- Admin dashboard (separate spec)
- Real-time analytics
- User/session tracking
- IP address storage
- Alerting or notifications
- Data export
- Migration from counter table
- Putting MCP behind CloudFront

## Open Questions

None.

Date: 2026-06-12
