.PHONY: chat
chat:
	@$(MAKE) -C mk/chat $(filter-out $@,$(MAKECMDGOALS))

ifneq ($(filter chat,$(MAKECMDGOALS)),)
.PHONY: master eval
master eval:
	@:
endif
