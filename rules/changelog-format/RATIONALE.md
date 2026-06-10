# Rationale

changelog-format enforces Keep a Changelog structure
across all project changelogs.

## Why a strict preamble is required

The preamble identifies the format to new contributors
immediately. Without it, changelogs drift into custom
formats that tooling can't parse.

## Why [Unreleased] is always present

Changes accumulate between releases. Without a staging
section, developers either forget to log changes or
add them directly under the last version heading —
both produce incorrect history.

## Why entries describe what, not how

Changelogs serve future readers asking "what changed."
Implementation details belong in commits and PRs. A
changelog entry that says "refactored the parser" tells
the reader nothing useful.

## Why no empty sections

Empty headings (### Fixed with nothing under it) add
noise. Only include sections that have entries. This
keeps changelogs scannable.
