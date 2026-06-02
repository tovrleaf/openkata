ROOT := ../..

.DEFAULT_GOAL := help

.PHONY: help master eval

help:
	@echo "Usage: make chat master|eval"

master:
	@cd $(ROOT) && kiro-cli chat --model claude-opus-4.6 --agent dojo-master --trust-all-tools

eval:
	@cd $(ROOT) && kiro-cli chat --model claude-opus-4.6 --agent kata-author --trust-all-tools
