---
name: grill-me
description: >
  Challenge a plan, spec, ADR, or design by interviewing the user
  relentlessly until all corners are covered. Identifies edge cases,
  questions assumptions, probes failure modes, and surfaces missing
  requirements. Use when the user wants to stress-test a plan, get
  grilled on their design, says "grill me", "challenge this", or
  "poke holes in this."
metadata:
  version: "1.0.0"
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

3. **Interview one question at a time** — Ask exactly ONE
   question per message. Never combine multiple questions.
   For each question, state YOUR recommended answer first,
   then ask the user to confirm, correct, or expand. Wait
   for a response before the next question. If a question
   can be answered by exploring the codebase, explore the
   codebase instead of asking.

4. **Continue until covered** — Keep going until all branches
   are resolved. Do not stop early. The goal is shared
   understanding, not speed.

5. **Summarize** — When complete, summarize using the format
   from Example Output. Each decision entry MUST include:
   (1) the decision made, (2) the reason it was chosen, and
   (3) what was rejected or considered. Do not list decisions
   without reasons.

6. **Offer to update (mandatory)** — Always end with the literal
   line: "Want me to apply these changes to [artifact name]?"
   where [artifact name] is the specific file being grilled
   (e.g., "spec.md", "design.md", "0003-use-postgres.md").
   This MUST be the final line of your output. Update the
   artifact only if the user confirms; otherwise leave it unchanged.

## Example Output

After grilling a spec, the summary looks like:

```markdown
## Decisions Made

1. **Auth strategy** — JWT with refresh tokens. Reason:
   stateless scaling requirement eliminates session affinity.
   Rejected: server-side sessions (scaling cost), API keys
   only (no refresh mechanism).
2. **Rate limiting** — Per-user at API gateway. Reason:
   prevents single-user abuse without penalizing shared IPs.
   Rejected: per-IP (punishes NAT users), no limiting
   (risk of runaway costs).
3. **Error format** — RFC 7807 Problem Details. Reason:
   standard format reduces client integration effort.
   Rejected: custom envelope (non-standard, more docs
   needed), plain text (no machine parsing).

## Changes to Apply

- Add "Rate Limiting" section to spec with per-user strategy
- Update Non-goals: remove "session management" (now in scope
  as JWT refresh)
- Add dependency: API gateway must support per-user rate limits

Want me to apply these changes to spec.md?
```

Every decision entry MUST follow this exact pattern:
`**Topic** — Decision. Reason: why. Rejected: alternatives.`
Do not omit the reason or the rejected alternatives.

## Boundaries

- DOES read specs, ADRs, designs, plans, and codebase
- DOES challenge assumptions and identify gaps
- DOES update the artifact only when user explicitly confirms
- Does NOT modify files without asking first
- Does NOT implement code changes — only updates plans
