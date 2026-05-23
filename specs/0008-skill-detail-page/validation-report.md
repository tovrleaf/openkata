# Validation: Skill Detail Page

## Requirements Check

### Navigation

- [x] Nav bar: brand (OPENKATA.dev) left — verified in
  `skill_detail.templ` line 8: `<a href="/" class="nav-brand">OPENKATA.dev</a>`
- [x] "Skills" link right — verified in `skill_detail.templ`
  line 10: `<a href="/skills/" class="nav-link">Skills</a>`

### URL structure

- [x] `/skills/{name}` displays latest released version —
  verified in `handlers.go` line 119: redirects to latest
  version URL
- [x] `/skills/{name}/{version}` displays specific version —
  verified in `handlers.go` line 111
- [ ] `/skills/{name}/raw/{filepath}` raw file (latest) —
  FAILED: no route exists for this pattern; only
  `/skills/{name}/{version}/raw/{filepath}` is implemented
  (line 94)
- [x] `/skills/{name}/{version}/raw/{filepath}` raw file
  (versioned) — verified in `handlers.go` line 94
- [x] Version in URL has no `v` prefix — verified:
  `gitVersions()` strips the `v` prefix from tags
- [x] Only released (tagged) versions served (production) —
  verified: S3 loader reads from versioned prefixes populated
  by deploy pipeline; local dev mode reads working tree
  (acceptable for dev)
- [ ] Version dropdown shows "v1.3.0 (Latest)" for latest —
  FAILED: dropdown shows `v1.3.0` without "(Latest)" suffix
  (`skill_detail.templ` line 29)
- [x] Version dropdown lists older versions — verified in
  template loop over `skill.Versions`
- [x] Selecting version navigates to that version — verified
  in JS: `window.location.href = '/skills/' + ... + '/' + this.value`

### Header

- [x] 技 character at two-line height — verified: CSS
  `.hero-kanji` at `4.5rem` font-size
- [ ] "Skill" label beside kanji — FAILED: no "Skill" label
  text exists; the skill name is placed directly in the
  `hero-text` div without a preceding label
- [ ] Skill name below the label — FAILED: skill name is
  beside the kanji (in `hero-text`), not below a "Skill"
  label
- [x] Download count with arrow and download button —
  verified: `{ fmt.Sprintf("%d", skill.Downloads) } ↓` and
  `.tar.gz` link
- [x] Tags with color-coded style — verified: `SplitTags`
  and `TagClass` produce badge classes
- [x] No cp command — verified: no clipboard/copy
  functionality present

### Tabs

- [x] Four tabs: Overview, Files, Changelog,
  Acknowledgments — verified in template (Acknowledgments
  conditionally shown)
- [x] Acknowledgments tab hidden when file doesn't exist —
  verified: `if skill.Acknowledgments != ""`

### Overview tab (default)

- [x] Description text displayed first — verified:
  `detail-description` paragraph before docs
- [x] SKILL.md body rendered as markdown — verified:
  `renderMarkdown()` called on SKILL.md content
- [x] Strip YAML frontmatter — verified: `stripFrontmatter()`
  removes `---` delimited frontmatter
- [x] Strip first `# Heading` — verified: `stripFirstH1()`
  removes first `<h1>` from rendered output
- [x] External links open in new window — verified:
  `addTargetBlankToExternalLinks()` adds `target="_blank"`
- [x] Relative links navigate to Files tab — verified: JS
  intercepts clicks on non-http links in overview panel,
  calls `switchTab('files')` and `highlightFile(href)`

### Files tab

- [x] File tree with connectors (`├──`, `└──`, `│`) —
  verified in `buildFileTree()` / `renderNode()` helper
- [ ] Directories collapsed by default — FAILED: directories
  render with `<details ... open>` attribute, meaning they
  are expanded by default (helpers.go line 142)
- [x] Clicking directory expands/collapses — verified:
  `<details>` element provides native toggle
- [x] Excluded files: tile.json, .tesslignore, evals/,
  CHANGELOG.md, references/ACKNOWLEDGMENTS.md — verified in
  `isExcludedFile()`
- [x] Clicking file displays content below tree — verified:
  JS `showFile()` shows `file-viewer` div
- [x] File content header: name left, toggle buttons right —
  verified in template: `file-viewer-header` with
  `file-viewer-name` and `file-viewer-actions`
- [x] Toggle buttons: Preview / Code / Raw — verified in
  template
- [ ] Preview: rendered markdown (default for .md) — PARTIAL:
  JS checks for `fileContents['__rendered__' + currentFile]`
  but this key is never populated in the Go handler;
  `FileContents` map only stores raw text, so preview falls
  back to `<pre><code>` display for .md files
- [x] Code: source in code block — verified: JS renders
  `<pre><code>` with escaped content
- [x] Raw: opens in new browser tab — verified:
  `<a ... target="_blank" rel="noopener">Raw</a>` with href
  to raw endpoint
- [x] Non-markdown files default to Code view — verified: JS
  `showFile()` calls `setView('code')` for non-.md files

### Changelog tab

- [x] Render CHANGELOG.md as markdown — verified:
  `renderMarkdown(cl)` in handler
- [ ] Show only entries from viewed version and prior —
  FAILED: no version-based filtering of changelog content;
  the entire CHANGELOG.md is rendered regardless of which
  version is being viewed
- [x] "No changelog available." when missing — verified in
  template

### Acknowledgments tab

- [x] Render references/ACKNOWLEDGMENTS.md as HTML —
  verified: `renderMarkdown(ack)` in handler
- [x] Hide tab when file doesn't exist — verified:
  conditional rendering in template

## Out of Scope Verified

- Confirmed: no Tessl registry data (quality score, eval
  results) present
- Confirmed: no search or filtering within the page
- Confirmed: no editing or contributing from the page
- Confirmed: no rule detail pages built
- Confirmed: unreleased skills handling — production uses S3
  versioned prefixes (only deployed releases)

## Issues Found

1. **Missing `/skills/{name}/raw/{filepath}` route** — the
   spec requires a raw file URL without version (serving
   latest), but only the versioned raw route is implemented.

2. **Version dropdown missing "(Latest)" label** — spec
   requires the latest version to display as "v1.3.0
   (Latest)" but implementation shows only "v1.3.0".

3. **Missing "Skill" label beside kanji** — spec says the
   header should have a "Skill" label beside the 技
   character with the skill name below it. Implementation
   places the skill name directly beside the kanji with no
   "Skill" label.

4. **Directories expanded by default** — spec says
   "Directories are collapsed by default" but the
   `<details open>` attribute causes them to render expanded.

5. **Markdown preview broken for Files tab** — the JS looks
   for pre-rendered HTML in `fileContents['__rendered__' +
   path]` but the Go handler never populates these keys.
   Markdown files in the file viewer will display as raw text
   in "Preview" mode instead of rendered HTML.

6. **Changelog not filtered by version** — spec requires
   showing only entries from the viewed version and prior.
   The implementation renders the entire CHANGELOG.md
   regardless of which version is selected.
