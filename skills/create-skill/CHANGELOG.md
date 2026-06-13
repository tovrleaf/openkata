# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.3.1] - 2026-06-13

### Added

- Enforce version 1.0.0 for new skills and Unreleased changelog
  convention

## [1.3.0] - 2026-06-11

### Added

- Portability check in validation step for distributable
  skills
- Common Failures guidance: use NEVER/MUST language for
  weak enforcement

### Changed

- Frontmatter rules now generic (project decides metadata
  structure)
- Acknowledgments step documents compact one-liner format

### Removed

- Version-banned-from-frontmatter rationale (moved to
  project conventions)

## [1.2.1] - 2026-05-26

### Changed

- Streamlined SKILL.md for lower token cost and 95%+ review
  score
- Condensed writing rules, example scenario, and common
  failures sections
- Simplified investigation and design steps

## [1.2.0] - 2026-05-24

### Added

- Imperative form writing rule replacing "explain the why"
- Structured validation report format with trigger analysis
- Skill design checklist: consistency, validation, and token
  budget sections
- Metadata tags in frontmatter

### Changed

- Boundaries section now mandatory (was recommended)
- Conciseness rule tightened to single-sentence form

## [1.1.0] - 2026-05-05

### Added

- Boundaries section as a recommended skill convention
- Boundaries checklist in skill design checklist
- Frontmatter rules: only name and description allowed
- Writing rule: complete examples over scattered snippets
- Writing rule: common failures must be non-obvious

### Changed

- Conciseness rule now explicitly says not to explain
  concepts the model already knows
- Removed version field from frontmatter (ADR 0005)

## [1.0.0] - 2026-04-28

### Added

- Repo investigation step before asking questions
- Progressive disclosure design guidance (SKILL.md, references/,
  scripts/, assets/)
- Trigger description writing rules
- Validation with representative positive and negative prompts
- Quality checklist for finalizing skills
