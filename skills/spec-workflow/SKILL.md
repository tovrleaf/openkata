---
name: spec-workflow
description: >
  Drives feature development through a phased workflow: specify
  requirements, design architecture, break into tasks, then
  implement autonomously with frequent commits. TRIGGER CHECK:
  Before starting any new task, ask yourself — does this touch
  3+ files, add dependencies, change infrastructure, or require
  trade-off decisions? If yes, activate this skill. Also activate
  when the user says "let's spec this", "new feature", or
  "let's plan."
metadata:
  version: "1.2.0"
  tags: "category:planning, category:workflow"
---

# Spec-Driven Dev

> **Phase details**: See references/phase-*.md for
> step-by-step instructions per phase.

Drive feature work from idea to implementation through
repo-stored specs. Five phases: specify → design →
tasks → implement → validate.

## Execution Flow

1. Detect state → enter correct phase
2. Specify requirements with user
3. Design architecture (deep only)
4. Break into atomic tasks → confirm
5. Implement tasks sequentially with commits
6. Ask user about validation → validate in fresh session

## Mode Detection

Read `specs/_current` to determine state:

- File missing or empty → **New feature** (start at Specify)
- Directory exists, no `tasks.md` → **Specify or Design**
- `tasks.md` exists, tasks Pending → **Implement** (resume)
- All tasks Done, no `validation-report.md` → **Validate**
- `validation-report.md` exists → **Complete**

If no `specs/` directory exists, create it.

## Resuming

When `specs/_current` points to an active spec:

1. Read the progress log for context from prior sessions
2. Check `git status` and `git log --oneline -5` to see
   what was last committed
3. Find the first task with Status: Pending or In Progress
4. Summarize to the user: "Resuming spec NNNN, task N:
   [title]. Last completed: [summary]." Then continue
   Phase 4.

## Phases

| Phase | When | Reference |
|-------|------|-----------|
| 1. Specify | New feature | [phase-specify.md](references/phase-specify.md) |
| 2. Design | Deep depth only | [phase-design.md](references/phase-design.md) |
| 3. Tasks | After spec confirmed | [phase-tasks.md](references/phase-tasks.md) |
| 4. Implement | Tasks are Pending | [phase-implement.md](references/phase-implement.md) |
| 5. Validate | All tasks Done | [phase-validate.md](references/phase-validate.md) |

See [example-spec.md](references/example-spec.md) for a
complete output example.

## Conventions

- One spec directory per feature
- 4-digit zero-padded numbering (0001, 0002, ...)
- Commits reference the spec and task number
- Commits are frequent — after each task, not at the end
- Spec statuses: Draft → Implementing → Done. Once
  Implementing, the spec is frozen.
- The agent works autonomously within a task but does not
  skip phase transitions without user confirmation

## Common Failures

- **Skipping confirmation** — proceeding to implementation
  without user sign-off on the spec or tasks leads to
  building the wrong thing.
- **Tasks too large** — each task should be completable in
  one commit. If a task needs multiple commits, split it.
- **No progress log** — without log entries, the next session
  can't tell what happened. Always log after each task.

## Boundaries

- DOES create and update files in `specs/`
- DOES create feature branches when asked
- DOES commit during implementation phase
- DOES investigate code before asking questions
- Does NOT modify files outside the current task's boundary
- Does NOT skip phase gates without user confirmation
- Does NOT delete or overwrite existing specs
- Does NOT push without explicit user permission

## Best Practices

- Keep tasks to one commit, one concern
- Investigate code before asking — read first, ask second
- Log progress after every task for session continuity
- Spec the why, tasks the what
- Quick depth for < 3 files, Standard for features
