---
title: How we use tessl
description: Where the tessl CLI fits in how skills are reviewed, scored and released, which files it adds, and what we do not use it for.
order: 6
---

# How we use tessl

[Tessl](https://tessl.io/) is a command-line tool for skills in the
Agent Skills format. We use it to check and score skills, and
optionally to publish them to its registry. It is a development tool.
No workflow in this repository runs it, and nothing a user installs
depends on it.

## The commands

| Command | What we use it for |
|---------|--------------------|
| `tessl skill lint <dir>` | Checks the skill's structure. Run before and after changes. |
| `tessl skill review <dir>` | Scores the skill and lists suggestions and warnings. |
| `tessl skill review --optimize <dir>` | Proposes edits. We apply the ones that help. |
| `tessl skill import --workspace openkata <dir>` | Sets a new skill up for tessl, once. |
| `tessl tile publish --bump patch` | Publishes to the registry, an optional release step. |

The linter checks Markdown links in `SKILL.md` as file paths relative
to the skill folder, so a link to a file outside the skill fails
the lint. The optimize command asks questions interactively. When
that prompt cannot be answered, the proposed diff is applied by
hand.

The CLI is listed as a prerequisite in the contributing guide, and
`make check` tests that it is installed.

## What it adds to a skill

- **`.tessl-plugin/plugin.json`** is made by `tessl skill import`. It
  holds the name `openkata/<skill>`, the description, the `skills`
  list, the `private` flag, tessl's own `version`, and the `score`.
- **`.tesslignore`** lists what is left out of the published package,
  such as `CHANGELOG.md` and `evals/`.

At the repository root, `tessl.json` holds the workspace settings.

The `score` field is filled in by our review workflow after a review,
so scores can be read without running tessl. The `version` in
`plugin.json` is tessl's version of the published package, and
`tessl tile publish --bump` manages it. It is expected to differ from
the version of the skill itself, which comes from its git tag.

These files are development-only. They are left out when a skill is
installed or downloaded, and changelogs do not mention them.

## The quality gate

A skill has to score 95% or more on `tessl skill review`. The review
workflow for a skill is:

1. Lint, and fix structural problems.
2. Review, and read the score, suggestions and warnings.
3. Optimize, and apply the changes that make sense.
4. Lint again, to confirm the skill is still valid.
5. Record the final score in `plugin.json`.
6. Run our own skill design checklist. A skill without a
   `## Boundaries` section fails it.

To bring several skills up to the bar, the optimize workflow runs
the review and the evals for each one, up to three rounds, and
stops and reports the final score if a skill is still below 95%.

The score is the quality score that the site shows on a skill's
Benchmarks tab. [How we score skills](scoring.md) explains it next to
effectiveness.

## What we do not use it for

**Evals.** Tessl has its own scenario evals, run with
`tessl eval run`. We do not use them for scoring, because a mismatch
in the model makes their results unreliable. Our evals run on a local
runner, described in [how we use evals](evals.md). The decision
record that proposed using tessl evals,
[ADR 0012](https://github.com/tovrleaf/openkata/blob/main/docs/adr/0012-use-tessl-evals-and-registry-for-distribution.md),
is still marked as proposed.

**CI.** Reviews and scores are produced by hand, on a developer's
machine, and committed with the change. Publishing a skill to the
registry is manual as well.

## The registry

At release time, after the git tag, the release workflow asks
whether to publish the skill to the registry too. If yes, it runs
`tessl tile publish --bump patch` for the skill, to the `openkata`
workspace. A skill's registry page is at
`https://tessl.io/registry/openkata/<name>`.

The `private` field of `plugin.json` records whether a skill is
published. A skill with `private: false` counts as published, and
the Benchmarks tab then adds a "View on Tessl Registry" link. The
quality score always links to the registry page. Publishing to the
registry is separate from publishing to the artifact bucket, which
the git tag triggers, as described in
[where artifacts are stored](storage.md).

## Working with agents

The project rule for agents is to ask before running tessl commands.
The one exception is the optimize workflow, which runs them without
asking because you started it for that purpose.
