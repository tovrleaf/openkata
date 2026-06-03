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
  version: "1.0.2"
  tags: "category:planning, category:workflow"
---

# Spec-Driven Dev

Drive feature work from idea to implementation through
repo-stored specs. One skill, five phases: specify → design
→ tasks → implement → validate.

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

## Phase 1: Specify

Before asking about requirements, complete these in order:

1. **Determine depth** — Ask the user:
   - **Quick** — small fix or change. Write a `brief.md`
     (a few sentences capturing scope) then skip to
     implementation. No full spec needed.
   - **Standard** — typical feature. `spec.md` + `tasks.md`.
   - **Deep** — complex feature. `spec.md` + `design.md` +
     `tasks.md`.

2. **Create the directory** — Next sequential 4-digit number:
   `specs/NNNN-slug/`. Check existing directories.

3. **Ask about branching** — "Would you like me to create a
   feature branch? Suggested name:
   `feature/NNNN-slug`"
   Do not create the branch without confirmation. If the
   branch already exists, switch to it instead of creating.

Then collect requirements:

4. **Investigate relevant code** — Read files related to
   the feature before asking questions. Check existing
   patterns, dependencies, and constraints. Don't ask the
   user what you can look up yourself.

5. **Ask targeted questions:**
   - What does the user want to build?
   - What does success look like?
   - What is out of scope?
   - Are there open questions that block implementation?

6. **Write spec.md** — Use the spec.md template from
   [spec-templates](references/spec-templates.md). Set
   status to `Draft` and depth to the chosen level.

7. **Set active** — Write the directory name to
   `specs/_current`.

8. **Resolve open questions** — If the spec has an
   Open Questions section with unresolved items, walk
   through each one with the user and get a decision.
   Update the spec to reflect the answers (move decisions
   into Requirements or remove the question). Do not
   proceed with open questions remaining.

9. **Confirm** — Show the spec to the user. Do not proceed
   to the next phase without confirmation.

## Phase 2: Design (Deep only)

1. **Investigate the codebase** — Read relevant files,
   existing patterns, dependencies.

2. **Write design.md** — Use the design.md template from
   [spec-templates](references/spec-templates.md).

3. **Confirm** — Show the design to the user. Do not proceed
   without confirmation.

## Phase 3: Tasks

1. **Break down the spec** (and design if present) into
   ordered tasks. Each task should be completable in one
   commit and touch one concern. If you need multiple
   commits, split the task.

2. **Write tasks.md** — Use the tasks.md template from
   [spec-templates](references/spec-templates.md).

3. **Confirm** — Show the task breakdown. Do not start
   implementing without confirmation.

## Phase 4: Implement

For each task with Status: Pending:

1. Update task status to `In Progress`
2. Read the progress log for notes from earlier tasks
3. Read the task's goal, boundary, and done-when criteria
4. Build the implementation
5. Run tests and verify against done-when criteria. If
   tests fail, fix before proceeding. For non-testable
   criteria (visual, behavioral), describe what was
   checked and the observed result in the progress log.
6. Update task status to `Done`
7. Add a progress log entry with date and summary
8. Commit:
   ```text
   type(scope): description

   Part of specs/NNNN-slug task N.
   ```
9. Move to next Pending task

When all tasks are Done:
- Ask the user: "All tasks complete. Would you like to
  run validation (recommended for a fresh-perspective
  review)?"
- If yes → proceed to Phase 5 (Validate)
- If no → continue below

Mark complete:
- Update `spec.md` status to `Done`
- Clear `specs/_current`
- Inform the user the feature is complete and ask whether
  to push the branch and open a PR

## Phase 5: Validate

A fresh agent session is recommended for validation to
avoid builder bias.

1. **Read spec.md** — load requirements and acceptance
   criteria only. Do not read tasks.md or the progress log.

2. **Review the implementation** — read the code changes
   (use `git diff main` or the feature branch diff).

3. **Write validation-report.md** — Use the template from
   [spec-templates](references/spec-templates.md):
   - Check each requirement against the implementation
   - Verify out-of-scope items were not built
   - Note any issues found

4. **Report** — show the validation report to the user.
   If issues are found, the feature returns to Phase 4
   for fixes.

## Conventions

- One spec directory per feature
- 4-digit zero-padded numbering (0001, 0002, ...)
- Commits reference the spec and task number
- Commits are frequent — after each task, not at the end
- The agent works autonomously within a task but does not
  skip phase transitions without user confirmation

## Boundaries

- DOES create and update files in `specs/`
- DOES create feature branches when asked
- DOES commit during implementation phase
- Does NOT modify files outside the current task's boundary
- Does NOT skip phase gates without user confirmation
- Does NOT delete or overwrite existing specs

## Example Scenario

User: "Let's build the design system for the website"

Agent detects spec-worthy work (new dependency, multiple
concerns), asks depth → Deep. Follows phases 1–5: creates
`specs/0001-design-system/`, collects requirements, writes
design, breaks into tasks, implements with commits, then
validates in a fresh session.

## Common Failures

- **Skipping confirmation** — proceeding to implementation
  without user sign-off on the spec or tasks leads to
  building the wrong thing.
- **Tasks too large** — each task should be completable in
  one commit. If a task needs multiple commits, split it.
- **No progress log** — without log entries, the next session
  can't tell what happened. Always log after each task.
