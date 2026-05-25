# Challenge the Prism Deployment Architecture

## Problem Description

The platform engineering team at a SaaS company has drafted a deployment architecture document for Prism, a new real-time analytics platform (see `inputs/deployment-design.md`). The document is due for stakeholder sign-off at the end of the week, and the engineering lead has asked for a rigorous independent challenge before it goes forward.

The architecture document records decisions about cloud provider, orchestration, storage, messaging, and caching — but the team worked quickly and the rationale behind each choice is thin. Several decisions were made implicitly, and alternative approaches were not discussed.

Your job is to challenge every decision in the document systematically. For each decision, surface the trade-offs, alternatives that were not considered, risks, and any assumptions that are baked in without being stated. Make sure nothing important is glossed over before the architecture is locked in.

## Output Specification

Write your complete challenge analysis to a file called `challenge-report.md`. The report should walk through each decision area, provide your own assessment, and end with a structured summary of what was settled and what changes the team should make to the document.

Do not modify `inputs/deployment-design.md`.
