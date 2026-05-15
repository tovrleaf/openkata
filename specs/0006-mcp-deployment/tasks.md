---
title: MCP Server Deployment — Tasks
---

# Tasks

## Phase 1: Infrastructure

- [x] Create S3 bucket `openkata-artifacts` (block public)
- [x] Create DynamoDB table `openkata-downloads`
- [x] Create IAM role `openkata-mcp-role` (S3 read,
      DynamoDB read/write)
- [x] Create Lambda function `openkata-mcp` with role
- [x] Create Function URL for `openkata-mcp`
- [x] ~~Set reserved concurrency (10) on both Lambdas~~
      Skipped: account concurrency limit too low
- [x] Update `openkata-web-role`: add S3 read, DynamoDB
      read/write
- [x] Update `openkata-ci` role: add S3 PutObject +
      ListObjects, Lambda UpdateFunctionCode for
      `openkata-mcp`
- [x] Script: `infra/create-mcp-stack.sh`

## Phase 2: Publish workflow

- [x] CI workflow: `.github/workflows/publish.yaml`
      (triggers on skill/rule tag push)
- [x] Publish existing 4 tagged skills to S3
- [x] Verify versions.json is correct

## Phase 3: MCP server rewrite

- [x] Remove local filesystem reading (os.ReadDir)
- [x] Remove stdio mode
- [x] Remove `target_dir` parameter
- [x] Add S3 client (read skills, versions.json)
- [x] Add DynamoDB client (read/write counts)
- [x] Cache versions.json in memory on init
- [x] Implement `list_skills` from cache + DynamoDB
- [x] Implement `list_rules` from cache + DynamoDB
- [x] Implement `install_skill` (S3 read + manifest
      generation + DynamoDB increment)
- [x] Implement `install_rule` (same)
- [x] Implement `skill_versions` (S3 ListObjects)
- [x] Implement `rule_versions` (S3 ListObjects)
- [x] Add `WithStateLess(true)` to server config
- [x] Add Lambda adapter (httpadapter.NewV2)
- [ ] Add tags parsing from metadata.tags

## Phase 4: Web downloads

- [ ] Implement `/skills/:name/archive` handler
- [x] Add S3 client to web server
- [x] Add DynamoDB client to web server
- [x] Display download counts on skills listing page
- [ ] Increment counter on download

## Phase 5: Deploy and verify

- [x] Deploy MCP Lambda
- [ ] CI workflow: `.github/workflows/deploy-mcp.yaml`
- [x] Test `list_skills` via curl
- [x] Test `install_skill` via curl
- [ ] Test `skill_versions` via curl
- [ ] Test web download route
- [x] Verify DynamoDB counters increment
