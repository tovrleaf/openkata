# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.1.0] - 2026-06-03

### Added

- Execution Flow summary for quick-glance overview
- Best Practices section with positive guidance
- Concrete example spec in references/example-spec.md

### Changed

- Split phase details into individual reference files for
  reduced token cost (phase-specify, phase-design,
  phase-tasks, phase-implement, phase-validate)
- SKILL.md reduced from ~160 to ~95 lines

## [1.0.2] - 2026-06-03

### Added

- Explicit step to resolve open questions before confirming spec

## [1.0.1] - 2026-06-02

### Fixed

- Prompt user for validation before clearing spec as done
  (Phase 5 was silently skipped in practice)

## [1.0.0] - 2026-05-24

### Added

- Five-phase workflow: specify → design → tasks → implement
  → validate
- Depth adaptation: Quick (brief.md), Standard (spec +
  tasks), Deep (spec + design + tasks)
- Mode detection from file state for session resumability
- Validation phase with fresh-agent recommendation
- Boundary and Depends annotations on tasks
- Progress log forward-reading between tasks
- References: spec-templates.md with all file templates
