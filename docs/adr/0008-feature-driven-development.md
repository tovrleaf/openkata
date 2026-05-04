---
status: PROPOSED
date: 2026-05-03
authors: [niko.kivela]
---

# 0008. Repo-stored feature specs

## Context

As the project grows beyond skills and rules into a website,
MCP server, and design system, feature requirements and design
decisions need a persistent home. Chat history doesn't survive
across sessions. Without stored specs, every new session starts
from zero.

## Decision Drivers

- Feature specs must survive across sessions
- An agent must be able to read the current feature state and
  resume without re-explaining
- Separate files per phase to minimize context window usage
- No external tools or issue trackers — everything in the repo
- The format must adapt to feature size

## Decision

Store feature specs in `.specs/` as numbered directories with
separate files per phase.

### Directory structure

```text
.specs/
├── _current                      # Active feature directory name
└── 0001-design-system/
    ├── spec.md                   # Story, requirements, acceptance
    ├── design.md                 # Architecture, approach (Deep)
    └── tasks.md                  # Ordered task breakdown
```

### Depth adaptation

The depth determines which files are created:

- **Quick** — no spec directory. Small fix or config change.
- **Standard** — `spec.md` + `tasks.md`. Typical feature work.
- **Deep** — `spec.md` + `design.md` + `tasks.md`. Complex or
  cross-cutting feature needing a design phase.

### File formats

**spec.md** — what and why:

```markdown
---
status: Draft | Ready | In Progress | Done
depth: Quick | Standard | Deep
---

# Feature Title

## Story
As a ..., I want ..., so that ...

## Requirements
- Acceptance criteria
- Out of scope

## Open Questions
- Unresolved items that block or change implementation
```

**design.md** — how (Deep only):

```markdown
# Design: Feature Title

## Architecture
- Approach, file paths, components affected

## Decisions
- Key trade-offs and rationale

## Dependencies
- External systems, libraries, other features
```

**tasks.md** — work breakdown:

```markdown
# Tasks: Feature Title

## Tasks

### 1. Task name
- **Status**: Pending | In Progress | Done
- **Goal**: What this achieves
- **Done when**: Observable criteria

## Progress Log
- [Date] Task N: Summary, decisions made
```

### Numbering

Directories use 4-digit zero-padded numbers (`0001-name/`)
matching the ADR convention. Check existing directories to
determine the next number.

## Alternatives Considered

### GitHub Issues + project board

- **Pros:** Standard tooling, visibility for collaborators
- **Cons:** Context lives outside the repo, agents cannot
  easily read/write issues, requires network access
- **Rejected because:** the agent needs to read and update
  feature state as files.

### Spec-kit (github/spec-kit)

- **Pros:** Mature ecosystem (92K stars), same separate-file
  pattern, CLI tooling
- **Cons:** Requires Python + uv, opinionated `.specify/`
  structure, heavy for a small project
- **Rejected because:** adds a language dependency to a Go
  project. The separate-file pattern is adopted without the
  tooling.

### Single file per feature

- **Pros:** All context in one place, no file management
- **Cons:** Gets long for complex features, agent loads full
  story when it only needs the current task
- **Rejected because:** separate files let the agent load
  only what it needs per phase.

## Consequences

### Positive

- Feature state survives across sessions and context resets
- Agent loads only the file relevant to the current phase
- Depth adapts to feature size
- Progress log creates an audit trail
- Follows the most common convention (spec-kit, OpenSpec)

### Negative

- More files per feature than a single-file approach
- No built-in visualization (no kanban board)

### Neutral

- Spec files are committed to the repo and visible in
  git history

## References

- [spec-kit](https://github.com/github/spec-kit) — separate
  spec.md + plan.md + tasks.md per feature
- [OpenSpec](https://openspec.dev/) — proposal.md + design.md
  + tasks.md
