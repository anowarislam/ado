# ==============================================================================
# Documentation Targets
# ==============================================================================
# MkDocs documentation build and serving

# ------------------------------------------------------------------------------
# Variables
# ------------------------------------------------------------------------------
MKDOCS ?= mkdocs
PIP ?= pip
DOCS_PORT ?= 8000
DOCS_DEPS := mkdocs-material mkdocs-minify-plugin pillow cairosvg

# ------------------------------------------------------------------------------
# Targets
# ------------------------------------------------------------------------------
.PHONY: docs.install docs.build docs.serve docs.deploy docs.clean docs.check

docs.install: ## Install MkDocs and dependencies
	$(call log_info,"Installing MkDocs dependencies...")
	@$(PIP) install $(DOCS_DEPS)
	$(call log_success,"MkDocs dependencies installed")

docs.build: _docs-prep ## Build documentation site
	$(call log_info,"Building documentation...")
	@$(MKDOCS) build --strict
	$(call log_success,"Documentation built in $(SITE_DIR)/")

docs.serve: _docs-prep ## Serve documentation locally (http://localhost:$(DOCS_PORT))
	$(call log_info,"Starting documentation server on http://localhost:$(DOCS_PORT)")
	@$(MKDOCS) serve --dev-addr localhost:$(DOCS_PORT)

docs.deploy: docs.build ## Deploy documentation to GitHub Pages
	$(call log_info,"Deploying documentation to GitHub Pages...")
	@$(MKDOCS) gh-deploy --force
	$(call log_success,"Documentation deployed")

docs.clean: ## Clean built documentation
	$(call log_info,"Cleaning documentation artifacts...")
	@rm -rf $(SITE_DIR)/
	@rm -f $(DOCS_DIR)/changelog.md
	$(call log_success,"Documentation artifacts cleaned")

docs.check: _docs-prep ## Check documentation for errors
	$(call log_info,"Checking documentation...")
	@$(MKDOCS) build --strict 2>&1 | grep -E "(WARNING|ERROR)" && exit 1 || true
	$(call log_success,"Documentation check passed")

# ------------------------------------------------------------------------------
# Internal Targets
# ------------------------------------------------------------------------------
.PHONY: _docs-prep
_docs-prep:
	@cp CHANGELOG.md $(DOCS_DIR)/changelog.md 2>/dev/null || true
