---
name: grill-me
description: >
  Challenge a plan, spec, ADR, or design by interviewing the user
  relentlessly until all corners are covered. Use when the user
  wants to stress-test a plan, get grilled on their design, says
  "grill me", "challenge this", or "poke holes in this."
metadata:
  tags: "category:review, category:architecture"
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

6. **Offer to update** — Ask the user if they want you to
   apply the decisions to the artifact. If yes, update it.
   If no, leave it unchanged.

## Example Output

After grilling a spec, the summary looks like:

```markdown
## Decisions Made

1. **Auth strategy** — JWT with refresh tokens, not sessions.
   Reason: stateless scaling requirement.
2. **Rate limiting** — Per-user, not per-IP. Applied at API
   gateway level.
3. **Error format** — RFC 7807 Problem Details. No custom
   envelope.

## Changes to Apply

- Add "Rate Limiting" section to spec with per-user strategy
- Update Non-goals: remove "session management" (now in scope
  as JWT refresh)
- Add dependency: API gateway must support per-user rate limits
```

## Boundaries

- DOES read specs, ADRs, designs, plans, and codebase
- DOES challenge assumptions and identify gaps
- DOES update the artifact only when user explicitly confirms
- Does NOT modify files without asking first
- Does NOT implement code changes — only updates plans
