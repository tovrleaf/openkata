# Phase 1: Specify

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
   [spec-templates](spec-templates.md). Set
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
