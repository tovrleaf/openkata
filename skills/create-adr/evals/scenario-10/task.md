# Documenting the Event Streaming Infrastructure Decision

## Problem/Feature Description

DataPulse is a data ingestion platform that collects events from hundreds of external data sources — IoT sensors, SaaS webhooks, and partner data feeds — and writes them to a PostgreSQL analytical store. The system currently writes events synchronously on the HTTP request path. This design worked at low volume, but as data source partners have scaled up their throughput, p99 latency has climbed to unacceptable levels and the API is becoming a bottleneck. Under write pressure, the database's connection pool exhausts and ingestion requests start failing.

The engineering team has decided to decouple event ingestion from event persistence by introducing a message streaming layer between the ingest API and the write workers. The team has done a round of internal research and benchmarking across the major options in this space. They have settled on a specific streaming technology and want to capture the decision — including the research that went into it — before beginning the implementation sprint.

The project files are under `inputs/`. There are no existing architectural decision records in this codebase, so this will be the first ADR. Create it in the standard location.

## Output Specification

- Create a new ADR documenting the message streaming technology decision, placed in the correct location within `inputs/` with a properly formatted filename.
- The ADR must be complete: all required sections filled with real, substantive content.
- Include any references, links, or external resources you consulted during your research directly in the ADR.
- Write a `streaming-decision-summary.md` at the root of your working directory containing: the ADR filename, the chosen technology, and the key technical reason it was selected over the alternatives.
