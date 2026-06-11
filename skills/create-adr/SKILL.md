---
name: create-adr
description: >
  Detects architectural decisions in conversation and guides creation of
  Architecture Decision Records. Activate when the user is making, discussing,
  or debating: technology choices (languages, frameworks, databases, cloud
  services), structural patterns (monolith vs microservices, event-driven,
  CQRS, API design), cross-cutting conventions (error handling, logging,
  authentication, API versioning), hard-to-reverse decisions (data models,
  public API contracts, infrastructure topology), or deviations from existing
  standards. Also activate when the user explicitly asks to create an ADR.
metadata:
  version: "1.2.0"
  tags: "category:documentation, category:architecture"
---

# Create ADR

## ADR Lifecycle

```text
PROPOSED → ACCEPTED → (SUPERSEDED or DEPRECATED)
```

- New ADRs always start as `PROPOSED` — including ADRs that
  supersede an existing one. Never set a new ADR to `ACCEPTED`.
- For material changes to a decision, write a new ADR that
  supersedes the old one. The new ADR's status is `PROPOSED`.
  Update the old ADR's status to `SUPERSEDED` and link to the
  new one.
- Minor corrections (typos, missing details) can be edited
  in place.

## Workflow

1. **Detect** — An architectural decision is in play.
2. **Propose** — Ask the user: "This looks like an architectural decision.
   Would you like to record it as an ADR?" If this is the first ADR in the
   session, also ask: "Would you like me to walk through the details one
   question at a time, or should I draft it from what I already know?"
3. **Explore the codebase** — Before asking the user anything, gather what
   you can on your own:
   - Read existing ADRs in `docs/adr/` for related or superseded decisions
   - Check the tech stack (package.json, go.mod, requirements.txt, etc.)
   - Find code patterns related to the decision area
   - Identify affected files and existing conventions

   Don't ask the user what you can look up yourself.
4. **Gather context** — Collect what you couldn't find in the codebase.
   For each question, provide your recommended answer with justification
   based on what you found in step 3. Let the user confirm or correct.
   Collect: the core decision, key decision drivers, rationale,
   alternatives with pros/cons/rejection reasons, consequences
   (positive, negative, neutral), non-goals, and relevant links.
5. **Generate** — Create the ADR using the template at
   [assets/adr-template.md](assets/adr-template.md). Fill every section with
   real content — do not leave placeholder text. The ADR must include:
   - YAML frontmatter with `status: PROPOSED`, `date`, and `authors`.
     Always set status to `PROPOSED` — never `ACCEPTED` or `DECIDED`.
     This is non-negotiable even when the user has already decided.
   - Context explaining the problem and forces at play
   - Decision drivers as a prioritized list
   - The decision in active voice
   - At least two alternatives with pros, cons, and rejection reasons
   - Consequences split into positive, negative, and neutral
   - Non-goals section if anything is explicitly out of scope (omit if none)
   - Reversibility section if it's useful to note how to undo the decision
     or what would trigger reconsideration (omit if none)
   - References section listing any links, related ADRs, or resources
     discovered during research (omit if none)
6. **Place the file** — Save to `docs/adr/NNNN-<slug>.md` where:
   - `NNNN` is the next sequential number (zero-padded to 4 digits)
   - `<slug>` is a lowercase, hyphenated summary of the decision
   - Check existing files in `docs/adr/` to determine the next number
   - Example: `docs/adr/0003-use-postgresql-for-persistence.md`
7. **Confirm** — Show the user the generated ADR and ask if any adjustments
   are needed before finalizing.
8. **Find blind spots** — When the ADR is finalized, check if
   `grill-with-docs` exists. If it does, ask: "Want me to
   grill this ADR against your domain language?" If yes,
   activate it. If `grill-with-docs` doesn't exist, try
   `grill-me` instead. If neither exists, skip silently.

## Quality self-check (E.C.A.D.R.)

Before finalizing an ADR, verify:

- **E**xplicit problem statement — Context makes the problem unambiguous
- **C**omprehensive options analysis — at least 2 alternatives with honest
  pros/cons
- **A**ctionable decision — specific enough to act on, in active voice
- **D**ocumented consequences — positive, negative, and neutral impacts
- **R**eferences included — any links from research are listed

For sections that cannot be filled from available data, insert investigation
prompts: `[INVESTIGATE: description of what needs follow-up]`

## Match depth to complexity

Omit optional sections (Non-goals, Reversibility, References)
when they add no information. A two-paragraph ADR for a
straightforward choice is better than a bloated one.

When a decision directly maps to code changes, add an
Implementation Plan section describing affected paths and
patterns to follow.

## Example Scenario

User: "Should we use Postgres or DynamoDB for the order service?"

1. Skill detects a technology choice decision
2. Asks: "This looks like an architectural decision. Want an ADR?"
3. Reads existing ADRs, checks go.mod/package.json for current DB
4. Drafts ADR with both options, pros/cons, and a recommendation
5. Saves to `docs/adr/0007-use-postgresql-for-order-service.md`

