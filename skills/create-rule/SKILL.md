---
name: create-rule
version: 1.0.0
description: >
  Creates rule files that define always-on constraints for agent
  sessions. Investigates the repo for existing conventions, writes
  a RULE.md with clear structure, and validates against real usage.
  Use when the user wants to create a rule, add a convention,
  enforce a standard, establish a project-wide constraint, define
  coding style, standardize formatting across generated files,
  notices inconsistent output across sessions, or is setting up a
  new project's conventions.
---

# Create Rule

Create rules — always-on constraints applied to every agent
session. Unlike skills, rules don't need triggers. They're loaded
at session start and apply to all work.

## Workflow

1. **Gather intent** — Understand what the rule should enforce.
   Collect:
   - A short name (lowercase, hyphenated)
   - What conventions or constraints to enforce
   - Any style guides or standards to draw from

2. **Investigate the repo** — Before asking questions, search
   for facts that reduce ambiguity:
   - Existing rules and conventions
   - Style guides, linter configs, or formatting standards
   - Patterns in existing files that reveal implicit conventions

   Don't ask the user what you can look up yourself.

3. **Clarify** — Ask only questions that materially affect the
   rule. Push until these are clear:
   - Which file types or contexts the rule applies to
   - Hard requirements vs preferences
   - Any exceptions or edge cases

4. **Write RULE.md** — Structure:

   ```markdown
   # Rule Name

   Brief description of what this rule enforces.

   ## Section

   - Convention one
   - Convention two
   ```

   Writing principles:
   - **Be specific.** "Use `-` for unordered lists" is
     enforceable. "Use consistent formatting" is not.
   - **State conventions, not explanations.** The agent doesn't
     need to know why a convention exists to follow it.
   - **Group by concern.** Separate sections for structure,
     spacing, emphasis, etc.
   - **Keep it short.** Rules are loaded every session — every
     line costs tokens on every interaction.

5. **Validate** — Check the rule against existing files:
   - Pick 2–3 representative files in the repo
   - Verify the conventions are consistent with what's already
     there
   - Note any conflicts with existing patterns

6. **Confirm** — Show the user the created rule and ask if
   adjustments are needed.

## Example Scenario

User: "Every model formats markdown differently, I want
consistency."

1. Skill checks repo for `.editorconfig`, linter configs,
   existing rules — finds none
2. Asks: "Which conventions matter most — headings, lists,
   code blocks, line length?"
3. Reads 2–3 existing markdown files to extract implicit
   patterns
4. Creates `markdown-consistency/RULE.md` with specific,
   enforceable conventions grouped by concern — see
   [example-rule.md](references/example-rule.md) for the
   finished output
5. Validates against existing files — no conflicts

## Common Failures

- **Too vague to enforce** — "use consistent formatting" is
  unenforceable. "Use `-` for unordered lists" is clear.
- **Conflicts with existing tooling** — the rule says one thing,
  the linter config says another. Check for `.editorconfig`,
  linter configs, and formatter settings first.
- **Too long** — rules load every session. A 200-line rule costs
  tokens on every interaction, even when irrelevant.
- **Overlaps with existing rules** — check `.agents/rules/` for
  rules that already cover the same concern.

A good rule is enforceable by reading it literally. If an agent
has to interpret intent, the rule is too vague.

## Quality References

Before finalizing, check the rule against:

- [Rule design checklist](references/rule-design-checklist.md)
- [Rule validation](references/rule-validation.md)
- [Rule token optimization](references/rule-token-optimization.md)
