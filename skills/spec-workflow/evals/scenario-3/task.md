# Resuming Implementation of a Partially-Complete Spec

## Problem/Feature Description

You are an AI developer agent using the spec-workflow skill. A developer has asked you to continue working on a feature. The spec at `specs/0002-notifications/` has tasks.md with 3 tasks — task 1 is Done, tasks 2 and 3 are Pending. The spec.md status is "Implementing". The `specs/_current` file points to this spec.

You should describe exactly what steps you would take to implement task 2 (the next pending task), in the exact order you'd perform them. Include every file modification, status update, log entry, and commit you would make. Be explicit about the sequence.

Do not write implementation code — describe the process you would follow step by step.

## Output Specification

Produce a numbered list of the exact steps you would take to complete task 2, including:
- What you read first
- What status changes you make and when
- What you add to the progress log and when
- What your commit message looks like
- What you do after the commit (move to next task or ask user)
