---
status: Done
depth: Standard
---

# Mobile Navigation

## Story

As a mobile visitor, I want the navigation to work on
small screens without overlapping or breaking, so I can
access all pages.

## Requirements

### Hamburger Menu

- Calligraphic brush-stroke 三 (san) SVG as menu toggle
- Inline SVG with `fill: currentColor` for theme support
- Centered in the nav bar (absolute positioning)
- Visible only below 768px breakpoint
- Tapping toggles a fixed dropdown panel with nav links
- Panel overlays content (position fixed, z-index 999)
- Links centered in the panel
- Tapping a link navigates and closes the menu
- Tapping the icon again closes the menu

### Responsive Behavior

- Below 768px: hide `.nav-links`, show 三 button
- Above 768px: show `.nav-links`, hide 三 button
- Nav bar height identical at all breakpoints
- Search icon and theme switcher always visible
- No nav layout overrides on mobile (same padding,
  position, height as desktop)

### Menu Panel

- Position fixed, top: 84px
- Full-width, opaque background
- Border-bottom separator
- No box-shadow
- Does not push content down

## Out of Scope

- Desktop nav changes
- Bottom navigation bar
- Slide-out side drawer
- Search overlay mobile adjustments

Date: 2026-06-12
