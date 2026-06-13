ROOT := ../..

.DEFAULT_GOAL := help

.PHONY: help web mcp

help:
	@echo "Usage: make deploy web|mcp"

web:
	@cd $(ROOT) && ./scripts/deploy-web.sh

mcp:
	@cd $(ROOT) && ./scripts/deploy-mcp.sh
