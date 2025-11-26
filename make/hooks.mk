# ==============================================================================
# Git Hooks and Pre-commit Targets
# ==============================================================================
# Git hooks management and pre-commit integration

# ------------------------------------------------------------------------------
# Variables
# ------------------------------------------------------------------------------
HOOKS_DIR := .githooks
PRECOMMIT ?= pre-commit

# ------------------------------------------------------------------------------
# Git Hooks Targets
# ------------------------------------------------------------------------------
.PHONY: hooks.install hooks.uninstall hooks.status

hooks.install: _check-git ## Install git hooks for conventional commits
	$(call log_info,"Installing git hooks...")
	@git config core.hooksPath $(HOOKS_DIR)
	$(call log_success,"Git hooks installed from $(HOOKS_DIR)/")
	@echo "  Commits will be validated for conventional commit format."

hooks.uninstall: ## Uninstall git hooks
	$(call log_info,"Uninstalling git hooks...")
	@git config --unset core.hooksPath || true
	$(call log_success,"Git hooks uninstalled")

hooks.status: ## Show git hooks status
	@echo "Hooks path: $$(git config core.hooksPath || echo 'default (.git/hooks)')"
	@echo ""
	@echo "Available hooks in $(HOOKS_DIR)/:"
	@ls -la $(HOOKS_DIR)/ 2>/dev/null || echo "  (none)"

# ------------------------------------------------------------------------------
# Pre-commit Targets
# ------------------------------------------------------------------------------
.PHONY: precommit.install precommit.uninstall precommit.run precommit.update

precommit.install: ## Install pre-commit framework hooks
	$(call log_info,"Installing pre-commit hooks...")
	@command -v $(PRECOMMIT) >/dev/null 2>&1 || { \
		echo "Installing pre-commit..."; \
		pip install pre-commit; \
	}
	@$(PRECOMMIT) install
	@$(PRECOMMIT) install --hook-type commit-msg
	$(call log_success,"Pre-commit hooks installed")

precommit.uninstall: ## Uninstall pre-commit hooks
	$(call log_info,"Uninstalling pre-commit hooks...")
	@$(PRECOMMIT) uninstall || true
	$(call log_success,"Pre-commit hooks uninstalled")

precommit.run: ## Run pre-commit on all files
	$(call log_info,"Running pre-commit on all files...")
	@$(PRECOMMIT) run --all-files

precommit.update: ## Update pre-commit hooks to latest versions
	$(call log_info,"Updating pre-commit hooks...")
	@$(PRECOMMIT) autoupdate
	$(call log_success,"Pre-commit hooks updated")
