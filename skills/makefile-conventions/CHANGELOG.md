# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.1] - 2026-06-11

### Changed

- Replaced catch-all `%: @:` pattern with conditional
  explicit subcommand declarations throughout

### Added

- Warning against placing `%: @:` catch-all in root Makefile
  (Common Failures, Structure conventions, and reference doc)
- RATIONALE.md explaining design decisions

## [1.0.0] — 2026-05-24

### Added

- SKILL.md with workflow for creating and organizing Make
  targets
- Modular Makefile structure reference with delegation pattern
- Conventions for naming, help system, and script delegation
