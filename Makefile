.PHONY: skills

skills: ## List all skills with type, version, and change status
	@./scripts/list-skills.sh

.PHONY: rules

rules: ## List all rules with type and version
	@./scripts/list-rules.sh