---
name: spec-driven-dev
description: >
  Drives feature development through a phased workflow: specify
  requirements, design architecture, break into tasks, then
  implement autonomously with frequent commits. Use when the
  user wants to build a feature, start a new feature, plan
  and implement something, spec out a feature, or says
  "let's spec this" or "new feature."
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

## Phase 1: Specify

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

4. **Collect requirements** — Ask targeted questions:
   - What does the user want to build?
   - What does success look like?
   - What is out of scope?
   - Are there open questions that block implementation?

4. **Write spec.md** — Use the spec.md template from
   [spec-templates](references/spec-templates.md). Set
   status to `Draft` and depth to the chosen level.

5. **Set active** — Write the directory name to
   `specs/_current`.

6. **Confirm** — Show the spec to the user. Do not proceed
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
   ordered tasks.

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
5. Run tests
5. Verify against done-when criteria
6. Update task status to `Done`
7. Add a progress log entry with date and summary
8. Commit:
   ```text
   type(scope): description

   Part of specs/NNNN-slug task N.
   ```
9. Move to next Pending task

When all tasks are Done:
- Update `spec.md` status to `Done`
- Clear `specs/_current`
- Inform the user the feature is complete

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

1. Skill checks `specs/_current` — empty, new feature
2. Asks depth — user says Deep
3. Creates `specs/0001-design-system/`
4. Asks about branch — user says yes
5. Collects requirements, writes `spec.md`, confirms
6. Investigates codebase, writes `design.md`, confirms
7. Breaks into tasks, writes `tasks.md`, confirms
8. Implements task by task with commits
9. User starts a fresh agent session for validation
10. Validator reads spec.md and code diff, writes report
11. Clears `_current` when validated

## Common Failures

- **Skipping confirmation** — proceeding to implementation
  without user sign-off on the spec or tasks leads to
  building the wrong thing.
- **Tasks too large** — each task should be completable in
  one commit. If a task needs multiple commits, split it.
- **No progress log** — without log entries, the next session
  can't tell what happened. Always log after each task.
