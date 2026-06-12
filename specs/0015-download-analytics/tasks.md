# Tasks: Download Analytics

## Tasks

### 1. Infrastructure: create events table and update IAM
- **Status**: Pending
- **Goal**: Add DynamoDB table creation to infra script,
  update IAM policies for both Lambdas
- **Boundary**: `infra/create-mcp-stack.sh`,
  `infra/iam-mcp-role-policy.json`,
  `infra/iam-web-role-mcp-policy.json`,
  `infra/README.md`
- **Key files**: existing table creation in
  `create-mcp-stack.sh` (Step 2)
- **Depends**: None
- **Done when**: Script includes events table creation,
  IAM policies grant PutItem on new table, README
  documents the table
- **Verify**: `bash -n infra/create-mcp-stack.sh`

**GATE: User must run infra scripts manually before
proceeding to task 2.**

### 2. Create shared analytics package
- **Status**: Pending
- **Goal**: Create `internal/analytics` package with
  `RecordDownload` function and client parsing logic
- **Boundary**: `internal/analytics/analytics.go`,
  `internal/analytics/analytics_test.go`
- **Key files**: `cmd/openkata-web/handlers.go`
  (incrementDownload), `cmd/openkata-mcp/main.go`
  (incrementCount)
- **Depends**: 1 (table must exist)
- **Done when**: Package compiles, has RecordDownload
  function that PutItems to events table, ParseClient
  function with tests
- **Verify**: `go test ./internal/analytics/...`

### 3. Integrate analytics in web server
- **Status**: Pending
- **Goal**: Call RecordDownload from the archive handler
  with source=web, parsed client, version, and country
- **Boundary**: `cmd/openkata-web/handlers.go`
- **Key files**: `handleArchive` function (line ~1073),
  `incrementDownload` function (line ~1337)
- **Depends**: 2
- **Done when**: Archive downloads write events to new
  table AND continue incrementing old counter; build
  passes
- **Verify**: `go build -o bin/openkata-web ./cmd/openkata-web/`

### 4. Integrate analytics in MCP server
- **Status**: Pending
- **Goal**: Call RecordDownload from installArtifact
  with source=mcp, best-effort client detection,
  version, and empty country
- **Boundary**: `cmd/openkata-mcp/main.go`
- **Key files**: `installArtifact` function (line ~287),
  `incrementCount` function (line ~480)
- **Depends**: 2
- **Done when**: MCP installs write events to new table
  AND continue incrementing old counter; build passes
- **Verify**: `go build -o bin/openkata-mcp ./cmd/openkata-mcp/`

### 5. Add tests
- **Status**: Pending
- **Goal**: Test client parsing logic, test event
  recording with mock DynamoDB client
- **Boundary**: `internal/analytics/analytics_test.go`
- **Depends**: 2
- **Done when**: Tests cover all client parsing cases
  and RecordDownload behavior
- **Verify**: `go test -race ./internal/analytics/...`

## Progress Log
