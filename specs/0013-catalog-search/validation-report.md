# Validation: Catalog Search

## Requirements Check

### Search Icon

- [x] Inline SVG in nav bar — verified in `nav.templ`:
  inline `<svg>` with viewBox matching Wikimedia
  Search_Icon.svg pattern
- [x] Uses `fill: currentColor` — verified: SVG element
  has `fill="currentColor"` attribute
- [x] Aligned right next to theme switcher — verified:
  search button and theme-switcher are siblings inside
  `.nav-actions` with `margin-left: auto`
- [x] Always visible on all pages (including mobile) —
  verified: `Nav()` is included in all page templates
  (Layout wraps all pages), no responsive hide rules on
  `.search-btn`
- [x] Clicking opens the search overlay — verified:
  `openBtn.addEventListener('click', openSearch)` in
  `search.js`

### Overlay

- [x] Full-screen overlay with semi-transparent dark
  backdrop — verified: `.search-overlay` uses
  `position: fixed; inset: 0`, backdrop uses
  `rgba(0, 0, 0, 0.6)`
- [x] Text input auto-focused on open — verified:
  `searchInput.focus()` called in `openSearch()`
- [x] Live filtering: results update as user types —
  verified: `searchInput.addEventListener('input', ...)`
  calls `renderOverlayResults`
- [x] Max 7 results shown — verified: `var max = 7` in
  `renderOverlayResults`, slices to `results.slice(0, max)`
- [x] Compact format: name + text type badge — verified:
  renders `search-overlay__item-name` +
  `search-overlay__item-type` (type is singular:
  `item.type` which is `type.slice(0, -1)` = skill/rule/
  profile)
- [x] "View all N results" link at bottom navigates to
  `/catalog/?q=...` — verified: renders
  `search-overlay__view-all` with
  `href="/catalog/?q=" + q`
- [x] Matched substring highlighted with red background —
  verified: `<mark class="search-highlight">` with CSS
  `background: rgba(220, 38, 38, 0.25)` (red at 25%
  opacity)
- [x] Tag-filter tokens highlight the matched tag itself
  — verified: `highlightTag()` wraps entire tag in
  `<mark class="search-highlight">` when tag token
  matches
- [x] Clicking a result navigates to detail page —
  verified: each result is an `<a href="item.path">`
  where path is `'/' + type + '/' + name + '/'`
- [x] Pressing Enter navigates to `/catalog/?q=...` —
  verified: `searchInput.addEventListener('keydown', ...)`
  checks `e.key === 'Enter'` and sets
  `window.location.href`
- [x] Closes on Escape — verified:
  `document.addEventListener('keydown', ...)` checks
  `e.key === 'Escape'`
- [x] Closes on clicking backdrop — verified:
  `backdrop.addEventListener('click', closeSearch)`
- [x] Closes on close button — verified:
  `closeBtn.addEventListener('click', closeSearch)`
- [x] No keyboard shortcuts to open (icon click only) —
  verified: no keydown listener that opens the overlay;
  only `openBtn.addEventListener('click', openSearch)`

### Catalog Page

- [x] Server renders HTML shell (search input + empty
  results container) — verified: `catalog.templ` renders
  `<input id="catalog-input">` and
  `<div id="catalog-results">` with no server-side
  content
- [x] No nav link — reachable only via search — verified:
  `nav.templ` has no link to `/catalog/`
- [x] JS fetches `/static/versions.json` eagerly on page
  load — verified: catalog page branch calls
  `fetchData(function() { ... })` immediately on load
  (not gated by overlay open)
- [x] Search input with same live filtering as overlay —
  verified: `catalogInput.addEventListener('input', ...)`
  calls `renderCatalogResults` which uses same `search()`
  function
- [x] Reads `?q=` query param on load to pre-fill and
  show results — verified:
  `new URLSearchParams(window.location.search).get('q')`
  sets `catalogInput.value` and calls
  `renderCatalogResults(q)`
- [x] Empty query shows all artifacts alphabetically —
  verified: `search()` with empty tokens returns
  `catalogData.slice()` which was sorted by
  `name.localeCompare(b.name)` in `parseData()`
