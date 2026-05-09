---
title: Infrastructure for Lambda Deployment
status: Draft
depth: Standard
created: 2026-05-09
---

# Infrastructure for Lambda Deployment

## Goal

Deploy the openkata-web server to AWS Lambda using CDK (Go)
so the site is publicly accessible. Everything as code,
deployable from a developer machine.

## Requirements

1. CDK stack in Go defines all AWS resources
2. Lambda runs the openkata-web binary (arm64, provided.al2023)
3. Lambda Function URL provides public HTTPS access
4. Static assets embedded in the binary (go:embed)
5. Single command deploys: `make deploy`
6. Infrastructure lives in `infra/` directory
7. No custom domain (use Function URL directly)

## Out of Scope

- Custom domain and Route 53
- CloudFront distribution (add with domain later)
- DynamoDB install counters
- CI/CD pipeline (GitHub Actions)
- MCP server integration (separate future work)

## Open Questions

- None — ADR 0007 covers the architecture decisions.

## Acceptance Criteria

- `make deploy` succeeds from a configured machine
- CloudFront URL returns the landing page HTML
- Static assets (CSS, htmx.js) load correctly
- Cold start under 500ms
