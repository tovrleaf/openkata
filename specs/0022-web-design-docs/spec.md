---
status: Draft
depth: Quick
---

# Web Page Design Documentation

## Story

As the frontend-developer agent, I need reference
documents describing the composition and behavior of
each page type on openkata.dev so I can make template
changes without trial-and-error guidance from the user.

## Requirements

### Design Reference Files

- Create `docs/design/` directory for page blueprints
- Document each page type: listing, detail, catalog,
  home, getting-started
- For each page describe: component hierarchy, link
  behavior, conditional sections, third-party branding
  rules

### Skill Detail Page

- Header section: kanji, name, version select, nav
- Meta section: downloads, archive link, tags
- Tab bar: which tabs appear conditionally
- Benchmarks section: tessl score placement, logo
  sizing, link targets, effectiveness table, scenarios
- External links: always open in new tab, use registry
  URL pattern

### Listing Pages

- Summary row: number, name, badges, version
- Effectiveness badge: when shown, format, position
- Expanded body: description, per-model breakdown
- Tag links: route to catalog with query param

### Agent Integration

- frontend-developer profile updated to reference
  design docs before modifying templates
- Design docs are living references (not archived
  after implementation)

## Constraints

- Docs describe current state, not aspirational design
- Keep each page doc under 100 lines
- Co-located in `docs/design/`, not in .agents/

## Acceptance Criteria

1. `docs/design/` exists with page-type documents
2. Skill detail page fully documented (components,
   conditionals, link behavior)
3. Listing page documented (badges, expand behavior)
4. frontend-developer profile references design docs
5. Agent can implement a new detail page section
   without asking where to put it

## Out of Scope

- CSS token system (covered by design-system rule)
- Visual mockups or screenshots
- New page types not yet built

Date: 2026-06-22
