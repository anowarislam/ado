# Makefile Reference

This project uses a modular Makefile system organized under the `make/` directory. Each file handles a specific concern, making it easy to extend and maintain.

## Quick Start

```bash
# First-time setup
make setup

# Build the binary
make build

# Run all tests
make test

# Full CI pipeline locally
make ci

# Show all available targets
make help
```

## Structure

```
Makefile              # Main entry point with composite targets
make/
├── common.mk        # Shared variables, colors, logging helpers
├── go.mk            # Go build, test, lint targets
├── python.mk        # Python lab targets
├── docker.mk        # Container build and management
├── docs.mk          # MkDocs documentation
├── hooks.mk         # Git hooks and pre-commit
└── help.mk          # Organized help display
```

## Targets by Category

### Build

| Target | Description |
|--------|-------------|
| `build` | Build the project (alias for `go.build`) |
| `go.build` | Build the ado binary with version info |
| `go.build.all` | Build for all platforms via goreleaser |
| `go.install` | Install binary to `$GOPATH/bin` |

### Test

| Target | Description |
|--------|-------------|
| `test` | Run all tests (Go + Python) |
| `test.cover` | Run all tests with coverage |
| `go.test` | Run Go tests with race detector |
| `go.test.cover` | Run Go tests with coverage report |
| `go.test.verbose` | Run Go tests with verbose output |
| `go.test.race` | Run Go tests with race detector |
| `go.bench` | Run Go benchmarks |
| `py.test` | Run Python lab tests |
| `py.test.cover` | Run Python tests with coverage |
| `py.test.verbose` | Run Python tests verbose |
| `py.test.watch` | Run Python tests in watch mode |

### Lint

| Target | Description |
|--------|-------------|
| `lint` | Run all linters (Go + Python) |
| `fmt` | Format all code |
| `check` | Quick syntax/lint check (no tests) |
| `go.fmt` | Format Go sources |
| `go.fmt.check` | Check Go formatting (fails if not formatted) |
| `go.vet` | Run go vet on all packages |
| `go.lint` | Run golangci-lint (if installed) |
| `go.tidy` | Sync Go module dependencies |
| `py.lint` | Lint Python code with ruff |
| `py.lint.fix` | Auto-fix Python lint issues |
| `py.fmt` | Format Python code with ruff |
| `py.type` | Run mypy type checking (if installed) |

### Docker

| Target | Description |
|--------|-------------|
| `docker.build` | Build Docker image for current platform |
| `docker.build.multi` | Build multi-arch image (amd64 + arm64) |
| `docker.run` | Run container with `meta info` |
| `docker.push` | Push image to ghcr.io registry |
| `docker.clean` | Remove local Docker images |
| `docker.lint` | Lint Dockerfile with hadolint |

### Documentation

| Target | Description |
|--------|-------------|
| `docs.install` | Install MkDocs and dependencies |
| `docs.build` | Build documentation site |
| `docs.serve` | Serve docs locally at http://localhost:8000 |
| `docs.deploy` | Deploy to GitHub Pages |
| `docs.clean` | Clean built documentation |
| `docs.check` | Check documentation for errors |

### Git Hooks

| Target | Description |
|--------|-------------|
| `hooks.install` | Install git hooks for conventional commits |
| `hooks.uninstall` | Remove git hooks |
| `hooks.status` | Show hooks status |
| `precommit.install` | Install pre-commit framework |
| `precommit.uninstall` | Uninstall pre-commit hooks |
| `precommit.run` | Run pre-commit on all files |
| `precommit.update` | Update pre-commit hooks |

### Dependencies

| Target | Description |
|--------|-------------|
| `deps` | Verify all dependencies are installed |
| `go.deps` | Download Go dependencies |
| `go.deps.update` | Update Go dependencies |
| `go.deps.graph` | Show Go dependency graph |
| `py.install` | Install Python lab package |
| `py.install.dev` | Install Python dev dependencies only |

### Cleanup

| Target | Description |
|--------|-------------|
| `clean` | Clean all build artifacts |
| `go.clean` | Remove Go build artifacts |
| `py.clean` | Remove Python virtual environment |
| `py.clean.cache` | Remove Python cache files only |
| `docker.clean` | Remove local Docker images |
| `docs.clean` | Remove documentation site |

### Info

| Target | Description |
|--------|-------------|
| `help` | Show organized help message |
| `help.all` | Show all targets (unorganized) |
| `version` | Show project version info |
| `info` | Show project information |
| `py.info` | Show Python environment info |

## Composite Targets

These high-level targets combine multiple sub-targets:

```bash
# First-time project setup
make setup        # hooks.install + go.deps + py.install

# Run all validations
make validate     # lint + test + docs.build

# Full CI pipeline locally
make ci           # validate + build

# Quick syntax check (no tests)
make check        # go.fmt.check + go.vet

# Pre-commit checks
make pre-commit   # check + test
```

## Colored Output

The Makefile system includes colored terminal output for better readability:

- **Cyan** `▸` - Info messages (starting a task)
- **Green** `✓` - Success messages (task completed)
- **Yellow** `⚠` - Warning messages
- **Red** `✗` - Error messages
- **Blue** `===` - Section headers

Example output:

```
▸ Building ado...
✓ Built: /path/to/ado

▸ Running Go tests...
✓ All Go tests passed

=== Validation Complete ===
✓ All checks passed!
```

## Build Variables

The build system automatically captures version info:

| Variable | Source | Example |
|----------|--------|---------|
| `VERSION` | `git describe --tags` | `v1.0.2-3-g1234567` |
| `COMMIT` | `git rev-parse --short HEAD` | `1234567` |
| `BUILD_TIME` | Current UTC time | `2024-01-15T10:30:00Z` |

These are embedded into the binary via ldflags and accessible via `ado meta info`.

## Overriding Variables

You can override default variables:

```bash
# Use a different Go binary
GO=go1.23 make go.build

# Use a different Python
PYTHON=python3.12 make py.test

# Use a different Docker tag
DOCKER_TAG=v2.0.0 make docker.build

# Change docs port
DOCS_PORT=9000 make docs.serve
```

## Adding New Targets

To add new targets:

1. **Create a new `.mk` file** in `make/` for a new category
2. **Add it to the main Makefile** with `include make/newfile.mk`
3. **Use the logging helpers** from `common.mk`:

```makefile
# In make/newfile.mk
.PHONY: new.target
new.target: ## Description for help system
	$(call log_info,"Starting task...")
	@your-command-here
	$(call log_success,"Task completed")
```

The `## Description` comment makes the target appear in `make help`.

## CI Integration

The Makefile is designed to work locally and in CI:

```yaml
# In GitHub Actions
- name: Run tests
  run: make test

- name: Build
  run: make build

- name: Full validation
  run: make validate
```

## Troubleshooting

### "command not found"

Some targets check for required tools. Install missing dependencies:

```bash
# Go
brew install go  # or download from https://go.dev

# Python
brew install python@3.12

# Docker
brew install --cask docker

# Pre-commit
pip install pre-commit

# MkDocs
pip install mkdocs-material
```

### Tests fail in CI but pass locally

Ensure you're running with the same flags:

```bash
# Run with race detector (like CI)
make go.test.race
```

### Docker build fails

Check Docker is running:

```bash
docker version
make docker.lint  # Validate Dockerfile syntax
```
