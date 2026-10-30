---
title: Skills, rules and profiles
description: What each kind of OpenKata artifact is, how they differ, and the kata vocabulary.
order: 2
---

# Skills, rules and profiles

OpenKata distributes three kinds of artifact. Each is a folder of
markdown files that an agent can read. Skills follow the
[Agent Skills](https://agentskills.io/) standard, so they work
with any agent that supports it. The three differ in when they
apply and what they are for.

## At a glance

| | Skill | Rule | Profile |
|---|---|---|---|
| Kata term | kata | dojo kun | sensei profile |
| What it is | A multi-step workflow agents follow | An always-on constraint applied to every session | An agent role with scoped permissions |
| Main file | `SKILL.md` | `RULE.md` | `<name>.md` |
| Folder | `skills/<name>/` | `rules/<name>/` | `profiles/<name>/` |

Every artifact folder also holds a `CHANGELOG.md` and a
`RATIONALE.md`. The rationale explains why the artifact is the way
it is, because practices should carry their reasoning and not only
the rule.

## Skills

A skill is a codified practice: a workflow an agent follows to get
a kind of task done the same way every time. The `description` in
its frontmatter tells the agent when the skill applies, and the
body is the workflow.

This is the start of the frontmatter of the `grill-me` skill:

```yaml
name: grill-me
description: >
  Challenge a plan, spec, ADR, or design by interviewing the user
  relentlessly until all corners are covered. ...
metadata:
  version: "1.0.1"
  tags: "category:review, category:architecture"
```

Besides `SKILL.md`, a skill folder can hold `references/` for
detail the workflow reads when it needs it (the `spec-workflow`
skill keeps its phase details there), and `evals/` with the
scenarios used to score the skill.

## Rules

A rule is a constraint that applies all the time, not a workflow
that runs on request. Each convention in a rule is written so it
can be checked literally.

The `git-naming` rule, for example, sets branch names to
`type/short-description` in lowercase kebab-case, and requires an
`Assisted-by` trailer on every AI-assisted commit. An agent does
not choose when to apply it; it applies to every commit.

## Profiles

A profile defines a role. It limits what an agent may touch and
which rules and skills it works with. Profiles are plain markdown
files with these sections:

- **Constraints**: what the agent must and must not do
- **Companions**: the skills the role requires or may use
- **Scope**: what it may read, and the only places it may write
- **Design intent**: why the role is shaped this way

The `spec-planner` profile can read the whole codebase but may
write only under `specs/`. It plans a feature and stops after the
task breakdown, and it never implements.

## How they fit together

The three types work as layers. The `spec-planner` profile says it
must follow the `markdown-style` and `git-naming` rules, so
those rules apply to everything it writes and commits. It uses the
`spec-workflow` skill for the planning steps and stops after phase
3. The profile sets the role and the limits, the rules set the
standards, and the skill supplies the workflow.

## The kata vocabulary

OpenKata names its parts after the martial art that inspired it:
kata are codified sequences of movements, practiced until they are
mastered. The vocabulary is an overlay on the technical terms,
which stay valid, and it applies across the docs and the tooling.

| Concept | Kata term | Description |
|---------|-----------|-------------|
| Skills | **kata** | Codified practices agents follow |
| Rules | **dojo kun** | Always-on constraints applied to every session |
| Profiles (roles) | **sensei profiles** | Agent role definitions with scoped permissions |
| Prompts | **kata forms** | Standardized templates for commits, PRs, reviews |
| ADRs | **dojo records** | Architecture decisions preserved for the school |
| MCP server | **the dojo** | Where kata are served and practiced |
| Framework | **the ryu** (school) | The complete system of practices |

"The dojo serves kata" is the short version: the MCP server is the
dojo, and the skills, rules and profiles it serves are the
practices. The reasoning behind the vocabulary is in the
[vocabulary decision](https://github.com/tovrleaf/openkata/blob/main/docs/adr/0004-adopt-kata-vocabulary.md).

## Browse the artifacts

Everything described here is listed on the site: see the
[skills](/skills/), [rules](/rules/) and [profiles](/profiles/).
