---
status: Done
depth: Shallow
---

# Skill Detail Page

## Story

As a developer browsing openkata.dev, I want to view a
skill's full documentation, files, and changelog so I can
decide whether to install it and understand how it works.

## Requirements

### Navigation

- Nav bar: brand (OPENKATA.dev) left, "Skills" link right
  (similar style to tidyrepo.vercel.app navigation)

### URL structure

- `/skills/{name}` — displays latest released version
- `/skills/{name}/{version}` — displays specific version
  (e.g., `/skills/commit-conventions/1.2.0`)
- `/skills/{name}/raw/{filepath}` — raw file (latest)
- `/skills/{name}/{version}/raw/{filepath}` — raw file
  (versioned)
- Version in URL has no `v` prefix
- Only released (tagged) versions are served; unreleased
  changes in version control are not visible
- Version dropdown top-right corner of page: shows
  "v1.3.0 (Latest)" and lists older versions; selecting
  navigates to that version

### Header

- 技 character at two-line height (same style as listing page)
  with "Skill" label beside it
- Skill name below the label
- Below name: `42 ↓ | .tar.gz` (download count with arrow,
  pipe, download button)
- Below that: tags, same color-coded style as listing page
- No cp command

### Tabs

Four tabs: Overview, Files, Changelog, Acknowledgments.

### Overview tab (default)

- Description text (from frontmatter) displayed first
- SKILL.md body rendered as markdown below description,
  starting after the first `# Heading` (which repeats the
  skill name)
- Strip YAML frontmatter from rendered output
- External links open in a new window (`target="_blank"`)
- Relative links to bundled files (e.g.,
  `references/commit-format.md`) navigate to the Files tab
  and display that file

### Files tab

- Display file tree mimicking `tree` command output style
  (indentation with `├──`, `└──`, `│` connectors)
- Directories are collapsed by default; clicking expands
  to show contained files
- Exclude: tile.json, .tesslignore, evals/, CHANGELOG.md,
  references/ACKNOWLEDGMENTS.md
- Clicking a file name displays its content below the tree
- File content header: file name left-aligned, toggle
  buttons right-aligned ("Preview" / "Code" / "Raw")
- "Preview": rendered markdown (default for .md files)
- "Code": syntax-highlighted source in a code block
- "Raw": opens the raw file in a new browser tab
- Non-markdown files default to "Code" view

### Changelog tab

- Render CHANGELOG.md as markdown
- Show entries from the viewed version and all prior versions
  (e.g., viewing v1.2.0 shows 1.2.0, 1.1.0, 1.0.0 — not
  1.3.0)
- If no changelog exists, show "No changelog available."

### Acknowledgments tab

- Render references/ACKNOWLEDGMENTS.md as HTML (markdown
  to HTML, same as Overview and Changelog)
- If none exists, hide the tab entirely

## Out of Scope

- Tessl registry data (quality score, eval results)
- Search or filtering within the page
- Editing or contributing from the page
- Rule detail pages (separate spec)
- Unreleased skills (no tagged version = not shown on site)

## Open Questions

- None
