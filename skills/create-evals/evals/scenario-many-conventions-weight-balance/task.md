## Problem/Feature Description

You have a skill called `format-changelog` that enforces
12 specific conventions when generating changelogs:
structure (heading, preamble, unreleased section, version
ordering), entry format (dash prefix, no empty sections,
blank lines between sections), versioning rules (semver,
dates, breaking = major), and quality rules (entries
describe what not how, no vague entries). All 12
conventions are equally important.

Generate evaluation scenarios that would test whether an
agent using this skill produces correct output. You need
to create criteria checklists that capture the diverse
conventions without making any single criterion trivially
small-weighted.

## Output Specification

Create 3–5 scenario directories inside
`skills/format-changelog/evals/`. Each directory must
contain:

- `scenario.json` with a description field
- `criteria.json` with a weighted scoring checklist
  (weights summing to exactly 100)
- `task.md` with a realistic task prompt

Criteria must be specific enough for binary pass/fail
scoring. Find a sensible grouping strategy for the 12
conventions so individual criteria weights remain in the
meaningful range (6–15 points each).
