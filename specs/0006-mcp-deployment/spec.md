---
status: Done
depth: Standard
---

# MCP Server Deployment

## Story

As an agent user, I want to discover, list, and install
skills/rules remotely via MCP so that I can use kata without
manually copying files.

## Requirements

- MCP server runs on Lambda with Streamable HTTP transport
  (stateless mode)
- Skills and rules served from S3 (only tagged releases)
- `versions.json` in S3 as lightweight index
- CI workflow publishes to S3 on tag creation
- Download counting via DynamoDB atomic counters
- Web server serves tar.gz downloads from S3
- Both Lambdas have S3 read + DynamoDB read/write
- Rate limiting via Lambda reserved concurrency
- MCP tools: `list_skills`, `list_rules`, `install_skill`,
  `install_rule`, `skill_versions`, `rule_versions`
- Tags in `metadata.tags` frontmatter (comma-separated)

## Out of Scope

- Auth / private skills
- Skill dependencies
- Quality scoring
- Local stdio mode (removed)

## Open Questions

- None
