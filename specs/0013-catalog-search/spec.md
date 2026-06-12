---
status: Done
depth: Standard
---

# Catalog Search

## Story

As a visitor, I want a global search that finds skills,
rules, and profiles by name, description, or tag, so I
can quickly locate artifacts without browsing each
listing page.

## Requirements

### Search Icon

- Inline SVG (from Wikimedia Search_Icon.svg) in the
  nav bar, aligned right next to the theme switcher
- Uses `fill: currentColor` to inherit theme colors
- Always visible on all pages (including mobile)
- Clicking opens the search overlay

### Overlay

- Full-screen overlay with semi-transparent dark backdrop
- Text input auto-focused on open
- Live filtering: results update as the user types
- Max 7 results shown, compact format: artifact name +
  text type badge (`skill` / `rule` / `profile`)
- "View all N results" link at bottom navigates to
  `/catalog/?q=...`
- Matched substring highlighted with red background
- Tag-filter tokens highlight the matched tag itself
- Clicking a result navigates to that artifact's detail
  page
- Pressing Enter navigates to `/catalog/?q=...`
- Closes on Escape, clicking backdrop, or close button
- No keyboard shortcuts to open (icon click only)

### Catalog Page (`/catalog/`)

- Server renders HTML shell (search input + empty
  results container); no nav link — reachable only via
  search
- JS fetches `/static/versions.json` eagerly on page
  load and populates results client-side
- Search input with same live filtering as overlay
- Reads `?q=` query param on load to pre-fill and show
  results
- Empty query shows all artifacts alphabetically by name
  (full catalog)
- Results: flat list ranked by match score (no grouping
  by type)
- Each result: name + description (with highlight) +
  tags (with highlight) + text type badge
- No download counts or version numbers in results
- No version numbers anywhere in catalog/overlay results
- Matched substrings highlighted with red background

### Search Logic

- Data source: `/static/versions.json` (name,
  description, tags) — served as embedded static asset
- On non-catalog pages: lazy-fetch on first overlay open,
  cached in JS variable
- On `/catalog/`: eager fetch on page load
- Query split by spaces into tokens, all AND'd
- Token with colon (no spaces) → exact tag match
  (e.g., `category:conventions`)
- Token without colon → case-insensitive substring
  search across name + description + tags
- Scoring: name match (highest) > tag match >
  description match (lowest)
- Results ordered by total score descending
- Purely client-side — vanilla JS, no htmx, no server
  endpoint

### Tag Links

- Tags on listing pages and detail pages become `<a>`
  links
- Visual style unchanged from current spans; add cursor
  pointer + subtle hover effect
- Navigate to `/catalog/?q=category:conventions` (the
  literal tag value)

## Out of Scope

- Full-text search over SKILL.md/RULE.md/PROFILE.md
  file contents
- Server-side search or indexing
- Search history or suggestions/autocomplete
- Fuzzy matching (exact substring matching only)
- Pagination of catalog results
- Mobile nav redesign (separate task)
- Keyboard shortcuts to open search
- Download counts or versions in search results

## Open Questions

None.

Date: 2026-06-11