- [x] Results: flat list ranked by match score (no
  grouping by type) — verified: single sorted array
  rendered sequentially
- [x] Each result: name + description + tags + text type
  badge — verified: `renderCatalogResults` renders
  `catalog-item__name`, `catalog-item__desc`,
  `catalog-item__tags`, `catalog-item__type`
- [x] No download counts in results — verified: no
  download count rendered in catalog or overlay items
- [x] No version numbers anywhere in catalog/overlay
  results — verified: version field not rendered in
  either `renderOverlayResults` or
  `renderCatalogResults`
- [x] Matched substrings highlighted with red
  background — verified: same `highlightText()` and
  `search-highlight` class used

### Search Logic

- [x] Data source: `/static/versions.json` (name,
  description, tags) — verified: `fetch('/static/versions.json')`
  in `fetchData()`; embed.go embeds `versions.json`
- [x] On non-catalog pages: lazy-fetch on first overlay
  open — verified: `openSearch()` calls `fetchData()`
  only if `!catalogData && !fetching`
- [x] On `/catalog/`: eager fetch on page load —
  verified: catalog page branch immediately calls
  `fetchData(cb)` without waiting for user interaction
- [x] Query split by spaces into tokens, all AND'd —
  verified: `tokenize()` splits by `/\s+/`;
  `scoreItem()` returns 0 if any token has no match
  (short-circuit `return 0`)
- [x] Token with colon → exact tag match — verified:
  `isTagToken()` checks `indexOf(':') > 0`; matching
  uses `tags[j] === token.toLowerCase()` (exact)
- [x] Token without colon → case-insensitive substring
  across name + description + tags — verified:
  `item.name.toLowerCase().indexOf(lower)`,
  `item.tags.toLowerCase().indexOf(lower)`,
  `item.description.toLowerCase().indexOf(lower)`
- [x] Scoring: name match (highest) > tag match >
  description match — verified: name = 100, tag = 50,
  description = 25
- [x] Results ordered by total score descending —
  verified: `results.sort(function(a, b) { return b.score - a.score; })`
- [x] Purely client-side — vanilla JS, no htmx, no
  server endpoint — verified: `search.js` is vanilla JS
  IIFE; no server search endpoint; `handleCatalog` only
  renders the shell

### Tag Links

- [x] Tags on listing pages become `<a>` links —
  verified: `skills.templ`, `rules.templ`,
  `profiles.templ` all render tags as
  `<a href="/catalog/?q=tag" class="badge ...">`
- [x] Tags on detail pages become `<a>` links —
  verified: `skill_detail.templ`, `rule_detail.templ`,
  `profile_detail.templ` all render tags as
  `<a href="/catalog/?q=tag" class="badge ...">`
- [x] Visual style unchanged — uses same badge classes —
  verified: `TagClass()` returns `"badge badge-green"`,
  `"badge badge-orange"`, `"badge badge-purple"`; CSS
  `a.badge` has `text-decoration: none`
- [x] Cursor pointer + subtle hover effect — verified:
  `a.badge` inherits pointer from being an anchor;
  `a.badge:hover { opacity: 0.8 }` provides hover
  effect
- [x] Navigate to `/catalog/?q=category:conventions`
  (the literal tag value) — verified: href uses
  `fmt.Sprintf("/catalog/?q=%s", tag)` — the raw tag
  value

## Out of Scope Verified

- Confirmed: no full-text search over SKILL.md/RULE.md/
  PROFILE.md file contents (search only uses
  versions.json fields)
- Confirmed: no server-side search or indexing (purely
  client-side JS)
- Confirmed: no search history or
  suggestions/autocomplete
- Confirmed: no fuzzy matching (uses `indexOf` for
  exact substring)
- Confirmed: no pagination of catalog results (all
  results rendered)
- Confirmed: no mobile nav redesign
- Confirmed: no keyboard shortcuts to open search
- Confirmed: no download counts or versions in search
  results

## Issues Found

None. All requirements are implemented correctly and
tests pass (`go test -count=1 ./cmd/openkata-web/...`
exits 0).

Date: 2026-06-12
