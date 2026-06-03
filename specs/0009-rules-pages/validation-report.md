# Validation: Rules Pages

Spec: specs/0009-rules-pages/spec.md
Date: 2026-06-03
Validator: Kiro (spec-validator)

## Requirements Check

| # | Requirement | Status | Evidence |
|---|-------------|--------|----------|
| 1.1 | 則 kanji at two-line height | PASS | rules.templ: `<span class="hero-kanji">則</span>` in hero section |
| 1.2 | "Rules" heading | PASS | rules.templ: `<h1 class="cursor">Rules</h1>` |
| 1.3 | Correct subtitle | PASS | `"Always-on constraints applied to every agent session."` |
| 1.4 | Vertical stack with same component as skills | PASS | Uses `artifact-list`/`artifact-entry` classes (shared with skills) |
| 1.5 | Sequential number | PASS | `fmt.Sprintf("%02d", i+1)` in `artifact-number` span |
| 1.6 | Name | PASS | `artifact-name` span with `r.Name` |
| 1.7 | Download count with ↓ | PASS | `artifact-downloads` span with `↓` suffix |
| 1.8 | Version | PASS | `artifact-version` span with `v` prefix |
| 1.9 | +/- toggle icon | PASS | `artifact-toggle-icon` span (CSS-driven via details/summary) |
| 1.10 | Tags with color-coding | PASS | `SplitTags`/`TagClass` applied per tag |
| 1.11 | "Open" button linking to /rules/:name/ | PASS | `btn-open` class, `href="/rules/{name}/"` |
| 1.12 | Expand/collapse description | PASS | Uses `<details>` element for toggle |
| 1.13 | Bold first sentence styling | PASS | `<strong>{ FirstSentence(...) }</strong>` with `<br/>` separator |
| 1.14 | Empty state message | PASS | `"No rules available yet."` when `len(rules) == 0` |
| 1.15 | Alphabetical sort | PASS | `loadArtifactList` sorts by name alphabetically |
| 2.1 | rule_detail.templ exists | PASS | File present with `RuleDetailPage` component |
| 2.2 | /rules/{name} route | PASS | handlers.go `len(parts) == 1` branch |
| 2.3 | /rules/{name}/{version} route | PASS | handlers.go `len(parts) == 2` branch |
| 2.4 | /rules/{name}/raw/{filepath} route | PASS | handlers.go `len(parts) >= 3 && parts[1] == "raw"` branch |
| 2.5 | /rules/{name}/{version}/raw/{filepath} route | PASS | handlers.go `len(parts) >= 4 && parts[2] == "raw"` branch |
| 2.6 | Version dropdown | PASS | `<select class="version-select">` when `len(rule.Versions) > 1` |
| 2.7 | Download count | PASS | `detail-downloads` span in rule_detail.templ |
| 2.8 | .tar.gz link | PASS | `detail-download` link to `/rules/{name}/archive` |
| 2.9 | Tags | PASS | `artifact-tags` div with `TagClass` color-coding |
| 2.10 | Tabs (Overview, Files, Changelog, Acknowledgments) | PASS | `TabBar` with conditional Acknowledgments tab |
| 2.11 | data-artifact-path="rules" on main | PASS | `<main data-artifact-path="rules">` |
| 3.1 | Description + RULE.md rendered as markdown | PASS | Description paragraph + `templ.Raw(rule.Docs)` |
| 3.2 | Strip frontmatter, skip first heading | PASS | `renderMarkdown` calls `stripFrontmatter` and `stripFirstH1` |
| 3.3 | External links in new window | PASS | `addTargetBlankToExternalLinks` adds `target="_blank"` |
| 3.4 | Relative links navigate to Files tab | PASS | detail.js `highlightFile` intercepts relative links |
| 4.1 | File tree with connectors | PASS | `renderNode` generates `├──`/`└──`/`│` connectors |
| 4.2 | Directories collapsed by default | PASS | `<details>` without `open` attribute |
| 4.3 | Excludes CHANGELOG.md and references/ACKNOWLEDGMENTS.md | PASS | `isExcludedFile` returns true for both |
| 4.4 | Click file shows content | PASS | `file-tree-link` href targets `file-block` details |
| 4.5 | Preview/Code/Raw toggle | PASS | `file-mode-btn` buttons in file_viewer.templ |
| 5.1 | Render CHANGELOG.md | PASS | `templ.Raw(rule.Changelog)` in changelog panel |
| 5.2 | Filter by version | PASS | `filterChangelogByVersion` keeps viewed version and prior |
| 6.1 | Acknowledgments tab only if file exists | PASS | Conditional `if rule.Acknowledgments != ""` for tab and panel |
| 7.1 | /rules/{name}/archive route | PASS | handlers.go archive branch calls `handleArchive` |
| 7.2 | /rules/{name}/archive/{version} route | PASS | `len(parts) >= 3` passes version to `handleArchive` |
| 8.1 | "Rules" link in nav on all pages | PASS | nav.templ: `<a href="/rules/" class="nav-link">Rules</a>` |
| 8.2 | Shared Nav() component | PASS | nav.templ defines Nav(), used via `@Nav()` in 7 templates |
| 8.3 | No active-state highlighting | PASS | No active state CSS for nav-link |
| 9.1 | All .skill-* renamed to .artifact-* | PASS | Zero matches for `.skill-` in CSS and templ files |
| 9.2 | Both skills and rules use .artifact-* | PASS | 23 `.artifact-*` rules in style.css, both templates use them |
| 10.1 | Shared nav.templ | PASS | File exists with `Nav()` component |
| 10.2 | Shared tabs.templ | PASS | File exists with `TabBar()` component |
| 10.3 | Shared file_viewer.templ | PASS | File exists with `FileViewer()` component |
| 10.4 | skill_detail.templ uses shared components | PASS | Uses `@Nav()`, `@TabBar(...)`, `@FileViewer(...)` |
| 10.5 | No inline JavaScript in skill_detail.templ | PASS | No `<script>` tags in templ files (JS in layout only) |
| 10.6 | detail.js uses data-path attribute | PASS | `main.dataset.artifactPath` drives version select URL |
| 10.7 | buildFileTree accepts artifactType param | FAIL | Signature is `buildFileTree(files, skillName, version)` — no artifactType. FileViewer handles type via Raw link URL. |
| 11.1 | RuleDetail struct exists (not reusing SkillDetail) | PASS | types.go: separate `RuleDetail` struct |
| 11.2 | Shared loadArtifactList helper | PASS | `loadArtifactList(ctx, artifactType)` used by skills, rules, profiles |
| 11.3 | Shared detail loader parameterized by artifact type | PASS | `loadArtifactDetailLocal(artifactType, ...)` and `loadArtifactDetailS3(...)` |
| 11.4 | gitVersions accepts artifact type | PASS | `gitVersions(artifactType, name string)` |
| 11.5 | Interface for shared detail loading | FAIL | No Go interface — uses conversion from `*SkillDetail` to `*RuleDetail` |
| 12.1 | Table-driven handler tests | PASS | `TestHandleRulesRouting` with named subtests |
| 12.2 | Listing 200 | PASS | `{name: "listing", path: "/rules/", wantCode: 200}` |
| 12.3 | Detail 200 | PASS | `{name: "detail latest", path: "/rules/test-rule", wantCode: 200}` |
| 12.4 | Version 200 | PASS | `{name: "detail specific version", path: "/rules/test-rule/1.0.0", wantCode: 200}` |
| 12.5 | Raw 200 | PASS | `{name: "raw file", ...wantCode: 200, wantBody: "#!/bin/bash"}` |
| 12.6 | Unknown 404 | PASS | `{name: "unknown rule", path: "/rules/no-such-rule", wantCode: 404}` |
| 12.7 | Test shared helpers | PASS | `TestLoadRuleDetailLocalFileContents`, `TestGitVersions` |
| 13.1 | Build verification passes | PASS | `templ generate && go build && go test -race` all succeed |
| 14.1 | Tag rules before deployment | SKIP | Task 8 intentionally not done |
| 14.2 | Regenerate versions.json | SKIP | Task 8 intentionally not done |

## Out-of-Scope Check

- Tag filtering or search: Not present ✓
- Sorting controls: Not present ✓
- Pagination: Not present ✓
- No unexpected features observed beyond spec requirements

## Issues Found

1. **buildFileTree signature** (10.7): Spec requires
   `buildFileTree` to accept `artifactType` string
   parameter. Current signature is
   `buildFileTree(files []string, skillName, version string)`.
   The `FileViewer` component handles artifact type for
   URL generation, so there's no functional bug — the
   tree HTML itself is type-agnostic. Low severity.

2. **No Go interface for detail loading** (11.5): Spec
   says "Interface for shared detail loading: shared
   internal helper populates via interface that both
   SkillDetail and RuleDetail implement." Implementation
   uses `loadArtifactDetailLocal` returning
   `*SkillDetail`, then converts to `*RuleDetail` via
   field copy. Achieves the goal (shared code, separate
   types) but doesn't use a Go `interface`. Low severity
   — the shared parameterized helper is present.

## Verdict

PARTIAL PASS

Two minor architectural deviations from spec (no
interface, buildFileTree signature) but all functional
requirements are met, tests pass, and the build is
clean. No blockers for deployment.
