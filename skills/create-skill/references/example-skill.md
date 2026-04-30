# Example Skill

A finished skill for running the test suite:

```markdown
---
name: run-tests
description: >
  Runs the project's test suite and reports results. Use when
  the user wants to run tests, check if tests pass, verify
  changes don't break anything, or asks about test failures.
---

# Run Tests

Run the full test suite, surface failures clearly, and suggest
fixes when the cause is obvious.

## Workflow

1. **Detect the test runner** — Check package.json scripts,
   Makefile targets, or pyproject.toml for the test command.
   Prefer `npm test`, `make test`, or `pytest` in that order.

2. **Run the suite** — Execute the detected command. Capture
   stdout and stderr.

3. **Report results** — If all tests pass, confirm with a
   one-line summary. If tests fail, list each failing test
   with its error message and the file:line reference.

4. **Suggest fixes** — For failures with an obvious cause
   (import error, missing env var, typo), propose a concrete
   fix. For ambiguous failures, ask the user before changing
   anything.

## Conventions

- Never modify test files to make tests pass.
- Run the full suite unless the user explicitly scopes to a
  subset.
```

This example shows:

- A description with natural trigger phrases and situations
- A tight four-step workflow
- Conventions that constrain behavior without rigid rules
