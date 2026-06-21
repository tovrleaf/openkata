---
status: Draft
depth: Standard
---

# Quality Lint

## Story

As the site maintainer, I want a static quality gate
that blocks publishing structurally incomplete skills
so every published skill meets Open Kata's conventions.

## Requirements

### Lint Tool

- New binary or subcommand: `cmd/openkata-lint/` or
  `openkata-eval lint`
- Static analysis of SKILL.md — no LLM calls
- Pass/fail per criterion, overall pass/fail
- Exit code 0 = pass, 1 = fail

### Criteria (initial set)

- Has YAML frontmatter with `name` and `description`
- Description is non-empty
- Has at least one `#` heading
- Has execution flow / workflow section
- Has boundaries (does/does not) section
- Under 500 lines
- No broken relative links (referenced files exist)
- Has metadata with `version` and `tags`

### Gate Invocation Points

1. **CI on every push** — runs lint on all changed
   skills; fails the check if any don't pass
2. **Publish CI job** — runs lint before S3 upload;
   blocks publish on failure
3. **`publish-tile.sh`** — runs lint before tessl
   publish; exits early on failure

### Output

- Terminal: list of criteria with ✓/✗ per item
- Exit code for CI integration
- Optional `--json` flag for machine-readable output

### Not Displayed to Users

- All published skills pass by definition
- No score badge on the site
- The lint is an internal guardrail only

## Constraints

- No network or LLM dependencies
- Criteria owned and versioned in this repo
- Must run fast (< 1s per skill)
- Replaces reliance on `tessl skill review` for
  quality enforcement (tessl review/optimize remain
  available as developer tools)

## Acceptance Criteria

1. `openkata-lint skills/create-adr` passes for all
   current published skills
2. Removing a required section causes lint failure
3. CI job runs lint on push
4. Publish CI blocks upload on lint failure
5. `publish-tile.sh` exits early on lint failure
6. `--json` output parseable

## Out of Scope

- LLM-based optimization suggestions (use tessl for
  that)
- Linting rules or profiles
- Displaying lint results to site users
- Auto-fixing lint failures

## Open Questions

None.

Date: 2026-06-20
