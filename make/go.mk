GO ?= go
GO_ENV ?= env GOCACHE=$(PWD)/.gocache
GO_BUILD_CMD ?= ./cmd/ado

# Packages to test (exclude main package which has no testable code)
GO_TEST_PKGS ?= ./cmd/ado/... ./internal/...

.PHONY: go.fmt go.tidy go.test go.build go.clean go.vet go.test.cover

go.fmt: ## Format Go sources
	@$(GO) fmt ./...

go.tidy: ## Sync Go module dependencies
	@$(GO) mod tidy

go.test: ## Run Go tests
	@$(GO_ENV) $(GO) test $(GO_TEST_PKGS)

go.test.cover: ## Run Go tests with coverage report
	@$(GO_ENV) $(GO) test -cover $(GO_TEST_PKGS)

go.vet: ## Run go vet on all packages
	@$(GO_ENV) $(GO) vet $(GO_TEST_PKGS)

go.build: ## Build the ado binary
	@$(GO_ENV) $(GO) build $(GO_BUILD_CMD)

go.clean: ## Remove repo-local Go build cache
	@rm -rf .gocache
