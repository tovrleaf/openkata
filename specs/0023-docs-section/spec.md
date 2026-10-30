---
status: Done
depth: Deep
---

# Documentation Section

## Story

As a maintainer of openkata.dev, I want a documentation engine
that renders Markdown pages into a docs section with
navigation, so that guides, reference, architecture notes and
contributor instructions can be written and previewed locally
before the section goes public.

This spec builds the engine and moves the existing
architecture docs into it. Writing the pages and making the
section public are spec 0024.

## Requirements

### Content Source and Location

- Documentation pages are Markdown files under `web/docs/`,
  next to `web/static/`
- `web/docs/` is a Go package (`embed.go`, same pattern as
  `web/static/embed.go`) so the content is compiled into the
  `openkata-web` binary, ready for go-live
- In local development the server reads the files from disk
  so edits show without a rebuild, as it does for
  `web/static/`
- Pages use frontmatter for `title`, optional `description`,
  and optional `order`; existing `stripFrontmatter` handling
  is reused

### Sections

- The engine supports four sections: Guides, Reference,
  Architecture and Contributing
- Each section has an `index.md`; in this spec the Guides,
  Reference and Contributing indexes are stubs
- The Architecture section contains the two moved documents
  below
- The page list for each section is spec 0024

There is a single live version of the docs (current `main`);
there is no per-release snapshot or version selector.

### Existing Architecture Docs

- `docs/agent-workflow.md` and `docs/infrastructure-diagram.md`
  move into the embedded tree as
  `web/docs/architecture/agent-workflow.md` and
  `web/docs/architecture/infrastructure-diagram.md`, with
  frontmatter added
- They are moved, not copied: no copy remains under `docs/`,
  and the files in `web/docs/` are the only source
- Nothing in the repository references the old paths today;
  the move must leave no dangling references

### Visibility

- `/docs/` is registered only in local development, in the
  same block as `/design-system/` and `/stats/`; production
  serves nothing under `/docs/`
- There is no "Docs" link in the main or mobile nav
- Existing production JavaScript is not changed

### Links

- Relative links between documentation pages resolve to the
  matching `/docs/...` URL
- Links to anything outside `web/docs/`, including ADRs, are
  absolute GitHub URLs on `main`; a relative link that leaves
  `web/docs/` is an error
- External links open in a new tab

### Diagrams and Preformatted Blocks

- Fenced `mermaid` blocks are shown as their source in a
  `<pre>` and marked with a class, so rendering can be added
  later without changing the server
- Preformatted blocks, such as the ASCII infrastructure
  diagram, scroll inside their own container instead of
  widening the page

### Routing and Rendering

- `/docs/` serves the documentation home with the four
  sections
- `/docs/<section>/` serves a section index and
  `/docs/<section>/<page>/` serves a page
- A page path without a trailing slash redirects to the slash
  form; `/docs/search-index.json` is exempt
- Unknown paths return 404; requests can only reach
  documentation pages (no path traversal)
- Markdown renders through the existing goldmark pipeline
  (tables, unsafe HTML for trusted repo content)

### Content Validation

The loader is strict. Any of these fails the whole docs site,
naming the file and the problem:

- Malformed frontmatter
- Two files mapping to the same URL
- A section directory without `index.md`
- A relative link that does not resolve to a docs page

Locally a failure returns 500 with the error for every
`/docs/` request; in CI the real-tree test fails.

### Ordering

- Sections and pages sort by frontmatter `order`
- Pages without `order` sort after all ordered pages, then by
  file name

### Navigation and Layout

- Section and page views show a breadcrumb trail at the top,
  `Docs > Section > Page`; every item but the last links to
  its page, and the home page has none
- Pages show a sidebar listing every section and its pages;
  the current page is highlighted
- On narrow screens the sidebar collapses behind a control
  and the page content remains readable without horizontal
  scroll
- Pages use the existing design tokens, themes and layout, so
  all themes render correctly
- Each page has previous and next links within its section

### Search Index

- The server serves `/docs/search-index.json`: page titles,
  paths, sections, descriptions and headings built from the
  docs content
- It is data, not a query endpoint
- The search overlay does not use it in this spec

### Quality

- Handler and loader behavior is covered by table-driven tests
  using `httptest` and `fstest.MapFS` fixtures
- A test runs the loader against the real embedded tree
- `go build` has no warnings and committed `_templ.go` files
  are in sync with `.templ` sources

## Out of Scope

- Writing guide, reference and contributing pages (spec 0024)
- Public route, nav link and search overlay integration
  (spec 0024)
- Rendering Mermaid diagrams
- Publishing ADRs; they stay in `docs/adr/` and are read on
  GitHub
- "Edit on GitHub" links on pages
- Versioned (per-release) docs and a version selector
- Generating reference pages from code (MCP tool list, CLI
  flags)
- A server-side search endpoint
- Rewriting the existing Getting Started page
- Documentation for agents only (spec 0022's `docs/design/`
  blueprints stay internal and are not published)
- Comments, feedback widgets or analytics specific to docs

## Rejected Approaches

- **Static site generator (Hugo, MkDocs)**: adds a second
  build and deploy pipeline for one section, and ADR 0007
  keeps the site server-rendered by the Go binary
- **Serving `docs/` directly**: mixes public pages with
  internal files and needs an exclude list to hide agent-only
  documents
- **Reading Markdown from S3 at request time**: pages would
  update independently of releases, but this adds latency and
  an AWS dependency for content that changes with the code
- **Publishing ADRs**: they record decisions for maintainers,
  not end users, and GitHub already serves them. Publishing
  them needed a second embedded tree and links resolved
  across both trees
- **Shipping engine and content together**: one large PR, and
  the engine waits for the slowest page
- **Hiding the section in production with only the nav link
  removed**: the route and search overlay would still expose
  half-written pages
- **Vendoring Mermaid now**: about 3 MB, more than the
  current 1.2 MB of static assets, plus theme re-render code,
  for four diagrams in one file

## Open Questions

None.
