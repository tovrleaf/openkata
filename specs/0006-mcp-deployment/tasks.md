---
title: MCP Server Deployment — Tasks
---

# Tasks

## Phase 1: Infrastructure

- [ ] Create S3 bucket `openkata-artifacts` (block public)
- [ ] Create DynamoDB table `openkata-downloads`
- [ ] Create IAM role `openkata-mcp-role` (S3 read,
      DynamoDB read/write)
- [ ] Create Lambda function `openkata-mcp` with role
- [ ] Create Function URL for `openkata-mcp`
- [ ] Set reserved concurrency (10) on both Lambdas
- [ ] Update `openkata-web-role`: add S3 read, DynamoDB
      read/write
- [ ] Update `openkata-ci` role: add S3 PutObject +
      ListObjects, Lambda UpdateFunctionCode for
      `openkata-mcp`
- [ ] Script: `infra/create-mcp-stack.sh`

## Phase 2: Publish workflow

- [ ] CI workflow: `.github/workflows/publish.yaml`
      (triggers on skill/rule tag push)
- [ ] Publish existing 4 tagged skills to S3
- [ ] Verify versions.json is correct

## Phase 3: MCP server rewrite

- [ ] Remove local filesystem reading (os.ReadDir)
- [ ] Remove stdio mode
- [ ] Remove `target_dir` parameter
- [ ] Add S3 client (read skills, versions.json)
- [ ] Add DynamoDB client (read/write counts)
- [ ] Cache versions.json in memory on init
- [ ] Implement `list_skills` from cache + DynamoDB
- [ ] Implement `list_rules` from cache + DynamoDB
- [ ] Implement `install_skill` (S3 read + manifest
      generation + DynamoDB increment)
- [ ] Implement `install_rule` (same)
- [ ] Implement `skill_versions` (S3 ListObjects)
- [ ] Implement `rule_versions` (S3 ListObjects)
- [ ] Add `WithStateLess(true)` to server config
- [ ] Add Lambda adapter (httpadapter.NewV2)
- [ ] Add tags parsing from metadata.tags

## Phase 4: Web downloads

- [ ] Add download route to web server
- [ ] Add S3 client to web server
- [ ] Add DynamoDB client to web server
- [ ] Implement `/skills/:name/archive` handler
- [ ] Increment counter on download
- [ ] Display download counts on website (future page)

## Phase 5: Deploy and verify

- [ ] Deploy MCP Lambda
- [ ] CI workflow: `.github/workflows/deploy-mcp.yaml`
- [ ] Test `list_skills` via curl
- [ ] Test `install_skill` via agent
- [ ] Test `skill_versions` via curl
- [ ] Test web download route
- [ ] Verify DynamoDB counters increment
