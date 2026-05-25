# Parallel Skill Evaluation Runner

## Problem/Feature Description

An AI platform team needs to validate three newly authored skills before they can be merged into the production tile registry. Each skill must be exercised by an isolated evaluation run — the runs must not share filesystem state, and the team wants all three to execute concurrently rather than waiting for one to finish before starting the next.

The repository hosts all three skills. Previously the team ran evaluations sequentially by switching branches, which took three times as long and occasionally left the working tree dirty when an eval crashed mid-run. The new approach should give each evaluation its own isolated workspace and launch all three in the background simultaneously, collecting results only after every eval has finished.

Write a single shell script `run-parallel-evals.sh` that automates this end-to-end: workspace preparation, parallel execution, and cleanup once all runs have completed. Use `tessl eval run skills/<skill-name>/` as the evaluation command for each skill.

## Output Specification

Produce a shell script `run-parallel-evals.sh` that:

- Sets up isolated workspaces for three evaluations: `skill-auth`, `skill-payments`, and `skill-notifications`
- Launches all three evaluations in the background simultaneously
- Waits for all three to complete before exiting
- Reports the outcome once all runs finish
- Cleans up all created workspaces and their branches after the runs complete
