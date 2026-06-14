.PHONY: check
check:
	@./scripts/check-prereqs.sh

.PHONY: test
test:
	@go test ./...

.PHONY: changelog
changelog:
	@./scripts/generate-changelog.sh

.PHONY: dev
dev:
	@lsof -ti:8080 | xargs kill 2>/dev/null || true
	@$(shell go env GOPATH)/bin/air

.PHONY: versions
versions:
	@go run ./cmd/generate-versions/ --local

.PHONY: badges
badges: ## Update README badges with current counts
	@./scripts/update-readme-badges.sh

.PHONY: stats-fetch
stats-fetch: ## Fetch analytics data to .local/stats/
	@go run ./cmd/stats-fetch/