---
status: PROPOSED
date: 2026-05-03
authors: [niko.kivela]
---

# 0008. Feature-driven development with repo-stored specs

## Context

As the project grows beyond skills and rules into a website,
MCP server, and design system, we need a structured way to
plan and execute feature work. Without it, features are built
ad-hoc — requirements live in chat history, design decisions
are implicit, and there is no record of what was planned vs
what was built.

The workflow must support AI-assisted development: an agent
should be able to read the current feature state, understand
what phase it is in, and continue working autonomously with
frequent commits.

## Decision Drivers

- Feature specs must survive across sessions — chat history
  does not
- An agent must be able to resume work without re-explaining
  the feature
- The process must adapt to feature size — a bug fix should
  not require a full spec cycle
- No external tools or issue trackers — everything in the repo
- Separate files per phase to minimize context window usage
  during implementation

## Decision

Adopt a feature-driven workflow where each feature is a
numbered directory in `.specs/` containing separate files
for each phase. The workflow progresses through phases:
specify → design → tasks → implement.

### Directory structure

```text
.specs/
├── _current                    # Active feature directory name
└── 001-design-system/
    ├── spec.md                 # Story, requirements, acceptance criteria
    ├── design.md               # Architecture, file paths, approach
    └── tasks.md                # Ordered task breakdown
```

### Depth adaptation

The depth determines which files are created:

- **Quick** — no spec directory. Small fix or config change.
  Jump straight to implementation.
- **Standard** — `spec.md` + `tasks.md`. Typical feature
  work. Story and requirements flow directly into tasks
  without a separate design phase.
- **Deep** — `spec.md` + `design.md` + `tasks.md`. Complex
  or cross-cutting feature. Full design phase with
  architecture decisions before task breakdown.

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

### 2. Next task
...

## Progress Log
- [Date] Task N: Summary, decisions made
```

### Implementation loop

For each task, the agent works autonomously:

1. Read the current task from `tasks.md`
2. Build the implementation
3. Run tests
4. Verify against the "done when" criteria
5. Commit with a reference to the feature and task
6. Update task status and progress log
7. Move to next task

Commits are frequent — after each task, not at the end of
the feature.

### Branch strategy

One branch per feature: `feature/{number}-{slug}`. Created
when the feature enters implementation.

### "Where am I?" capability

The skill reads `.specs/_current`, opens the active feature
directory, determines the current phase and next action, and
guides the user or continues autonomously.

## Alternatives Considered

### GitHub Issues + project board

- **Pros:** Standard tooling, visibility for collaborators
- **Cons:** Context lives outside the repo, agents cannot
  easily read/write issues, requires network access
- **Rejected because:** the agent needs to read and update
  feature state as files.

### Spec-kit (github/spec-kit)

- **Pros:** Mature ecosystem (92K stars), extensive extensions,
  CLI tooling, same separate-file pattern
- **Cons:** Requires Python + uv, opinionated `.specify/`
  structure, heavy for a small project
- **Rejected because:** adds a language dependency to a Go
  project. The phased workflow and separate-file pattern are
  adopted without the tooling.

### BMAD Method

- **Pros:** Scale-adaptive, 12+ specialized agents, complete
  lifecycle coverage
- **Cons:** Framework-level complexity, requires npm, 34+
  workflows
- **Rejected because:** designed for multi-agent orchestration
  at enterprise scale. The scale-adaptive depth idea is
  adopted without the framework.

### Single file per feature

- **Pros:** All context in one place, no file management
- **Cons:** Gets long for complex features, agent loads full
  story when it only needs the current task, wastes context
  window
- **Rejected because:** separate files let the agent load
  only what it needs per phase. spec-kit and OpenSpec both
  use separate files for this reason.

## Consequences

### Positive

- Feature state survives across sessions and context resets
- Agent loads only the file relevant to the current phase
- Depth adapts to feature size — no overhead for small changes
- Progress log creates an audit trail of decisions
- Follows the most common convention (spec-kit, OpenSpec)
- No external dependencies or tooling required

### Negative

- More files per feature than a single-file approach
- No built-in visualization (no kanban board)
- Requires discipline to keep spec and tasks in sync

### Neutral

- The workflow skill is local to this project initially but
  could become distributable if the pattern proves useful
- Spec files are committed to the repo and visible in
  git history

## References

- [spec-kit](https://github.com/github/spec-kit) — phased
  specify → plan → tasks → implement workflow, separate files
- [BMAD Method](https://github.com/bmad-code-org/BMAD-METHOD)
  — scale-adaptive depth, "what's next?" guidance
- [OpenSpec](https://openspec.dev/) — generates proposal.md +
  design.md + tasks.md
- [spec-engineer](https://github.com/villetakanen/asdlc-io)
  — spec-anchored development, same-commit rule
- [user-story-clarifier](https://github.com/n-n-code/n-n-code-skills)
  — story card format, Definition of Ready checklist
