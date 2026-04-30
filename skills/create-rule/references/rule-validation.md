# Rule Validation

Use this to verify a rule works against real files before
shipping.

## Goal

Prove that the rule:

- Is enforceable by reading it literally
- Does not conflict with existing tooling or rules
- Matches patterns already present in the repo
- Does not create contradictions with itself

## Minimum Validation

Check at least 2–3 representative files:

- **conforming file** — already follows the conventions.
  Confirms the rule matches existing patterns.
- **non-conforming file** — violates at least one convention.
  Confirms the rule would catch it.
- **edge case file** — unusual structure or mixed patterns.
  Confirms the rule handles ambiguity.

## Conflict Checks

Before finalizing, verify:

- No overlap with existing rules in `.agents/rules/`
- No contradiction with linter configs, `.editorconfig`,
  or formatter settings
- No contradiction between sections within the rule itself

## Common Failures

- Convention too vague — agent interprets instead of follows
- Convention conflicts with tooling — agent and linter fight
- Convention too specific — breaks on edge cases
- Rule too long — token cost exceeds value
- Rule overlaps another rule — agent gets contradictory
  instructions
