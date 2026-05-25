# Documenting the Caching Strategy Decision

## Problem/Feature Description

Prism Analytics serves executive dashboards to enterprise customers. The reporting API aggregates multi-dimensional business metrics and exposes them through a REST API consumed by the dashboard frontend. As the customer base has grown, API response times have climbed. Profiling sessions consistently point to repeated, expensive data retrieval as the bottleneck — the same aggregated results being computed on every request with no reuse between calls.

The engineering team has evaluated approaches to reduce the load and has settled on an approach. They need the decision formally recorded before the implementation sprint begins.

The project files are in `inputs/`. There are no existing architectural decision records in this project. Create the ADR documenting the caching strategy decision and save it in the standard location within `inputs/`.

## Output Specification

- Create a new ADR documenting the caching strategy decision, placed in the correct location within `inputs/` with a properly formatted filename.
- The ADR must be complete with substantive content — no placeholder text.
- Write a `caching-decision-summary.md` at the root of your working directory containing: the ADR filename, the chosen caching approach, and one sentence on why the alternatives were not selected.
