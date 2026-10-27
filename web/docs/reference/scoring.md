---
title: How we score skills
description: The two scores a skill can carry, effectiveness and the tessl quality score, how each is produced and recorded, and where they are shown.
order: 5
---

# How we score skills

A skill can carry two scores, and they answer different questions.
**Effectiveness** says whether an agent that follows the skill does
the right things. The **quality score** says how well formed and
well written the skill is, according to a review by tessl. Only
skills have scores. Rules and profiles do not.

## The two scores

| | Effectiveness | Quality score |
|---|---|---|
| Question | Does an agent using the skill do what the skill asks? | How good is the skill itself, by tessl's review? |
| Produced by | The eval runner: an agent attempts, a judge scores | `tessl skill review` on the skill's files |
| Result | A percentage for each model | One percentage for the skill |
| Recorded in | `skills/<name>/evals/results/<model>/<time>.json` | `score` in `skills/<name>/.tessl-plugin/plugin.json` |
| Target | 95% overall | 95% or more |

Both are snapshots. Nothing recomputes them, so after a skill
changes, the evals and the review are run again and the new values
are committed. The review workflow records the score in
`plugin.json` so it can be read without running tessl again.

## Effectiveness

Each scenario is scored as the share of its criteria points that the
judge passed, and a model's effectiveness is the average over all
the scenarios, rounded to a whole percent. [How we use
evals](evals.md) describes the scenarios, the judge and the scoring
in detail.

Two rules decide which results count:

- For each model, only the newest report is used.
- A report counts only if its `skill_version` is the skill's current
  version. After a new release the old results are ignored until the
  evals have been run again for the new version.

A skill can have results for several models. The model name shown
is the label given with `--model-label` when the evals were run.

## Quality score

The quality score is the percentage that `tessl skill review`
reports for the skill. The review workflow runs the review, applies
the worthwhile suggestions, and writes the final percentage into
`score` in the skill's `plugin.json`. The Benchmarks tab describes
the score as a structural quality measure covering formatting,
metadata and the structure of `SKILL.md`.

The same file has a `private` flag. A skill that is not private is
treated as published on the Tessl registry, and the Benchmarks tab
then links to its registry page at
`https://tessl.io/registry/openkata/<name>`. A score of zero is
treated as no score.

## Where the scores are shown

Both scores appear only for skills that have them. A skill without
evals or without a review simply shows no scores, and a skill with
neither has no Benchmarks tab.

- **On the listing**, a skill with eval results shows its best
  effectiveness and the model that earned it, for example `96%`
  with a model name. If two models tie, one of them is named.
  Expanding the entry lists the score of every model.
- **On the skill's page**, the Benchmarks tab holds:
  - the **Quality score**, linked to the registry page, with a
    **View on Tessl Registry** link when the skill is published
  - the **Effectiveness** table, one row for each model
  - the **Scenarios**, for each model, with the score and the
    description of each scenario and a bar that is green from 95%,
    yellow from 70% and red below that

[Using the catalog](../guides/using-the-catalog.md) shows how to open
the tab.

## How the scores reach the site

The website reads the scores from `versions.json`, the index
described in [where artifacts are stored](storage.md):

| Score | In the index | Written by |
|-------|--------------|------------|
| Effectiveness | `models`: label, effectiveness, scenarios per model | The repository mode only |
| Quality score | `tessl_score`, and `published` | Both modes |

The index built from the repository, with `make versions`, reads the
eval reports in each skill's `evals/results/` and the skill's
`plugin.json`. The index built from the bucket by the publish
workflow reads `plugin.json` from the released version, and does not
read the eval reports.

## Keeping scores current

The scores are a record of the last time someone ran the tools. A
typical round for a skill is:

1. Change `SKILL.md`.
2. Run the evals until the overall result is 95% or more, and the
   review until its score is 95% or more.
3. Commit the new eval report and `plugin.json` with the change.
4. Release the skill, and rebuild the index.

A skill that is new, or changed since its last run, shows no
effectiveness until its evals have been run for the current version.
