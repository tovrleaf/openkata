# Rationale

## Why update the glossary inline, not batched

Terms resolved mid-session lose context if batched.
By the end, the user may not remember why "account"
was changed to "customer." Immediate writes preserve
the reasoning while it's fresh.

## Why a three-gate test for ADRs

Without gates, every decision becomes an ADR and the
docs/ directory becomes noise. The three criteria
(hard to reverse, surprising, real trade-off) filter
for decisions that actually help future readers.

## Why delegate to create-adr instead of creating directly

The ADR format, numbering, and lifecycle are owned by
create-adr. Duplicating that logic creates drift.
Delegation keeps each skill focused on its own concern.

## Why separate from grill-me

grill-me is deliberately read-only — it challenges
without modifying. grill-with-docs writes files during
the session. Merging them would compromise grill-me's
"safe to run on anything" property.
