---
status: Draft
depth: Standard
---

# Rules Pages

## Story

As a visitor, I want to browse and read rules on
openkata.dev so I can understand the always-on constraints
available for agent sessions.

## Requirements

### Listing Page (`/rules/`)

- Hero section: 則 kanji at two-line height with "Rules"
  heading and subtitle "Always-on constraints applied to
  every agent session."
- List all rules in a vertical stack (same component as
  skills listing)
- Each entry shows: sequential number, name, download
  count with ↓, version, +/- toggle icon
- Tags displayed as second row, color-coded by prefix
  (same scheme as skills: category green, tool orange,
  language purple)
- "Open" button right-aligned in tags row, links to
  detail page (`/rules/:name/`)
- Clicking entry expands/collapses description
- Description styled same as skills: bold first sentence,
  rest after line break
- Empty state: "No rules available yet."
- Sorted alphabetically

### Detail Page (`/rules/:name`)

- URL structure:
  - `/rules/{name}` — latest released version
  - `/rules/{name}/{version}` — specific version
  - `/rules/{name}/raw/{filepath}` — raw file (latest)
  - `/rules/{name}/{version}/raw/{filepath}` — raw file
    (versioned)
  - Version in URL has no `v` prefix
- Header: 則 kanji with "Rule" label, rule name below
- Version dropdown (when multiple versions exist)
- Download count + .tar.gz archive link
- Tags row (color-coded, same as listing)
- Tabs: Overview, Files, Changelog, Acknowledgments
  (Acknowledgments tab only if
  `references/ACKNOWLEDGMENTS.md` exists)

### Overview tab (default)

- Description text displayed first
- RULE.md body rendered as markdown (strip frontmatter,
  skip first `# Heading`)
- External links open in new window
- Relative links to bundled files navigate to Files tab

### Files tab

- File tree with `├──`/`└──`/`│` connectors
- Directories collapsed by default
- Exclude: CHANGELOG.md, references/ACKNOWLEDGMENTS.md
  (reuse existing `isExcludedFile`)
- Click file → show content below tree
- Preview/Code/Raw toggle (Preview default for .md)
- Raw opens in new browser tab

### Changelog tab

- Render CHANGELOG.md as markdown
- Filter entries: show viewed version and all prior
- "No changelog available." if missing

### Acknowledgments tab

- Render references/ACKNOWLEDGMENTS.md as markdown
- Hide tab entirely if file doesn't exist

### Archive/Download

- `/rules/{name}/archive` — download latest .tar.gz
- `/rules/{name}/archive/{version}` — download specific
  version
- Archive routing already exists in `handleRules`; wired
  up and functional

### Navigation

- Add "Rules" link to nav bar on all pages, to the right
  of "Skills"
- Extract nav into a shared `Nav()` templ component in
  `nav.templ` (currently duplicated across 6 templates)
- No active-state highlighting on nav links

### CSS Refactor

- Rename `.skill-*` CSS classes to `.artifact-*` for
  shared listing/detail components
- Update all templates referencing old class names
- Both skills and rules use `.artifact-*` classes
- Prerequisite: spec 0010 (implement first)

### Template Architecture

- Shared components in individual files:
  - `nav.templ` — Nav()
  - `tabs.templ` — TabBar(), TabPanel()
  - `file_viewer.templ` — FileTree(), FileBlocks()
- `skill_detail.templ` and `rule_detail.templ` compose
  from shared components
- No inline JavaScript — move detail page JS to
  `web/static/js/detail.js`
- Version select uses `data-path` attribute for navigation
  URL (not hardcoded to `/skills/`)
- `buildFileTree` accepts `artifactType` string parameter

### Code Architecture

- New `RuleDetail` struct (do not reuse `SkillDetail`)
- Interface for shared detail loading:
  shared internal helper populates via interface that
  both `SkillDetail` and `RuleDetail` implement
- Extract shared `loadArtifactList(ctx, artifactType)`
  from `loadSkillsList` / `loadRulesList`
- Extract shared detail loader parameterized by
  artifact type (local + S3 variants)
- Parameterize `gitVersions` to accept artifact type
  prefix (currently hardcoded to `skills/`)

### Data

- Rules data sourced from versions.json (same as skills)
- Tags present in versions.json and RULE.md frontmatter
- Download counts from DynamoDB (same as skills)
- Git tags: `rules/<name>/v<version>`
- `generate-versions` already handles rules (no changes)

### Testing

- Table-driven handler tests (same pattern as skills)
- Test routing: listing 200, detail 200, version 200,
  raw file 200, unknown rule 404
- Tests use `t.TempDir()` + fixture files (existing
  pattern)
- Test extracted shared helpers

### Release

- Tag all rules before deployment
  (`rules/<name>/v1.0.0`)
- Rules not visible in production without tags

## Out of Scope

- Tag filtering or search
- Sorting controls
- Pagination

## Open Questions

- None
