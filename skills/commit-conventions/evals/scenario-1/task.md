# Sprint Branch Planning

## Background

Your team is kicking off a two-week sprint. Four pieces of work have been assigned
and each needs its own git branch. The project uses GitHub and the team tracks work
in an issue tracker — some of these tasks are linked to open issues.

Here are the four work items for the sprint:

1. **Issue #456** — Customers are reporting intermittent 500 errors when the session
   token expires during a long-running request. Needs investigation and a patch.

2. **Issue #461** — Product wants a new "saved searches" panel added to the
   dashboard so users can rerun recent queries without retyping them.

3. No issue — The CI pipeline's Node.js version is pinned to v16; it needs to be
   bumped to v20 and the workflow YAML updated accordingly.

4. No issue — The `README.md` and inline API docs are out of date following last
   quarter's auth refactor; they need to be rewritten to match the current behavior.

## Your Task

Propose a git branch name for each of the four work items above.

Write your proposals to a file named `branch-plan.md`. For each item include the
work item description and the branch name you've chosen. Briefly note any naming
decision that wasn't obvious.
