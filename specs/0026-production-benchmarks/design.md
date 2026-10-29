# Design: Benchmarks and Version List in Production

Status: decided. The sandbox model fix is dropped; the other
decisions are confirmed.

## Where the data comes from

Eval results are files in the skill's folder,
`evals/results/<model>/<time>.json`, committed with the skill. The
publish workflow uploads the whole artifact folder at the tag, so the
results of a released version are already in the bucket under
`skills/<name>/<version>/evals/results/`. Installs and archives leave
`evals/` out, but it is stored.

**Decided: read the results from the bucket.** The index generator, in
its bucket mode, reads the results from the released version's folder,
as it already reads `.tessl-plugin/plugin.json`. This keeps one rule
everywhere: what the site shows for a version is what was committed with
that version. The alternative, committing a copy of the index, would
need a manual step after every release and is what lags today.

## The release order

A result counts only when its `skill_version` equals the skill's current
version, and the runner takes that from the newest released changelog
heading. So for a release the order is:

1. update the changelog and bump the frontmatter version
2. run the evals, which write a report for the new version
3. commit the report with the release commit
4. tag

The release steps in the docs and in the `openkata-ryu-release` skill do
not say this. Without it a release carries only reports for the previous
version, and the index shows no effectiveness.

## Shared code

The index is defined twice today, in `cmd/generate-versions` and in the
web server, and the eval parsing exists only in the generator.

- **`internal/versionindex`**: the types of `versions.json`, used by the
  generator and the web server, so both agree on the shape.
- **`internal/evalresults`**: pure functions that take the report files
  of a skill, as paths and contents, and return the effectiveness by
  model for a given version. The generator's two modes only differ in how
  they collect the files: from disk, or from the bucket. The selection
  rules (newest report per model, matching `skill_version`, rounding)
  live in one place.

## The web server

- **A bucket interface.** The server calls a concrete S3 client in 12
  places. Define a small interface with `GetObject` and `ListObjectsV2`,
  keep the real client for production, and use a fake in tests.
- **Benchmarks.** The production artifact loader reads `tessl_score`,
  `published` and `models` from the index entry, the way the local loader
  does, through one shared function.
- **Version list.** The production loader lists the version folders under
  `<type>/<name>/` and sorts them newest first with the existing
  `compareVersions`, instead of showing only the current version.

**Decided: benchmarks only for the latest version.** The index holds one
entry per artifact, the latest version. The local loader shows that
entry's scores on any version's page, which breaks the rule that a result
of another version is never shown. In production the Benchmarks tab and
the quality score appear only when the page is for the version in the
index, and are hidden on older versions.

## The sandbox model

The sandbox starts `kiro-cli` without `--model`, so a result filed under
a model may have run with another.

**Decided: leave it.** The sandbox is going to be removed, so fixing the
model flag is not worth doing. Until it is removed, a result from a
sandboxed scenario may have run with a different model than its label,
and the evals reference page says so.

## Checking published versions

A script compares the pushed tags with the bucket, using the raw-file
route's behavior as the check, and lists tags that were never pushed and
versions that are missing, so neither goes unnoticed again.

**Decided: add the script, and run it by hand.** A scheduled workflow
would need more bucket access than the CI role has.

## Not in this spec

- Making the catalog search read the live index. It reads the copy
  bundled in the binary, which lags.
- Computing scores. Evals and reviews stay a local step.
