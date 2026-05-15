---
title: MCP Server Deployment
status: In Progress
depth: Standard
created: 2026-05-10
adr: 0010-deploy-mcp-to-lambda
---

# MCP Server Deployment

Deploy the OpenKata MCP server to AWS Lambda so agents
can discover, list, and install skills/rules remotely.

## Requirements

1. MCP server runs on Lambda with Streamable HTTP
   transport (stateless mode)
2. Skills and rules served from S3 (only tagged releases)
3. `versions.json` in S3 as lightweight index
4. CI workflow publishes to S3 on tag creation
5. Download counting via DynamoDB atomic counters
6. Web server serves tar.gz downloads from S3
7. Both Lambdas have S3 read + DynamoDB read/write
8. Rate limiting via Lambda reserved concurrency
9. MCP tools: `list_skills`, `list_rules`,
   `install_skill`, `install_rule`, `skill_versions`,
   `rule_versions`
10. Tags in `metadata.tags` frontmatter (comma-separated)

## Out of scope

- Auth / private skills
- Skill dependencies
- Quality scoring
- Local stdio mode (removed)

## Technical decisions

- Separate Lambda: `openkata-mcp`
- S3 bucket: `openkata-artifacts`
- DynamoDB table: `openkata-downloads`
- Region: `eu-north-1`
- Stateless: `WithStateLess(true)` in mcp-go
- `versions.json` cached in Lambda memory on cold start
- Version history resolved via S3 `ListObjects`
- Web downloads: `/skills/<name>/download?v=<version>`
- Install returns file contents + manifest with checksums
- Agent handles update logic (checksum comparison)
