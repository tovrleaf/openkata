---
status: Implementing
depth: Deep
---

# Static Asset Delivery and Diagram Rendering

## Story

As a maintainer of openkata.dev, I want static assets cached and
large vendored files kept out of the Lambda binary, and a
decided way to show diagrams in the docs, so that pages load
fast and the docs can grow without bloating the deploy.

## Background

- `/static/` is served by the `openkata-web` Lambda from
  `web/static` embedded in the binary. The files are 1.3 MB of
  the 35 MB linux/arm64 binary
- CloudFront distribution `E1X978YU1XMI3Q` has two origins, the
  web and MCP Lambda Function URLs. Its default behavior and
  `/mcp*` both use cache policy
  `4135ea2d-6df8-44a3-9df3-4b5a84be39ad`, the managed
  *CachingDisabled* policy (confirmed on 2026-11-16)
- The Go code sets no `Cache-Control` on static responses, and
  `embed.FS` files carry no modification time, so browsers get
  no validators or cache hints
- Every static file request therefore invokes the Lambda
- Spec 0023 shows `mermaid` blocks as source because vendoring
  Mermaid (about 3 MB) would more than triple the static assets

## Requirements

### Part 1: Static asset delivery

- `web/static` stays embedded in the Lambda binary. With no
  Mermaid library, no large vendored file is coming, so S3 is
  not used for assets
- A `/static/*` CloudFront behavior caches static files at the
  edge, so the Lambda only sees the first request for each URL
- Static responses carry `Cache-Control` and an `ETag`, so a
  browser can revalidate with a `304`
- URLs that templates write carry a per-file content hash
  (`?v=<hash>`). These files are cached for a year, and a
  changed file reaches visitors on the next page load
- Files without a hash in their URL reach visitors within 5
  minutes of a deploy. These are CSS `@import`s, files that
  scripts fetch and web manifest icons
- Local development serves `web/static` from disk, as today,
  and an edited file shows on the next reload
- Any infrastructure change is scripted under `infra/` like
  the existing scripts

### Part 2: Diagram rendering

- Diagrams are SVG files, drawn by hand or by an agent, next to
  the docs page that uses them. The files are the source: nothing
  in the build generates them. They are linked with Markdown
  image syntax and descriptive alt text
- The docs renderer inlines a linked local `.svg` file into the
  page, so the SVG's CSS variables follow the site theme
- Colors in the SVG use the site's design tokens, with fallback
  values so the file reads on GitHub in both color schemes
- The four `mermaid` diagrams in
  `web/docs/architecture/agent-workflow.md` are redrawn this way
- The strict docs loader rejects a page that links an SVG file
  that does not exist
- No Mermaid remains: the SVG files hold no Mermaid source, and
  the renderer's special handling of `mermaid` blocks from spec
  0023 is removed. A `mermaid` block renders like any other code
  block
- [visual-explainer](https://github.com/nicobailon/visual-explainer)
  may be used to draft a diagram. It is not a dependency of
  the site or a required step for contributors

The comparison of the four diagram options and the asset origin
decision are recorded in [design.md](design.md).

## Out of Scope

- Moving HTML pages or the docs Markdown to S3
- Changing the Lambda or CloudFront architecture beyond a cache
  behavior
- Serving static files from S3
- Adopting visual-explainer as a project dependency or
  workflow requirement
- Redrawing the ASCII diagrams in `text` blocks
- Generating diagrams from Mermaid source. Spec 0030 does that
  once Mermaid v13 is stable
