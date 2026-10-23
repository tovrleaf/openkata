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
  - Skills, rules and profiles: what each is and how they
    differ, with the kata vocabulary (dojo kun, sensei
    profiles, the dojo, the ryu)
  - Installing artifacts through the MCP server
  - Using the catalog
- **Reference**
  - How MCP discovery works: the `list_*` tools and tag
    filter, and where the index lives (`versions.json`)
  - Where skills, rules and profiles are stored and how
    `versions.json` is generated
  - Downloads: how downloads are recorded and counted
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

### Go-Live

- `/docs/` moves out of the local-only block and is served in
  production from the embedded tree, with the `Site` built
  once and cached
- A "Docs" link is added to the main nav and the mobile nav
  menu
- `search.js` loads `/docs/search-index.json` and shows docs
  results in the overlay; `/catalog/` filters them out

## Open Questions

- The stats dashboard (`/stats/`) is local-only. Should the
  downloads page describe it, or only the public download
  counts?
- How much of the S3 bucket layout and AWS setup should be
  public?
- Is rendered Mermaid needed at go-live, or is the source view
  enough? Spec 0025 evaluates diagram rendering and static
  asset delivery; decide whether go-live waits for it
