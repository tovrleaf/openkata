# Phase 5: Validate

A fresh agent session is recommended for validation to
avoid builder bias.

1. **Read spec.md** — load requirements and acceptance
   criteria only. Do not read tasks.md or the progress log.

2. **Review the implementation** — read the code changes
   (use `git diff main` or the feature branch diff).

3. **Write validation-report.md** — Use the template from
   [spec-templates](spec-templates.md):
   - Check each requirement against the implementation
   - Verify out-of-scope items were not built
   - Note any issues found

4. **Report** — show the validation report to the user.
   If issues are found, the feature returns to Phase 4
   for fixes.
