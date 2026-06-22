---
status: Draft
depth: Standard
---

# Public Tessl Skills

## Story

As the site maintainer, I want to publish skills to the
tessl registry so they're discoverable beyond openkata.dev
and link to the registry from the detail page.

## Requirements

### Publishing Workflow

- Each skill has `.tessl-plugin/plugin.json` with
  `"private": true`
- Publishing = set `"private": false` then run
  `tessl tile publish`
- Once public, stays public (no flip-back)
- New script: `scripts/publish-tile.sh` accepts skill
  name, sets private to false if needed, runs publish
- Makefile target: `make publish-tile SKILL=name`
- Tessl publish is part of the standard release
  workflow — run alongside tag + S3 publish
- Publishing remains manual (per ADR 0012), not in CI

### Registry Link on Detail Page

- Skill detail page shows "View on Tessl Registry"
  link inside the Benchmarks tab (not in header)
- Registry URL: `https://tessl.io/registry/openkata/<name>`
- `generate-versions` reads `plugin.json` `"private"`
  field → writes `"published": true/false` to
  `versions.json`
- Unpublished skills: no link shown

### Scope

- Skills only (rules and profiles not published to
  tessl)
- All distributable skills will eventually be public

## Constraints

- Tessl remains a developer tool and distribution
  channel — not a data source for the site
- `plugin.json` is tessl's file — only modify the
  `"private"` field
- `.tessl-plugin/` excluded from user installations
- All `.tesslignore` files must include `evals/` before
  publishing (already fixed across all skills)

## Acceptance Criteria

1. `make publish-tile SKILL=create-adr` sets private
   to false and runs `tessl tile publish`
2. `generate-versions` reads published status from
   `plugin.json` into `versions.json`
3. Skill detail page shows registry link for published
   skills
4. Unpublished skills show no registry link
5. Script is idempotent (re-running on already-public
   skill just republishes)

## Out of Scope

- Automatic publishing via CI
- Publishing rules or profiles
- Pulling data from tessl registry into the site
- Unpublishing skills

## Open Questions

None.

Date: 2026-06-20
