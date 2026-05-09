# Tasks: Website Foundation

## Tasks

### 1. Add templ dependency
- **Status**: Done
- **Goal**: Add the templ module to go.mod
- **Boundary**: `go.mod`, `go.sum`
- **Depends**: None
- **Done when**: `go mod tidy` succeeds with templ imported

### 2. Create the landing page template
- **Status**: Done
- **Goal**: Write the base layout and home page in templ
- **Boundary**: `cmd/openkata-web/templates/`
- **Depends**: 1
- **Done when**: `templ generate` produces `.go` files
  without errors

### 3. Create the web server
- **Status**: Done
- **Goal**: HTTP server that renders the home template and
  serves static files
- **Boundary**: `cmd/openkata-web/main.go`,
  `cmd/openkata-web/handlers.go`
- **Depends**: 2
- **Done when**: `go run ./cmd/openkata-web` starts and
  `curl localhost:8080` returns HTML with the project name

### 4. Add static assets
- **Status**: Done
- **Goal**: Create CSS placeholder and vendor htmx.min.js
- **Boundary**: `web/static/`
- **Depends**: None
- **Done when**: CSS file exists, htmx.min.js exists, both
  are referenced in the layout template

### 5. Configure air for hot reload
- **Status**: Done
- **Goal**: Air config that watches .go, .templ, and .css
  files, runs templ generate before rebuild
- **Boundary**: `web/.air.toml`
- **Depends**: 3
- **Done when**: Running `air -c web/.air.toml` starts the
  server and rebuilds on file changes

### 6. Add make dev target
- **Status**: Done
- **Goal**: Single command to start local development
- **Boundary**: `Makefile`
- **Depends**: 5
- **Done when**: `make dev` starts air, server runs on
  localhost:8080, changing a .templ file triggers rebuild

## Progress Log

- [2026-05-09] Tasks 1-4: Created web server with templ
  templates, handlers, static assets. Server verified with
  curl returning full HTML.
- [2026-05-09] Tasks 5-6: Added air config and make dev
  target. templ CLI installed at ~/go/bin/templ.
