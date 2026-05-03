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
- One file per feature to keep context together and reduce
  file management overhead

## Decision

Adopt a feature-driven workflow where each feature is a single
markdown file in `.features/` that progresses through phases:
story → design → tasks → implementation.

### Feature file structure

```text
.features/
├── _current              # Filename of active feature
└── 001-design-system.md  # One file per feature
```

### Feature file format

```markdown
---
status: Specifying | Designing | In Progress | Done
depth: Quick | Standard | Deep
---

# Feature Title

## Story
As a ..., I want ..., so that ...

## Requirements
- Collected from conversation with user
- Acceptance criteria

## Design
- Architecture, file paths, approach
- Out of scope

## Tasks
### 1. Task name
- **Status**: Pending | In Progress | Done
- **Goal**: What this achieves
- **Done when**: Observable criteria

## Progress Log
- [Date] Task N: Summary, decisions made
```

### Depth adaptation

The skill determines depth from the feature's scope:

- **Quick** — small fix or config change. No story or design
  phase. Jump straight to implementation with a commit
  message that explains the change.
- **Standard** — typical feature. Story → design → tasks →
  implement. Most work falls here.
- **Deep** — complex or cross-cutting feature. Story →
  detailed requirements gathering → design with ADR
  reference → tasks → implement with verification steps.

### Implementation loop

For each task, the agent works autonomously:

1. Read the task spec
2. Build the implementation
3. Run tests
4. Verify against the "done when" criteria
5. Commit with a reference to the feature and task
6. Update task status and progress log
7. Move to next task

Commits are frequent — after each task, not at the end of
the feature. This provides checkpoints for debugging and
review.

### Branch strategy

One branch per feature: `feature/{number}-{slug}`. Created
when the feature enters the "In Progress" phase.

### "Where am I?" capability

The skill can read `.features/_current`, open the active
feature file, determine the current phase and next action,
and guide the user or continue autonomously.

## Alternatives Considered

### GitHub Issues + project board

- **Pros:** Standard tooling, visibility for collaborators,
  integrates with PRs
- **Cons:** Context lives outside the repo, agents cannot
  easily read/write issues, specs and tasks are split across
  systems, requires network access
- **Rejected because:** the agent needs to read and update
  feature state as files. Issues are a UI, not a development
  artifact.

### Spec-kit (github/spec-kit)

- **Pros:** Mature ecosystem, 90K+ stars, extensive extension
  catalog, CLI tooling
- **Cons:** Requires Python + uv, opinionated directory
  structure (`.specify/`), heavy for a small project, not
  aligned with Agent Skills format
- **Rejected because:** adds a language dependency (Python)
  to a Go project and imposes a framework where a single
  skill suffices. The phased workflow idea is adopted without
  the tooling.

### BMAD Method

- **Pros:** Scale-adaptive, 12+ specialized agents, complete
  lifecycle coverage
- **Cons:** Framework-level complexity, requires npm, 34+
  workflows is far more than needed, specialized agent
  personas don't map to single-agent usage
- **Rejected because:** designed for multi-agent orchestration
  at enterprise scale. The scale-adaptive depth idea is
  adopted without the framework.

### Separate spec and task files

- **Pros:** Clean separation of concerns, each file has one
  purpose
- **Cons:** More files to manage, context split across files,
  agent must read multiple files to understand state
- **Rejected because:** a single file keeps all context
  together. The agent reads one file and knows everything
  about the feature.

## Consequences

### Positive

- Feature state survives across sessions and context resets
- Agent can resume autonomously by reading one file
- Depth adapts to feature size — no overhead for small changes
- Progress log creates an audit trail of decisions
- No external dependencies or tooling required
- Frequent commits provide rollback points

### Negative

- Feature files can grow long for complex features
- No built-in visualization (no kanban board, no burndown)
- Single-file format may need restructuring if features
  become very large

### Neutral

- The workflow skill is local to this project initially but
  could become distributable if the pattern proves useful
- Feature files are committed to the repo and visible in
  git history

## References

- [spec-kit](https://github.com/github/spec-kit) — phased
  specify → plan → tasks → implement workflow
- [BMAD Method](https://github.com/bmad-code-org/BMAD-METHOD)
  — scale-adaptive depth, "what's next?" guidance
- [spec-engineer](https://github.com/villetakanen/asdlc-io)
  — spec-anchored development, same-commit rule
- [user-story-clarifier](https://github.com/n-n-code/n-n-code-skills)
  — story card format, Definition of Ready checklist
