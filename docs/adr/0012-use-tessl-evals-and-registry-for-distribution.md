---
status: PROPOSED
date: 2026-05-18
authors: [niko.kivela]
---

# 0012. Use Tessl evals and registry for distribution

## Context

ADR 0006 adopted Tessl for skill quality (lint, review, optimize)
but explicitly excluded publishing and distribution. Since then,
the project has grown to 8 distributable skills and the dojo MCP
server is no longer the only distribution path.

Tessl now offers scenario-based evals that measure whether a skill
actually improves agent behavior — not just whether it's well-
structured. The registry provides discoverability, security review,
and an Impact score that differentiates skills for potential users.

Without evals, we can only verify that skills read well. With evals,
we can prove they work. Without the registry, users must find skills
via GitHub or the dojo. With the registry, skills are searchable and
versioned with quality signals.

## Decision Drivers

- Must prove skills work, not just that they're well-written
- Must catch regressions when skill content changes
- Must get automated security review on published skills
- Must provide discoverability beyond the dojo and GitHub
- Evals are server-side only — cost and time are factors
- Manual publishing workflow preferred over CI automation

## Decision

We will use Tessl evals and the Tessl registry for distributable
skills. Specifically:

- **Evals:** Generate scenarios with `tessl scenario generate` and
  run them with `tessl eval run` before publishing. Evals are
  committed to the repo in `skills/<name>/evals/` but excluded
  from the published tile via `.tesslignore`.
- **Publishing:** Publish distributable skills to the `openkata`
  workspace on the Tessl registry using `tessl tile publish`.
  Publishing is manual — triggered during the release workflow,
  not automated via CI.
- **When to run evals:** Only for distributable skills, only when
  skill content has changed meaningfully, only before publishing.
  Not on every commit.
- **Versioning:** Skill versions (CHANGELOG.md) and tile versions
  (tile.json) remain independent per ADR 0006.

This supersedes the "no publishing" clause in ADR 0006. The dojo
MCP server remains an alternative distribution channel.

## Alternatives Considered

### Continue without evals or registry

- **Pros:** No server-side dependency, no cost, simpler workflow
- **Cons:** No proof skills work, no regression detection, limited
  discoverability, no security review
- **Rejected because:** As skills are published publicly, quality
  signals and discoverability matter for adoption.

### Automate publishing via CI (GitHub Actions)

- **Pros:** Every merge publishes automatically, no manual step
- **Cons:** Publishes on every change regardless of readiness,
  evals run on every push (cost), less control over what ships
- **Rejected because:** Skills are updated infrequently and
  batched intentionally. Manual publishing gives full control
  over timing and quality gates.

### Run evals for all skills including local

- **Pros:** Consistent quality bar across everything
- **Cons:** Local skills are repo-bound, not distributed, and
  change frequently during development. Eval cost and time not
  justified for internal-only artifacts.
- **Rejected because:** Evals prove value for consumers. Local
  skills have no external consumers.

## Consequences

### Positive

- Regression detection when skills are edited
- Impact score proves skills work for potential users
- Security review on publish catches issues automatically
- Discoverability via the Tessl registry

### Negative

- Server-side eval dependency — can't run locally
- Eval runs take ~15 minutes per skill
- Publishing adds a step to the release workflow

### Neutral

- Evals are committed to the repo but excluded from tiles
- The dojo remains available as an alternative install path
- Tile versions stay independent from skill versions

## References

- [Tessl eval documentation](https://docs.tessl.io/evaluate/evaluate-skill-quality-using-scenarios)
- [Tessl registry](https://tessl.io/registry)
- ADR 0006 — Adopt Tessl as skill quality toolchain
- ADR 0005 — Version distributed artifacts
