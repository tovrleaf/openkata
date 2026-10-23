---
status: Draft
depth: Standard
---

# Static Asset Delivery and Diagram Rendering

## Story

As a maintainer of openkata.dev, I want static assets cached and
large vendored files kept out of the Lambda binary, and a
decided way to show diagrams in the docs, so that pages load
fast and the docs can grow without bloating the deploy.

## Background

- `/static/` is served by the `openkata-web` Lambda from
  `web/static` embedded in the binary (1.2 MB today)
- CloudFront has one origin, the Lambda Function URL, and its
  cache policy ID is believed to be the managed
  *CachingDisabled* policy (`infra/create-cloudfront.sh`); this
  must be confirmed on the live distribution
- The Go code sets no `Cache-Control` on static responses, and
  `embed.FS` files carry no modification time, so browsers get
  no validators or cache hints
- Every static file request therefore invokes the Lambda
- Spec 0023 shows `mermaid` blocks as source because vendoring
  Mermaid (about 3 MB) would more than triple the static assets

## Requirements

### Part 1: Static asset delivery

- Confirm the live CloudFront cache policy and record it in
  the spec
- Static responses carry an explicit cache policy, and a
  change to an asset reaches visitors without a stale-cache
  window longer than the agreed limit
- `/static/*` is cached at the edge so the Lambda only sees
  misses
- Decide whether large vendored files are served from S3
  behind CloudFront or stay embedded, with the decision and
  its trade-offs recorded in the design. Points to settle:
  version skew between templates and assets, hashed file names
  or invalidation, local development, and behavior when the
  asset origin is unavailable
- Any infrastructure change is scripted under `infra/` like
  the existing scripts

### Part 2: Diagram rendering

Compare these options against the docs' four existing
diagrams and decide:

- **Mermaid, client-side:** keeps diagrams as text in
  Markdown; needs a vendored or hosted JS file and theme
  re-render code
- **Inline SVG drawn with design tokens:** follows the theme
  with no JS; diagrams are edited as SVG. Candidate authoring
  aid: [visual-explainer](https://github.com/nicobailon/visual-explainer),
  an MIT agent skill that writes hand-drawn SVG pages. It is an
  authoring tool, not a runtime library, and is not a
  dependency of the site
- **Pre-rendered SVG from Mermaid source at build time:**
  text source stays in the repo; output is static
- **Source view only (spec 0023 behavior)**

The decision criteria are: readable in both themes, works
without JavaScript, bytes added to the page and the binary,
how a contributor edits a diagram, and the number of diagrams
the docs will realistically have after spec 0024.

## Out of Scope

- Moving HTML pages or the docs Markdown to S3
- Changing the Lambda or CloudFront architecture beyond cache
  behaviors and an optional asset origin
- Adopting visual-explainer as a project dependency or
  workflow requirement

## Open Questions

- Is the live cache policy really *CachingDisabled*?
- Does the 0024 go-live need rendered diagrams, or can this
  spec land after it?
- If S3 is chosen, is it a new bucket or a prefix in
  `openkata-artifacts`?
