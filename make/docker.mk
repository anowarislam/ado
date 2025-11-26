# ==============================================================================
# Docker Targets
# ==============================================================================
# Container build and management targets

# ------------------------------------------------------------------------------
# Variables
# ------------------------------------------------------------------------------
DOCKER ?= docker
DOCKER_IMAGE := ghcr.io/anowarislam/ado
DOCKER_TAG ?= $(VERSION)
DOCKERFILE := Dockerfile
DOCKERFILE_GORELEASER := goreleaser.Dockerfile

# Build arguments
DOCKER_BUILD_ARGS := \
	--build-arg VERSION=$(VERSION) \
	--build-arg COMMIT=$(COMMIT) \
	--build-arg BUILD_TIME=$(BUILD_TIME)

# ------------------------------------------------------------------------------
# Targets
# ------------------------------------------------------------------------------
.PHONY: docker.build docker.build.multi docker.run docker.push docker.clean docker.lint

docker.build: _check-docker ## Build Docker image for current platform
	$(call log_info,"Building Docker image...")
	@$(DOCKER) build \
		$(DOCKER_BUILD_ARGS) \
		-t $(DOCKER_IMAGE):$(DOCKER_TAG) \
		-t $(DOCKER_IMAGE):latest \
		-f $(DOCKERFILE) \
		.
	$(call log_success,"Image built: $(DOCKER_IMAGE):$(DOCKER_TAG)")

docker.build.multi: _check-docker ## Build multi-arch Docker image (requires buildx)
	$(call log_info,"Building multi-arch Docker image...")
	@$(DOCKER) buildx build \
		$(DOCKER_BUILD_ARGS) \
		--platform linux/amd64,linux/arm64 \
		-t $(DOCKER_IMAGE):$(DOCKER_TAG) \
		-t $(DOCKER_IMAGE):latest \
		-f $(DOCKERFILE) \
		.
	$(call log_success,"Multi-arch image built: $(DOCKER_IMAGE):$(DOCKER_TAG)")

docker.run: docker.build ## Run Docker container
	$(call log_info,"Running container...")
	@$(DOCKER) run --rm $(DOCKER_IMAGE):$(DOCKER_TAG) meta info

docker.push: _check-docker ## Push Docker image to registry
	$(call log_info,"Pushing image to registry...")
	@$(DOCKER) push $(DOCKER_IMAGE):$(DOCKER_TAG)
	@$(DOCKER) push $(DOCKER_IMAGE):latest
	$(call log_success,"Image pushed: $(DOCKER_IMAGE):$(DOCKER_TAG)")

docker.clean: ## Remove local Docker images
	$(call log_info,"Removing local images...")
	@$(DOCKER) rmi $(DOCKER_IMAGE):$(DOCKER_TAG) 2>/dev/null || true
	@$(DOCKER) rmi $(DOCKER_IMAGE):latest 2>/dev/null || true
	$(call log_success,"Docker images cleaned")

docker.lint: _check-docker ## Lint Dockerfile with hadolint
	$(call log_info,"Linting Dockerfile...")
	@command -v hadolint >/dev/null 2>&1 || { \
		$(DOCKER) run --rm -i hadolint/hadolint < $(DOCKERFILE); \
		exit 0; \
	}
	@hadolint $(DOCKERFILE)
	$(call log_success,"Dockerfile lint passed")
