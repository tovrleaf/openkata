## Problem/Feature Description

You have a skill called `greet-user` that generates
personalized greeting messages. The skill has three
workflow steps: detect the user's locale, select a
greeting template, and output the formatted message.

The skill is located at `skills/greet-user/SKILL.md`
(provided in inputs). Generate evaluation scenarios
that test whether an agent using this skill produces
correct output across different situations.

## Output Specification

Create 3–5 scenario directories inside
`skills/greet-user/evals/`. Each directory must contain:

- `scenario.json` with a description field
- `criteria.json` with a weighted scoring checklist
  (weights summing to exactly 100)
- `task.md` with a realistic task prompt

Include input fixtures only where they add necessary
project context.
