## Problem/Feature Description

You have a skill called `deploy-service` that orchestrates
service deployments through a 7-step workflow: verify
prerequisites, run tests, build artifacts, tag the release,
push to registry, deploy to staging, and promote to
production. It has a `references/` directory with
detailed runbook docs and a `scripts/` directory with
deployment scripts.

The skill is located at `skills/deploy-service/SKILL.md`
(provided in inputs along with its references and scripts).
Generate evaluation scenarios that thoroughly test whether
an agent produces correct output across the full range
of situations this skill handles.

## Output Specification

Create 3–5 scenario directories inside
`skills/deploy-service/evals/`. Each directory must
contain:

- `scenario.json` with a description field (add
  `"include": ["./inputs"]` when fixtures are present)
- `criteria.json` with a weighted scoring checklist
  (weights summing to exactly 100)
- `task.md` with a realistic, self-contained task prompt

Include input fixtures where the scenario needs existing
project context to be concrete. Scenarios must vary in
complexity and cover different aspects of the skill's
workflow.
