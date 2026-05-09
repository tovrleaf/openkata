# Design: Website Foundation

## Architecture

```text
cmd/openkata-web/
├── main.go              # HTTP server, route registration
├── handlers.go          # Route handlers (home page)
└── templates/
    ├── layout.templ     # Base HTML layout (head, body, scripts)
    └── home.templ       # Landing page content

web/
├── static/
│   ├── css/
│   │   └── style.css   # Placeholder stylesheet
│   └── js/
│       └── htmx.min.js # htmx library (vendored)
└── .air.toml            # Air configuration for hot reload
```

- `cmd/openkata-web/` — the binary, following existing
  `cmd/openkata-mcp/` convention
- `web/` — static assets served by the Go server. Separate
  from `cmd/` so assets aren't mixed with Go source.
- Templates live in `cmd/openkata-web/templates/` because
  templ generates Go code alongside the templates.

## Decisions

- **Vendor htmx** — no CDN dependency, works offline, one
  less external request. Copy `htmx.min.js` into static.
- **Filesystem serving in dev** — `http.Dir("web/static")`
  so CSS/JS changes are instant without rebuild. Production
  will use `embed.FS` later (out of scope).
- **Air config in `web/`** — keeps dev tooling config near
  the assets it watches, not in project root.
- **Templates in cmd/** — templ generates `.go` files next
  to `.templ` files. Keeping them in the binary's package
  avoids import complexity.

## Dependencies

- `github.com/a-h/templ` — HTML templating (new dependency)
- `air` — dev tool, not a Go dependency (installed via
  `go install`)
- htmx — vendored JS file, no Go dependency
