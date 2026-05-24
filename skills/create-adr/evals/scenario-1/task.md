# Updating the API Communication Strategy

## Problem/Feature Description

A fintech startup originally built its internal services to communicate over REST. That decision was documented in the team's architectural records eighteen months ago. Since then, the product has grown: new services have been added, teams complain about over-fetching data, and the mobile team is struggling with multiple round-trips for each screen load.

Leadership has decided to adopt GraphQL for client-facing APIs going forward, while keeping REST for internal service-to-service calls. This is a material change from the existing decision, and the engineering team wants the architectural record updated correctly — preserving the history of why REST was originally chosen, and clearly documenting the new direction and its reasoning.

The project files are in `inputs/`. The existing architectural records are in `inputs/docs/adr/`.

## Output Specification

- Create a new ADR file that documents the decision to adopt GraphQL for client-facing APIs.
- Update the existing REST ADR file to reflect its new status and cross-reference the new decision.
- Write a `change-summary.md` at the root of your working directory describing: which files were created or modified, the sequence number assigned to the new ADR, and what status the old ADR was changed to.
