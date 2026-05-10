---
title: Custom Domain with CloudFront
status: Draft
depth: Standard
created: 2026-05-10
---

# Custom Domain with CloudFront

## Goal

Serve openkata.dev (and www.openkata.dev) via CloudFront
in front of the existing Lambda Function URL. TLS via ACM.

## Requirements

1. openkata.dev serves the website
2. www.openkata.dev redirects to openkata.dev
3. Route 53 hosted zone for openkata.dev
4. ACM certificate covering both apex and www
5. CloudFront distribution with custom domain
6. All infrastructure scripted — no manual AWS console work
7. Agent does not run AWS commands directly; scripts are
   provided for the user to execute

## Out of Scope

- MCP path routing (future)
- Cache invalidation automation
- WAF / rate limiting

## Manual Steps (user)

1. Point registrar nameservers to Route 53 (after hosted
   zone is created)
2. Run scripts in order
3. Wait for DNS propagation

## Acceptance Criteria

- `curl https://openkata.dev` returns the landing page
- `curl https://www.openkata.dev` redirects to openkata.dev
- TLS certificate is valid and auto-renewing (ACM)
- All resources created via scripts in `infra/`
