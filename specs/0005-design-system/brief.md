---
title: Design System
status: Done
depth: Quick
created: 2026-05-10
---

# Design System

Token-based CSS design system for openkata.dev with three
switchable themes and technical code-style decorations.

## Themes

Three themes, switchable via a single circle button that
cycles through them. Choice persists in localStorage.

- **japandi** — clean default. White background, dark text,
  red accent. Japanese + Scandinavian minimalism.
- **hokusai** — warm. Parchment background inspired by
  Hokusai's Great Wave woodblock print. Deep blue ink,
  vermillion red accent, soft blue borders.
- **indigo** — dark. Deep indigo background, cream text,
  red accent.

## Color Palette

Raw palette colors (available in all themes for badges):
- Red `#d42b1e` — primary accent, links
- Purple `#6b3a6b` — experimental badge
- Green `#2d3d2d` — released badge
- Yellow `#f5e03b` — pending badge
- Orange `#e8943b` — warning badge
- Black `#0a0a0a` — ink

Each theme overrides semantic tokens (ink, surface, border,
accent) while raw palette colors remain constant.

## Typography

- **Body:** Noto Sans JP — clean, supports Japanese
  characters for future kanji use
- **Monospace:** Noto Sans Mono — same family as body,
  used for headings, code, skill names

All headings use monospace. This reinforces the technical
feel without needing decorative fonts.

## Code-Style Decorations

Techniques that make the site feel like source code:

- `// ` prefix on h2 headings (muted color, via CSS
  `::before`)
- `/* */` wrapping for comment-style descriptions
  (`.comment` class)
- Blinking underscore cursor on hero title (`.cursor`
  class, CSS animation)
- Monospace for all headings and the nav brand

## Considered and Discarded

- `---` as section dividers — browser `<hr>` default
  border was impossible to fully remove across browsers.
  Dropped.
- `/* */` around nav brand — felt cluttered, removed.
- Block cursor `█` — too dominant, replaced with `_`.
- Red blinking cursor — too bright, changed to muted.
- Two-color "Open" + "Kata" brand — unnecessary
  complexity, kept single color.
- Dark indigo as default — user preferred light default.
- Vermillion/gold/gray palette (first attempt) — replaced
  with richer 6-color palette.
- CDN fonts — rejected per ADR 0007 (vendor everything).
  Currently using Google Fonts for preview; will vendor
  before production.

## File Structure

```
web/static/css/
├── tokens.css              — structural tokens (spacing,
│                             typography, layout)
├── themes/
│   ├── japandi.css          — default color palette
│   ├── hokusai.css         — warm theme
│   └── indigo.css          — dark theme
├── base.css                — reset, body, links
├── typography.css          — headings, code, cursor, hr
├── components/
│   ├── card.css            — card component
│   └── theme-switcher.css  — circle toggle button
├── layouts/
│   └── grid.css            — card grid
└── style.css               — imports + page-specific
```

## Theme Switcher

Single circle button showing the next theme's primary
color. Click cycles: japandi → hokusai → indigo → japandi.
Implemented in `web/static/js/theme.js` using
`data-theme` attribute on `<html>` and localStorage.
