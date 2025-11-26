# ==============================================================================
# ado Makefile
# ==============================================================================
# Main entry point - includes all sub-makefiles from make/
#
# Usage:
#   make help        Show all available targets
#   make setup       First-time project setup
#   make build       Build the binary
#   make test        Run all tests
#   make ci          Run full CI pipeline locally
#
# See make/*.mk for individual target definitions.
# ==============================================================================

.DEFAULT_GOAL := help

# ------------------------------------------------------------------------------
# Include Sub-Makefiles
# Order matters: common.mk must be first as it defines shared variables
# ------------------------------------------------------------------------------
include make/common.mk
include make/go.mk
include make/python.mk
include make/docker.mk
include make/docs.mk
include make/hooks.mk
include make/help.mk

# ==============================================================================
# Composite Targets
# ==============================================================================
# High-level targets that combine multiple sub-targets

# ------------------------------------------------------------------------------
# Setup
# ------------------------------------------------------------------------------
.PHONY: setup
setup: hooks.install go.deps py.install ## First-time project setup
	$(call log_header,"Setup Complete")
	@echo ""
	@echo "Next steps:"
	@echo "  make build    Build the binary"
	@echo "  make test     Run all tests"
	@echo "  make help     Show all targets"
	@echo ""

# ------------------------------------------------------------------------------
# Build
# ------------------------------------------------------------------------------
.PHONY: build
build: go.build ## Build the project (alias for go.build)

# ------------------------------------------------------------------------------
# Test
# ------------------------------------------------------------------------------
.PHONY: test
test: go.test py.test ## Run all tests (Go + Python)

.PHONY: test.cover
test.cover: go.test.cover py.test.cover ## Run all tests with coverage

# ------------------------------------------------------------------------------
# Lint
# ------------------------------------------------------------------------------
.PHONY: lint
lint: go.fmt.check go.vet py.lint ## Run all linters

.PHONY: fmt
fmt: go.fmt py.fmt ## Format all code

# ------------------------------------------------------------------------------
# Validation (mirrors CI)
# ------------------------------------------------------------------------------
.PHONY: validate
validate: lint test docs.build ## Run all validations (lint + test + docs)
	$(call log_header,"Validation Complete")
	$(call log_success,"All checks passed!")

.PHONY: ci
ci: validate build ## Run full CI pipeline locally
	$(call log_header,"CI Pipeline Complete")
	$(call log_success,"Ready to push!")

# ------------------------------------------------------------------------------
# Quick Checks
# ------------------------------------------------------------------------------
.PHONY: check
check: go.fmt.check go.vet ## Quick syntax/lint check (no tests)
	$(call log_success,"Quick check passed")

.PHONY: pre-commit
pre-commit: check test ## Run before committing
	$(call log_success,"Pre-commit checks passed")

# ==============================================================================
# Info
# ==============================================================================
.PHONY: info
info: version ## Show project information
	@echo ""
	@echo "Directories:"
	@echo "  Root:     $(PROJECT_ROOT)"
	@echo "  Docs:     $(DOCS_DIR)"
	@echo "  Build:    $(BUILD_DIR)"
	@echo ""
	@echo "Binaries:"
	@ls -la $(GO_BINARY) 2>/dev/null || echo "  (not built yet - run 'make build')"
