.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help
	@echo "Open Kata - Make Commands"
	@echo "========================="
	@echo ""
	@echo "Catalog:"
	@echo "  make skills  - List all skills with type and version"
	@echo "  make rules   - List all rules with type and version"
	@echo "  make adrs    - List all architecture decision records"
	@echo ""
	@echo "Development:"
	@echo "  make check     - Check development prerequisites"
	@echo "  make changelog - Generate root CHANGELOG.md"

.PHONY: skills
skills: ## List all skills with type, version, and change status
	@./scripts/list-skills.sh

.PHONY: rules
rules: ## List all rules with type and version
	@./scripts/list-rules.sh

.PHONY: check
check: ## Check development prerequisites
	@./scripts/check-prereqs.sh

.PHONY: adrs
adrs: ## List all architecture decision records
	@./scripts/list-adrs.sh

.PHONY: changelog
changelog: ## Generate root CHANGELOG.md from artifact changelogs
	@./scripts/generate-changelog.sh
