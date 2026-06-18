# Tasks: Skill Detail Page

## Tasks

### 1. Add markdown rendering with goldmark
- **Status**: Done
- **Goal**: Replace the `renderMarkdown` passthrough with
  proper markdown-to-HTML using goldmark. Support headings,
  lists, code blocks, links, bold/italic.
- **Boundary**: `cmd/openkata-web/handlers.go`
- **Depends**: None
- **Done when**: Markdown renders as HTML, external links
  get `target="_blank"`, frontmatter is stripped

### 2. Update SkillDetail type and data loading
- **Status**: Done
- **Goal**: Add tags, version list, and file content map to
  the data model. Filter excluded files. Load content from
  git tags (local) or S3 (prod). Only serve released versions.
- **Boundary**: `cmd/openkata-web/templates/types.go`,
  `cmd/openkata-web/handlers.go`
- **Depends**: 1
- **Done when**: Handler serves versioned skill data with
  filtered file list and per-file content

### 3. Add versioned routing
- **Status**: Done
- **Goal**: Handle `/skills/{name}`, `/skills/{name}/{version}`,
  and `/skills/{name}/{version}/raw/{filepath}` routes.
  Latest redirects to highest released version.
- **Boundary**: `cmd/openkata-web/handlers.go`
- **Depends**: 2
- **Done when**: All URL patterns resolve correctly, raw
  serves plain text in new tab

### 4. Rewrite skill detail template — header
- **Status**: Done
- **Goal**: 技 character, skill name, `0 ↓ | .tar.gz` line,
  color-coded tags, version dropdown top-right
- **Boundary**: `cmd/openkata-web/templates/skill_detail.templ`
- **Depends**: 2
- **Done when**: Header renders with all elements, dropdown
  navigates between versions

### 5. Rewrite skill detail template — tabs and Overview
- **Status**: Done
- **Goal**: Four tabs (Overview, Files, Changelog,
  Acknowledgments). Overview shows description then rendered
  body (first heading stripped). Relative links navigate to
  Files tab and auto-open the file.
- **Boundary**: `cmd/openkata-web/templates/skill_detail.templ`
- **Depends**: 1, 4
- **Done when**: Overview renders correctly, relative links
  switch to Files tab

### 6. Implement Files tab with tree and file viewer
- **Status**: Done
- **Goal**: Tree display with `├──`/`└──` connectors,
  collapsible directories, file viewer with Preview/Code/Raw
  toggle. Non-markdown defaults to Code.
- **Boundary**: `cmd/openkata-web/templates/skill_detail.templ`
- **Depends**: 2, 5
- **Done when**: Tree renders, directories expand, files
  display with three view modes

### 7. Implement Changelog and Acknowledgments tabs
- **Status**: Done
- **Goal**: Changelog filtered to viewed version and prior.
  Acknowledgments tab hidden when file doesn't exist.
- **Boundary**: `cmd/openkata-web/templates/skill_detail.templ`,
  `cmd/openkata-web/handlers.go`
- **Depends**: 1, 5
- **Done when**: Changelog shows correct entries per version,
  Acknowledgments tab absent when empty

### 8. Add Skills link to navigation
- **Status**: Done
- **Goal**: Add "Skills" to nav bar on all pages
- **Boundary**: `cmd/openkata-web/templates/layout.templ`
- **Depends**: None
- **Done when**: Nav shows Skills link, links to `/skills/`

### 9. CSS styling
- **Status**: Done
- **Goal**: Style header, tabs, file tree, file viewer,
  version dropdown, code blocks. Match existing design system.
- **Boundary**: `web/static/css/style.css`
- **Depends**: 4, 5, 6, 7
- **Done when**: Page looks consistent with listing page

## Progress Log

- [2026-06-18] All tasks confirmed complete. Skill detail
  page live with goldmark rendering, versioned routing,
  file tree viewer, changelog filtering, and full styling.
