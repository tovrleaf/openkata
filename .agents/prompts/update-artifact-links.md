# Update Artifact Links

Scan all skills for well-known file paths they create or
maintain, then update `fileArtifactMap` in
`cmd/openkata-web/handlers.go` to keep cross-links current.

## Workflow

1. Read each `skills/*/SKILL.md` and look for file paths
   the skill creates or maintains (e.g., `docs/context/`,
   `docs/adr/`, config files)
2. Compare findings against the current `fileArtifactMap`
   in `cmd/openkata-web/handlers.go`
3. Add any missing entries, remove stale ones
4. Rebuild: `go build -o bin/openkata-web ./cmd/openkata-web/`
