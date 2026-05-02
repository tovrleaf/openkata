# Makefile Structure

## Directory layout

```text
project/
├── Makefile              # Entry point, includes mk/*.mk
└── mk/
    ├── test.mk           # Delegates to mk/test/Makefile
    ├── deploy.mk         # Delegates to mk/deploy/Makefile
    ├── test/
    │   └── Makefile      # Subcommands: sh, yaml, md, html
    └── deploy/
        └── Makefile      # Subcommands: staging, production
```

## Main Makefile

The root Makefile includes all domain modules and provides
the top-level help:

```makefile
.DEFAULT_GOAL := help

include mk/test.mk
include mk/deploy.mk

.PHONY: help
help: ## Show this help
	@echo "Project Name - Make Commands"
	@echo "============================"
	@echo ""
	@echo "Testing:"
	@echo "  make test       - Show test help"
	@echo "  make test sh    - Run shellcheck"
	@echo "  make test yaml  - Run yamllint"
	@echo ""
	@echo "Deployment:"
	@echo "  make deploy     - Show deploy help"
	@echo "  make deploy production - Deploy to production"
```

## Domain module (mk/test.mk)

Each `.mk` file delegates to a subdirectory Makefile:

```makefile
.PHONY: test

test:
	@$(MAKE) -C mk/test $(filter-out $@,$(MAKECMDGOALS))

test-%:
	@$(MAKE) -C mk/test $*

%:
	@:
```

The `filter-out` pattern passes subcommands through:
`make test sh` calls `$(MAKE) -C mk/test sh`.

The catch-all `%: @:` prevents "No rule to make target"
errors for subcommand arguments.

## Domain Makefile (mk/test/Makefile)

Contains the actual targets and domain-specific help:

```makefile
.DEFAULT_GOAL := help

.PHONY: help sh yaml

help:
	@echo "Testing"
	@echo "======="
	@echo ""
	@echo "Commands:"
	@echo "  make test sh   - Run shellcheck"
	@echo "  make test yaml - Run yamllint"

sh:
	@./scripts/test-shell.sh

yaml:
	@yamllint .
```

## Simple targets

Not every target needs the delegation pattern. Simple,
standalone targets go directly in the main Makefile or
a flat `.mk` file:

```makefile
.PHONY: check
check: ## Check development prerequisites
	@./scripts/check-prereqs.sh

.PHONY: clean
clean: ## Remove build artifacts
	@rm -rf dist/ build/
```

Use the modular pattern when a domain has multiple
subcommands. Use flat targets when there is only one
command for the concern.
