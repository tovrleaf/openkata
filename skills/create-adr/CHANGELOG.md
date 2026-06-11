# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.2.1] - 2026-06-11

### Added

- RATIONALE.md explaining design decisions

## [1.2.0] - 2026-05-24

### Added

- Workflow step 8: offer to run grill-me skill on finalized
  ADRs to find blind spots

## [1.1.0] - 2026-04-30

### Added

- Example scenario showing end-to-end ADR creation flow
- Gotcha: never leave placeholder text in sections
- references/ACKNOWLEDGMENTS.md for source attribution

### Changed

- Moved References to references/ACKNOWLEDGMENTS.md
- Merged Common Failures into Gotchas to reduce overlap with
  E.C.A.D.R. quality checklist

### Removed

- "When to suggest" section (duplicated the description)
- Inline References section (moved to ACKNOWLEDGMENTS.md)

## [1.0.2] - 2026-04-26

### Fixed

- Applied markdown-consistency rule: added language specifier to
  code block, blank lines around headings in changelog

## [1.0.1] - 2026-04-25

### Changed

- Trimmed intro paragraphs redundant with description
- Compacted "When to suggest" list for conciseness

### Fixed

- CHANGELOG.md formatting (added missing category header)

## [1.0.0] - 2026-04-12

### Added

- Customized Nygard ADR template with YAML frontmatter, decision
  drivers, structured alternatives, and categorized consequences
- Optional Non-goals, Reversibility, and References sections
- Workflow: codebase exploration before asking questions
- Workflow: recommend answers with justification
- Workflow: one-at-a-time mode for first ADR in session
- E.C.A.D.R. quality self-check with INVESTIGATE markers
- ADR lifecycle: PROPOSED → ACCEPTED → SUPERSEDED/DEPRECATED
