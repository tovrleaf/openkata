# Tasks: MCP Server Deployment

## Tasks

### 1. Create infrastructure
- **Status**: Done
- **Goal**: S3 bucket, DynamoDB table, IAM roles, Lambda
  function, and Function URL for the MCP server
- **Boundary**: `infra/create-mcp-stack.sh`
- **Depends**: None
- **Done when**: Script creates all resources, Lambda has
  a Function URL

### 2. Create publish workflow
- **Status**: Done
- **Goal**: CI workflow that publishes skills/rules to S3
  on tag push, generates versions.json
- **Boundary**: `.github/workflows/publish.yaml`
- **Depends**: 1
- **Done when**: Tagging a skill triggers upload to S3 and
  versions.json is correct

### 3. Rewrite MCP server for Lambda
- **Status**: Done
- **Goal**: Remove local filesystem, add S3/DynamoDB
  clients, implement all MCP tools, add Lambda adapter
- **Boundary**: `cmd/openkata-mcp/`
- **Depends**: 1
- **Done when**: All tools work via Function URL (list,
  install, versions)

### 4. Add tags parsing
- **Status**: Done
- **Goal**: Parse `metadata.tags` from skill/rule metadata
- **Boundary**: `cmd/openkata-mcp/`
- **Depends**: 3
- **Done when**: Tags appear in list_skills/list_rules
  responses

### 5. Web download routes
- **Status**: Done
- **Goal**: Archive download handler with DynamoDB counter
  increment
- **Boundary**: `cmd/openkata-web/handlers.go`
- **Depends**: 1
- **Done when**: `/skills/:name/archive` serves tar.gz and
  increments download count

### 6. Deploy and verify
- **Status**: Done
- **Goal**: Deploy MCP Lambda, create deploy workflow,
  verify all tools end-to-end
- **Boundary**: `.github/workflows/deploy-mcp.yaml`
- **Depends**: 3
- **Done when**: All MCP tools respond correctly via curl,
  DynamoDB counters increment

## Progress Log

- [2026-05-10] Tasks 1-3 completed. Infrastructure created,
  publish workflow working, MCP server deployed and verified
  via curl.
- [2026-06-18] Tasks 4-6 confirmed complete. Tags parsing
  live in MCP server, archive routes serving downloads with
  analytics, deploy workflow running on main merge.
