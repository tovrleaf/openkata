# Rationale

## Why phases are split into reference files

The main SKILL.md stays under 130 lines by deferring
phase details to `references/phase-*.md`. This keeps
the execution flow scannable while giving each phase
full step-by-step instructions. An agent reads the
phase file only when entering that phase — saving
context tokens for the actual work.

## Why mode detection uses a _current pointer

`specs/_current` lets the agent resume without scanning
every spec directory. One file read determines the active
spec. Without it, the agent would need to walk all spec
dirs looking for pending tasks — expensive in large repos.

## Why phases have gates

Each phase transition requires user confirmation. Without
gates, the agent builds the wrong thing (confirmed by
real-world testing). The cost of one "does this look
right?" question is far less than rebuilding after a
misunderstood requirement.

## Why tasks must be one commit

Multi-commit tasks create ambiguous resume points. When
a session dies mid-task, the progress log and git history
must agree on what's done. One task = one commit = one
log entry makes this unambiguous.

## Why progress logs exist

Agent sessions have finite context. The progress log
is the handoff mechanism between sessions — it records
what was done, what was decided, and what to do next.
Without it, resuming requires the user to re-explain
context that was already established.

## Why there are two depth modes

Quick (< 3 files, no design phase) and Standard (full
5-phase). Most changes don't need architecture design.
Forcing every change through five phases wastes time
on trivial work. The trigger check in the description
helps agents self-select the right depth.
