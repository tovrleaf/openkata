.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help
	@echo "Open Kata - Make Commands"
	@echo "========================="
	@echo ""
	@echo "Catalog:"
	@printf "  \033[36m%-12s\033[0m %s\n" "skills" "List all skills with type and version"
	@printf "  \033[36m%-12s\033[0m %s\n" "rules" "List all rules with type and version"
	@printf "  \033[36m%-12s\033[0m %s\n" "adrs" "List all architecture decision records"
	@printf "  \033[36m%-12s\033[0m %s\n" "specs" "List all feature specs with status"
	@echo ""
	@echo "Development:"
	@printf "  \033[36m%-12s\033[0m %s\n" "check" "Check development prerequisites"
	@printf "  \033[36m%-12s\033[0m %s\n" "changelog" "Generate root CHANGELOG.md"
	@printf "  \033[36m%-12s\033[0m %s\n" "dev" "Start local dev server with hot reload"

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

.PHONY: specs
specs: ## List all feature specs with status
	@./scripts/list-specs.sh

.PHONY: changelog
changelog: ## Generate root CHANGELOG.md from artifact changelogs
	@./scripts/generate-changelog.sh

.PHONY: dev
dev: ## Start local dev server with hot reload
	@lsof -ti:8080 | xargs kill 2>/dev/null || true
	@$(shell go env GOPATH)/bin/air
