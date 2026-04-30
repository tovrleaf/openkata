# Skill Validation

Use this when a skill needs prompt-level validation before
shipping.

## Goal

Prove that the skill:

- Triggers for obvious requests
- Triggers for paraphrased requests
- Does not trigger for nearby but unrelated work
- Gives another agent enough detail to act without guessing

## Minimum Prompt Set

Write at least 3 prompts:

- **positive-obvious** — direct request using the most likely
  trigger words
- **positive-paraphrased** — same job, different wording
- **negative-adjacent** — close enough to confuse a weak
  description, but should not load the skill. The best
  negatives are near-misses: queries that share keywords or
  domain with the skill but need something different. "Write
  a fibonacci function" is too easy as a negative for a
  deploy skill — "set up a staging environment" is a real
  near-miss that tests discrimination.

Add more prompts only when the surface area is large or the
user explicitly wants deeper validation.

## Comparison Modes

Choose the lightest comparison that answers the risk:

- **manual simulation** — read the prompt against the skill
  and judge whether the trigger and workflow would work
- **before vs after** — compare the current skill against the
  revised skill
- **trigger wording A vs B** — use when the main risk is
  activation quality rather than body content

## When Validation Can Be Skipped

Skip only when all of these are true:

- The edit does not change the trigger surface
- The edit does not change the workflow meaning
- The edit does not add or remove important resources

If any of those changed, run at least a lightweight prompt
simulation.

## Review Checklist

For each prompt, record:

- Should the skill trigger?
- Which words or phrases should cause activation?
- Which part of the body should guide the next step?
- Where could an agent take a shortcut or misread?

Common failures:

- Description too vague to trigger
- Description so broad that it over-triggers
- Description summarizes workflow, tempting the agent to skip
  the body
- Body assumes repo facts it never tells the reader to discover
- Examples are longer than the rules they clarify
