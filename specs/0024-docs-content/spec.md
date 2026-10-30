---
status: Draft
depth: Standard
---

# Documentation Content and Go-Live

## Story

As a visitor to openkata.dev, I want a documentation section
with guides, reference, architecture notes and contributor
instructions so that I can learn how OpenKata works, install
and use its artifacts, and contribute without reading the
repository source.

Builds on spec 0023, which provides the docs engine.

## Requirements

### Pages

Each section gets at least these pages, written by hand:

- **Guides**
  - Manifesto: a short page that points to `MANIFESTO.md` in the
    repository root, which stays the single source, so
    developers can find it from the docs
  - Skills, rules and profiles: what each is and how they
    differ, with the kata vocabulary (dojo kun, sensei
    profiles, the dojo, the ryu)
  - Installing artifacts through the MCP server
  - Using the catalog
- **Reference**
  - How MCP discovery works: the `list_*` tools and tag
    filter, and where the index lives (`versions.json`)
  - Where skills, rules and profiles are stored and how
    `versions.json` is generated, including the S3 bucket
    layout; the infrastructure is public as code, so nothing
    is held back
  - Downloads: how downloads are recorded and counted, and
    the public download counts shown in the catalog. The
    local-only stats dashboard (`/stats/`) is not documented
  - How we use evals, including the local eval runner and
    its Docker sandbox
  - How we score skills: effectiveness and the tessl score
  - How we use tessl: reviews, evals and the registry
  - Versioning and `.openkata-manifest.json`
- **Architecture**
  - Our agent workflow (moved in spec 0023)
  - The infrastructure diagram (moved in spec 0023)
- **Contributing**
  - Authoring skills, rules and profiles
  - Release workflow, one page each for a new version of the
    site, of a skill, rule or profile, and of the MCP server
  - Repository conventions

Pages may ship one section per PR while the section is still
local-only.

### Last-Updated Date

- Every docs page, including the docs home and the section
  indexes, shows when it was last updated at the bottom, below
  the previous/next links, for example `Last updated 25 October
  2026`
- The date is the author date of the last commit that touched the
  page's file, so an unrelated commit does not change it. It is
  not written in the page
- Locally the date is read from git. A page git has no date for,
  such as a new uncommitted file, shows `Last updated: not
  committed yet` in the same place and style, so the page looks
  the same everywhere
- In production the dates are stamped at deploy time from the full
  git history, embedded in the binary, and not committed; CI never
  pushes a commit for them. A deploy fails if any page has no
  date, so production never shows the placeholder

### Go-Live

- Decision, 2026-10-30: go-live does not wait for spec 0025 or for
  the benchmarks and version list of spec 0026 to show on the live
  site. Diagrams stay as they are until 0025 decides how they are
  rendered, and the benchmarks and the missing versions are handled
  separately
- The pages above did not wait for go-live; they were written while
  the section was local-only
- `/docs/` moves out of the local-only block and is served in
  production from the embedded tree, with the `Site` built
  once and cached
- A "Docs" link is added to the main nav and the mobile nav
  menu
- `search.js` loads `/docs/search-index.json` and shows docs
  results in the overlay; `/catalog/` filters them out. This is
  task 21, built after go-live

## Open Questions

None.
