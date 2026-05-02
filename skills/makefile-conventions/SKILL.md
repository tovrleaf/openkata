---
name: makefile-conventions
version: 1.0.0
description: >
  Structures Makefiles as a universal command interface using
  modular includes and self-documenting help. Use when the user
  wants to add a make target, create a Makefile, organize build
  commands, set up a make-based workflow, or asks how to
  structure make targets.
---

# Makefile Conventions

Use Make as the universal command interface. All project
operations are accessible via `make`, regardless of the
underlying implementation language.

## Workflow

1. **Investigate the project** — Before writing targets, check:
   - Does a Makefile already exist? Read it.
   - What scripts, tools, or build steps does the project use?
   - What language ecosystems are involved (Go, Node, Python)?
   - Is there an existing `mk/` directory?

2. **Decide placement** — If the target belongs to an existing
   concern group, add it there. If it starts a new concern,
   create a new modular makefile. See
   [makefile-structure](references/makefile-structure.md) for
   the modular pattern.

3. **Write the target** — Follow these conventions:

   - Declare every target `.PHONY` unless it produces a file
   - Add a `## Description` comment for the help system
   - Use tabs for indentation (Make requirement)
   - Keep targets short — delegate to scripts for complex logic
   - Use `@` prefix to suppress command echo for clean output
   - Quote all variable references: `"$(VAR)"`

4. **Update help** — Add the new target to the appropriate
   section in `make help` output. Group by concern:

   ```text
   Content Management:
     make live       - Manage live performances
     make media      - Manage media files

   Testing:
     make test sh    - Run shellcheck
     make test yaml  - Run yamllint

   Development:
     make hooks setup - Install pre-commit hooks
   ```

5. **Verify** — Run `make help` to confirm the target appears.
   Run the target to confirm it works.

## Conventions

### Structure

- One main `Makefile` at the project root
- Modular makefiles in `mk/` for each concern domain
- Domains with subcommands get `mk/<domain>/Makefile`
- Main Makefile includes all `mk/*.mk` files
- `.DEFAULT_GOAL := help` so bare `make` shows help

### Naming

- Target names are lowercase, hyphenated for multi-word
- Subcommands use space separation: `make test sh`,
  `make adr new`
- Variable names are UPPER_CASE

### Help system

- `make` or `make help` shows all available commands
- `make <domain>` shows domain-specific help
- Help output groups targets by concern with descriptions
- Consistent format: `make <command>` followed by description

### Delegation pattern

Makefiles are the interface, not the implementation. Complex
logic lives in scripts:

```makefile
deploy:
	@./scripts/deploy.sh "$(ENV)"
```

This lets you change the implementation language without
changing the command interface.

## Example Scenario

User: "Add a make target for running database migrations"

1. Reads existing Makefile — finds `mk/` structure with
   `test.mk`, `deploy.mk`
2. No existing `db.mk` — creates `mk/db.mk` with a `db`
   domain target
3. Creates `mk/db/Makefile` with `migrate`, `rollback`,
   `status` subcommands
4. Adds `include mk/db.mk` to main Makefile
5. Adds Database section to `make help` output
6. Verifies: `make db migrate` runs the migration script

## Common Failures

- **No help text** — targets without descriptions are
  undiscoverable. Always update `make help`.
- **Logic in Makefiles** — complex bash in a target is hard
  to debug and test. Delegate to scripts.
- **Missing .PHONY** — without it, Make checks for a file
  with the target name and may skip execution.
- **Spaces instead of tabs** — Make requires tabs for recipe
  indentation. This is the most common syntax error.
