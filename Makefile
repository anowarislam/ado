SHELL := /bin/bash
.DEFAULT_GOAL := help

include make/go.mk
include make/python.mk

.PHONY: test
test: go.test py.test ## Run all tests (Go + Python lab)

.PHONY: hooks.install
hooks.install: ## Install git hooks for conventional commits
	@git config core.hooksPath .githooks
	@echo "Git hooks installed. Commits will be validated for conventional format."

.PHONY: hooks.uninstall
hooks.uninstall: ## Uninstall git hooks
	@git config --unset core.hooksPath || true
	@echo "Git hooks uninstalled."

.PHONY: help
help: ## Show available make targets
	@echo "Available targets:"
	@grep -hE '^[a-zA-Z0-9_.-]+:.*##' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS=":.*##"} {printf "  %-18s %s\n", $$1, $$2}'
