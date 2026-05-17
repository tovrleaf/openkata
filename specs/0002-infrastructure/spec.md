---
status: Done
depth: Standard
---

# Infrastructure for Lambda Deployment

## Story

As a maintainer, I want the website deployed to AWS Lambda
so that the site is publicly accessible without managing
servers.

## Requirements

- CloudFormation template defines all AWS resources
- Lambda runs the openkata-web binary (arm64,
  provided.al2023)
- Lambda Function URL provides public HTTPS access
- Static assets embedded in the binary (go:embed)
- Single command deploys: `make deploy`
- Infrastructure lives in `infra/` directory
- No custom domain (use Function URL directly)

## Out of Scope

- Custom domain and Route 53
- CloudFront distribution (add with domain later)
- DynamoDB install counters
- CI/CD pipeline (GitHub Actions)
- MCP server integration (separate future work)

## Open Questions

- None — ADR 0007 covers the architecture decisions.
