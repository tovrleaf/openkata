---
status: Done
depth: Standard
---

# Profile Pages

## Story

As a visitor, I want to browse and read agent profiles
on openkata.dev so I can understand the available agent
roles and their scoped permissions.

## Requirements

### Listing Page (`/profiles/`)

- Hero section: 師 kanji at two-line height with "Profiles"
  heading and subtitle "Agent role definitions with scoped
  permissions."
- List all profiles in a vertical stack (same component
  as skills/rules listing)
- Each entry shows: sequential number, name, download
  count with ↓, version, +/- toggle icon
- Tags displayed as second row, color-coded by prefix
  (same scheme as skills/rules)
- "Open" button right-aligned in tags row, links to
  detail page (`/profiles/:name/`)
- Clicking entry expands/collapses description
- Description styled same as skills: bold first sentence,
  rest after line break
- Empty state: "No profiles available yet."
- Sorted alphabetically

### Detail Page (`/profiles/:name`)

- URL structure:
  - `/profiles/{name}` — latest released version
  - `/profiles/{name}/{version}` — specific version
  - `/profiles/{name}/raw` — raw file (latest)
  - `/profiles/{name}/{version}/raw` — raw file (versioned)
  - Version in URL has no `v` prefix
- Header: 師 kanji with "Profile" label, profile name below
- Version dropdown (when multiple versions exist)
- Download count + .tar.gz archive link
- Tags row (color-coded, same as listing)
- Tabs: Overview, Changelog
  (No Files tab — profiles are single files;
  no Acknowledgments — profiles have none)

### Overview tab (default)

- Description text displayed first
- Profile markdown body rendered (strip frontmatter if
  any, skip first `# Heading`)
- External links open in new window

### Changelog tab

- Render CHANGELOG.md as markdown if it exists
- "No changelog available." if missing

### Archive/Download

- `/profiles/{name}/archive` — download latest .tar.gz
- `/profiles/{name}/archive/{version}` — download
  specific version
- Archive routing already exists in `handleProfiles`;
  verify it works end-to-end locally

### Navigation

- Add "Profiles" link to nav bar, to the right of "Rules"

### Code Architecture

- Refactor `SkillDetail`, `RuleDetail` into a single
  `ArtifactDetail` concrete struct with an `Type` field
  ("skills", "rules", "profiles")
- Remove field-by-field copy in `loadRuleDetailVersion`
- Shared loaders return `*ArtifactDetail` directly
- Templates accept `ArtifactDetail` instead of
  type-specific structs
- Profile detail loader uses shared helpers where
  possible, with simplified handling for single-file
  structure (no file tree, no acknowledgments)
- Profile filename stays as `{name}.md` (not PROFILE.md)
  for drag-and-drop usability after download

### Deployment Verification

- Confirm publish workflow handles profiles correctly
- Profiles must be uploaded to S3 and visible in
  production after tagging
- Profile directory structure: `profiles/{name}/{name}.md`
  (matches git archive expectations)

### Testing

- Table-driven handler tests for profiles (same pattern
  as rules)
- Test routing: listing 200, detail 200, version 200,
  raw file 200, unknown profile 404
- Test archive download works locally
- Regression tests after ArtifactDetail refactor:
  - Skills listing still returns 200
  - Skill detail page still renders
  - Skill archive download works
  - Rules listing still returns 200
  - Rule detail page still renders
  - Rule archive download works
  - MCP listing (if applicable) still works
- All existing tests must pass unchanged or be
  updated to use new type

## Out of Scope

- Tag filtering or search
- Sorting controls
- Pagination
- File tree (profiles are single files)


