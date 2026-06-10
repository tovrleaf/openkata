# Grilling Session on Payment Service Integration Plan

## Problem/Feature Description

Your team is planning to integrate a third-party payment provider into an existing e-commerce platform. A junior engineer has written a rough design document describing how the payment service will interact with the order management and notification systems. The document outlines several domain concepts — "checkout", "payment", "transaction", "order" — and describes a flow where "the system handles payment confirmation and then sends a notification".

You have been asked to run a grilling session on this plan to surface unclear decisions, sharpen the domain language, and produce updated documentation as terms are resolved.

The codebase already has a glossary at `docs/context/GLOSSARY.md` that defines several key terms used across the platform. The plan needs to be stress-tested against this glossary, and any terms resolved during the session must be documented immediately.

## Output Specification

- `docs/context/GLOSSARY.md` — updated with any terms resolved during the grilling session (created or updated in-place)
- `session-log.md` — a log of each question asked and the answer provided, in the order they were raised

The session log must show each question on its own, with any recommendation the agent offered alongside it, followed by a sample answer. Simulate the user's responses yourself based on the plan — the goal is to show how the session would proceed. Continue the session until all major domain terms and design branches in the plan are covered.
