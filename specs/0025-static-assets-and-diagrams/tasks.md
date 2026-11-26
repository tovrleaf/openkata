# Tasks: Static Asset Delivery and Diagram Rendering

Part 1 tasks 1–3 and Part 2 tasks 4–7 are independent of each
other. Task 3 changes AWS, and the account owner runs it.

## Tasks

### 1. Static handler with hashes and cache headers
- **Status**: Done
- **Goal**: Static responses carry `ETag` and `Cache-Control`,
  and code can ask for a file's hashed URL
- **Boundary**: `internal/assets/` (new),
  `cmd/openkata-web/main.go`, tests
- **Depends**: None
- **Done when**:
  - `assets.URL` returns `/static/<path>?v=<hash>`
  - The handler sets `ETag` and answers `If-None-Match` with
    `304`
  - A matching `v` gets the one-year `immutable` header, and any
    other request gets `max-age=300`
  - Outside Lambda it serves from disk with `no-cache` and fresh
    hashes
- **Verify**: `go test ./internal/assets/ ./cmd/openkata-web/`
  and `curl -I` against the local server

### 2. Templates use hashed URLs
- **Status**: Done
- **Goal**: Every `/static/` link in a template carries its hash
- **Boundary**: `cmd/openkata-web/templates/*.templ`, generated
  `_templ.go` files, a new template test
- **Depends**: 1
- **Done when**:
  - No template writes a literal `/static/` URL
  - A test fails on a literal one, or on a path with no file
- **Verify**: `templ generate ./cmd/openkata-web/templates/`,
  `go test ./cmd/openkata-web/...`, and the local site renders
  with `?v=` on CSS, JS and images

### 3. CloudFront: cache `/static/*` at the edge
- **Status**: Done
- **Goal**: The Lambda only sees the first request for each
  static URL
- **Boundary**: `infra/add-static-cache-behavior.sh` (new),
  `infra/iam-admin-policy.json`, `infra/README.md`
- **Depends**: 1–2 deployed
- **Done when**:
  - The admin policy grants the three cache policy actions
  - The script finds or creates the `openkata-static` cache
    policy and adds the `/static/*` behavior on the `lambda`
    origin
  - A second run changes nothing
  - The README lists the script
- **Verify**: after the account owner runs it, two requests for
  the same hashed URL return `x-cache: Miss` then
  `x-cache: Hit from cloudfront`

### 4. Renderer inlines SVG diagrams
- **Status**: Done
- **Goal**: A Markdown image of a local `.svg` file shows as an
  inline SVG that follows the theme
- **Boundary**: `cmd/openkata-web/docs_render.go`, the docs
  loader, `web/static/css/components/docs.css`, tests
- **Depends**: None
- **Done when**:
  - A relative `.svg` image inside `web/docs` is replaced with
    the file's contents in a `<figure class="docs-diagram">`
  - Other images are unchanged
  - The strict loader rejects a link to a missing SVG
- **Verify**: `go test ./cmd/openkata-web/ -run 'Docs|Embedded'`

### 5. Redraw the agent workflow diagrams
- **Status**: Done
- **Goal**: The four diagrams show as pictures in both themes
- **Boundary**: `web/docs/architecture/agent-workflow.md`,
  `web/docs/architecture/agent-workflow/*.svg` (new)
- **Depends**: 4
- **Done when**:
  - Each `mermaid` block is replaced with an image link and alt
    text
  - Each SVG follows the conventions in `design.md`
  - Each SVG shows the same nodes and edges as its Mermaid block
- **Verify**: `go test ./cmd/openkata-web/ -run Embedded`, and a
  look at the page on the local server in both themes and on
  narrow screens

### 6. Remove the Mermaid source view
- **Status**: Done
- **Goal**: No Mermaid handling remains
- **Boundary**: `cmd/openkata-web/docs_render.go`,
  `cmd/openkata-web/docs_render_test.go`
- **Depends**: 5
- **Done when**: the `mermaid-source` rewrite and its tests are
  gone, and `mermaid` appears nowhere in `cmd/` or `web/`
- **Verify**: `go test ./cmd/openkata-web/...` and
  `grep -rn mermaid cmd web` prints nothing

### 7. Document how to add a diagram
- **Status**: Pending
- **Goal**: A contributor can add a diagram that follows the
  conventions
- **Boundary**: `web/docs/contributing/authoring.md`
- **Depends**: 4
- **Done when**: the page shows the image link, where the SVG
  file goes, and the SVG conventions from `design.md`
- **Verify**: `go test ./cmd/openkata-web/ -run Embedded`

## Progress Log

### 2026-11-16: Task 1

Added `internal/assets`. The handler sets the hash as `ETag` and
picks `Cache-Control` from the `v` query. `main.go` serves
`/static/` through `assets.Default` in both modes. On the local
server a file gets `no-cache`, a repeat with its ETag gets `304`,
and a directory gets `404`. Against the real `static.FS`, a
matching `v` gets one year and no `v` gets 5 minutes. Directory listings under `/static/` now return
`404`. `site.webmanifest` is still served as `text/plain`, as
before.

### 2026-11-16: Task 2

Every `/static/` link in the templates goes through `assets.URL`,
including the local-only chart script on the stats page. A new
test fails on a literal `/static/` URL in a template and on a link
to a missing file. A temporary bad link in `nav.templ` made it
fail as expected. On the local server, the home, docs and getting
started pages render every static link with `?v=`.

### 2026-11-26: Task 3

Added `infra/add-static-cache-behavior.sh`. The managed
*UseOriginCacheControlHeaders-QueryStrings* policy was dropped,
because it keys on `Host` and the Function URL rejects a foreign
`Host`. The script now creates a custom `openkata-static` policy
that keys on `v` alone, and the admin policy grants the actions
it needs. The behavior copies the default behavior's CloudFront
functions, so `www` still redirects. Against a mocked `aws`, the
first run created the policy and added the behavior after
`/mcp*`, and a second run left the behavior alone. The account
owner applied the admin policy and ran the script. On the live
site, a second request for `/static/favicon.svg` returned
`x-cache: Hit from cloudfront`, and `www` still redirects. Tasks
1–2 are not deployed yet, so these files are cached for the
policy's 300-second default. Check a hashed URL for the one-year
hit after this branch deploys.

### 2026-11-26: Task 4

A relative `.svg` image alone in its paragraph now renders as the
file's `<svg>` element inside `<figure class="docs-diagram">`,
with any XML declaration dropped. External, site-absolute and
non-SVG images render as before. A missing file, a path outside
`web/docs`, a file without an `<svg>` element, or a diagram
inside a sentence fails the build with `errDocsLink`. The loader
now allows a directory named after a page in the same section,
which is where the SVG files go. With inlining turned off, all
seven diagram tests failed.

### 2026-11-26: Task 5

The four `mermaid` blocks in `agent-workflow.md` are now image
links to SVG files in `architecture/agent-workflow/`, each with
alt text. Each diagram keeps the nodes and edges of its Mermaid
block. The two wide ones are laid out top to bottom, so they fit
a phone. The first draft filled the column and its text came out
at about twice the body size. Each SVG now sets a `width` of 8/7
of its viewBox width, and `design.md` says so. Screenshots of the
local page in both themes, at desktop width and in a 375px frame,
showed every diagram whole and readable. Opened on its own with
the OS in dark mode, as GitHub shows it, each SVG used its dark
fallbacks.

### 2026-11-26: Task 6

Removed the rewrite that marked `mermaid` code blocks as
`mermaid-source`, and its three tests. A `mermaid` block now
renders like any other code block. `grep -rn mermaid cmd web`
prints nothing.
