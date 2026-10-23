# Tasks: Documentation Section

## Tasks

### 1. Move architecture docs into the embedded tree
- **Status**: Done
- **Goal**: `web/docs/` exists as an embeddable package with
  the four sections and the two moved documents
- **Boundary**: `web/docs/**`, `docs/agent-workflow.md`,
  `docs/infrastructure-diagram.md`
- **Key files**: `web/static/embed.go`
- **Depends**: None
- **Done when**: `git mv` of both documents into
  `web/docs/architecture/` with `title`, `description` and
  `order` frontmatter; `web/docs/index.md` and an `index.md`
  stub per section with `title` and `order`; `embed.go`
  (package `docs`) exports `FS`; no references to the old
  paths remain
- **Verify**: `go build -o bin/openkata-web ./cmd/openkata-web/`
  and `grep -rn "docs/agent-workflow\|docs/infrastructure-diagram"
  --exclude-dir=.git --exclude-dir=specs .` returns nothing

### 2. Loader: frontmatter, model and ordering
- **Status**: Done
- **Goal**: `buildSite(fs.FS)` builds sections and pages with
  URLs and ordering, and fails strictly on structural errors
- **Boundary**: `cmd/openkata-web/docs.go`,
  `cmd/openkata-web/docs_test.go`
- **Key files**: `cmd/openkata-web/handlers.go`
  (`stripFrontmatter`), `rules/go-testing/RULE.md`
- **Depends**: None
- **Done when**: `Page`, `Section` and `Site` exist as in the
  design; the title falls back to the first `# ` heading;
  unordered pages sort after ordered ones, then by file name;
  malformed frontmatter, URL collisions and a missing section
  `index.md` return errors naming the file; table tests with
  `fstest.MapFS` cover each case
- **Verify**: `go test ./cmd/openkata-web/ -run Docs`

### 3. Loader: rendering, link rewriting and headings
- **Status**: Done
- **Goal**: Each page's body is rendered HTML with docs links
  rewritten, Mermaid blocks marked and headings collected
- **Boundary**: `cmd/openkata-web/docs.go`,
  `cmd/openkata-web/docs_test.go`
- **Key files**: `renderMarkdownFull`,
  `addTargetBlankToExternalLinks`, `stripFirstH1` in
  `cmd/openkata-web/handlers.go`
- **Depends**: 2
- **Done when**: goldmark with Table, unsafe HTML and auto
  heading IDs; relative `.md` links (with fragments) resolve
  to `/docs/...` URLs; links that leave the tree or match no
  page fail `buildSite` with a message pointing to absolute
  GitHub URLs; `language-mermaid` blocks get `mermaid-source`
  on the `<pre>`; `h2` and `h3` text lands in `Headings`;
  table tests cover each rule
- **Verify**: `go test ./cmd/openkata-web/ -run Docs`

### 4. Handler, routing and search index
- **Status**: Done
- **Goal**: `/docs/` serves pages locally, with 404, redirect
  and search index behavior from the design
- **Boundary**: `cmd/openkata-web/docs.go`,
  `cmd/openkata-web/docs_test.go`, `cmd/openkata-web/main.go`
- **Key files**: `handleStats` and the local-only block in
  `main.go`
- **Depends**: 3
- **Done when**: `handleDocs` is registered inside the
  `AWS_LAMBDA_FUNCTION_NAME == ""` block and rebuilds the site
  from `os.DirFS("web/docs")` per request; a `buildSite` error
  returns 500 with its text; unknown paths 404; known pages
  without a trailing slash 301; `/docs/search-index.json`
  returns the JSON index; `httptest` table tests cover home,
  section, page, 404, 301, traversal attempts (`..`, encoded
  `%2e%2e`, `//`) and the index body. The handler renders
  plain HTML until task 5 adds templates
- **Verify**: `go test ./cmd/openkata-web/...`

### 5. Templates: docs home and page with sidebar
- **Status**: Done
- **Goal**: Docs pages render inside the site layout with a
  sidebar and previous/next links
- **Boundary**: `cmd/openkata-web/templates/docs.templ`,
  `cmd/openkata-web/templates/docs_templ.go`,
  `cmd/openkata-web/docs.go`
- **Key files**: `cmd/openkata-web/templates/layout.templ`,
  `cmd/openkata-web/templates/skill_detail.templ` (prev/next)
- **Depends**: 4
- **Done when**: `DocsHome` lists the four sections; `DocsPage`
  shows the title, body, a `<details>` sidebar with every
  section and page and the current page highlighted, and
  previous/next links within the section; handler tests
  assert the highlight and prev/next links
- **Verify**: `templ generate ./cmd/openkata-web/templates/`
  then `go test ./cmd/openkata-web/...`

### 6. Styles and sidebar script
- **Status**: Done
- **Goal**: The docs layout works in both themes and on narrow
  screens
- **Boundary**: `web/static/css/components/docs.css`,
  `web/static/css/style.css`, `web/static/js/docs.js`,
  `cmd/openkata-web/templates/docs.templ`
- **Key files**: existing files in `web/static/css/components/`,
  `web/static/css/themes/`
- **Depends**: 5
- **Done when**: two-column layout on wide screens, one column
  on narrow ones; only design tokens are used; `pre` scrolls
  inside its container; `docs.js` opens the sidebar on wide
  screens with `matchMedia` and is loaded only by the docs
  template; checked by hand in japandi and indigo at phone and
  desktop widths, including the infrastructure diagram page
- **Verify**: run the server locally and open `/docs/`

### 7. Breadcrumb navigation
- **Status**: Done
- **Goal**: Every docs page shows a breadcrumb trail at the top
- **Boundary**: `cmd/openkata-web/docs_handler.go`,
  `cmd/openkata-web/docs_templates_test.go`,
  `cmd/openkata-web/templates/docs.templ`,
  `cmd/openkata-web/templates/docs_types.go`,
  `web/static/css/components/docs.css`
- **Key files**: `docsSidebar` and `docsPager` in
  `cmd/openkata-web/templates/docs.templ`
- **Depends**: 5, 6
- **Done when**: section index and page views show
  `Docs > Section > Page` above the layout, spanning the full
  width; every item but the last is a link and the last is
  plain text with `aria-current="page"`; the home page has no
  breadcrumb; the trail wraps on narrow screens without
  horizontal scroll; tests cover the home, section index and
  page trails
- **Verify**: `templ generate ./cmd/openkata-web/templates/`
  then `go test ./cmd/openkata-web/...`; check by hand in both
  themes at phone and desktop widths

### 8. Real-tree test and final checks
- **Status**: Pending
- **Goal**: CI fails if the committed docs tree is invalid
- **Boundary**: `cmd/openkata-web/docs_test.go`
- **Key files**: `web/docs/embed.go`
- **Depends**: 1, 3
- **Done when**: a test runs `buildSite(docs.FS)` and asserts
  no error, four sections, and non-empty titles; templ output
  is in sync; build and tests pass
- **Verify**: `templ generate ./cmd/openkata-web/templates/`,
  `go build -o bin/openkata-web ./cmd/openkata-web/`,
  `go test ./cmd/openkata-web/...`, and `git status` shows no
  changed `_templ.go` files

## Progress Log
