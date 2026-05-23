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
- Does the body avoid restating trigger criteria already
  covered by the description?
- Does the description answer both "what" and "when" in a
  single read? (completeness)
- Could this skill accidentally trigger instead of another
  skill in the same workspace? Are the trigger terms unique
  to this skill's domain? (distinctiveness)

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

## Consistency

- Does the name follow existing patterns (lowercase-hyphenated)?
- Does the description avoid leaking workflow steps that would
  let the agent skip the body?
- Is there a companion rule or skill that should be referenced?
- Does the CHANGELOG format match other artifacts?

## Validation

- Have you tested with at least one near-miss negative prompt?
- Could someone verify the output is correct without re-reading
  the whole skill?

## Token Budget

- Is the front-loaded cost justified for activation frequency?
- Could any section move to `references/` without losing
  workflow clarity?

## Boundaries

- Does the skill state what it DOES?
- Does it state what it Does NOT do?
- Are the boundaries specific enough to prevent scope creep?
- Would an agent know which files and actions are off-limits?
