# Tasks: Skills Listing Page

## Tasks

### 1. Add tag and text helpers
- **Status**: Done
- **Goal**: Split comma-separated tags, map prefixes to CSS
  classes, extract first sentence from descriptions
- **Boundary**: `cmd/openkata-web/templates/helpers.go`
- **Depends**: None
- **Done when**: `SplitTags()`, `TagClass()`,
  `FirstSentence()`, `AfterFirstSentence()` compile

### 2. Update skills template
- **Status**: Done
- **Goal**: Two-row summary (name + tags) with expandable
  description, "Open" button, bold first sentence with
  conditional rest-of-text display
- **Boundary**: `cmd/openkata-web/templates/skills.templ`
- **Depends**: 1
- **Done when**: Template renders correctly with
  expand/collapse

### 3. Add CSS for tags, button, and description
- **Status**: Done
- **Goal**: Border-only color-coded badges, solid "Open"
  button, description with `>` prefix, left border,
  Inconsolata font, bold first sentence
- **Boundary**: `web/static/css/style.css`,
  `web/static/css/tokens.css`
- **Depends**: None
- **Done when**: Visual styling matches requirements

### 4. Add Inconsolata font
- **Status**: Done
- **Goal**: Load Inconsolata from Google Fonts for
  description text
- **Boundary**: `cmd/openkata-web/templates/layout.templ`
- **Depends**: None
- **Done when**: Font loads and renders in descriptions

## Progress Log

- [2026-05-16] All tasks completed. Skills listing page
  renders with color-coded tags, expand/collapse
  descriptions, and "Open" button linking to detail pages.
- [2026-05-17] Refined description styling: Inconsolata
  medium font, `>` prefix with left border, bold first
  sentence, conditional text truncation.
