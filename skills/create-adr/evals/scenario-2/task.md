# Selecting a State Management Solution for the Dashboard

## Problem/Feature Description

A SaaS company's engineering team is building a React-based analytics dashboard. As the dashboard has grown, different engineers have reached for different state management approaches: some components use React's built-in `useState`/`useContext`, others have started using Zustand, and a recent pull request introduced Redux Toolkit. The inconsistency is causing bugs when state changes in one part of the app don't propagate correctly to others.

The team wants to commit to a single state management strategy across the dashboard. They've already had informal discussions about the trade-offs, but nothing is documented. The tech lead wants an ADR to align the team before the next sprint.

The project files are in `inputs/`. Explore them to understand the current state of things before producing the ADR. There is already one architectural record in the project.

## Output Specification

- Create a new ADR documenting the chosen state management approach.
- The ADR must be placed in the correct directory with the correct filename format.
- If any information needed to complete a section cannot be determined from the project files or general knowledge, use an appropriate inline marker to flag what needs follow-up.
- Write a `research-notes.md` file at the root of your working directory summarizing what you found in the codebase (files examined, patterns observed, existing ADRs found).
