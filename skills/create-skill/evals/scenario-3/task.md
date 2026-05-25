# Structured Code Review Skill

## Problem Description

An engineering team at a B2B software company has developed an internal code review process that goes beyond style linting. Reviews must check: that each PR has a linked ticket, that API contract changes are backward-compatible, that new endpoints follow the team's REST naming conventions, and that sensitive data (PII fields, secrets) is never logged. The team currently documents this process in an internal wiki, but reviewers apply it inconsistently because the page is long and easy to skim.

The engineering manager wants to package this structured review process as an agent skill so that engineers can trigger a consistent, thorough review without memorizing the checklist. The skill should guide the agent from first contact with the PR through to producing a structured review comment.

When creating the skill, make sure the guidance is direct enough that another agent following it wouldn't need to improvise — every step should tell the agent exactly what action to take.

## Output Specification

Create a skill package for this structured code review workflow. Place all output in a `code-review/` directory:

- `code-review/SKILL.md` — the complete skill definition

Produce a `skill-summary.md` in the working directory (not inside the skill folder) that notes: the skill name, the trigger phrases used in the description, and one potential pitfall that someone using this skill might encounter.
