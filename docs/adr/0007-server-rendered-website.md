---
status: ACCEPTED
date: 2026-05-03
authors: [niko.kivela]
---

# 0007. Server-rendered website with Go, templ, and htmx

## Context

Open Kata needs a public website to introduce the project, document
the MCP server, provide installation instructions, and let users
browse available skills and rules. The site must be cheap to host,
discoverable by search engines and AI agents, and maintainable
within the existing Go monorepo.

A separate portfolio site will share the same design system in the
future, so the visual foundation must be reusable.

## Decision Drivers

- Minimal hosting cost — ideally within AWS free tier
- No separate frontend framework or build toolchain
- Content derived from existing SKILL.md and RULE.md files
- Must support dynamic interactions (catalog browsing, filtering)
- Design system must be reusable across future sites
- Single repo, single language (Go)

## Decision

Build a server-rendered website using Go, templ, and htmx. Serve
it alongside the MCP server from a single Lambda behind CloudFront
with a custom domain.

### Architecture

One Lambda binary handles both website and MCP traffic, routed
by path:

- `openkata.dev/*` → website (templ + htmx)
- `openkata.dev/mcp` → MCP protocol handler

CloudFront sits in front for caching, TLS termination, and the
custom domain.

### Technology choices

- **Go + templ** — server-side HTML rendering, no separate
  frontend. Templ generates type-safe Go code from templates.
- **htmx** — dynamic interactions (filtering skills, loading
  details) via HTML fragment requests to the same Go server.
  No JavaScript framework needed.
- **Token-based CSS design system** — CSS custom properties
  for colors, typography, and spacing. Structured as
  `tokens.css`, `base.css`, `typography.css`, `components/`,
  `layouts/`. Reusable across future sites.

### Content pipeline

The binary reads `skills/*/SKILL.md` and `rules/*/RULE.md` at
startup, parses frontmatter (name, description), and builds the
catalog. The same discovery code the MCP server already uses is
shared between both.

### Install tracking

DynamoDB counter table tracks skill/rule installs. The MCP
`install_skill` and `install_rule` handlers atomically increment
a counter per artifact. The website reads counts for display.

### Discoverability

- `robots.txt` — standard crawl directives
- `llms.txt` — AI agent discovery
- `sitemap.xml` — generated from the skill/rule catalog
- Open Graph and meta tags per page, derived from SKILL.md
  frontmatter

### Infrastructure

- **Lambda Function URL** behind CloudFront
- **ACM certificate** for TLS on the custom domain
- **DynamoDB** single table for install counters
- **GitHub Actions** for CI/CD — build, test, deploy

### Local development

Single command: `make dev` runs `air`, which watches `.go`,
`.templ`, and `.css` files, regenerates templ, rebuilds, and
serves on `localhost:8080`. Filesystem-based asset serving in
dev, embedded FS in production.

### Project structure

```text
cmd/
├── openkata-mcp/    # Existing MCP server (becomes shared)
└── openkata-web/    # Website + MCP in single binary
web/
├── styles/
│   └── ds/          # Design system (tokens, base, typography)
├── templates/       # Templ templates
└── static/          # robots.txt, llms.txt, favicon
```

## Alternatives Considered

### Static site on S3 + CloudFront

- **Pros:** Cheapest possible, no Lambda, simple deployment
- **Cons:** Cannot support htmx interactions, no dynamic
  catalog browsing, MCP server needs separate infrastructure
- **Rejected because:** htmx requires a running server to
  return HTML fragments. Two separate deployments (S3 for
  site, Lambda for MCP) is more complex than one Lambda
  for both.

### Astro or Hugo static site generator

- **Pros:** Rich ecosystem, markdown-native, fast builds
- **Cons:** Introduces Node.js or a second language, cannot
  share Go code with MCP server, separate build toolchain
- **Rejected because:** adds a language and toolchain to a
  Go-only project. Templ keeps everything in Go.

### Separate Lambda for website and MCP

- **Pros:** Independent scaling and deployment
- **Cons:** Two Lambdas, two deployments, duplicated skill
  discovery code, more CloudFront routing complexity
- **Rejected because:** at this scale, one Lambda is simpler
  and cheaper. Can split later if needed.

## Consequences

### Positive

- Single binary, single deployment, single language
- Shared skill discovery code between website and MCP
- htmx provides dynamic interactions without a JS framework
- Design system is reusable for future sites
- Install tracking comes nearly free with DynamoDB
- Hosting cost within AWS free tier for expected traffic

### Negative

- Lambda cold starts affect first request (Go mitigates this —
  typically under 100ms)
- Server-rendered means no offline capability
- `air` is an additional dev dependency
- Templ is less mature than established template engines

### Neutral

- The design system skill and rule (future Open Kata artifacts)
  will be developed alongside this site
- MCP and website sharing a binary can be split later if
  scaling demands it

## References

- [templ](https://templ.guide) — Go HTML templating
- [htmx](https://htmx.org) — HTML-driven interactions
- [air](https://github.com/air-verse/air) — Go live reload
- ADR 0002 — MCP server for distribution
- ADR 0003 — Go as implementation language
