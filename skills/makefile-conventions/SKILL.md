---
name: makefile-conventions
description: >
  Structures Makefiles as a universal command interface using
  modular includes and self-documenting help. Use when the user
  wants to add a make target, create a Makefile, organize build
  commands, set up a make-based workflow, or asks how to
  structure make targets. Also activate before modifying any
  Makefile or mk/ file in this project.
metadata:
  version: "1.0.0"
  tags: "category:conventions, tool:makefile"
---

# Makefile Conventions

Use Make as the universal command interface. All project
operations are accessible via `make`, regardless of the
underlying implementation language.

## Workflow

1. **Investigate the project** — Before writing targets, check:
   - Does a Makefile already exist? Read it.
   - What scripts, tools, or build steps does the project use?
   - Is there an existing `mk/` directory?

2. **Decide placement** — If the target belongs to an existing
   concern group, add it there. If it starts a new concern,
   create a new modular makefile. See
   [makefile-structure](references/makefile-structure.md) for
   the full modular pattern.

3. **Write the target** — Follow these conventions:

   - Declare every target `.PHONY` unless it produces a file
   - Add a `## Description` comment for the help system
   - Delegate complex logic to scripts
   - Use `@` prefix for clean output
   - Quote variable references: `"$(VAR)"`

4. **Update help** — Add the new target to the appropriate
   concern group in `make help` output.

5. **Verify** — Run `make help` to confirm the target appears.
   Run the target to confirm it works.

## Complete Example

Root Makefile pattern — includes domain modules and
provides grouped help:

```makefile
.DEFAULT_GOAL := help

include mk/test.mk
include mk/deploy.mk

.PHONY: help
help: ## Show this help
	@echo "Project - Make Commands"
	@echo "======================"
	@echo ""
	@echo "Testing:"
	@echo "  make test       - Show test help"
	@echo "  make test sh    - Run shellcheck"
	@echo ""
	@echo "Development:"
	@echo "  make check      - Check prerequisites"
```

For the full modular delegation pattern (domain `.mk`
files, subdirectory Makefiles, catch-all targets), see
[makefile-structure](references/makefile-structure.md).

## Conventions

### Structure

- One main `Makefile` at the project root
- Modular makefiles in `mk/` for each concern domain
- Domains with subcommands get `mk/<domain>/Makefile`
- `.DEFAULT_GOAL := help` so bare `make` shows help

### Naming

- Targets: lowercase, hyphenated for multi-word
- Subcommands: space separation (`make test sh`)
- Variables: UPPER_CASE

### Help system

- `make` or `make help` shows all commands grouped by concern
- `make <domain>` shows domain-specific help
- When a domain has 2 or more subcommands, bare
  `make <domain>` must show subcommand help. Do not make it
  execute a default action.

### Delegation

Makefiles are the interface, not the implementation:

```makefile
deploy:
	@./scripts/deploy.sh "$(ENV)"
```

## Example Scenario

User: "Add a make target for running database migrations"

1. Reads existing Makefile — finds `mk/` structure
2. No existing `db.mk` — creates `mk/db.mk` and
   `mk/db/Makefile` with `migrate`, `rollback`, `status`
3. Adds `include mk/db.mk` to main Makefile
4. Adds Database section to `make help` output
5. Verifies: `make db migrate` works

## Common Failures

- **No help text** — targets without descriptions are
  undiscoverable. Always update `make help`.
- **Logic in Makefiles** — complex bash in a target is hard
  to debug. Delegate to scripts.
