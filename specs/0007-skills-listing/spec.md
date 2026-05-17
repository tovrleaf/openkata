---
status: Done
depth: Shallow
---

# Skills Listing Page

## Story

As a visitor, I want to browse available skills with their
tags and metadata so I can quickly find relevant kata and
open their detail pages.

## Requirements

- List all skills at `/skills/` in a vertical stack
- Each entry shows: number, name, download count, version,
  and a +/- toggle icon
- Tags displayed as a second row, always visible
- Tags color-coded by prefix: category (green), tool
  (orange), language (purple)
- Tags styled as border-only badges (transparent background,
  colored border and text, monospace, compact)
- "Open" button right-aligned in the tags row, solid black
  background with white text, links to detail page
- Clicking anywhere on the entry expands/collapses the
  description
- Tags remain in place when expanded/collapsed
- Description styled with `>` prefix, left border, and
  Inconsolata medium (500) font
- First sentence of description rendered bold (weight 900)
- If first sentence exceeds 200 characters, only show that;
  otherwise show rest of text after a line break
- Empty state shows "No skills available yet."
- Skill name is plain text, not a link

## Out of Scope

- Install command display
- Tag filtering or search
- Sorting controls
- Pagination

## Open Questions

- None
