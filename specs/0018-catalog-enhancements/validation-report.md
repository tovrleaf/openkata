# Validation Report — Spec 0018

Date: 2026-06-21
Validator: spec-validator

## Acceptance Criteria Results

| # | Criterion | Status | Notes |
|---|-----------|--------|-------|
| 1 | `openkata-eval` stores model label and skill version in result JSON | PASS | `JSONOutput.SkillVersion` set from `parseSkillVersion()`; `JSONOutputConfig.ModelLabel` set from `cfg.ModelLabel` (populated via `--model-label` flag). Both serialise to JSON. |
| 2 | Results stored at `skills/<name>/evals/results/<model>/<ts>.json` | PASS | `autoResolvePath()` constructs `<skillPath>/evals/results/<model>/<ts>.json`; called in `main.go` when no `--output` flag and mode is `modeSkill`, with `direct=true` so the path is used as-is. |
| 3 | `generate-versions --local` produces `models` field with effectiveness per model | PASS | `runLocal()` calls `loadEvalResults()` for each skill, which reads per-model result files, computes `int(overallPercentage + 0.5)`, and populates `artifactInfo.Models`. Field is omitted (`omitempty`) when no results exist. |
| 4 | `/skills/` listing shows best effectiveness score with model name | PASS | `loadArtifactList()` computes `BestScore`/`BestModel` by iterating models. `skills.templ` renders `<span class="artifact-effectiveness">{ fmt.Sprintf("%d%%", s.BestScore) } { s.BestModel }</span>` inside `if s.HasBenchmarks`. Test `TestSkillsListingEffectiveness/shows_badge_when_data_exists` verifies `96%` and `Opus 4.6` appear. |
| 5 | Listing rows expand to show per-model breakdown | PASS | `skills.templ` renders `<p class="artifact-models">` in the `<details>` body, iterating `s.Models` with label and percentage. Uses `<details>`/`<summary>` for expand behaviour. |
| 6 | Detail page shows model table | PASS | `skill_detail.templ` renders a `<table class="benchmark-table">` with Model, Effectiveness, and Scenarios columns inside `#panel-benchmarks`, conditioned on `len(skill.Models) > 0`. `SkillTabs()` only adds the tab when models or tessl score exist. Test `TestSkillDetailBenchmarksTab/shows_tab_when_data_exists` verifies the tab and model labels appear. |
| 7 | Skills without eval results show no effectiveness info | PASS | `HasBenchmarks` is only set `true` when `len(info.Models) > 0` in `loadArtifactList()`. Template gate `if s.HasBenchmarks` prevents rendering the badge. Benchmarks panel and tab are suppressed by `len(skill.Models) > 0 || skill.TesslScore > 0`. Tests `TestSkillsListingEffectiveness/omits_badge_when_no_data` and `TestSkillDetailBenchmarksTab/hides_tab_when_no_data` verify both cases. |
| 8 | Stale results (version mismatch) are not displayed | PASS | `loadEvalResults()` reads `skill_version` from each result JSON and skips any entry where `result.SkillVersion != currentVersion`. Only matching-version results are included in the returned `models` map. |

## Summary

All eight acceptance criteria are implemented and verified. The
implementation is coherent end-to-end: `openkata-eval` writes
labelled, versioned result files at the specified path;
`generate-versions --local` aggregates them with version
filtering; the web handlers propagate the model data from
`versions.json` into template structs; and the templates render
the listing badge, expandable breakdown, and detail table
correctly. Test coverage exists for criteria 4, 5, 6, and 7.

## Issues Found

None.
