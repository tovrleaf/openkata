# Choosing a Database for the Inventory Service

## Problem/Feature Description

A small e-commerce startup is building a new inventory management service. The engineering team has been debating whether to use PostgreSQL or MongoDB for the service's primary data store. Both options have been prototyped internally, and the team needs a permanent record of the decision to onboard new engineers and avoid revisiting it later.

The project is a Node.js application (see `inputs/package.json`). There is currently no architectural documentation — this will be the team's first ADR. The decision involves trade-offs around query flexibility, schema enforcement, operational complexity, and the team's existing expertise.

Your job is to produce a properly structured Architecture Decision Record capturing this technology choice. Use the information available in the project files to inform the context and decision, and document the trade-offs honestly.

## Output Specification

- Create the ADR as a Markdown file in the correct directory within the project (under `inputs/`), with a properly formatted filename.
- The ADR must include all appropriate sections filled with real content.
- Write a brief `decision-log.md` file at the top level of your working directory summarizing: which file was created, what directory was created (if any), and what the ADR sequence number was assigned.
