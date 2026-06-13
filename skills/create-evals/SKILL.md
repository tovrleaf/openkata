---
name: create-evals
description: >
  Generate diverse evaluation scenarios that test whether an agent using
  a skill produces correct output, including happy-path cases, edge cases,
  and failure scenarios across varying complexity levels. Produces structured
  scenario directories containing task prompts, weighted scoring checklists,
  and optional input fixtures for each scenario. Use when the user says
  "create evals", "generate eval scenarios", "write evals", "test this
  skill", or before publishing a skill.
metadata:
  version: "1.0.0"
  tags: "category:scaffolding, category:quality"
---

# Create Evals

Generate diverse eval scenarios for a skill. Each scenario
tests whether an agent following the skill produces the
expected output given a realistic task.

## Workflow

1. **Identify target** — Determine which skill to create
   evals for from the user's request.

2. **Read the skill** — Read the skill's SKILL.md and any
   `references/` to understand what behaviors it enforces.

3. **Plan scenarios** — Design 3–5 diverse scenarios:
   - At least one straightforward happy-path case
   - At least one edge case or tricky situation
   - Vary complexity (simple, moderate, challenging)
   - Cover different aspects of the skill's workflow

4. **Create scenario directories** — One directory per
   scenario in the skill's `evals/` directory. Use
   descriptive kebab-case names:
   ```text
   skill-name/
   └── evals/
       ├── scenario-happy-path/
       ├── scenario-edge-case/
       └── scenario-complex/
   ```

5. **Write files** — Each scenario directory contains:
   - `scenario.json` — Metadata (see structure below)
   - `criteria.json` — Weighted scoring checklist
   - `task.md` — Realistic task prompt
   - `inputs/` — Optional fixture files

6. **Verify** — Confirm all criteria weights sum to exactly
   100 per scenario.

7. **Report** — List the scenarios created with descriptions.

## File Formats

### scenario.json

```json
{
  "description": "Brief description of the scenario"
}
```

Add `"include": ["./inputs"]` when input fixtures exist.

### criteria.json

```json
{
  "context": "What the scenario tests",
  "type": "weighted_checklist",
  "checklist": [
    {
      "name": "criterion-name",
      "description": "Observable behavior to verify",
      "max_score": 15
    }
  ]
}
```

### task.md

```markdown
## Problem/Feature Description

Realistic scenario description.

## Output Specification

What files or artifacts the agent must produce.
```

## Criteria Rules

- Weights MUST sum to exactly 100 per scenario
- Weight guide:
  - 12–15: Core skill behaviors
  - 10–12: Important explicitly-emphasized rules
  - 6–8: Supporting rules that matter but aren't the focus
  - 4–6: Minor style or formatting details
- Each description must be specific enough for binary
  pass/fail scoring — no subjective criteria

## Task Rules

- Two sections: Problem/Feature Description and Output
  Specification
- Self-contained and realistic — a believable scenario
- Never reference the skill being tested
- Keep focused on the skill's domain

## Input Files

- Optional — only include when the scenario needs existing
  project context
- Keep minimal — just enough to make the scenario concrete
- Must be realistic (real-looking code, configs, docs)

## Example Scenario

For a commit-conventions skill:

```text
evals/
└── scenario-breaking-change/
    ├── scenario.json
    ├── criteria.json
    ├── task.md
    └── inputs/
        └── diff.patch
```

Where `criteria.json` might weight "includes breaking
change footer" at 15, "imperative mood header" at 12,
"scope present" at 10, etc., summing to 100.

## Boundaries

- DOES generate eval scenario files
- Does NOT run evals or judge output
- Does NOT modify the skill itself
- Does NOT publish or deploy

## Common Failures

- **Weights not summing to 100** — always verify arithmetic
  before writing criteria.json.
- **Criteria too vague** — "well-written commit message" is
  unjudgeable. "Header uses imperative mood" is binary.
- **Task leaks skill name** — the task must not hint at
  which skill is active or what behaviors are expected.
- **Scenarios all test the same aspect** — cover different
  parts of the workflow, not variations of one case.
- **Unrealistic inputs** — fixture files must look like
  real project artifacts, not toy examples.
