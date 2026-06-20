---
name: openkata-eval-runner
description: >
  Run eval scenarios for an OpenKata skill locally and iterate
  fixes until 95% threshold. Use when the user says "run evals",
  "check evals", "eval this skill", or when validating skill
  quality.
---

# Eval Runner

Run eval scenarios locally and iterate skill fixes until
the 95% threshold is met.

## Tool Split

| Tool | Use for |
|------|---------|
| `./bin/openkata-eval` | Scoring (runs locally, same model as users) |
| `tessl skill review --optimize` | Instruction polish after passing |
| `tessl eval run` | NOT used — model mismatch makes results unreliable |

## Workflow

1. **Identify the skill** — Determine which skill from the
   user's request. Resolve path (skills/<name>/).

2. **Check evals exist** — Verify `skills/<name>/evals/`
   has scenario directories. If missing, activate the
   `create-evals` skill first.

3. **Build the runner** — `make eval-local` builds the binary
   automatically if missing.

4. **Run evals** — Execute:
   ```bash
   ./bin/openkata-eval skills/<name>
   ```
   The overall average must be 95% or above.

5. **Report** — If passing (95%+), report the score.
   If failing, report which scenarios/criteria failed
   with reasons, and suggest specific SKILL.md fixes.

6. **Fix and retry** — If the user confirms fixes, apply
   them to SKILL.md, then re-run. Iterate until 95%+.

7. **Optimize** — Once passing locally, run:
   ```bash
   tessl skill review --optimize skills/<name>
   ```
   Apply any worthwhile suggestions, re-run local evals
   to confirm no regression.

8. **Commit** — Stage changes and commit.

## Single Scenario Debugging

To iterate on one failing scenario without running all:
```bash
./bin/openkata-eval skills/<name>/evals/scenario-X
```
Single scenario mode always exits 0 (debug tool).

## Boundaries

- DOES run evals, suggest and apply SKILL.md fixes
- DOES run tessl optimize for polish after passing
- Does NOT publish, tag, or release
- Does NOT use tessl eval run for scoring

## Gotchas

- Only distributable skills (`skills/`) have evals
- Conversational skills need `"sandbox": false` in
  scenario.json (no Docker required)
- Skills that use tools need `"sandbox": true` (default)
  and Docker running
