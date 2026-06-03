.PHONY: skills
skills:
	@$(MAKE) -C mk/skills $(filter-out $@,$(MAKECMDGOALS))

ifneq ($(filter skills,$(MAKECMDGOALS)),)
.PHONY: list status
list status:
	@:
endif
