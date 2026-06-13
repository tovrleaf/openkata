---
name: openkata-eval-runner
description: >
  Regenerate and run eval scenarios for an OpenKata skill.
  Use when the user says "regenerate evals", "run evals",
  "check evals", or when evals are stale.
---

# Eval Runner

Generate and run eval scenarios for distributable skills.
Iterate fixes until the 95% threshold is met.

## Workflow

1. **Identify the artifact** — Determine which skill from the
   user's request. Resolve path (skills/<name>/).
2. **Create worktree** — `git worktree add .worktrees/eval-<name> -b eval/<name>`.
   cd into it. If it already exists, reuse it.
3. **Regenerate scenarios** — Follow the `create-evals` skill
   to generate scenarios for `skills/<name>/`.
4. **Run evals** — Run `tessl eval run skills/<name>/ --variant with-context`.
   The "With context" average must be 95% or above.
5. **Report** — If passing (95%+), report the score. If failing,
   report which scenarios/checks failed and suggest specific
   SKILL.md fixes.
6. **Fix and retry** — If the user confirms fixes, apply them to
   SKILL.md, then go back to step 3. Iterate until 95%+.
7. **Commit** — Stage changes and commit in the worktree.
8. **Notify** — Tell the user the worktree is ready to merge:
   `eval/<name>`

## Boundaries

- DOES generate scenarios, run evals, suggest and apply
  SKILL.md fixes to improve scores
- Does NOT commit, tag, release, or publish — those are
  `openkata-ryu-release`'s job

## Gotchas

- Only distributable skills (`skills/`) have evals. Local
  skills (`.agents/skills/`) do not.
- If `create-evals` finds existing scenarios that already
  cover the skill's behaviors, it may skip generation.

