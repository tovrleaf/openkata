---
status: Implementing
depth: Standard
---

# Skill Effectiveness Display

## Story

As a developer browsing openkata.dev, I want to see how
effective each skill is per model so I can choose the
right skill and model pairing for my project.

## Requirements

### Eval Result Storage

- `openkata-eval` results committed to repo at:
  `skills/<name>/evals/results/<model_id>/<timestamp>.json`
- Already excluded from Files tab, archive downloads,
  MCP installs, and tessl publish (existing `evals/`
  filter)
- Result JSON includes: model ID, model label, skill
  version, overall percentage, per-scenario scores
- Add `--model-label` flag (e.g., "Sonnet 4.6") stored
  as `"model_label"` in result JSON
- Add `"skill_version"` to result JSON (read from
  CHANGELOG.md or git tag)

### versions.json Extension

`generate-versions` reads eval results and produces:

```json
{
  "create-adr": {
    "version": "1.3.0",
    "description": "...",
    "tags": "...",
    "models": {
      "claude-sonnet-4-20250514": {
        "label": "Sonnet 4.6",
        "effectiveness": 92
      },
      "claude-haiku-4-20250514": {
        "label": "Haiku 4",
        "effectiveness": 78
      }
    }
  }
}
```

- Only includes results where `skill_version` matches
  the current published version (no stale data)
- Picks latest result per model if multiple exist
- CI uses `generate-versions --local` (repo is checked
  out) then uploads `versions.json` to S3
- Skills without matching eval results: `models` field
  absent
- Scenario names and descriptions are read from the
  result JSON (no separate file parsing needed)

### Display — Listing Page

- Show best effectiveness score with model name:
  `96% claude-opus-4.6`
- Collapsible row expands to show all tested models:
  `claude-sonnet-4.5: 92% | claude-haiku-4: 78% | claude-opus-4.6: 96%`
- Skills without eval data show no effectiveness info
- Skills with data sort above those without (when
  sorting by relevant criteria in spec 0019)

### Display — Detail Page

New "Benchmarks" tab (shown only when eval data or
tessl data exists, same pattern as Acknowledgments tab).

Tab contents:

1. **Tessl section** — tessl quality score (from
   `plugin.json` `"score"` field) and "View on Tessl
   Registry" link (if skill is public)
2. **Effectiveness table** — per-model scores:

| Model | Effectiveness | Scenarios |
|-------|--------------|-----------|
| claude-opus-4.6 | 96% | 5/6 pass |
| claude-sonnet-4.5 | 92% | 5/6 pass |
| claude-haiku-4 | 78% | 4/6 pass |

3. **Scenario breakdown** — list of eval scenarios with
   per-model pass/fail indicators

- Tab hidden if no eval results AND skill not published
  to tessl (nothing to show)
- Rules and profiles have no Benchmarks tab

## Constraints

- No LLM calls at build or deploy time
- `versions.json` remains single source of truth for
  listings
- No new infrastructure
- Model label is human-readable; API model ID stored
  for reproducibility
- Only current-version eval results displayed (version
  mismatch → omit)
- Eval runs are manual, infrequent, developer-initiated

## Acceptance Criteria

1. `openkata-eval` stores model label and skill version
   in result JSON
2. Results stored at
   `skills/<name>/evals/results/<model>/<ts>.json`
3. `generate-versions --local` produces `models` field
   with effectiveness per model
4. `/skills/` listing shows best effectiveness score
   with model name
5. Listing rows expand to show per-model breakdown
6. Detail page shows model table
7. Skills without eval results show no effectiveness
   info
8. Stale results (version mismatch) are not displayed

## Out of Scope

- Token tracking (deferred)
- Multi-model eval pipeline automation
- Historical effectiveness trends
- Sorting (separate spec 0019)
- Publishing to tessl registry (separate spec 0020)
- Quality lint (separate spec 0021)

## Open Questions

None.

Date: 2026-06-20
