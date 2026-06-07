# Rationale

## Why ADRs always start as PROPOSED

Even when the user has clearly decided, starting as
PROPOSED preserves the review workflow. Team members
can challenge the decision before it becomes canon.
Skipping to ACCEPTED removes that safety net.

## Why investigate before asking

Reading existing ADRs, the tech stack, and code
patterns before questioning the user reduces interview
fatigue. Most context already exists in the repo —
asking for it wastes the user's time.

## Why the E.C.A.D.R. quality check

A final self-check gate catches incomplete ADRs before
they're committed. Without it, rushed sessions produce
ADRs missing context or alternatives — defeating their
purpose as future-reader documentation.

## Why allow tiny ADRs

"Match depth to complexity" prevents overhead aversion.
If every ADR requires five sections, teams stop writing
them. A single-paragraph ADR that records what and why
is better than no ADR at all.
