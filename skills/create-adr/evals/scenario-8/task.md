# Documenting the Order Storage Migration

## Problem/Feature Description

OrderFlow has been running its order storage on MongoDB since the early days of the product. The decision made sense at the time: partner schemas varied widely, the team had strong Mongo experience, and launch pressure favoured familiarity. Two years on, the landscape has changed considerably. Partner schemas have stabilised and are now consistent enough to express as a relational model. The operations team has been burned multiple times by the lack of multi-document transactions — compensating saga logic has introduced subtle consistency bugs that have required manual data remediation. A significant portion of the engineering effort in the last two quarters has gone toward working around MongoDB's consistency limitations rather than building features.

After an extended internal review, the engineering leadership has made a firm decision to migrate order storage to PostgreSQL. The pg package is already present in the codebase (it was originally pulled in for the analytics reporting service). The project files are under `inputs/`. The `docs/adr/` directory contains prior architectural decisions for this project — you should read these before creating any new documentation.

## Output Specification

- Create a new ADR documenting the decision to move order storage to PostgreSQL. Place it in the correct location within `inputs/` with a properly formatted filename. The existing ADRs in `docs/adr/` determine the correct sequence number — check before writing.
- Update the ADR in `inputs/` that originally documented the order storage choice to reflect that it has been superseded.
- Write a `migration-summary.md` file at the root of your working directory containing: the filename of the new ADR, the filename of the ADR that was superseded, and a one-paragraph description of why the team is making this change.
