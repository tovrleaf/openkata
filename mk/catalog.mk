.PHONY: rules
rules:
	@./scripts/list-rules.sh

.PHONY: profiles
profiles:
	@./scripts/list-profiles.sh

.PHONY: agents
agents:
	@./scripts/list-agents.sh

.PHONY: adrs
adrs:
	@./scripts/list-adrs.sh

.PHONY: specs
specs:
	@./scripts/list-specs.sh
