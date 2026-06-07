---
name: grill-with-docs
description: >
  Domain-driven grilling session that challenges a plan against
  the project's language, sharpens terminology, and writes
  documentation inline as decisions crystallise. Use when the
  user wants to stress-test a plan against domain language, says
  "grill with docs", "sharpen terms", "challenge against domain",
  or wants a grilling session that produces glossary updates and
  ADRs as it goes.
metadata:
  version: "1.0.0"
  tags: "category:review, category:architecture, category:documentation"
---

# Grill With Docs

Interview the user relentlessly about their plan while
maintaining domain documentation inline. Walk down each
branch of the design tree, resolving dependencies between
decisions one by one.

## Workflow

1. **Read the artifact** — Read the plan, spec, or design.
   Understand the full scope before asking anything.

2. **Explore the codebase** — Before asking questions, look
   for existing documentation:
   - `docs/context/GLOSSARY.md` — domain language
   - `docs/context/CODEBASE.md` — context boundaries
   - `docs/adr/` — existing architectural decisions
   - Code patterns related to the plan's domain

3. **Interview one question at a time** — Ask exactly ONE
   question per message. For each question, state YOUR
   recommended answer first. Wait for a response before
   continuing. If a question can be answered by exploring
   the codebase, explore instead of asking.

4. **Challenge against the glossary** — When the user uses
   a term that conflicts with `GLOSSARY.md`, call it out:
   "Your glossary defines 'X' as Y, but you seem to mean
   Z — which is it?"

5. **Sharpen fuzzy language** — When the user uses vague or
   overloaded terms, propose a precise canonical term.
   "You're saying 'account' — do you mean the Customer or
   the User?"

6. **Discuss concrete scenarios** — Stress-test domain
   relationships with specific scenarios that probe edge
   cases and force precision about boundaries.

7. **Cross-reference with code** — When the user states how
   something works, check whether the code agrees. Surface
   contradictions: "Your code does X, but you just said
   Y — which is right?"

8. **Update GLOSSARY.md inline** — When a term is resolved,
   update `docs/context/GLOSSARY.md` immediately. Do not
   batch. Create the file lazily if it doesn't exist. See
   [glossary-format.md](references/glossary-format.md).

9. **Update CODEBASE.md if needed** — When context
   boundaries or relationships are clarified, update
   `docs/context/CODEBASE.md`. Create lazily. See
   [codebase-format.md](references/codebase-format.md).

10. **Offer ADRs sparingly** — Only when all three criteria
    in [adr-criteria.md](references/adr-criteria.md) are
    met. When warranted, delegate to the `create-adr` skill.

11. **Continue until covered** — Keep going until all
    branches are resolved. Do not stop early.

12. **Summarize and offer to apply** — Summarize decisions
    made, then ask: "Want me to apply these changes to
    [artifact name]?"

## Boundaries

**DOES:**
- Interview relentlessly about a plan or design
- Update GLOSSARY.md inline as terms are resolved
- Update CODEBASE.md when boundaries are clarified
- Delegate to create-adr when ADR criteria are met
- Cross-reference claims against code

**Does NOT:**
- Implement code changes
- Create ADRs directly (delegates to create-adr)
- Replace grill-me for quick artifact reviews
- Modify files without resolving terms first

## Common Failures

- **Batching glossary updates** — write terms as they're
  resolved, not at the end. Users forget context.
- **Creating ADRs too eagerly** — most decisions don't
  need one. Apply the three-gate test strictly.
- **Accepting vague terms** — "the service handles it"
  is not precise. Push until the term has a one-sentence
  definition and an _Avoid_ list.
