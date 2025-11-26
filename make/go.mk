# ==============================================================================
# Go Targets
# ==============================================================================
# Go build, test, and development targets

# ------------------------------------------------------------------------------
# Variables
# ------------------------------------------------------------------------------
GO ?= go
GO_ENV ?= env GOCACHE=$(PROJECT_ROOT)/.gocache
GO_BUILD_CMD ?= ./cmd/ado
GO_BINARY ?= $(PROJECT_ROOT)/ado

# Build flags
GO_LDFLAGS := -s -w \
	-X github.com/anowarislam/ado/internal/meta.Version=$(VERSION) \
	-X github.com/anowarislam/ado/internal/meta.Commit=$(COMMIT) \
	-X github.com/anowarislam/ado/internal/meta.BuildTime=$(BUILD_TIME)

# Packages to test/vet (exclude main package which has no testable code)
GO_TEST_PKGS ?= ./cmd/ado/... ./internal/...

# Test flags
GO_TEST_FLAGS ?= -race -timeout 30s
GO_COVER_FLAGS ?= -coverprofile=coverage.out -covermode=atomic

# ------------------------------------------------------------------------------
# Build Targets
# ------------------------------------------------------------------------------
.PHONY: go.build go.build.all go.install

go.build: _check-go ## Build the ado binary
	$(call log_info,"Building $(PROJECT_NAME)...")
	@$(GO_ENV) $(GO) build -ldflags "$(GO_LDFLAGS)" -o $(GO_BINARY) $(GO_BUILD_CMD)
	$(call log_success,"Built: $(GO_BINARY)")

go.build.all: _check-go ## Build for all platforms (via goreleaser)
	$(call log_info,"Building for all platforms...")
	@command -v goreleaser >/dev/null 2>&1 || { \
		echo "goreleaser not installed. Install from https://goreleaser.com"; \
		exit 1; \
	}
	@goreleaser build --snapshot --clean
	$(call log_success,"All platform builds complete in dist/")

go.install: _check-go ## Install binary to GOPATH/bin
	$(call log_info,"Installing $(PROJECT_NAME)...")
	@$(GO_ENV) $(GO) install -ldflags "$(GO_LDFLAGS)" $(GO_BUILD_CMD)
	$(call log_success,"Installed to $$(go env GOPATH)/bin/$(PROJECT_NAME)")

# ------------------------------------------------------------------------------
# Test Targets
# ------------------------------------------------------------------------------
.PHONY: go.test go.test.cover go.test.verbose go.test.race go.bench

go.test: _check-go ## Run Go tests
	$(call log_info,"Running Go tests...")
	@$(GO_ENV) $(GO) test $(GO_TEST_FLAGS) $(GO_TEST_PKGS)
	$(call log_success,"All Go tests passed")

go.test.cover: _check-go ## Run Go tests with coverage report
	$(call log_info,"Running Go tests with coverage...")
	@$(GO_ENV) $(GO) test $(GO_TEST_FLAGS) $(GO_COVER_FLAGS) $(GO_TEST_PKGS)
	@$(GO) tool cover -func=coverage.out | tail -1
	$(call log_success,"Coverage report: coverage.out")

go.test.cover.html: go.test.cover ## Generate HTML coverage report
	$(call log_info,"Generating HTML coverage report...")
	@$(GO) tool cover -html=coverage.out -o coverage.html
	$(call log_success,"Coverage report: coverage.html")
	@echo "Open coverage.html in a browser to view the report"

go.test.cover.check: go.test.cover ## Check coverage meets threshold (80%)
	$(call log_info,"Checking coverage threshold (80%)...")
	@COVERAGE=$$($(GO) tool cover -func=coverage.out | grep total | awk '{print $$3}' | tr -d '%'); \
	if [ "$${COVERAGE%.*}" -lt 80 ]; then \
		echo "$(COLOR_RED)Coverage $${COVERAGE}% is below 80% threshold$(COLOR_RESET)"; \
		exit 1; \
	fi
	$(call log_success,"Coverage meets threshold")

go.test.verbose: _check-go ## Run Go tests with verbose output
	$(call log_info,"Running Go tests (verbose)...")
	@$(GO_ENV) $(GO) test -v $(GO_TEST_FLAGS) $(GO_TEST_PKGS)

go.test.race: _check-go ## Run Go tests with race detector
	$(call log_info,"Running Go tests with race detector...")
	@$(GO_ENV) $(GO) test -race $(GO_TEST_PKGS)
	$(call log_success,"No race conditions detected")

go.bench: _check-go ## Run Go benchmarks
	$(call log_info,"Running Go benchmarks...")
	@$(GO_ENV) $(GO) test -bench=. -benchmem $(GO_TEST_PKGS)

# ------------------------------------------------------------------------------
# Lint Targets
# ------------------------------------------------------------------------------
.PHONY: go.fmt go.fmt.check go.vet go.lint go.tidy

go.fmt: ## Format Go sources
	$(call log_info,"Formatting Go sources...")
	@$(GO) fmt ./...
	$(call log_success,"Go sources formatted")

go.fmt.check: ## Check Go code formatting (fails if not formatted)
	$(call log_info,"Checking Go formatting...")
	@test -z "$$(gofmt -l .)" || { \
		echo "$(COLOR_RED)Go files need formatting:$(COLOR_RESET)"; \
		gofmt -l .; \
		exit 1; \
	}
	$(call log_success,"Go formatting check passed")

go.vet: _check-go ## Run go vet on all packages
	$(call log_info,"Running go vet...")
	@$(GO_ENV) $(GO) vet $(GO_TEST_PKGS)
	$(call log_success,"go vet passed")

go.lint: _check-go ## Run golangci-lint (if installed)
	$(call log_info,"Running golangci-lint...")
	@command -v golangci-lint >/dev/null 2>&1 || { \
		echo "golangci-lint not installed. Falling back to go vet."; \
		$(MAKE) --no-print-directory go.vet; \
		exit 0; \
	}
	@golangci-lint run ./...
	$(call log_success,"golangci-lint passed")

go.tidy: _check-go ## Sync Go module dependencies
	$(call log_info,"Tidying Go modules...")
	@$(GO) mod tidy
	$(call log_success,"Go modules tidied")

# ------------------------------------------------------------------------------
# Dependency Targets
# ------------------------------------------------------------------------------
.PHONY: go.deps go.deps.update go.deps.graph

go.deps: _check-go ## Download Go dependencies
	$(call log_info,"Downloading Go dependencies...")
	@$(GO) mod download
	@$(GO) mod verify
	$(call log_success,"Go dependencies downloaded and verified")

go.deps.update: _check-go ## Update Go dependencies
	$(call log_info,"Updating Go dependencies...")
	@$(GO) get -u ./...
	@$(GO) mod tidy
	$(call log_success,"Go dependencies updated")

go.deps.graph: _check-go ## Show Go dependency graph
	@$(GO) mod graph

# ------------------------------------------------------------------------------
# Cleanup Targets
# ------------------------------------------------------------------------------
.PHONY: go.clean

go.clean: ## Remove Go build artifacts
	$(call log_info,"Cleaning Go build artifacts...")
	@rm -rf .gocache
	@rm -f $(GO_BINARY)
	@rm -f coverage.out
	$(call log_success,"Go build artifacts cleaned")
