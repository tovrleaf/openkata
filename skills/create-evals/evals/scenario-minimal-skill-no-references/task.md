## Problem/Feature Description

You have a small utility skill called `rename-files` that
batch-renames files using a pattern. It has only a
SKILL.md — no references, no scripts, no assets. The
skill's workflow is: detect files matching a glob, preview
the rename plan, and apply renames after confirmation.

The skill is located at `skills/rename-files/SKILL.md`
(provided in inputs). Generate evaluation scenarios for
this skill.

## Output Specification

Create 3–5 scenario directories inside
`skills/rename-files/evals/`. Each directory must contain:

- `scenario.json` with a description field
- `criteria.json` with a weighted scoring checklist
  (weights summing to exactly 100)
- `task.md` with a realistic task prompt

Only include input fixtures if the scenario genuinely
requires project context. Do not fabricate references
to files or directories that don't exist in the skill.
