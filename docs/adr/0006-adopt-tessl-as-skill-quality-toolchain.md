---
status: PROPOSED
date: 2026-04-25
authors: [niko.kivela]
---

# 0006. Adopt Tessl as the skill quality toolchain

## Context

OpenKata skills follow the Agent Skills format (ADR 0001) and are
distributed via the dojo MCP server (ADR 0002). However, the project
lacks tooling for evaluating and improving skill quality — validating
structure, reviewing content against best practices, and optimizing
instructions for agent consumption.

Currently, skill quality is verified manually by reading SKILL.md files.
There is no automated validation that a skill's structure is correct
and no standardized review process. As the number of kata grows and
external contributors join, this manual approach won't scale.

Tessl is a CLI toolchain purpose-built for the Agent Skills format. It
provides `skill lint`, `skill review`, and `skill review --optimize`
commands that directly address these gaps.

## Decision Drivers

- Must validate skill structure automatically before merge
- Must support review and optimization of skill content
- Must align with the Agent Skills format already adopted (ADR 0001)
- Must not interfere with the dojo as the sole distribution mechanism
- Low adoption cost — should work with existing SKILL.md files

## Decision

We will adopt Tessl as the standard toolchain for skill quality in
OpenKata. Specifically:

- `tessl skill lint` validates skill structure during development and CI
- `tessl skill review` evaluates skill quality and compliance
- `tessl skill review --optimize` suggests improvements to skill content

Each skill directory will contain a `tile.json` alongside its SKILL.md,
generated via `tessl skill import`. This file is required by Tessl for
its operations. The `tile.json` is a development-only artifact — the
dojo does not copy it to users when installing kata.

The `tessl.json` at the project root manages workspace configuration.

We will not use Tessl for publishing or distribution. The dojo MCP
server (ADR 0002) remains the sole distribution channel for kata.

## Alternatives Considered

### Manual validation only

- **Pros:** No external dependency, no tooling to learn
- **Cons:** No automated structure validation, quality depends entirely
  on reviewer diligence, no way to systematically optimize skill content
- **Rejected because:** Doesn't scale as the number of kata and
  contributors grows. Provides no feedback loop for improving skill
  quality.

### Custom linting and review scripts

- **Pros:** Full control, no external dependency
- **Cons:** Must build and maintain SKILL.md parsing, frontmatter
  validation, structure checks, and quality heuristics from scratch.
  Duplicates what Tessl already provides.
- **Rejected because:** Reinventing tooling that already exists for the
  exact format we use. Maintenance burden with no quality benefit.

## Consequences

### Positive

- Automated validation catches structural issues before merge
- Review and optimization provide a consistent quality bar for kata
- Works with existing SKILL.md files via `tessl skill import`
- Development-only tooling — no impact on how users consume kata

### Negative

- Adds an external CLI dependency to the development workflow
- Contributors need Tessl installed locally (or CI handles it)
- Each skill directory gains a `tile.json` file that duplicates some
  metadata from SKILL.md frontmatter (name, version, description)

### Neutral

- `tile.json` is a Tessl requirement, not distributed to users — the
  dojo skips it when installing kata into projects
- Tessl's vendored mode means no external dependencies are pulled into
  the project

## Versioning

Each skill has two version fields:

| Field | File | Managed by | Purpose |
|-------|------|-----------|---------|
| `version` in frontmatter | SKILL.md | `release` | Source of truth. Follows ADR 0005 semver scheme. |
| `version` in tile.json | tile.json | `tessl skill publish` | Tessl-internal. Defaults to `0.1.0` on import. |

The SKILL.md version is what users see and what the dojo uses. The
tile.json version is only relevant if publishing to the Tessl registry,
which we don't do (see Decision above). Tessl does not read or sync
the SKILL.md frontmatter version.

The tile.json version will remain at its default unless manually changed.
This is acceptable — tile.json is a development-only file not distributed
to users.

## Reversibility

Tessl adoption is reversible. The core artifacts are SKILL.md files
(ADR 0001), which exist independently of Tessl. Removing Tessl would
mean deleting `tile.json` files, `tessl.json`, and the `.tessl/`
directory. The skills themselves would remain fully functional.

Reconsider if Tessl development stalls, the Agent Skills ecosystem
converges on a different toolchain, or a better quality tool emerges.

## References

- [Tessl CLI documentation](https://docs.tessl.io)
- [Agent Skills specification](https://agentskills.io/specification)
- ADR 0001 — Use Agent Skills format
- ADR 0002 — Use MCP server for distribution
