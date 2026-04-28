---
name: create-rule
description: >
  Creates rule files that define always-on constraints for agent
  sessions. Investigates the repo for existing conventions, writes
  a RULE.md with clear structure, and validates against real usage.
  Use when the user wants to create a rule, add a convention,
  enforce a standard, establish a project-wide constraint, define
  coding style, or standardize formatting across generated files.
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
