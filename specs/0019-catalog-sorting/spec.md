---
status: Draft
depth: Standard
---

# Catalog Sorting

## Story

As a developer browsing openkata.dev, I want to sort
the catalog by name, type, downloads, tokens, version,
or release date so I can find skills relevant to my
needs.

## Requirements

### Sort Options

- **Name** (A–Z) — default
- **Type** (skills → rules → profiles) — catalog page
  only
- **Downloads** (high → low)
- **Version** (highest semver first)
- **Updated** (most recently released first)

### Implementation

- Sort controls: horizontal button group above the
  artifact list
- Active sort highlighted with accent style
- Query parameter: `?sort=name`, `?sort=downloads`,
  `?sort=version`, `?sort=updated`, `?sort=type`
- Server-side sorting, no client-side JS
- Catalog search (HTMX) integrates sort parameter —
  sort applies to filtered results only
- Changing sort preserves search query; changing search
  preserves sort

### Release Date (for "Updated" sort)

- Local mode: git tag timestamp via
  `git tag -l --format='%(creatordate:iso)'`
- S3 mode: S3 object `LastModified` timestamp
- `generate-versions` writes `"updated"` field (ISO
  timestamp) into `versions.json`
- Not displayed on listing pages — used only as sort
  key (release dates visible in changelog)

### Skills Without Data

When sorting by downloads or version, all skills have
data. No special handling needed.

## Constraints

- No client-side JS for sorting
- `versions.json` provides all data needed for sorting
- Sort does not persist across sessions (no cookies)

## Acceptance Criteria

1. Sort controls visible on `/skills/`, `/rules/`,
   `/profiles/`, `/catalog/`
2. Each sort option reloads with correct order
3. Active sort visually highlighted
4. Catalog page: sort + search work together
5. "Updated" sort uses release date from git tag / S3
6. Default sort is name A–Z

## Out of Scope

- User-configurable sort persistence
- Client-side (JS) sorting
- Sorting by tokens (deferred with token tracking)
- Sorting by effectiveness (separate spec 0018)

## Open Questions

None.

Date: 2026-06-20
