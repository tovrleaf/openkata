# Rule Design Checklist

Use this before finalizing a generated or revised rule.

## Problem and Scope

- Is the target problem concrete and repeatable?
- Is the rule enforcing one coherent concern, not several?
- Are out-of-scope cases obvious from the wording?
- Is it clear which file types or contexts the rule applies to?

## Convention Quality

- Is every convention specific enough to enforce literally?
- Could an agent follow each convention without interpreting
  intent?
- Are conventions stated as directives, not explanations?
- Are hard requirements distinguishable from preferences?
- Is each convention stated exactly once? No duplicates
  across sections.

## Packaging

- Is RULE.md sufficient on its own?
- If not, are extra details in `references/` instead of
  bloating the main file?
- Are conventions grouped by concern with clear section
  headings?

## Validation

- Has the rule been checked against 2–3 representative files?
- Are there conflicts with existing tooling (linters,
  formatters, `.editorconfig`)?
- Does the rule overlap with existing rules in `.agents/rules/`?
- Are edge cases or exceptions documented?

## Portability

- Is the wording generic unless the user explicitly asked for
  a repo-bound rule?
- Does the rule avoid environment assumptions it cannot
  justify?
- If the rule is repo-bound, does it say so plainly?

## Boundaries

- Is it clear where this rule's concern ends?
- Could this rule conflict with a future rule in an adjacent
  domain?

## Consistency

- Does the name follow existing patterns (lowercase-hyphenated)?
- Is there a companion skill that enforces this rule's
  conventions?
- Does the CHANGELOG format match other artifacts?

## Token Cost

- Is the rule under 100 lines? (Rules load every session.)
- Does every line earn its place?
- Are explanations of *why* removed in favor of *what*?
- Could any section move to `references/` without losing
  enforceability?
