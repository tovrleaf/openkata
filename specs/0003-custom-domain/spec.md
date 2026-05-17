---
status: Done
depth: Standard
---

# Custom Domain with CloudFront

## Story

As a visitor, I want to access the site at openkata.dev so
that it has a professional, memorable URL with TLS.

## Requirements

- openkata.dev serves the website
- www.openkata.dev redirects to openkata.dev
- Route 53 hosted zone for openkata.dev
- ACM certificate covering both apex and www
- CloudFront distribution with custom domain
- All infrastructure scripted — no manual AWS console work
- Agent does not run AWS commands directly; scripts are
  provided for the user to execute

## Out of Scope

- MCP path routing (future)
- Cache invalidation automation
- WAF / rate limiting

## Open Questions

- None
