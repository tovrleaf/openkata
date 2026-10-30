---
title: Authoring skills, rules and profiles
description: Where artifacts live, what each kind contains and how long it may be, how to create one, and the review it has to pass.
order: 1
---

# Authoring skills, rules and profiles

This page is for adding a skill, a rule or a profile to OpenKata. It
says where the files go, what each kind of artifact has to contain,
and what it must pass before it is merged. The
[skills, rules and profiles guide](../guides/skills-rules-profiles.md)
explains what the three kinds are for.

You need [Go](https://go.dev/dl/) 1.26 or newer, to build the MCP
server, and the
[Tessl CLI](https://docs.tessl.io/introduction-to-tessl/installation),
to review skills.

## Where artifacts live

Each kind can be distributable, which means it is published and
installed by users, or local, which means only this repository uses
it.

| Kind | Distributable | Local |
|------|---------------|-------|
| Skill | `skills/<name>/` | `.agents/skills/<name>/` |
| Rule | `rules/<name>/` | `.agents/rules/<name>/` |
| Profile | `profiles/<name>/<name>.md` | `.agents/profiles/<name>.md` |

Distributable artifacts are symlinked into `.agents/` so that the
repository can use them too. For example,
`.agents/skills/create-adr` points to `../../skills/create-adr`.
Local artifacts are not tagged or released.

A distributable artifact must work on its own, in any repository. It
must not name a local skill, depend on this repository's paths or
tools, or assume a particular agent.

### Names

Names are lowercase and hyphenated.

- Distributable skills have no prefix and follow a pattern:
  `create-<artifact>`, `review-<artifact>`,
  `<domain>-conventions` or `<name>-workflow`.
- Local skills and rules start with `openkata-`.
- Distributable rules follow a pattern such as `<domain>-style`,
  `<domain>-naming` or `<domain>-system`, for example
  `markdown-style` and `git-naming`.
- Profile names are lowercase, and one word is preferred.

## What an artifact contains

### A skill

```text
skills/<name>/
├── SKILL.md
├── CHANGELOG.md
├── RATIONALE.md
├── references/       detail that is loaded on demand
├── scripts/          executable code
├── assets/           templates and other resources
└── evals/            scenarios that score the skill
```

`SKILL.md` has frontmatter with a `name`, a `description`, and a
`metadata` block with the `version` and the `tags`:

```yaml
---
name: skill-name
description: >
  What the skill does. Use when the user wants to ...
metadata:
  version: "1.0.0"
  tags: "category:workflow"
---
```

- **The description is for activation.** It says what the skill does
  and when to use it, in words a user would say. The steps belong in
  the body, not in the description.
- **The body is procedural.** It is written in the imperative, is
  concise, and includes one complete example.
- **A `## Boundaries` section is required**, listing what the skill
  does and does not do. A `## Common Failures` section lists the
  mistakes an agent would otherwise make.
- **`SKILL.md` stays under 500 lines.** Longer detail moves to
  `references/`, and deterministic steps move to `scripts/`.
- A new skill starts at version `1.0.0`, never `0.x.y`.
- Sources the skill draws on go in `references/ACKNOWLEDGMENTS.md`,
  one line for each: the name and link, the license, what was
  adapted, and the version it was adopted in.

Distributable skills also have `evals/`, described in
[how we use evals](../reference/evals.md). A skill that has been set
up for tessl carries `.tessl-plugin/plugin.json` and a
`.tesslignore`, described in
[how we use tessl](../reference/tessl.md). Both are development
files and are not installed.

### A rule

```text
rules/<name>/
├── RULE.md
├── CHANGELOG.md
└── RATIONALE.md
```

`RULE.md` has a `name`, a `description` and `tags` in its frontmatter,
then a heading, a line saying what the rule enforces, and sections
grouped by concern, each a list of conventions.

- **Every convention must be enforceable by reading it literally.**
  "Use `-` for unordered lists" can be applied. "Use consistent
  formatting" cannot.
- **State the convention, not the reason.** An agent does not need to
  know why to follow it. The reasons go in `RATIONALE.md`.
- **`RULE.md` stays under 100 lines.** Rules are loaded in every
  session, so every line costs tokens on every interaction.
- Keep a rule to one concern, and check that no other rule covers it.

### A profile

```text
profiles/<name>/
├── <name>.md
├── CHANGELOG.md
└── RATIONALE.md
```

A profile is a short markdown file: a title, a one-line description
of the role, and these sections. The published profiles also have a
`description` and `tags` in their frontmatter.

- **Constraints**: what the agent must always or never do, and which
  rules it follows
- **Scope**: what it may read and the only places it may write, and
  when to stop and hand work to a person or another profile
- **Design intent**: how it should approach decisions in its domain

Scope is a boundary, not a task list. A scope that is too broad, such
as "the backend", covers half the repository. One that is too narrow
leaves the agent unable to do useful work. A constraint that applies
to every agent is a rule, not a profile constraint. The `create-profile`
skill suggests keeping a profile under 40 lines.

## Changelog and rationale

Every artifact has a `CHANGELOG.md` in the Keep a Changelog format,
with an `## [Unreleased]` section at the top. How a release adds a
version to it is described in
[versioning and the manifest](../reference/versioning.md). Changelogs
describe changes that users of the artifact see. They do not mention
the tessl files.

`RATIONALE.md` explains the design: why the artifact is the way it is
and what trade-offs were made. Each heading answers one "why"
question in a paragraph. The site shows it as the Rationale tab. It
is stored with each version, but it is left out of installs and
archives.

## Creating one

The `create-skill`, `create-rule` and `create-profile` skills lead
you through the work, and you can follow their steps by hand. In
short:

1. Say what the artifact should do and when it applies.
2. Look in the repository for existing artifacts and conventions that
   overlap.
3. Write the main file, following the sections above.
4. Add `CHANGELOG.md`, and `RATIONALE.md` where it applies.
5. Place it, and symlink it into `.agents/` if it is distributable.
6. Run the review below.

## The review

Nothing is committed until the review passes.

**Skills**

- The score on `tessl skill review` is 95% or more.
- The skill has a `## Boundaries` section.
- It passes the skill design checklist, which covers the problem and
  scope, how well it triggers, the workflow, packaging, portability,
  consistency and token budget.
- It has been tried with two or three realistic prompts that should
  trigger it, and at least one that should not.
- The evals pass at 95% or more.

**Rules**

- They pass the Quality Checklist in `create-rule`: every convention
  is specific enough to apply mechanically, none needs the agent to
  guess an intent, and sections are grouped by concern.
- The rule was checked against two or three existing files, and
  against any `.editorconfig`, linter or formatter settings, for
  conflicts.
- No existing rule covers the same concern.
- Rules do not use tessl.

**Profiles**

- The scope matches the real structure of the repository.
- It does not overlap an existing profile.
- Every rule it names exists.
- It says when to stop and hand off.

Tessl commands are run by hand. The project asks agents to confirm
with you before running them.
