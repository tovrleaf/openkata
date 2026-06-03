.DEFAULT_GOAL := help

include mk/skills.mk
include mk/chat.mk
include mk/catalog.mk
include mk/dev.mk

.PHONY: help
help: ## Show this help
	@echo "Open Kata - Make Commands"
	@echo "========================="
	@echo ""
	@echo "Catalog:"
	@printf "  \033[36m%-16s\033[0m %s\n" "skills" "Show skills help"
	@printf "  \033[36m%-16s\033[0m %s\n" "skills list" "List all skills with type and version"
	@printf "  \033[36m%-16s\033[0m %s\n" "skills status" "List skills with registry status"
	@printf "  \033[36m%-16s\033[0m %s\n" "rules" "List all rules with type and version"
	@printf "  \033[36m%-16s\033[0m %s\n" "profiles" "List all agent profiles"
	@printf "  \033[36m%-16s\033[0m %s\n" "agents" "List all Kiro agent configs"
	@printf "  \033[36m%-16s\033[0m %s\n" "adrs" "List all architecture decision records"
	@printf "  \033[36m%-16s\033[0m %s\n" "specs" "List all feature specs with status"
	@echo ""
	@echo "Development:"
	@printf "  \033[36m%-16s\033[0m %s\n" "check" "Check development prerequisites"
	@printf "  \033[36m%-16s\033[0m %s\n" "test" "Run all tests"
	@printf "  \033[36m%-16s\033[0m %s\n" "changelog" "Generate root CHANGELOG.md"
	@printf "  \033[36m%-16s\033[0m %s\n" "dev" "Start local dev server with hot reload"
	@printf "  \033[36m%-16s\033[0m %s\n" "chat master" "Start Kiro chat with dojo-master agent"
	@printf "  \033[36m%-16s\033[0m %s\n" "chat eval" "Start Kiro chat with kata-author for eval work"
	@printf "  \033[36m%-16s\033[0m %s\n" "versions" "Generate versions.json from local files"
	@printf "  \033[36m%-16s\033[0m %s\n" "deploy" "Deploy web server to AWS Lambda"
