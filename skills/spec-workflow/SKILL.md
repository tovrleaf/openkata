---
name: spec-workflow
description: >
  Drives feature development through a repo-stored, phased spec workflow:
  specify requirements, design architecture, break into atomic tasks, then
  implement autonomously with frequent per-task commits and a progress log
  for session continuity. Use when the user says "let's spec this", "new
  feature", or "let's plan", or when a change touches 3+ files, adds
  dependencies, alters infrastructure, or requires trade-off decisions.
  Distinct from general coding assistance: this skill manages a persistent
  spec directory lifecycle (Draft → Implementing → Done) with explicit
  phase gates and user confirmation checkpoints before each transition.
metadata:
  version: "1.2.1"
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

### Spec Directory Structure

```
specs/
  _current          # plain text: path to active spec dir, e.g. specs/0001-auth
  0001-auth/
    spec.md         # requirements (Draft → Implementing → Done)
    design.md       # architecture notes (deep depth only)
    tasks.md        # task list with Status fields
    progress.log    # append-only session log
    validation-report.md  # created during Phase 5
  0002-payments/
    ...
```

**Sample `specs/_current` content:**
```
specs/0001-auth
```

**Sample `specs/0001-auth/tasks.md` excerpt:**
```
## Task 0001-03 — Add JWT middleware
Status: Pending
Files: src/middleware/auth.ts, src/routes/index.ts
```

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

| Phase | When | Output | Reference |
|-------|------|--------|-----------|
| 1. Specify | New feature | `spec.md` with sections: Problem, Requirements, Constraints, Acceptance Criteria | [phase-specify.md](references/phase-specify.md) |
| 2. Design | Deep depth only | `design.md` with sections: Architecture, Key Components, Data Flow, Trade-offs | [phase-design.md](references/phase-design.md) |
| 3. Tasks | After spec confirmed | `tasks.md` with atomic, commit-sized tasks, each listing Status and Files | [phase-tasks.md](references/phase-tasks.md) |
| 4. Implement | Tasks are Pending | Code commits (one per task) + `progress.log` entries | [phase-implement.md](references/phase-implement.md) |
| 5. Validate | All tasks Done | `validation-report.md` summarising acceptance criteria results | [phase-validate.md](references/phase-validate.md) |

See [example-spec.md](references/example-spec.md) for a
complete output example.

## Rules

**Do:**
- One spec directory per feature; 4-digit zero-padded numbering (0001, 0002, …)
- Commit after each task (not at the end); reference the spec and task number
- Keep tasks to one commit, one concern — split if a task needs multiple commits
- Investigate code before asking — read first, ask second
- Log progress after every task for session continuity
- Spec statuses: Draft → Implementing → Done. Spec is frozen once Implementing.
- Spec the *why*; tasks the *what*. Quick depth for < 3 files, Standard for features.
- Create feature branches when asked; create and update files in `specs/`

**Do not:**
- Skip phase gates or proceed to implementation without user confirmation on spec or tasks
- Modify files outside the current task's boundary
- Delete or overwrite existing specs
- Push without explicit user permission

**Common failures:**
- **Skipping confirmation** — building without sign-off leads to building the wrong thing.
- **Tasks too large** — if a task needs multiple commits, split it.
- **No progress log** — without log entries, the next session can't resume correctly.
