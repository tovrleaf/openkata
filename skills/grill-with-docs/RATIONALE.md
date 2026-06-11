# Rationale

grill-with-docs is a domain-driven interview that
challenges plans against project language and writes
documentation (glossary, ADRs) inline as decisions
crystallise.

## Why the glossary is updated inline, not batched

Terms resolved mid-session lose context if batched.
By the end, the user may not remember why "account"
was changed to "customer." Immediate writes preserve
the reasoning while it's fresh.

## Why ADR creation uses a three-gate test

Without gates, every decision becomes an ADR and the
docs/ directory becomes noise. The three criteria
(hard to reverse, surprising, real trade-off) filter
for decisions that actually help future readers.

## Why the skill delegates to create-adr

The ADR format, numbering, and lifecycle are owned by
create-adr. Duplicating that logic creates drift.
Delegation keeps each skill focused on its own concern.

## Why this skill is separate from grill-me

grill-me is deliberately read-only — it challenges
without modifying. grill-with-docs writes files during
the session. Merging them would compromise grill-me's
"safe to run on anything" property.
