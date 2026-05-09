---
status: Done
depth: Deep
---

# Website Foundation

## Story

As a maintainer, I want the website project structure and
local dev environment working so that I can iterate on pages
with hot reload and build toward the full site incrementally.

## Requirements

- Go binary at `cmd/openkata-web/` that serves HTML
- Templ for server-side HTML rendering
- htmx included for future dynamic interactions
- `make dev` runs air with templ generation and auto-reload
- Landing page at `/` with project name and tagline
- Filesystem-based asset serving in dev mode
- CSS served from a static directory (no design system yet —
  just the structure for it)
- Project builds with `go build ./cmd/openkata-web/`

## Out of Scope

- Design system (tokens, typography, colors) — separate spec
- Skill/rule catalog pages
- MCP documentation
- DynamoDB install tracking
- AWS deployment (Lambda, CloudFront)
- Custom domain
- robots.txt, llms.txt, sitemap
- Production embedded FS

## Open Questions

- None — scope is clear enough to proceed
