# Design: Documentation Section

## Architecture

### Content layout

```text
web/docs/
  embed.go                 # package docs, embeds the pages
  index.md                 # /docs/
  guides/                  # one directory per section
    index.md
  reference/
    index.md
  architecture/
    index.md
    agent-workflow.md      # moved from docs/
    infrastructure-diagram.md
  contributing/
    index.md
```

- Each section directory has an `index.md` whose frontmatter
  sets the section `title` and `order`
- Page frontmatter: `title`, optional `description`, optional
  `order`. Pages without `title` use their first `# ` heading
- The URL of `guides/foo.md` is `/docs/guides/foo/`;
  `index.md` is the section root

### Embedding

- `web/docs/embed.go` embeds `*.md` and the four section
  directories, the same shape as `web/static/embed.go`
- The loader takes an `fs.FS`. Tests pass `fstest.MapFS`, local
  development passes `os.DirFS("web/docs")`, and the real-tree
  test passes the embedded FS. Production serves nothing yet,
  so the embedded FS is only read by that test until spec 0024

### Loading and model (`cmd/openkata-web/docs.go`)

```go
type Page struct {
    Section, URL, Title, Description string
    Order    int  // ordered pages first; see hasOrder
    hasOrder bool
    Headings []string
    Body     []byte // rendered HTML
}
type Section struct {
    Slug, Title string
    Order       int
    Index       *Page
    Pages       []*Page // ordered, then unordered by file name
}
type Site struct {
    Sections []*Section
    byURL    map[string]*Page // the only lookup used for requests
}

func buildSite(fsys fs.FS) (*Site, error)
```

- Frontmatter is parsed with `gopkg.in/yaml.v3`, already in
  `go.mod`; the existing `stripFrontmatter` is reused for the
  body
- `buildSite` renders every page up front, so link errors are
  found at load time, not on the request that hits them
- `buildSite` is strict and returns the first error, naming
  the file: malformed frontmatter, a URL collision, a section
  without `index.md`, or a relative link that does not resolve
- Pages without `order` sort after ordered pages, then by file
  name
- Locally the `Site` is rebuilt from `os.DirFS` on each
  request, so edits show without a rebuild. A `buildSite`
  error returns 500 with the error text
- Requests are resolved only through `byURL`. A path that is
  not a known page is a 404, so no file path ever comes from
  the request and traversal is impossible by construction

### Routing

- `mux.HandleFunc("/docs/", handleDocs)` in `main.go`, inside
  the existing `AWS_LAMBDA_FUNCTION_NAME == ""` block
- `/docs/` → home, `/docs/<section>/` → section index,
  `/docs/<section>/<page>/` → page,
  `/docs/search-index.json` → the search index
- Any other path without a trailing slash redirects (301) to
  the slash form if that form is a known page, otherwise 404

### Rendering

- One goldmark pipeline, configured as `renderMarkdownFull`
  configures it (Table extension, unsafe HTML for trusted repo
  content), plus `parser.WithAutoHeadingID()` for heading
  anchors
- Post-processing, using the existing helpers where they
  exist:
  - external links get `target="_blank"`
    (`addTargetBlankToExternalLinks`)
  - the first `<h1>` is stripped (`stripFirstH1`) because the
    template renders the title
  - relative links ending in `.md` (with an optional
    `#fragment`) are resolved against the current page's
    directory inside `web/docs/` and rewritten to `/docs/...`
    URLs. A link that leaves the tree or matches no page is a
    `buildSite` error whose message says to use an absolute
    GitHub URL for files outside `web/docs/`
  - `<pre><code class="language-mermaid">` gets the class
    `mermaid-source` on the `<pre>`, and the source stays as
    is
- `templates/docs.templ` provides `DocsPage(model)` and
  `DocsHome(model)` inside the existing `Layout` and `Nav`

### Navigation and layout

- No nav changes in this spec
- A breadcrumb trail (`<nav aria-label="Breadcrumb">` with an
  ordered list) sits above the two columns and spans the full
  width: `Docs` links to `/docs/`, the section links to its
  index, and the last item is the current page as plain text
  with `aria-current="page"`. A section index shows
  `Docs > Section`. The page model gets the trail as a list of
  titled links, built from the page's section
- The page is two columns: a sidebar of every section and its
  pages, and the content. The current page is highlighted, and
  previous/next links follow the section order
- The sidebar is a `<details>`. A small `docs.js`, loaded only
  by the docs template, opens it on wide screens with
  `matchMedia` and leaves it collapsed on narrow ones; without
  JavaScript it stays collapsed but every link works
- New `css/components/docs.css`, imported from `style.css`,
  uses only design tokens so both themes work. `pre` blocks get
  `overflow-x: auto` inside their container

### Search index

- `buildSite` collects `h2` and `h3` headings while rendering
- The index is JSON:
  `[{title, path, section, description, headings}]`, served at
  `/docs/search-index.json`. It is data, not a query endpoint
- `search.js` is not changed; spec 0024 connects the overlay

### Testing

- `docs_test.go`, table-driven, with `fstest.MapFS` fixtures:
  frontmatter parsing, ordering (including unordered pages),
  slug and URL rules, link rewriting with fragments, every
  strict-mode error, Mermaid marking, 404 and redirect
  behavior, traversal attempts, search index contents
- One test runs `buildSite` against the embedded `web/docs` FS
  and asserts: no error, the four sections exist, every title
  is non-empty. Strict mode already covers unique URLs and
  resolving links

### Files affected

- New: `web/docs/**`, `web/docs/embed.go`,
  `cmd/openkata-web/docs.go`, `cmd/openkata-web/docs_test.go`,
  `templates/docs.templ` (and its committed `_templ.go`),
  `web/static/js/docs.js`, `web/static/css/components/docs.css`
- Changed: `main.go` (local-only route),
  `web/static/css/style.css` (import)
- Moved: `docs/agent-workflow.md` and
  `docs/infrastructure-diagram.md` into `web/docs/architecture/`

## Decisions

- **Engine first, content in spec 0024.** The engine is
  reviewable on its own and does not wait for pages
- **Local-only route.** It follows the existing
  `/design-system/` and `/stats/` pattern; nothing half-written
  reaches production. Go-live moves the route out of the block
- **Embed now.** Production does not serve from it yet, but the
  real-tree test proves the `go:embed` patterns, and go-live is
  a route move
- **Resolve requests through an in-memory index, not the
  filesystem.** Safer than sanitizing paths, and the same index
  drives the sidebar, previous/next links and search index
- **Strict loader.** Trusted content with one author group; a
  loud error naming the file is the fastest feedback, and the
  same check is the CI gate
- **One tree, absolute links outside it.** Dropping ADRs leaves
  a single `fs.FS`, so link resolution stays inside it
- **Reuse goldmark and the existing helpers.** No new Markdown
  dependency; behavior matches skill and rule pages
- **Mermaid source only.** A class on the `<pre>` keeps the
  door open for client-side rendering without server changes
- **`<details>` sidebar with a JS enhancement.** Works without
  scripts and avoids a custom menu component. `docs.js` is new
  and loaded only on docs pages, so no existing script changes
- **Generated `_templ.go` files stay committed.** That is the
  existing repository rule

## Dependencies

- `github.com/yuin/goldmark` and `gopkg.in/yaml.v3`, already in
  `go.mod`
- Spec 0022's `docs/design/` stays out of the embedded tree
- No infrastructure or deploy script changes
