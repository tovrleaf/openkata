---
name: openkata-review-rule
description: >
  Reviews a rule against the rule design checklist, validates
  conventions are enforceable, and checks token cost. Use when
  the user creates a new rule, updates an existing rule, asks
  to review a rule, or says "review" in the context of rules.
---

# Review Rule

Review a rule against the design checklist. No external
tooling — rules are validated by reading and checking.

## Workflow

1. **Identify the rule** — Determine which rule to review.
   If ambiguous, list available rules and ask.

2. **Read the rule** — Read RULE.md in full.

3. **Run the checklist** — Evaluate each section of the
   rule design checklist from
   `skills/create-rule/references/rule-design-checklist.md`:
   - Problem and Scope
   - Convention Quality
   - Packaging
   - Validation
   - Portability
   - Token Cost

4. **Check enforceability** — For every convention in the
   rule, verify it is literally enforceable. If a convention
   requires interpretation ("keep it clean", "be consistent"),
   flag it.

5. **Check line count** — RULE.md must be under 100 lines.
   If over, identify what can move to `references/`.

6. **Check for duplicates** — Scan for conventions stated
   more than once across sections.

7. **Report** — Summarize findings per checklist section.
   Use ✅ for pass, ⚠️ for minor issues, ❌ for failures.

8. **Feed learnings back** — After every review, check
   whether the findings could improve `create-rule` or
   the rule design checklist itself. If a pattern keeps
   surfacing, ask the user if they want to update the
   relevant skill so future rules benefit.

## When Applying Fixes

- Fix the rule first, then update `create-rule` if the
  issue was a gap in the creation workflow.
- Don't add guidance to `create-rule` for one-off issues.
  Only update it when a pattern repeats across multiple
  rules.

## Example Scenario

User: "Review the bash-style rule."

1. Reads `rules/bash-style/RULE.md`
2. Evaluates against checklist: Problem/Scope ✅,
   Convention Quality ⚠️ (one subjective qualifier found)
3. Reports: 97 lines (under limit), one enforceability
   issue flagged, no duplicates

## Boundaries

**DOES:**
- Evaluate rules against the design checklist
- Check enforceability, line count, and duplicates
- Report findings and suggest improvements

**Does NOT:**
- Fix rules automatically without user confirmation
- Create new rules
- Modify `create-rule` without asking
