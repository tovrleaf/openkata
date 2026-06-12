# Tasks: Catalog Search

## Tasks

### 1. Serve versions.json as static asset
- **Status**: Done
- **Goal**: Add `versions.json` to the embedded static
  files so it's available at `/static/versions.json`
- **Boundary**: `web/static/embed.go`
- **Depends**: None
- **Done when**: Fetching `/static/versions.json` in
  the browser returns the JSON data; build passes
- **Verify**: `go build -o bin/openkata-web ./cmd/openkata-web/`

### 2. Add search icon to nav
- **Status**: Done
- **Goal**: Inline SVG search icon in the nav bar,
  right-aligned next to the theme switcher, using
  `fill: currentColor` for theme colors
- **Boundary**: `cmd/openkata-web/templates/nav.templ`,
  `web/static/css/style.css`
- **Depends**: None
- **Done when**: Icon visible on all pages, inherits
  theme color, positioned next to theme switcher
- **Verify**: `templ generate && go build -o bin/openkata-web ./cmd/openkata-web/`

### 3. Build search overlay
- **Status**: Done
- **Goal**: Clicking the search icon opens a full-screen
  overlay with backdrop, auto-focused input, close
  button. Closes on Escape, backdrop click, or close
  button.
- **Boundary**: `web/static/js/search.js`,
  `web/static/css/style.css`,
  `cmd/openkata-web/templates/nav.templ` (overlay HTML)
- **Depends**: 2
- **Done when**: Overlay opens/closes correctly, input
  is auto-focused, Escape and backdrop click close it
- **Verify**: Manual browser test

### 4. Implement search logic and overlay results
- **Status**: Done
- **Goal**: Lazy-fetch `versions.json` on first overlay
  open, filter with AND'd tokens (colon = exact tag,
  else substring on name/description/tags), score
  (name > tag > description), show top 7 results with
  name + type badge, highlight matched substrings in
  red, show "View all N results" link
- **Boundary**: `web/static/js/search.js`
- **Depends**: 1, 3
- **Done when**: Typing in overlay shows filtered,
  ranked, highlighted results; "View all" links to
  `/catalog/?q=...`; clicking result navigates to
  detail page; Enter navigates to catalog
- **Verify**: Manual browser test

### 5. Add /catalog/ page (server shell)
- **Status**: Done
- **Goal**: Add route handler and templ template for
  `/catalog/` — renders an HTML shell with search input
  and empty results container
- **Boundary**: `cmd/openkata-web/main.go`,
  `cmd/openkata-web/handlers.go`,
  `cmd/openkata-web/templates/catalog.templ`
- **Depends**: None
- **Done when**: `/catalog/` renders a page with search
  input; build passes
- **Verify**: `templ generate && go build -o bin/openkata-web ./cmd/openkata-web/`

### 6. Implement catalog page JS
- **Status**: Done
- **Goal**: On `/catalog/`, eagerly fetch `versions.json`,
  read `?q=` param, populate results with same search
  logic as overlay. Live filtering from the page's
  search input. Full result format: name + description
  + tags + type badge, all with highlights. Empty query
  shows all artifacts alphabetically.
- **Boundary**: `web/static/js/search.js`
- **Depends**: 4, 5
- **Done when**: Catalog page shows results from query
  param, live-filters, highlights matches, empty shows
  all alphabetically
- **Verify**: Manual browser test

### 7. Make tags clickable links
- **Status**: Done
- **Goal**: Convert tag `<span>` elements to `<a>` links
  on listing and detail pages. Same visual style, add
  hover effect. Navigate to `/catalog/?q=tag-value`.
- **Boundary**: `cmd/openkata-web/templates/skills.templ`,
  `cmd/openkata-web/templates/rules.templ`,
  `cmd/openkata-web/templates/profiles.templ`,
  `cmd/openkata-web/templates/skill_detail.templ`,
  `cmd/openkata-web/templates/rule_detail.templ`,
  `cmd/openkata-web/templates/profile_detail.templ`,
  `web/static/css/style.css`
- **Depends**: 5
- **Done when**: Clicking a tag navigates to catalog
  with that tag as query; hover effect visible; visual
  style unchanged otherwise
- **Verify**: `templ generate && go build -o bin/openkata-web ./cmd/openkata-web/`

### 8. Add tests
- **Status**: Done
- **Goal**: Test `/catalog/` handler serves the page,
  test `versions.json` is served from static. Verify
  tag links render as `<a>` elements with correct href.
- **Boundary**: `cmd/openkata-web/handlers_test.go`
- **Depends**: 1, 5, 7
- **Done when**: Tests pass covering catalog route,
  static versions.json serving, and tag link rendering
- **Verify**: `go test -race ./cmd/openkata-web/...`

## Progress Log

- 2026-06-12: All 8 tasks implemented — embed, icon, overlay, search logic, catalog page, catalog JS, tag links, tests
