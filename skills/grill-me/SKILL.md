---
name: grill-me
description: >
  Challenge a plan, spec, ADR, or design by interviewing the user
  relentlessly until all corners are covered. Use when the user
  wants to stress-test a plan, get grilled on their design, says
  "grill me", "challenge this", or "poke holes in this."
---

# Grill Me

Challenge the user's plan until every decision is explicit and
every corner case is addressed.

## Workflow

1. **Read the artifact** — Read the spec, ADR, design, or plan
   the user wants challenged. Understand the full scope before
   asking anything.

2. **Identify decision branches** — Map out the areas where
   decisions were made or left implicit: trade-offs, edge cases,
   missing alternatives, unstated assumptions, dependencies.

3. **Interview one question at a time** — Walk down each branch
   of the decision tree, resolving dependencies between
   decisions one by one. For each question:
   - Provide your recommended answer
   - Wait for the user to confirm, correct, or expand
   - If a question can be answered by exploring the codebase,
     explore the codebase instead of asking

4. **Continue until covered** — Keep going until all branches
   are resolved. Do not stop early. The goal is shared
   understanding, not speed.

5. **Summarize** — When complete, briefly list the decisions
   made and any changes the user should apply to the artifact.

## Boundaries

- DOES read specs, ADRs, designs, plans, and codebase
- DOES challenge assumptions and identify gaps
- Does NOT modify any files
- Does NOT implement changes — only identifies them
