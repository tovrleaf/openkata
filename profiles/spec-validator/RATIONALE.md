# Rationale

## Why never read tasks.md

Tasks represent implementation intent, not
requirements. A validator who reads tasks develops
sympathy for partial progress. "It's in the plan"
is not the same as "it's in the code." Blindness
to tasks preserves objectivity.

## Why report without editorializing

The validator's job is to observe, not prescribe.
"Requirement X is not met" is useful. "You should
do Y to fix X" oversteps — the implementer may
have context the validator lacks.

## Why check for out-of-scope work

Scope creep is invisible without explicit checking.
Features built beyond the spec may seem helpful but
indicate the implementer deviated from the plan.
Flagging them surfaces drift early.
