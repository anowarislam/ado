# ==============================================================================
# Help System
# ==============================================================================
# Organized help display for all targets

# ------------------------------------------------------------------------------
# Help Target
# ------------------------------------------------------------------------------
.PHONY: help
help: ## Show this help message
	@echo ""
	@echo "$(COLOR_BOLD)$(PROJECT_NAME)$(COLOR_RESET) - Development Commands"
	@echo ""
	@echo "$(COLOR_BOLD)Usage:$(COLOR_RESET)"
	@echo "  make $(COLOR_CYAN)<target>$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_BOLD)Quick Start:$(COLOR_RESET)"
	@echo "  make setup      $(COLOR_CYAN)# First-time setup$(COLOR_RESET)"
	@echo "  make build      $(COLOR_CYAN)# Build the binary$(COLOR_RESET)"
	@echo "  make test       $(COLOR_CYAN)# Run all tests$(COLOR_RESET)"
	@echo "  make ci         $(COLOR_CYAN)# Run full CI locally$(COLOR_RESET)"
	@echo ""
	@$(MAKE) --no-print-directory _help-section SECTION="Build" PATTERN="go\.(build|clean)|build$$"
	@$(MAKE) --no-print-directory _help-section SECTION="Test" PATTERN="(go|py)\.test|^test$$"
	@$(MAKE) --no-print-directory _help-section SECTION="Lint" PATTERN="(go|py)\.(vet|lint|fmt)|^lint$$"
	@$(MAKE) --no-print-directory _help-section SECTION="Validation" PATTERN="^(validate|ci|deps)$$"
	@$(MAKE) --no-print-directory _help-section SECTION="Docker" PATTERN="^docker\."
	@$(MAKE) --no-print-directory _help-section SECTION="Documentation" PATTERN="^docs\."
	@$(MAKE) --no-print-directory _help-section SECTION="Hooks" PATTERN="^(hooks|precommit)\."
	@$(MAKE) --no-print-directory _help-section SECTION="Cleanup" PATTERN="clean$$"
	@$(MAKE) --no-print-directory _help-section SECTION="Info" PATTERN="^(version|help)$$"
	@echo ""

# Internal: Print a help section
.PHONY: _help-section
_help-section:
	@echo "$(COLOR_BOLD)$(SECTION):$(COLOR_RESET)"
	@grep -hE "^[a-zA-Z0-9_.-]+:.*## " $(MAKEFILE_LIST) 2>/dev/null | \
		grep -E "$(PATTERN)" | \
		sort | \
		awk 'BEGIN {FS = ":.*## "}; {printf "  $(COLOR_CYAN)%-18s$(COLOR_RESET) %s\n", $$1, $$2}'
	@echo ""

# Alternative simple help
.PHONY: help.all
help.all: ## Show all targets (unorganized)
	@echo "$(COLOR_BOLD)All Available Targets:$(COLOR_RESET)"
	@echo ""
	@grep -hE "^[a-zA-Z0-9_.-]+:.*## " $(MAKEFILE_LIST) 2>/dev/null | \
		sort | \
		awk 'BEGIN {FS = ":.*## "}; {printf "  $(COLOR_CYAN)%-20s$(COLOR_RESET) %s\n", $$1, $$2}'
