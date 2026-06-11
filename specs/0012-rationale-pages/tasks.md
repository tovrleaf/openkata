# Tasks: RATIONALE.md Support

## Tasks

### 1. Add RATIONALE.md to isExcludedFile and detail loaders
- **Status**: Done
- **Goal**: Exclude RATIONALE.md from the Files tab but
  NOT from archives, load it for detail page rendering
  (same inclusion/exclusion pattern as CHANGELOG.md)
- **Boundary**: `cmd/openkata-web/handlers.go`
- **Depends**: None
- **Done when**: RATIONALE.md excluded from Files tab,
  included in archive downloads, loaded and rendered as
  HTML in detail struct, build passes
- **Verify**: `go build -o bin/openkata-web ./cmd/openkata-web/`

### 2. Add Rationale tab to detail page templates
- **Status**: Done
- **Goal**: Show "Rationale" tab on skill, rule, and
  profile detail pages when RATIONALE.md exists
- **Boundary**: `cmd/openkata-web/templates/skill_detail.templ`,
  `cmd/openkata-web/templates/rule_detail.templ`,
  `cmd/openkata-web/templates/profile_detail.templ`
- **Depends**: 1
- **Done when**: Tab appears only when rationale content
  exists, positioned after Files and before Changelog,
  renders markdown
- **Verify**: `templ generate && go build -o bin/openkata-web ./cmd/openkata-web/`

### 3. Write RATIONALE.md for spec-workflow
- **Status**: Done
- **Goal**: Create the first RATIONALE.md explaining
  spec-workflow's design decisions — this is the first
  of many; all skills and rules should eventually have
  one
- **Boundary**: `skills/spec-workflow/RATIONALE.md`
- **Depends**: None
- **Done when**: File documents key decisions (reference
  split, mode detection placement, phase gates, progress
  log, token economics)
- **Verify**: File exists and renders locally

### 4. Update skill conventions
- **Status**: Done
- **Goal**: Document RATIONALE.md as an optional
  convention in openkata-skill-conventions
- **Boundary**: `.agents/skills/openkata-skill-conventions/SKILL.md`
- **Depends**: None
- **Done when**: Convention documented with guidance on
  when to include a rationale file
- **Verify**: Read the updated file

### 5. Add tests
- **Status**: Done
- **Goal**: Test that RATIONALE.md is excluded from
  Files tab, included in archive, and loaded for detail
  rendering
- **Boundary**: `cmd/openkata-web/handlers_test.go`
- **Depends**: 1, 2
- **Done when**: Tests verify exclusion from Files tab,
  inclusion in archive, and presence in detail struct
- **Verify**: `go test -race ./cmd/openkata-web/...`

## Progress Log

- 2026-06-07: All tasks implemented
