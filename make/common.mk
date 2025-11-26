# ==============================================================================
# Common Variables and Utilities
# ==============================================================================
# Shared configuration for all Makefiles

# ------------------------------------------------------------------------------
# Shell Configuration
# ------------------------------------------------------------------------------
SHELL := /bin/bash
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

# ------------------------------------------------------------------------------
# Project Info
# ------------------------------------------------------------------------------
PROJECT_NAME := ado
PROJECT_ROOT := $(shell pwd)
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

# ------------------------------------------------------------------------------
# Directories
# ------------------------------------------------------------------------------
BIN_DIR := $(PROJECT_ROOT)
BUILD_DIR := $(PROJECT_ROOT)/dist
DOCS_DIR := $(PROJECT_ROOT)/docs
SITE_DIR := $(PROJECT_ROOT)/site

# ------------------------------------------------------------------------------
# Terminal Colors (if supported)
# ------------------------------------------------------------------------------
TERM_SUPPORTS_COLOR := $(shell tput colors 2>/dev/null || echo 0)
ifeq ($(shell test $(TERM_SUPPORTS_COLOR) -ge 8 && echo yes),yes)
    COLOR_RESET := $(shell tput sgr0)
    COLOR_BOLD := $(shell tput bold)
    COLOR_RED := $(shell tput setaf 1)
    COLOR_GREEN := $(shell tput setaf 2)
    COLOR_YELLOW := $(shell tput setaf 3)
    COLOR_BLUE := $(shell tput setaf 4)
    COLOR_CYAN := $(shell tput setaf 6)
else
    COLOR_RESET :=
    COLOR_BOLD :=
    COLOR_RED :=
    COLOR_GREEN :=
    COLOR_YELLOW :=
    COLOR_BLUE :=
    COLOR_CYAN :=
endif

# ------------------------------------------------------------------------------
# Output Helpers
# ------------------------------------------------------------------------------
define log_info
	@printf "$(COLOR_CYAN)▸$(COLOR_RESET) %s\n" $(1)
endef

define log_success
	@printf "$(COLOR_GREEN)✓$(COLOR_RESET) %s\n" $(1)
endef

define log_warn
	@printf "$(COLOR_YELLOW)⚠$(COLOR_RESET) %s\n" $(1)
endef

define log_error
	@printf "$(COLOR_RED)✗$(COLOR_RESET) %s\n" $(1)
endef

define log_header
	@printf "\n$(COLOR_BOLD)$(COLOR_BLUE)=== %s ===$(COLOR_RESET)\n" $(1)
endef

# ------------------------------------------------------------------------------
# Dependency Checks
# ------------------------------------------------------------------------------
define check_tool
	@command -v $(1) >/dev/null 2>&1 || { \
		printf "$(COLOR_RED)✗$(COLOR_RESET) Required tool '$(1)' not found. $(2)\n"; \
		exit 1; \
	}
endef

# Check Go is installed
.PHONY: _check-go
_check-go:
	$(call check_tool,go,Please install Go from https://go.dev/dl/)

# Check Python is installed
.PHONY: _check-python
_check-python:
	$(call check_tool,python3,Please install Python 3.12+)

# Check Docker is installed
.PHONY: _check-docker
_check-docker:
	$(call check_tool,docker,Please install Docker from https://docs.docker.com/get-docker/)

# Check Git is installed
.PHONY: _check-git
_check-git:
	$(call check_tool,git,Please install Git)

# ------------------------------------------------------------------------------
# Common Targets
# ------------------------------------------------------------------------------
.PHONY: version
version: ## Show project version info
	@echo "Project:    $(PROJECT_NAME)"
	@echo "Version:    $(VERSION)"
	@echo "Commit:     $(COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"

.PHONY: clean
clean: go.clean py.clean docs.clean ## Clean all build artifacts
	$(call log_success,"All build artifacts cleaned")

.PHONY: deps
deps: _check-go _check-python ## Verify all dependencies are installed
	$(call log_info,"Checking Go dependencies...")
	@go mod verify
	$(call log_success,"All dependencies verified")
