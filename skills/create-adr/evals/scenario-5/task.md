# Documenting the Authentication Strategy Overhaul

## Problem/Feature Description

PayStream's internal API is consumed exclusively by other services within the platform — payment processors, fraud detection, and the webhook delivery engine all authenticate to it programmatically. The current authentication mechanism relies on server-side sessions backed by Redis, which was a reasonable choice when the API had a handful of trusted callers. Now the platform is expanding: third-party fintech partners will soon integrate directly, the number of service-to-service callers is expected to triple, and the infrastructure team is planning a multi-region deployment.

After a technical review, the team has agreed to move away from session-based authentication toward stateless token-based authentication for machine-to-machine communication. The API codebase is under `inputs/`. The project already has an ADR directory with prior decisions recorded — check it before creating anything new.

## Output Specification

- Create a new ADR documenting the authentication strategy decision, placed in the correct location within `inputs/` with a properly formatted filename.
- The ADR must be complete: all required sections filled with substantive content, no placeholder text remaining.
- Write an `auth-decision-summary.md` file at the root of your working directory containing: the filename of the ADR created, the names of any source files that would need to change as part of this decision, and a one-sentence summary of the architectural change.
