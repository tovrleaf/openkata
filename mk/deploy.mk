.PHONY: deploy
deploy:
	@$(MAKE) -C mk/deploy $(filter-out $@,$(MAKECMDGOALS))

ifneq ($(filter deploy,$(MAKECMDGOALS)),)
.PHONY: web mcp
web mcp:
	@:
endif
