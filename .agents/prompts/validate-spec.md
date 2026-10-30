# Validate Spec

Check a built feature against its spec with an agent that has no
builder bias, and write a validation report.

Run this after all tasks of a spec are built and before opening a
pull request.

## Workflow

1. **Resolve the spec** — The user may name one: a bare number
   (`23` → `0023`), words from the slug (`docs sec`), or a full
   directory name. Scan `specs/` for matches. One match: use it.
   Several: list them and ask. None: list the specs and ask. With
   no hint, use the path in `specs/_current`, then the session
   context.
2. **Confirm the spec** — Read `specs/NNNN-slug/spec.md` to check
   it exists. Update `specs/_current` to point to it.
3. **Get the diff** — Run `git diff main...HEAD`. If the work is
   already merged, use the commit range the user gives.
4. **Spawn a fresh agent** — Use a new agent with no session
   context. Give it the instructions in
   `profiles/spec-validator/spec-validator.md` and pass only:
   - the path to `specs/NNNN-slug/spec.md`
   - the raw diff
   - the paths of the relevant source files
   - the validation report template in
     `skills/spec-workflow/references/spec-templates.md`

   Do not pass `tasks.md`, `design.md`, the progress log or any
   other session context.
5. **Report** — The agent writes
   `specs/NNNN-slug/validation-report.md`. Show the result to the
   user. If issues are found, the feature returns to the build
   phase; do not fix them as part of validation.

## Constraints

- The validator reads the spec and the code, never the task plan.
- The validator reports issues and does not fix them.
- Commit the report with the feature so the check is on record.
