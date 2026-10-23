---
status: Draft
depth: Deep
---

# Documentation Section

## Story

As a visitor to openkata.dev, I want a documentation section
with guides, reference, architecture notes, and contributor
instructions so that I can learn how OpenKata works, install
and use its artifacts, and contribute without reading the
repository source.

## Requirements

### Content Source and Location

- Documentation pages are Markdown files under `web/docs/`,
  next to `web/static/`
- `web/docs/` is a Go package (`embed.go`, same pattern as
  `web/static/embed.go`) so the content is compiled into the
  `openkata-web` binary; production ships only the binary
- In local development the server may read the files from
  disk so edits show without a rebuild, as it does for
  `web/static/`
- Pages use frontmatter for `title`, optional `description`,
  and optional `order`; existing `stripFrontmatter` handling
  is reused

### Sections

The first version has four sections, each with an index page
and at least the pages listed:

- **Guides**
  - Skills, rules and profiles: what each is and how they
    differ, with the kata vocabulary (dojo kun, sensei
    profiles, the dojo, the ryu)
  - Installing artifacts through the MCP server
  - Using the catalog
- **Reference**
  - How MCP discovery works: the `list_*` tools and tag
    filter, and where the index lives (`versions.json` in S3)
  - Where skills, rules and profiles are stored: the S3
    bucket layout and how `versions.json` is generated
  - Downloads: how downloads are recorded and counted, and
    what the stats dashboard shows
  - How we use evals, including the local eval runner and
    its Docker sandbox
  - How we score skills: effectiveness and the tessl score
  - How we use tessl: reviews, evals and the registry
  - Versioning and `.openkata-manifest.json`
- **Architecture**
  - Our agent workflow
  - The infrastructure diagram
  - The published ADRs
- **Contributing**
  - Authoring skills, rules and profiles
  - Release workflow, one page each for a new version of the
    site, of a skill, rule or profile, and of the MCP server
  - Repository conventions

Reference pages are written by hand. There is a single live
version of the docs (current `main`); there is no per-release
snapshot or version selector.

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

### ADRs

- ADRs stay in `docs/adr/` as the source of truth, so existing
  repository links, tooling and agent instructions keep
  working
- They are published under the Architecture section from
  those files, with no hand-edited copy
- No generated content is committed to the repository
- A clean checkout builds and serves exactly the ADR set that
  `docs/adr/` contains, and a build can never silently ship
  without them
- The mechanism (embedding `docs/adr/` directly or generating a
  copy that is gitignored) is decided in the design

### Diagrams

- Fenced `mermaid` blocks render as diagrams, because
  `docs/agent-workflow.md` already uses them and the site has
  no Mermaid support today
- Diagrams follow the active theme and stay readable in light
  and dark themes; if scripts do not run, the diagram source
  remains visible
- Preformatted blocks such as the ASCII infrastructure diagram
  scroll inside their own container instead of widening the
  page

### Routing and Rendering

- `/docs/` serves the documentation home with the four
  sections
- `/docs/<section>/` serves a section index and
  `/docs/<section>/<page>/` serves a page
- Unknown paths return 404; requests can only reach embedded
  documentation files (no path traversal)
- Markdown renders through the existing goldmark pipeline
  (tables, external links open in a new tab)
- Relative links between documentation pages resolve to the
  matching `/docs/...` URL

### Navigation and Layout

- A "Docs" link is added to the main nav and to the mobile
  nav menu
- Pages show a sidebar listing every section and its pages,
  built from the file tree and frontmatter `order`; the
  current page is highlighted
- On narrow screens the sidebar collapses behind a control
  and the page content remains readable without horizontal
  scroll
- Pages use the existing design tokens, themes and layout, so
  all themes render correctly
- Each page has previous and next links within its section

### Search

- Documentation pages are findable through the existing
  search overlay
- The server builds a small client-side index of page titles,
  descriptions and headings from the embedded content and the
  overlay loads it; no new search endpoint is added

### Quality

- Handler behavior is covered by table-driven tests using
  `httptest` and fixtures in `testdata/`
- `go build` has no warnings and committed `_templ.go` files
  are in sync with `.templ` sources

## Out of Scope

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

- **Static site generator (Hugo, MkDocs)** — adds a second
  build and deploy pipeline for one section, and ADR 0007
  keeps the site server-rendered by the Go binary
- **Serving `docs/` directly** — mixes public pages with
  internal files and requires an exclude list to hide
  agent-only documents
- **Reading Markdown from S3 at request time** — pages would
  update independently of releases, but this adds latency and
  an AWS dependency for content that changes with the code
- **Moving ADRs to `web/docs/architecture/`** — one copy and no
  build step, but every link, tool and agent instruction that
  points at `docs/adr/` would need updating

## Open Questions

None.
