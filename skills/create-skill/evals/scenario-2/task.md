# Observability Pipeline Onboarding Skill

## Problem Description

A platform engineering team at a mid-sized SaaS company has standardized on OpenTelemetry for distributed tracing, metrics, and log collection. New engineers and on-call responders regularly need to instrument services and configure the collector pipeline — a process that involves a dozen configuration options, several environment-specific quirks, and detailed runbook steps for common failure modes.

Currently this knowledge lives in a long Confluence page that nobody reads. The team wants to package the onboarding and instrumentation workflow as a reusable agent skill so that engineers can ask their AI assistant to walk them through it. The skill needs to cover: identifying what a service is missing, generating the OTel configuration boilerplate, and pointing the engineer toward the right troubleshooting reference for their environment (Kubernetes vs. bare metal vs. Docker Compose).

The team has emphasized that the skill should be maintainable: the core workflow steps should be easy to scan, while the environment-specific runbooks, edge cases, and exhaustive configuration references should be kept separate so they don't drown out the essential procedure.

## Output Specification

Create a skill package for this OpenTelemetry onboarding workflow. Place all output in a `otel-onboarding/` directory:

- `otel-onboarding/SKILL.md` — the complete skill definition (core workflow only)
- Any supporting documentation or reference material the skill needs (use appropriately named subdirectories)
- Any scripts the skill relies on for deterministic or fragile operations

Produce a brief `packaging-notes.md` in the working directory (not inside the skill folder) that explains what content you placed in the main file versus supporting files, and why.
