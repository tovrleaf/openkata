# DevOps Agent Profile

## Problem Description

A platform engineering team manages a monorepo that includes application code, infrastructure-as-code, CI configuration, deployment scripts, and database provisioning. The head of engineering wants to introduce a dedicated agent to handle all the "infrastructure and operations" work — a broad remit that the team has struggled to pin down.

Internally the team has been calling this the "DevOps agent," but they haven't agreed on exactly what it should own. Some engineers think it should cover everything from build pipelines to database provisioning; others argue it should be more focused. The agent should have a clearly bounded scope — specific enough to do meaningful, targeted work, but not so narrowly defined that it's useless for most day-to-day infrastructure tasks.

Create a profile for this DevOps role. Explore the project structure to understand what directories and configuration files exist, then decide on a scope that is appropriately specific without being over-constrained. The project contains the following top-level directories: `.github/`, `infra/`, `scripts/deploy/`, `scripts/db/`, `src/`, `docs/`.

## Output Specification

Create an agent profile for the DevOps role. Save it as a markdown file in the appropriate profiles location for this project. Choose an appropriately short, lowercase name for the file.
