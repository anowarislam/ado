# ==============================================================================
# Python Lab Targets
# ==============================================================================
# Python lab development, testing, and linting targets

# ------------------------------------------------------------------------------
# Variables
# ------------------------------------------------------------------------------
PYTHON ?= python3
VENV := $(PROJECT_ROOT)/.venv
VENV_BIN := $(VENV)/bin
PIP := $(VENV_BIN)/pip
PYTEST := $(VENV_BIN)/pytest
RUFF := $(VENV_BIN)/ruff
LAB_DIR := $(PROJECT_ROOT)/lab/py

# Python version check
PYTHON_MIN_VERSION := 3.12

# ------------------------------------------------------------------------------
# Environment Targets
# ------------------------------------------------------------------------------
.PHONY: py.env py.env.check py.install py.install.dev

py.env: ## Create Python virtual environment
	$(call log_info,"Creating Python virtual environment...")
	@if [ ! -d "$(VENV)" ]; then \
		$(PYTHON) -m venv "$(VENV)"; \
		$(PIP) install --upgrade pip; \
	fi
	$(call log_success,"Virtual environment ready: $(VENV)")

py.env.check: ## Check Python version
	$(call log_info,"Checking Python version...")
	@$(PYTHON) -c "import sys; v=sys.version_info; exit(0 if v >= (3,12) else 1)" || { \
		echo "$(COLOR_RED)Python $(PYTHON_MIN_VERSION)+ required$(COLOR_RESET)"; \
		$(PYTHON) --version; \
		exit 1; \
	}
	$(call log_success,"Python version OK: $$($(PYTHON) --version)")

py.install: py.env ## Install lab package and requirements
	$(call log_info,"Installing lab package...")
	@"$(PIP)" install -e "$(LAB_DIR)[dev]"
	$(call log_success,"Lab package installed")

py.install.dev: py.env ## Install development dependencies only
	$(call log_info,"Installing dev dependencies...")
	@"$(PIP)" install ruff pytest pytest-cov
	$(call log_success,"Dev dependencies installed")

# ------------------------------------------------------------------------------
# Test Targets
# ------------------------------------------------------------------------------
.PHONY: py.test py.test.verbose py.test.cover py.test.watch

py.test: py.env ## Run Python lab tests
	$(call log_info,"Running Python tests...")
	@if [ ! -x "$(PYTEST)" ]; then \
		echo "pytest not found. Run 'make py.install' first."; \
		exit 1; \
	fi
	@if [ -d "$(LAB_DIR)" ]; then \
		"$(PYTEST)" "$(LAB_DIR)" -q; \
	else \
		echo "No lab/py directory present."; \
	fi
	$(call log_success,"Python tests passed")

py.test.verbose: py.env ## Run Python tests with verbose output
	$(call log_info,"Running Python tests (verbose)...")
	@"$(PYTEST)" "$(LAB_DIR)" -v

py.test.cover: py.env ## Run Python tests with coverage
	$(call log_info,"Running Python tests with coverage...")
	@"$(PYTEST)" "$(LAB_DIR)" --cov="$(LAB_DIR)" --cov-report=term-missing

py.test.watch: py.env ## Run Python tests in watch mode
	$(call log_info,"Running Python tests in watch mode...")
	@"$(PYTEST)" "$(LAB_DIR)" --watch

# ------------------------------------------------------------------------------
# Lint Targets
# ------------------------------------------------------------------------------
.PHONY: py.lint py.lint.fix py.fmt py.type

py.lint: py.env ## Lint Python code with ruff
	$(call log_info,"Linting Python code...")
	@if [ ! -x "$(RUFF)" ]; then \
		echo "ruff not found. Run 'make py.install' first."; \
		exit 1; \
	fi
	@"$(RUFF)" check "$(LAB_DIR)"
	$(call log_success,"Python lint passed")

py.lint.fix: py.env ## Lint and auto-fix Python code
	$(call log_info,"Linting and fixing Python code...")
	@"$(RUFF)" check --fix "$(LAB_DIR)"
	$(call log_success,"Python lint fixes applied")

py.fmt: py.env ## Format Python code with ruff
	$(call log_info,"Formatting Python code...")
	@"$(RUFF)" format "$(LAB_DIR)"
	$(call log_success,"Python code formatted")

py.type: py.env ## Run type checking with mypy (if installed)
	$(call log_info,"Running type checks...")
	@command -v mypy >/dev/null 2>&1 || { \
		echo "mypy not installed. Skipping type checks."; \
		exit 0; \
	}
	@mypy "$(LAB_DIR)"
	$(call log_success,"Type checks passed")

# ------------------------------------------------------------------------------
# Cleanup Targets
# ------------------------------------------------------------------------------
.PHONY: py.clean py.clean.cache

py.clean: ## Remove Python virtual environment
	$(call log_info,"Cleaning Python environment...")
	@rm -rf "$(VENV)"
	@rm -rf "$(LAB_DIR)/__pycache__"
	@rm -rf "$(LAB_DIR)/.pytest_cache"
	@rm -rf "$(LAB_DIR)/*.egg-info"
	$(call log_success,"Python environment cleaned")

py.clean.cache: ## Remove Python cache files only
	$(call log_info,"Cleaning Python cache...")
	@find "$(LAB_DIR)" -type d -name "__pycache__" -exec rm -rf {} + 2>/dev/null || true
	@find "$(LAB_DIR)" -type f -name "*.pyc" -delete 2>/dev/null || true
	$(call log_success,"Python cache cleaned")

# ------------------------------------------------------------------------------
# Info Targets
# ------------------------------------------------------------------------------
.PHONY: py.info

py.info: py.env ## Show Python environment info
	@echo "Python: $$($(PYTHON) --version)"
	@echo "Venv: $(VENV)"
	@echo "Lab dir: $(LAB_DIR)"
	@echo ""
	@echo "Installed packages:"
	@"$(PIP)" list --format=freeze | head -20
