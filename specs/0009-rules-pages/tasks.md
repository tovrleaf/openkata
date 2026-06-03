# Tasks: Rules Pages

Prerequisite: spec 0010 (CSS rename) must be complete.

## Tasks

### 1. Extract shared Nav component
- **Status**: Done
- **Goal**: Replace duplicated nav markup with a shared
  `Nav()` templ component
- **Boundary**: `cmd/openkata-web/templates/nav.templ`
  (new), all `.templ` files with inline nav (6 files)
- **Depends**: 0010 complete
- **Done when**: All pages render nav from `Nav()`,
  "Rules" link appears to the right of "Skills",
  `templ generate` + build passes
- **Verify**: `go build -o bin/openkata-web ./cmd/openkata-web/`

### 2. Extract detail page JS to external file
- **Status**: Done
- **Goal**: Move inline JS from `skill_detail.templ` to
  `web/static/js/detail.js`, parameterized via `data-`
  attributes (no hardcoded `/skills/` paths)
- **Boundary**: `web/static/js/detail.js` (new),
  `cmd/openkata-web/templates/skill_detail.templ`,
  `cmd/openkata-web/templates/layout.templ` (script tag)
- **Depends**: 1
- **Done when**: Skill detail page functions identically
  (tabs, version select, file tree, preview/code toggle),
  version select uses `data-path` attribute, no inline
  `<script>` in skill_detail.templ
- **Verify**: Manual test of skill detail page interactions

### 3. Extract shared template components
- **Status**: Done
- **Goal**: Create reusable templ components for detail
  pages: tab bar, file tree, file blocks
- **Boundary**: `cmd/openkata-web/templates/tabs.templ`
  (new), `cmd/openkata-web/templates/file_viewer.templ`
  (new), `cmd/openkata-web/templates/skill_detail.templ`
  (refactor to compose from shared)
- **Depends**: 2
- **Done when**: `skill_detail.templ` uses shared
  components, skill detail page renders identically,
  `templ generate` + build passes
- **Verify**: `go build -o bin/openkata-web ./cmd/openkata-web/`

### 4. Extract shared handler helpers
- **Status**: Done
- **Goal**: Refactor handlers into shared helpers:
  `loadArtifactList`, shared detail loader (interface),
  parameterized `gitVersions`, `buildFileTree` with
  `artifactType` param
- **Boundary**: `cmd/openkata-web/handlers.go`,
  `cmd/openkata-web/templates/types.go` (add `RuleDetail`,
  interface)
- **Depends**: 3
- **Done when**: `loadSkillsList` and `loadSkillDetail*`
  delegate to shared helpers, existing skill routes work
  unchanged, `RuleDetail` type exists, tests pass
- **Verify**: `go test ./cmd/openkata-web/...`

### 5. Add rules listing handler and template
- **Status**: Done
- **Goal**: Rewrite `/rules/` listing to match skills
  listing style (numbered, expandable, tags, "Open"
  button, 則 kanji hero)
- **Boundary**: `cmd/openkata-web/templates/rules.templ`,
  `cmd/openkata-web/handlers.go` (update `handleRules`
  to use shared helper)
- **Depends**: 4
- **Done when**: `/rules/` shows all rules with number,
  name, downloads, version, tags, expand/collapse,
  "Open" links to `/rules/:name/`, sorted alphabetically
- **Verify**: `go build -o bin/openkata-web ./cmd/openkata-web/ && go test ./cmd/openkata-web/...`

### 6. Add rule detail handler and template
- **Status**: Done
- **Goal**: Build `/rules/:name` and `/rules/:name/:version`
  detail pages with Overview, Files, Changelog,
  Acknowledgments tabs
- **Boundary**: `cmd/openkata-web/templates/rule_detail.templ`
  (new), `cmd/openkata-web/handlers.go` (extend
  `handleRules` routing, add detail loading via shared
  helper)
- **Depends**: 5
- **Done when**: `/rules/bash-style` renders detail page
  with RULE.md content, file tree, changelog, version
  dropdown, raw file serving works at
  `/rules/:name/:version/raw/:filepath`
- **Verify**: `go build -o bin/openkata-web ./cmd/openkata-web/ && go test ./cmd/openkata-web/...`

### 7. Add handler tests for rules
- **Status**: Pending
- **Goal**: Table-driven tests for rules listing and
  detail routing (200, 404 cases, raw file serving)
- **Boundary**: `cmd/openkata-web/handlers_test.go`
- **Depends**: 6
- **Done when**: Tests cover: rules listing 200, rule
  detail 200, specific version 200, raw file 200,
  unknown rule 404, shared helper tests pass
- **Verify**: `go test -race ./cmd/openkata-web/...`

### 8. Tag rules and verify end-to-end
- **Status**: Pending
- **Goal**: Tag all rules with v1.0.0, regenerate
  versions.json, verify pages work locally
- **Boundary**: git tags, `web/static/versions.json`
- **Depends**: 7
- **Done when**: All four rules tagged, `make versions`
  updates versions.json with correct versions, listing
  and detail pages display real data locally
- **Verify**: `make versions` + visual check

## Progress Log

- 2026-06-03: Task 6 done — created rule_detail.templ
  mirroring skill detail (則 kanji, tabs, file viewer,
  changelog, acknowledgments). Extended handleRules with
  detail routing (/rules/:name, /rules/:name/:version,
  raw file serving). Added loadRuleDetailVersion using
  shared helpers. templ generate + build + tests pass.
- 2026-06-03: Task 5 done — rewrote rules.templ to match
  skills listing style: kanji hero (則), numbered expandable
  entries, tags, version, downloads, "Open" button linking
  to /rules/:name/. Handler already uses shared helper.
  templ generate + build + tests pass.
- 2026-06-03: Task 4 done — extracted shared helpers:
  loadArtifactList, loadDownloadCounts, loadArtifactDetailLocal,
  loadArtifactDetailS3. Parameterized gitVersions with
  artifactType. Added RuleDetail type. All three list loaders
  delegate to loadArtifactList. Build + tests pass.
- 2026-06-03: Task 3 done — extracted TabBar and FileViewer
  shared components into tabs.templ and file_viewer.templ,
  added TabDef type to types.go, refactored skill_detail.templ
  to compose from shared components, templ generate + build
  + tests pass.
- 2026-06-03: Task 2 done — extracted inline JS from
  skill_detail.templ to web/static/js/detail.js,
  parameterized via data-artifact-path attribute,
  added script tag to layout.templ, build + tests pass.
- 2026-06-03: Task 1 done — created `nav.templ` with
  shared `Nav()` component (Skills + Rules links),
  replaced inline nav in all 6 templates, templ generate
  + build + tests pass.
