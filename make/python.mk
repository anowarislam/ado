PYTHON ?= python3
VENV ?= .venv
PIP ?= $(VENV)/bin/pip
PYTEST ?= $(VENV)/bin/pytest
LAB_DIR ?= lab/py

.PHONY: py.env py.install py.test py.clean py.lint

py.env: ## Create Python virtual environment for lab work
	@if [ ! -d "$(VENV)" ]; then $(PYTHON) -m venv "$(VENV)"; fi

py.install: py.env ## Install lab package and requirements
	@"$(PIP)" install -e "$(LAB_DIR)[dev]"

py.test: py.env ## Run Python lab tests with pytest
	@if [ ! -x "$(PYTEST)" ]; then \
		echo "pytest not found in $(VENV); install it (e.g., make py.install)"; exit 1; \
	fi
	@if [ -d "$(LAB_DIR)" ]; then \
		"$(PYTEST)" "$(LAB_DIR)"; \
	else \
		echo "No lab/py directory present."; \
	fi

py.lint: py.env ## Lint Python lab code with ruff
	@if [ ! -x "$(VENV)/bin/ruff" ]; then \
		echo "ruff not found in $(VENV); run 'make py.install' first"; exit 1; \
	fi
	@"$(VENV)/bin/ruff" check "$(LAB_DIR)"

py.clean: ## Remove Python virtual environment
	@rm -rf "$(VENV)"
