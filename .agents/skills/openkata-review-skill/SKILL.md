---
name: openkata-review-skill
description: >
  Lints, reviews, and optimizes a skill using the Tessl toolchain. Validates
  structure, evaluates quality, and applies improvements. Use when the user
  creates a new skill, updates an existing skill, asks to review or optimize
  a skill, or says "lint", "review", or "optimize" in the context of skills.
---

# Review Skill

Run the Tessl quality pipeline on a skill — lint, review, optimize. Works
on both new and updated skills in `.agents/skills/` or `skills/`.

## Workflow

1. **Identify the skill** — Determine which skill to review from the
   user's request. Resolve the directory path. If ambiguous, list
   available skills and ask.

2. **Lint** — Run `tessl skill lint <skill-directory>`. Fix any
   structural issues before proceeding.

3. **Review** — Run `tessl skill review <skill-directory>`. Read the
   full output — scores, suggestions, and warnings. This context is
   needed for changelogs and commits.

4. **Optimize** — Run `tessl skill review --optimize <skill-directory>`.
   Read the diff and suggested changes. Apply improvements that make
   sense. If the interactive prompt blocks, apply the diff manually.

5. **Re-lint** — Run `tessl skill lint <skill-directory>` again to
   confirm the optimized skill is still valid.

6. **Report** — Summarize what changed, the before/after review score,
   and any remaining warnings.

7. **Persist score** — After review, update the skill's
   `.tessl-plugin/plugin.json` with the final score:
   read the JSON, set the `"score"` field to the numeric
   percentage, and write it back. This keeps scores
   queryable without re-running tessl.

8. **Quality checklist** — If `skills/create-skill/` and
   `skills/create-skill/references/skill-design-checklist.md` both
   exist, run the checklist against the reviewed skill. Report
   findings. A missing `## Boundaries` section is a blocking
   failure — do not pass the review without it. If either path
   is missing, skip silently.

9. **Feed learnings back** — After every review, check whether the
   insights could improve `create-skill` or the skill design
   checklist. If a pattern keeps surfacing (e.g., a common mistake,
   a better convention), ask the user if they want to update the
   skill creation workflow so future skills benefit.

## When Applying Improvements

- **Generalize from feedback.** Don't overfit fixes to specific cases.
  A skill will be used across many different prompts — ensure changes
  apply broadly, not just to the example that surfaced the issue.
- **Look for repeated patterns.** If the agent keeps writing the same
  helper script or taking the same multi-step approach when using a
  skill, bundle that into `scripts/` so every future invocation
  benefits.
- **Explain the why.** When adding instructions, explain reasoning
  instead of using heavy-handed `ALWAYS`/`NEVER` rules. Agents respond
  better to understanding why something matters.
- **Keep it lean.** Remove instructions that aren't pulling their
  weight. If something isn't improving outputs, cut it.

## Example Scenario

User: "Review the create-adr skill."

1. Runs `tessl skill lint skills/create-adr` — passes
2. Runs `tessl skill review skills/create-adr` — 99%
3. Runs `tessl skill review --optimize` — no changes needed
4. Runs checklist — all sections pass
5. Reports: 99%, no blocking issues, one orphaned file warning

## Boundaries

**DOES:**
- Run tessl lint, review, and optimize on skills
- Run the skill design checklist
- Suggest improvements and feed learnings back

**Does NOT:**
- Create new skills
- Publish or release skills
- Modify `create-skill` without asking

## Gotchas

- Tessl's optimize prompt is interactive — apply diffs manually when
  the prompt can't be answered.
- Markdown links in SKILL.md are validated by the linter as file paths
  relative to the skill directory. Use plain text references for files
  outside the skill folder.
