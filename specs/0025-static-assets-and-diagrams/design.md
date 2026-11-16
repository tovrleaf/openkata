# Design: Static Asset Delivery and Diagram Rendering

## Architecture

```text
Browser ─► CloudFront E1X978YU1XMI3Q
             ├─ /static/*  ─► Lambda openkata-web   (UseOriginCacheControlHeaders-QueryStrings)
             ├─ /mcp*      ─► Lambda openkata-mcp   (CachingDisabled)
             └─ default    ─► Lambda openkata-web   (CachingDisabled)
```

### Edge cache

- A new `/static/*` behavior sends requests to the existing
  `lambda` origin, with compression on and `GET` and `HEAD` only
- It uses the managed cache policy
  *UseOriginCacheControlHeaders-QueryStrings*
  (`4cc15a8a-d715-48a4-82b8-cc0b614638fe`). All query strings
  are in the cache key, so `?v=` works at the edge. TTLs come
  from the response's `Cache-Control`, between 0 seconds and one
  year
- The policy also keys on cookies and `Host`. The site sets no
  cookies, and `www` redirects to the apex, so neither splits
  the cache
- A made-up query string skips the cache and reaches the Lambda,
  as every request does today
- New script: `infra/add-static-cache-behavior.sh`, idempotent
  like `add-mcp-to-cloudfront.sh`. It adds the behavior with
  `GetDistributionConfig` and `UpdateDistribution`, which the
  admin policy already grants. No IAM change is needed

A custom policy that keys on `v` alone was rejected. It would
need cache policy actions added to the admin policy, for a gain
that only matters against deliberate cache busting.

### Static handler

A new handler in `cmd/openkata-web` replaces the bare
`http.FileServerFS`:

- In Lambda it serves the embedded `static.FS`. At startup it
  hashes every file once (SHA-256, first 10 hex characters).
  Hashing 1.3 MB takes milliseconds
- It sets `ETag` to the hash and serves with `http.ServeContent`,
  which answers `If-None-Match` with `304`
- `Cache-Control`:
  - `?v=` matches the file's hash:
    `public, max-age=31536000, immutable`
  - No `v`, or an old one: `public, max-age=300`
- Outside Lambda it serves `web/static` from disk, hashes the
  file on each request and sends `Cache-Control: no-cache`, so
  an edited file shows on the next reload

### Content-hash URLs

- A new package `internal/assets` holds the hash table and
  exposes `URL(path string) string`, which returns
  `/static/<path>?v=<hash>`. A path with no file returns the URL
  without `?v=`
- Templates change from `href="/static/css/style.css"` to
  `href={ assets.URL("css/style.css") }`
- A test walks every `.templ` file and fails on a `/static/`
  link that does not go through `assets.URL`, or that names a
  missing file
- Files that other static files reference keep plain URLs and
  the 5-minute cache. These are the CSS `@import`s in
  `style.css`, `versions.json` and the icons in
  `site.webmanifest`. `search.js` fetches `versions.json`

## Diagram rendering

### Markdown

```markdown
![Plan, implement and validate a feature. A failed validation
goes back to implement.](agent-workflow/feature-development.svg)
```

The SVG file lives in a folder named after the page:
`web/docs/architecture/agent-workflow/feature-development.svg`.

### Renderer

- `docs_render.go` gets a goldmark node renderer for images. A
  relative `.svg` destination that resolves inside `web/docs` is
  replaced with the file's contents, wrapped in
  `<figure class="docs-diagram">`
- Any other image renders as goldmark does today
- The strict loader reports a missing SVG file like a broken
  link
- `web/docs` is already embedded for the docs, so the SVG files
  ship in the binary. Each diagram adds about 3 KB

### SVG conventions

- Each SVG sets a `viewBox` and no fixed `width` or `height`.
  CSS sets `max-width: 100%`
- It has a `<title>` and `role="img"` with `aria-labelledby`
- Colors come only from tokens with fallbacks, in a `<style>`
  block: `fill: var(--color-surface-raised, #fff)`
- An `@media (prefers-color-scheme: dark)` rule in the same block
  sets dark fallbacks for GitHub. On the site the tokens win
- Text uses `font-family: inherit`
- IDs inside the SVG (markers, title) carry a prefix unique to
  the diagram, because several SVGs share one page
- No Mermaid source is kept in the file

### Mermaid source view

`docs_render.go` rewrites `mermaid` code blocks to
`<pre class="mermaid-source">` (spec 0023). With the four blocks
gone nothing uses it, so the rewrite and its tests in
`docs_render_test.go` are removed.

## Decisions

### Where static files are served from

| | Embedded, cached at the edge | S3 prefix in `openkata-artifacts` |
|---|---|---|
| Deploy size | Every deploy ships all 1.3 MB | Only changed files upload |
| Lambda hits | First request per URL | None |
| Version skew | None: binary and assets match | Assets upload first; old files stay |
| Origin unavailable | Lambda down means the site is down anyway | Pages render without CSS or JS |
| Infrastructure | One behavior, one cache policy | Origin, OAC, bucket policy, cache policy, sync step |

Chosen: embedded. Without a Mermaid library no large file is
coming, and 1.3 MB is about 4% of the 35 MB binary. Edge caching
with hashed URLs removes almost all Lambda hits for a fraction of
the S3 setup.

### How a changed file reaches visitors

| | Site version `?v=` | Per-file hash `?v=` | Hashed file names |
|---|---|---|---|
| Unchanged files after a release | Downloaded again | Kept | Kept |
| Build changes | None | Hash table at startup | Rename and manifest |
| CSS `@import` and JS `fetch` | Not covered | Not covered | Need rewriting |

Chosen: per-file hash in the query string. Only changed files
get a new URL, and file names stay readable. Files that other
static files reference keep a 5-minute cache.

### Diagrams

The four options were compared against the four `mermaid`
diagrams in `agent-workflow.md`.

| | Mermaid in the browser | Inline SVG with tokens | Mermaid rendered at build | Source view |
|---|---|---|---|---|
| Both themes | Needs re-render code | Yes, through tokens | Fixed colors; needs rewriting | Yes |
| Without JavaScript | No | Yes | Yes | Yes |
| Bytes | About 3 MB of JS | About 3 KB per diagram | 10–30 KB per diagram | None |
| Editing | Edit text | Move coordinates, or ask an agent | Edit text, run Node and Chromium | Edit text |
| Reads as a diagram | Yes | Yes | Yes | No |

Chosen: inline SVG with tokens, kept in separate files. It is the
only option that follows the theme and works without JavaScript
with no new tooling.

SVG written straight into the Markdown was rejected: the raw page
fills with coordinates, and GitHub strips inline `<svg>`.

More diagrams are expected, and hand-drawn SVG costs the most to
edit. Mermaid v13, in alpha since November 2026, renders SVG in
Node without a browser, which removes the Chromium cost of
build-time rendering. Spec 0030 adopts it once it is stable. Its
output fits the same pipeline: an SVG file that the renderer
inlines.

## Dependencies

- Spec 0023: docs renderer and strict loader
- AWS: one cache behavior, added by the account owner with the
  admin credentials
