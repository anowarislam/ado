SHELL := /bin/bash
.DEFAULT_GOAL := help

include make/go.mk
include make/python.mk

# =============================================================================
# Testing
# =============================================================================

.PHONY: test
test: go.test py.test ## Run all tests (Go + Python lab)

# =============================================================================
# Validation (mirrors CI)
# =============================================================================

.PHONY: validate
validate: lint test ## Run all validations (lint + test)

.PHONY: lint
lint: go.vet go.fmt.check py.lint ## Run all linters

.PHONY: go.fmt.check
go.fmt.check: ## Check Go code formatting (fails if not formatted)
	@echo "Checking Go formatting..."
	@test -z "$$(gofmt -l .)" || (echo "Go files need formatting:" && gofmt -l . && exit 1)

.PHONY: ci
ci: validate go.build ## Run full CI pipeline locally

# =============================================================================
# Pre-commit hooks
# =============================================================================

.PHONY: hooks.install
hooks.install: ## Install git hooks for conventional commits
	@git config core.hooksPath .githooks
	@echo "Git hooks installed. Commits will be validated for conventional format."

.PHONY: hooks.uninstall
hooks.uninstall: ## Uninstall git hooks
	@git config --unset core.hooksPath || true
	@echo "Git hooks uninstalled."

.PHONY: precommit.install
precommit.install: ## Install pre-commit framework hooks
	@command -v pre-commit >/dev/null 2>&1 || (echo "Installing pre-commit..." && pip install pre-commit)
	@pre-commit install
	@pre-commit install --hook-type commit-msg
	@echo "Pre-commit hooks installed."

.PHONY: precommit.run
precommit.run: ## Run pre-commit on all files
	@pre-commit run --all-files

# =============================================================================
# Documentation
# =============================================================================

.PHONY: docs.install
docs.install: ## Install MkDocs and dependencies
	@pip install mkdocs-material mkdocs-minify-plugin pillow cairosvg

.PHONY: docs.build
docs.build: ## Build documentation site
	@cp CHANGELOG.md docs/changelog.md
	@mkdocs build --strict
	@echo "Documentation built in site/"

.PHONY: docs.serve
docs.serve: ## Serve documentation locally (http://localhost:8000)
	@cp CHANGELOG.md docs/changelog.md
	@mkdocs serve

.PHONY: docs.clean
docs.clean: ## Clean built documentation
	@rm -rf site/
	@rm -f docs/changelog.md
	@echo "Cleaned documentation build artifacts"

# =============================================================================
# Help
# =============================================================================

.PHONY: help
help: ## Show available make targets
	@echo "Available targets:"
	@echo ""
	@echo "Testing:"
	@grep -hE '^[a-zA-Z0-9_.-]+:.*## Run' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS=":.*##"} {printf "  %-20s %s\n", $$1, $$2}'
	@echo ""
	@echo "Validation:"
	@grep -hE '^(validate|lint|ci|go\.fmt\.check):.*##' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS=":.*##"} {printf "  %-20s %s\n", $$1, $$2}'
	@echo ""
	@echo "Hooks:"
	@grep -hE '^(hooks\.|precommit\.)[a-zA-Z0-9_.-]+:.*##' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS=":.*##"} {printf "  %-20s %s\n", $$1, $$2}'
	@echo ""
	@echo "All targets:"
	@grep -hE '^[a-zA-Z0-9_.-]+:.*##' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS=":.*##"} {printf "  %-20s %s\n", $$1, $$2}'
