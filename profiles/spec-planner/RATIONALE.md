# Rationale

spec-planner is a planning-only agent that uses the
spec-workflow skill's first three phases to produce
specs and task breakdowns, then hands off to an
implementer.

## Why the planner stops after Phase 3

Planning and implementation require different modes
of thinking. Combining them in one agent dilutes
both. Stopping at tasks forces a clean handoff and
prevents the planner from making implementation
assumptions.

## Why the planner reads everything but writes only to specs/

Good plans require full codebase context — patterns,
conventions, existing implementations. But a planner
that modifies code is dangerous. Read-all write-narrow
gives context without risk.

## Why tasks must be independently implementable

If a task requires re-reading the full spec to
understand, the breakdown failed. Each task should
contain enough context for an implementer to act
without the full picture. This enables parallel work
and session resumption.
