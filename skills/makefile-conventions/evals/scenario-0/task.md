# Build Command Interface for a New Data Pipeline Project

## Problem Description

A data engineering team has just scaffolded a new Python project called `pipeflow` that processes and transforms large datasets. The project already has several scripts in a `scripts/` directory that handle common operations — running tests, linting the codebase, building distribution packages, and generating documentation. Currently, every developer has to remember the exact script paths and arguments to run anything, which leads to inconsistency and onboarding friction.

The team lead wants a single, consistent command interface so that anyone can run any project operation without knowing which script does what. The interface should be self-documenting — running it with no arguments should display all available commands organized by category. There are three main areas of concern: `test` (unit tests, coverage), `lint` (style checks, type checking), and `build` (packaging, docs). Each of these has two or more distinct subcommands.

The existing scripts are real and should be invoked by the command interface rather than duplicating their logic. The project has no existing Makefile.

## Output Specification

Create a complete Makefile-based command interface for the `pipeflow` project. The solution should include:

- A root `Makefile` at the project root
- Any supporting files needed for the modular command structure
- Supporting scripts can be assumed to already exist at `scripts/test-unit.sh`, `scripts/test-coverage.sh`, `scripts/lint-style.sh`, `scripts/lint-types.sh`, `scripts/build-package.sh`, and `scripts/build-docs.sh`

Demonstrate that the help system works by capturing the output of running `make` with no arguments and writing it to `make-help-output.txt`. Also capture `make test` (with no subcommand) to `make-test-help-output.txt` and `make lint` to `make-lint-help-output.txt`.
