# Challenge the JWT Migration ADR

## Problem Description

The engineering team is preparing to sign off on ADR-014 (see `inputs/adr.md`), which proposes migrating the authentication system from session-based auth to JWT. Before the ADR is ratified, the tech lead has asked for a rigorous challenge of every decision and assumption in it.

The current authentication implementation lives under `inputs/src/`. It is the live system the ADR is proposing to migrate away from, so it provides important context about the current behaviour and configuration.

Conduct a thorough challenge of the ADR, covering every decision, trade-off, and open question. Produce a challenge report as `challenge-report.md`.

Do not modify `inputs/adr.md` or any files under `inputs/src/`.
