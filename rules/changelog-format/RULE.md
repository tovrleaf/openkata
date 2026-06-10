---
name: changelog-format
description: >
  Keep a Changelog format for all CHANGELOG.md files.
  Applied when creating or updating any changelog.
tags: category:conventions, language:markdown
---

# Changelog Format

Every CHANGELOG.md follows Keep a Changelog strictly.

## Structure

```markdown
# Changelog

All notable changes to this project will be documented
in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - 2026-05-19

### Added

- Initial release
```

## Rules

- Start with `# Changelog` heading
- Include the two-line preamble (notable changes + format)
- `## [Unreleased]` section always present, above versions
- Versions in descending order (newest first)
- Version format: `## [MAJOR.MINOR.PATCH] - YYYY-MM-DD`
- Use only these section headings under each version:
  Added, Changed, Deprecated, Removed, Fixed, Security
- Each entry starts with `- ` (dash, space)
- Entries describe what changed, not how or why
- No empty sections — omit headings with no entries
- Blank line between sections

## Entry Quality

- Bad: "Fixed bug"
- Good: "Fixed crash when uploading files larger than 10MB"
- Bad: "Updated dependencies"
- Good: "Upgraded goldmark to v1.7 for CommonMark compliance"

## Versioning

- Breaking changes → MAJOR bump
- New features → MINOR bump
- Bug fixes → PATCH bump
- First release is always 1.0.0
