---
status: Draft
depth: Standard
---

# Benchmarks and Version List in Production

## Story

As a visitor to openkata.dev, I want a skill's page to show its
benchmarks and every released version, as the site does when it
runs locally, so that I can judge a skill's quality and pick a
version.

## Background

Checked against the live site and the code on 2026-10-27:

- A skill page on openkata.dev has the tabs Overview, Files,
  Rationale, Changelog and Acknowledgments. There is no Benchmarks
  tab, and the version control shows only the current version, for
  example `v1.0.1 (Latest)`
- The skills listing shows no effectiveness badge
- The repository has 13 skills and the live listing shows 11
- In production the artifact page is built by
  `loadArtifactDetailS3` in `cmd/openkata-web/handlers.go`. It
  reads only `version`, `description` and `tags` from the index,
  and sets the version list to the current version alone. It never
  sets the tessl score, the published flag or the eval results
- Those are set only by `loadArtifactDetailLocal`, which runs when
  the site is started outside Lambda, so the local site shows the
  Benchmarks tab and all versions
- The publish workflow builds `versions.json` with
  `generate-versions` in its S3 mode, which writes the tessl score
  but not the eval results (`models`). Only the `--local` mode
  writes them
- The listing page reads `models` from the index in production too,
  so with no `models` there it shows no badge
- The copy of `versions.json` bundled in the website lists 5 rules
  while the README badge counts 6, and 1 of its 13 skills has eval
  results

## Requirements

### Index

- The `versions.json` that the publish workflow uploads includes,
  for each skill that has them, the tessl score, the published flag
  and the eval results by model, in the same shape the local mode
  writes
- The eval results come from the skill's own files at the released
  version, so a result for another version is never shown. The
  design decides where they are read from

### Artifact page

- The production artifact page shows the Benchmarks tab when the
  skill has a tessl score or eval results, with the same content as
  the local site, and no tab when it has neither
- The version control lists every released version of the artifact,
  newest first, with the newest marked as the latest, and choosing
  one opens that version's page
- Rules and profiles list all their versions too

### Listing

- The skills listing shows the best effectiveness score and its
  model for skills that have eval results, as the local site does

### Versions in the bucket

- Every released version of an artifact, that is every git tag of
  the form `<type>/<name>/v<version>`, has a folder in the bucket.
  Versions that are missing are published

### Quality

- The production loaders are covered by tests with a fake bucket
  that holds several versions and some skills without benchmarks
- Local and production show the same data for the same files
- The Reference pages that say which mode of `generate-versions`
  writes which fields (where artifacts are stored, how we score
  skills) are updated to match the change

## Out of Scope

- Generating eval results or tessl scores; that stays a local step
- Sorting or filtering the catalog (spec 0019)
- Changing how scores are computed

## Open Questions

- Where do the eval results come from in production: the
  `evals/results/` folder that is already uploaded with each
  released version, or a committed copy of the index?
- The sandbox does not pass `--model` to `kiro-cli`, so a result
  filed under a model name may have run with another model. Is
  that settled before the scores are shown publicly?
- Why does the live listing show 11 of the 13 skills? Are two not
  released?
