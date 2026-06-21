## Task 0018-01 — Add model-label and skill-version to openkata-eval

Status: Done
Files: cmd/openkata-eval/main.go, cmd/openkata-eval/config.go, cmd/openkata-eval/output.go, cmd/openkata-eval/output_test.go

Add `--model-label` flag to openkata-eval CLI. Read skill
version from CHANGELOG.md (latest `## [x.y.z]` entry).
Write both fields to JSON output:
- `"skill_version"` at top level
- `"model_label"` inside `config` object

Update `JSONOutput` struct and `writeJSONOutput` function.
Add test for version parsing from CHANGELOG.md.

## Task 0018-02 — Output eval results to nested directory

Status: Done
Files: cmd/openkata-eval/output.go, cmd/openkata-eval/output_test.go

Change default `--output` path resolution when running
against a skill. If no explicit `--output` is given and
the input is a skill path, write to:
`<skill>/evals/results/<model_id>/<timestamp>.json`

Where `model_id` is the `--model` flag value and
`timestamp` is `2006-01-02T150405` format. Create
directories as needed.

## Task 0018-03 — Extend generate-versions to read eval results

Status: Done
Files: cmd/generate-versions/main.go

In local mode, after scanning skills:
1. For each skill, glob `evals/results/*/` directories
2. In each model dir, find the latest `.json` file
3. Parse it: check `skill_version` matches the skill's
   current version (from git tag)
4. Extract: model ID (directory name), `model_label`,
   `overall_percentage`, scenario results
5. Write into the skill's entry as a `models` map

Skip skills with no results or version mismatch.

## Task 0018-04 — Add models and scenarios to versions.json schema

Status: Done
Files: cmd/generate-versions/main.go

Extend `artifactInfo` struct to include:
```go
type modelInfo struct {
    Label         string           `json:"label"`
    Effectiveness int              `json:"effectiveness"`
    Scenarios     []scenarioInfo   `json:"scenarios"`
}
type scenarioInfo struct {
    Name        string `json:"name"`
    Description string `json:"description"`
    Pass        bool   `json:"pass"`
}
```

Add `Models map[string]modelInfo` to `artifactInfo`
(omitempty so skills without data have no field).

## Task 0018-05 — Update SkillEntry type for effectiveness data

Status: Done
Files: cmd/openkata-web/templates/types.go

Add fields to `SkillEntry`:
- `BestScore int` (0 if no data)
- `BestModel string` (label of best-scoring model)
- `HasBenchmarks bool`

Add `BenchmarkModel` struct for per-model data:
- `Label string`
- `Effectiveness int`
- `Scenarios []BenchmarkScenario`

Add to `ArtifactDetail`:
- `Models []BenchmarkModel`
- `TesslScore int`
- `Published bool`

## Task 0018-06 — Load effectiveness data into skill listings

Status: Done
Files: cmd/openkata-web/handlers.go

Update `loadArtifactList` to parse the `models` field
from versions.json. Compute best score and best model
label. Populate `BestScore`, `BestModel`, `HasBenchmarks`
on each `SkillEntry`.

## Task 0018-07 — Display effectiveness on listing page

Status: Done
Files: cmd/openkata-web/templates/skills.templ

In the skill listing row:
- Show `96% claude-opus-4.6` badge after downloads
  (only if `HasBenchmarks` is true)
- In the collapsible `<details>` body, show per-model
  breakdown line: `model: NN% | model: NN%`

## Task 0018-08 — Add Benchmarks tab to skill detail page

Status: Done
Files: cmd/openkata-web/templates/skill_detail.templ, cmd/openkata-web/handlers.go

Add conditional "Benchmarks" tab (hidden if no data):
1. Tessl section: show score from `plugin.json` and
   registry link if published
2. Effectiveness table: model, effectiveness %, N/M
   scenarios pass
3. Scenario breakdown: list with per-model pass/fail

Load benchmark data in the detail handler from
versions.json `models` field. Read `plugin.json` score
and published status.

## Task 0018-09 — Update handlers_test.go for effectiveness

Status: Done
Files: cmd/openkata-web/handlers_test.go

Add test cases:
- Listing page renders effectiveness badge when data
  exists
- Listing page omits badge when no data
- Detail page shows Benchmarks tab with model table
- Detail page hides Benchmarks tab when no data
- Stale version results are not displayed

## Task 0018-10 — Generate templ and verify build

Status: Done
Files: cmd/openkata-web/templates/*_templ.go

Run `templ generate ./cmd/openkata-web/templates/` and
verify:
- `go build -o bin/openkata-web ./cmd/openkata-web/`
- `go test ./cmd/openkata-web/...`
- `go test ./cmd/generate-versions/...`
- `go test ./cmd/openkata-eval/...`
