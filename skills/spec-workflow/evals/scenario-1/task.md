# Break Down an Approved Feature Spec into Implementation Tasks

## Problem/Feature Description

A developer on your team has finished writing the specification for a password reset feature and the spec has been approved. The spec is stored at `specs/0001-password-reset/spec.md` and `specs/_current` already points to it.

Before any coding begins, the team needs a clear breakdown of the work into discrete, ordered implementation tasks. The developer will tackle one task at a time and commit after each one — so each task should represent exactly the right amount of work for a single focused commit. Tasks should be ordered so that anything with dependencies on earlier work comes later.

Your job is to produce the implementation task list so the developer can start coding immediately after reviewing it. Do not write any implementation code — only the task breakdown file.

## Output Specification

- Create the tasks breakdown file in the appropriate location within the spec directory at `specs/0001-password-reset/`
- Each task must have enough detail for a developer to pick it up and start working without further clarification
- The task list should cover all the requirements in the spec
- Do not skip writing the progress tracking section at the bottom of the file
