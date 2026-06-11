# ADR Criteria

Only offer to create an ADR when ALL THREE are true:

1. **Hard to reverse** — the cost of changing your mind
   later is meaningful
2. **Surprising without context** — a future reader will
   wonder "why did they do it this way?"
3. **The result of a real trade-off** — there were genuine
   alternatives and you picked one for specific reasons

If any is missing, skip the ADR.

## What qualifies

- Architectural shape (monorepo, event sourcing, etc.)
- Integration patterns between contexts
- Technology choices that carry lock-in (database, auth
  provider, deployment target)
- Boundary and scope decisions (who owns what data)
- Deliberate deviations from the obvious path
- Constraints not visible in code (compliance, SLAs)
- Rejected alternatives when rejection is non-obvious

## What does NOT qualify

- Easy to reverse — just reverse it when needed
- Not surprising — nobody will wonder why
- No real alternative — "we did the obvious thing"
- Library choices (unless lock-in is significant)
- General best practices everyone follows

## When criteria are met

Delegate to the `create-adr` skill. Do not duplicate its
workflow — just pass the decision context and let it handle
the format, numbering, and file placement.
