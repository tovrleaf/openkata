---
status: ACCEPTED
date: 2026-05-04
authors: [niko.kivela]
---

# 0009. Feature development workflow

## Context

ADR 0008 defines where and how feature specs are stored. This
ADR defines the workflow — how the agent and human collaborate
to move a feature from idea to implementation.

The workflow must support autonomous agent execution with
human gates at key decision points.

## Decision Drivers

- Agent should work autonomously within a task but not skip
  phases without human confirmation
- Frequent commits for debugging and rollback
- The workflow must be resumable — an agent in a new session
  should pick up where the last one left off
- Branch-per-feature for clean history

## Decision

Adopt a phased workflow: specify → design → tasks → implement
→ validate. The agent drives each phase, the human confirms
transitions.

### Phases

1. **Specify** — Collect requirements from the user. Write
   `spec.md` with story, acceptance criteria, and out of
   scope. Get user confirmation before proceeding.

2. **Design** (Deep only) — Investigate the codebase, propose
   architecture, write `design.md`. Get user confirmation.

3. **Tasks** — Break the spec (and design if present) into
   ordered tasks in `tasks.md`. Each task has a goal and
   observable "done when" criteria.

4. **Implement** — For each task, autonomously:
   1. Read the task spec
   2. Build the implementation
   3. Run tests
   4. Verify against "done when" criteria
   5. Commit with feature/task reference
   6. Update task status and progress log
   7. Move to next task

5. **Validate** — A fresh agent session is recommended to
   avoid builder bias. Read `spec.md` requirements and the
   code diff. Write `validation-report.md` checking each
   requirement against the implementation. If issues are
   found, return to Implement for fixes.

### Branch strategy

One branch per feature: `feature/{number}-{slug}`. Created
when the feature enters the Implement phase.

### Resumability

The agent reads `specs/_current` to find the active feature,
then reads `spec.md` and `tasks.md` to determine the current
phase and next action:

- No `tasks.md` → still in Specify or Design phase
- `tasks.md` exists, all Pending → ready to start implementing
- `tasks.md` with some Done → resume from first Pending task
- All tasks Done, no `validation-report.md` → ready to validate
- `validation-report.md` exists → feature complete

### Commit conventions

During implementation, commits reference the feature:

```text
feat(website): add templ layout templates

Part of specs/0001-website task 3.
```

Commits are frequent — after each task, not at the end.

### Human gates

The agent does not proceed past these points without user
confirmation:

- spec.md → tasks.md (requirements are complete)
- design.md → tasks.md (architecture is approved, Deep only)
- Before first implementation commit (task plan is approved)

Within the implementation loop, the agent works autonomously.

## Alternatives Considered

### Fully autonomous (no gates)

- **Pros:** Faster, no waiting for human
- **Cons:** Agent may build the wrong thing, no course
  correction until the end
- **Rejected because:** human oversight at phase transitions
  catches misunderstandings early.

### Gate every task

- **Pros:** Maximum control
- **Cons:** Constant interruptions, defeats the purpose of
  autonomous execution
- **Rejected because:** task-level autonomy with phase-level
  gates is the right balance.

### No workflow — just build

- **Pros:** No overhead for small changes
- **Cons:** Complex features get built ad-hoc, requirements
  are lost
- **Rejected because:** the depth-adaptive model (ADR 0008)
  already handles this — Quick depth skips the workflow
  entirely.

## Consequences

### Positive

- Agent can resume from any point by reading file state
- Human stays in control at decision points
- Frequent commits provide rollback points
- Branch-per-feature keeps main clean

### Negative

- Phase transitions require human presence
- Workflow overhead for features that turn out to be simpler
  than expected

### Neutral

- The workflow will be codified as a local skill once tested
  on a real feature
- Quick-depth features bypass this workflow entirely

## References

- ADR 0008 — spec storage and file formats
- [spec-kit](https://github.com/github/spec-kit) — specify →
  plan → tasks → implement phased workflow
- [BMAD Method](https://github.com/bmad-code-org/BMAD-METHOD)
  — scale-adaptive depth, "what's next?" resumability
- [MUSUBI](https://github.com/nahisaho/MUSUBI) — validation
  reports per feature, review gates between phases
- [Superpowers](https://github.com/obra/superpowers) —
  two-stage review (spec compliance + code quality),
  fresh-agent validation
- [cc-sdd](https://github.com/gotalab/cc-sdd) — boundary
  annotations on tasks, implementation notes propagating
  forward between tasks
- [get-shit-done](https://github.com/gsd-build/get-shit-done)
  — discuss phase, auto-detect next step, context engineering
- [don-cheli-sdd](https://github.com/doncheli/don-cheli-sdd)
  — worktree isolation, drift detection between spec and code
- [agent-os](https://github.com/buildermethods/agent-os) —
  standards discovery and injection into agent context
- [shotgun](https://github.com/shotgun-sh/shotgun) —
  codebase-aware research before specification, staged PRs
- [specswarm](https://github.com/MartyBonacci/specswarm) —
  quality scoring (0-100), tech stack drift prevention
