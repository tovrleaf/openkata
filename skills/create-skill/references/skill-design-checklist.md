# Skill Design Checklist

Use this before finalizing a generated or revised skill.

## Problem and Scope

- Is the target problem concrete and repeatable?
- Is the skill solving one coherent job rather than several?
- Is the audience clear?
- Are out-of-scope cases obvious from the wording?

## Trigger Quality

- Does the description say what the skill does?
- Does it say when to use it?
- Does it include phrases a user would actually say?
- Does it avoid vague language like "helps with" or "handles"?
- If over-triggering is a risk, does the description narrow
  scope clearly?

## Workflow Quality

- Are the steps in the correct order?
- Does the skill investigate local context before asking
  avoidable questions?
- Are decision points and defaults explicit?
- Are required inputs and expected outputs stated?
- Are external dependencies named only when necessary?

## Packaging

- Is SKILL.md sufficient on its own?
- If not, are extra details in `references/` instead of
  bloating the main file?
- Are `scripts/` included only for deterministic or fragile
  tasks?
- Are `assets/` included only when they materially improve
  execution?

## Validation

- Are there concrete examples or scenarios?
- Does the skill define what success looks like?
- Does it describe how to catch obvious failure modes?
- If assumptions were needed, are they stated explicitly?

## Portability

- Is the wording generic unless the user explicitly asked for
  a repo-bound skill?
- Does the skill avoid environment assumptions it cannot
  justify?
- If the skill is repo-bound, does it say so plainly?
