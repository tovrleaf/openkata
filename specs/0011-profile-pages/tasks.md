# Tasks: Profile Pages

## Tasks

### 1. Restructure profiles into directories
- **Status**: Pending
- **Goal**: Move flat profile files into directory
  structure: `profiles/{name}/{name}.md`. Add
  CHANGELOG.md to each. Update symlinks in
  `.agents/profiles/` to point to new paths. Update
  `generate-versions` to scan the new structure.
- **Boundary**: `profiles/` directory,
  `.agents/profiles/` symlinks,
  `cmd/generate-versions/main.go`
- **Depends**: None
- **Done when**: Each profile is at
  `profiles/{name}/{name}.md` with a CHANGELOG.md,
  symlinks updated, `make versions` works
- **Verify**: `make versions` produces correct output,
  symlinks resolve correctly

### 2. Refactor to single ArtifactDetail struct
- **Status**: Pending
- **Goal**: Replace SkillDetail, RuleDetail with a single
  ArtifactDetail struct. Update all templates and handlers.
- **Boundary**: `cmd/openkata-web/templates/types.go`,
  `cmd/openkata-web/templates/skill_detail.templ`,
  `cmd/openkata-web/templates/rule_detail.templ`,
  `cmd/openkata-web/templates/file_viewer.templ`,
  `cmd/openkata-web/handlers.go`,
  `cmd/openkata-web/handlers_test.go`
- **Depends**: None
- **Done when**: Single `ArtifactDetail` struct used
  everywhere, no more SkillDetail/RuleDetail types,
  all existing tests pass, skills and rules pages
  render identically
- **Verify**: `go test -race ./cmd/openkata-web/...`

### 3. Add Profiles link to Nav
- **Status**: Pending
- **Goal**: Add "Profiles" link to the shared Nav component
- **Boundary**: `cmd/openkata-web/templates/nav.templ`
- **Depends**: 1
- **Done when**: "Profiles" appears to the right of
  "Rules" in nav on all pages
- **Verify**: `templ generate && go build -o bin/openkata-web ./cmd/openkata-web/`

### 4. Update profiles listing page
- **Status**: Pending
- **Goal**: Rewrite listing to match skills/rules style
  (師 kanji, numbered expandable entries, tags, Open)
- **Boundary**: `cmd/openkata-web/templates/profiles.templ`
- **Depends**: 3
- **Done when**: `/profiles/` shows all profiles with
  same style as skills/rules listing
- **Verify**: `templ generate && go build -o bin/openkata-web ./cmd/openkata-web/`

### 5. Add profile detail loader
- **Status**: Pending
- **Goal**: Implement `loadProfileDetailVersion` using
  shared helpers for single-file structure
- **Boundary**: `cmd/openkata-web/templates/types.go`,
  `cmd/openkata-web/handlers.go`
- **Depends**: 1, 2
- **Done when**: Loader returns profile data with
  rendered markdown, version list, downloads
- **Verify**: `go build -o bin/openkata-web ./cmd/openkata-web/ && go test ./cmd/openkata-web/...`

### 6. Add profile detail page handler and template
- **Status**: Pending
- **Goal**: Build `/profiles/:name` detail page with
  Overview and Changelog tabs
- **Boundary**: `cmd/openkata-web/templates/profile_detail.templ` (new),
  `cmd/openkata-web/handlers.go` (extend handleProfiles)
- **Depends**: 4, 5
- **Done when**: `/profiles/spec-planner` renders with
  profile content, version dropdown, download link,
  raw file serving at `/profiles/:name/:version/raw`
- **Verify**: `templ generate && go build -o bin/openkata-web ./cmd/openkata-web/ && go test ./cmd/openkata-web/...`

### 7. Verify archive download works
- **Status**: Pending
- **Goal**: Confirm `/profiles/{name}/archive` serves
  a valid .tar.gz locally
- **Boundary**: verification only
- **Depends**: 6
- **Done when**: Downloading archive returns valid
  tar.gz containing the profile markdown
- **Verify**: `curl -o /tmp/test.tar.gz localhost:8080/profiles/spec-planner/archive && tar -tzf /tmp/test.tar.gz`

### 8. Add handler tests for profiles
- **Status**: Pending
- **Goal**: Table-driven tests for profile routes
- **Boundary**: `cmd/openkata-web/handlers_test.go`
- **Depends**: 6
- **Done when**: Tests cover listing 200, detail 200,
  version 200, raw 200, unknown 404, archive 200
- **Verify**: `go test -race ./cmd/openkata-web/...`

### 9. Verify deployment pipeline
- **Status**: Pending
- **Goal**: Confirm publish workflow handles profile
  directories correctly (git archive extracts files,
  S3 upload works)
- **Boundary**: `.github/workflows/publish.yaml`,
  `cmd/generate-versions/main.go`
- **Depends**: 1
- **Done when**: Manual workflow_dispatch for a
  profile tag succeeds and profile appears in S3
- **Verify**: Check GitHub Actions + S3 contents

## Progress Log
