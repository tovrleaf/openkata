# Rationale

create-adr detects architectural decisions in
conversation and guides creation of Architecture
Decision Records with a structured lifecycle.

## Why ADRs always start as PROPOSED

Even when the user has clearly decided, starting as
PROPOSED preserves the review workflow. Team members
can challenge the decision before it becomes canon.
Skipping to ACCEPTED removes that safety net.

## Why the skill investigates the codebase before asking

Reading existing ADRs, the tech stack, and code
patterns before questioning the user reduces interview
fatigue. Most context already exists in the repo —
asking for it wastes the user's time.

## Why the E.C.A.D.R. quality check exists

A final self-check gate catches incomplete ADRs before
they're committed. Without it, rushed sessions produce
ADRs missing context or alternatives — defeating their
purpose as future-reader documentation.

## Why tiny single-paragraph ADRs are allowed

If every ADR requires five sections (Context, Drivers,
Options, Decision, Consequences), the process feels
heavy and people stop writing them entirely. Allowing
a single paragraph that records what was decided and
why removes that friction. A three-sentence ADR is
infinitely more useful than no ADR at all.
