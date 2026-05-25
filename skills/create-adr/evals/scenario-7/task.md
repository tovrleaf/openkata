# Documenting the Frontend State Management Decision

## Problem/Feature Description

ClearLedger's finance operations dashboard has grown from a handful of views to over thirty screens covering transaction review, reconciliation, approval workflows, and reporting. The codebase currently uses React's local component state — each feature manages its own data fetching and caching independently. As the application has scaled, the team has run into recurring pain points: sibling components that share data fetch it redundantly, filter state set in one panel doesn't propagate to dependent charts, and a performance audit found the same API endpoint being called six times on a typical page load.

The frontend lead has made a decision on a state management approach and wants it formally recorded before the team begins the migration sprint. The project files are under `inputs/`. There are no existing architectural records in this project, so this will be the first ADR.

## Output Specification

- Create a new ADR documenting the state management decision, saved in the correct location within `inputs/` with a properly formatted filename.
- The ADR must be complete: all sections filled with real, substantive content — no placeholder text remaining.
- Write a `state-decision-notes.md` file at the root of your working directory containing: the ADR filename created, the chosen state management approach, and a brief explanation of why the options you didn't choose were rejected.
