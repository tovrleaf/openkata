# Phase 4: Implement

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
