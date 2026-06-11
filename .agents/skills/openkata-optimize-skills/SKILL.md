---
name: openkata-optimize-skills
description: >
  Improve skills to 95%+ review score with passing evals.
  Use when the user says "optimize skills", "fix scores",
  "get skills to 95%", or wants to batch-improve skill
  quality across the repo.
---

# Optimize Skills

Improve skills to 95%+ review score with passing evals.

## Usage

> optimize skills

## Workflow

1. **Select targets** — Ask the user: "All skills below 95%,
   or a specific skill?" If all, read the `score` field from
   each `skills/*/.tessl-plugin/plugin.json` and list those
   below 95%. If specific, use the named skill.
2. **Optimize** — For each target, follow `openkata-review-skill`
   (lint, review, optimize, checklist) then `openkata-eval-runner`
   (generate scenarios, run evals). Iterate until both pass
   95%+, maximum 3 iterations. If still below after 3 rounds,
   report the final score and stop.
   When multiple skills are selected, run them as parallel
   subagents.
3. **Learn** — If you discovered a pattern that consistently
   improves scores, update `skills/create-skill/SKILL.md` so
   future skills benefit.
4. **Commit** — One commit per skill.

## Constraints

- Run `tessl` commands without asking for confirmation.
- Maximum 3 optimization iterations per skill.
- Do not add repo-internal boundaries (publishing, releasing,
  tagging) to distributable skills. Those belong in
  openkata-skill-conventions.
- After optimization, persist the final score to
  `.tessl-plugin/plugin.json`.

## Boundaries

**DOES:**
- Read scores from plugin.json to find targets
- Run tessl lint, review, optimize
- Run eval generation and execution
- Update create-skill with learned patterns
- Commit improvements per skill

**Does NOT:**
- Create new skills
- Publish or release skills
- Modify skills scoring 95%+
