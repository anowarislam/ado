# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`ado` is a CLI binary (Go) with a Python lab for prototyping. Go ships to users; Python is for R&D only.

- **Module**: `github.com/anowarislam/ado`
- **Go**: 1.23+ with Cobra CLI framework
- **Python Lab**: `lab/py/` (not required at runtime)

## Quick Commands

```bash
# Go
make go.build            # Build ./ado binary
make go.test             # Run tests
make go.vet              # Run go vet
go test -v -run TestName ./internal/config/...  # Single test

# Test Coverage (80% minimum enforced by CI)
make go.test.cover       # Run tests with coverage report
make go.test.cover.check # Verify 80% threshold met
make go.test.cover.html  # Generate HTML coverage report

# Python Lab
make py.install          # Install lab package
make py.test             # Run pytest
make py.lint             # Run ruff

# Validation (mirrors CI)
make test                # Go + Python tests
make validate            # Lint + test + docs build
make ci                  # Full CI pipeline locally
make docker.test         # Test GoReleaser Dockerfile (catches release issues)

# Setup
make hooks.install       # Install commit hooks
make help                # Show all available targets
```

## Architecture

```
cmd/ado/<command>/    → Cobra commands, each exports NewCommand()
internal/             → Shared packages (meta, config, ui)
lab/py/               → Python prototypes (promote to Go when stable)
make/                 → Makefile modules (go.mk, python.mk, docker.mk, etc.)
```

**Command wiring**: `cmd/ado/root/root.go` registers all subcommands via `AddCommand()`.

**Build metadata**: `internal/meta/info.go` - Version/Commit/BuildTime set via ldflags in `.goreleaser.yaml`.

**Config resolution**: `internal/config/paths.go` - XDG_CONFIG_HOME → ~/.config/ado → ~/.ado

## Commit Format

```
<type>[(scope)][!]: <description>
```
Types: `feat|fix|docs|style|refactor|perf|test|build|ci|chore`

Breaking changes: add `!` (e.g., `feat!:`)

## Key Patterns

- **Table-driven tests** with `t.Run` subtests
- **Error wrapping**: `fmt.Errorf("context: %w", err)`
- **Spec-driven**: Check `docs/commands/*.md` before implementing

## Release System

Fully automated via release-please + GoReleaser. PRs with conventional commits → Release PR created → merge triggers build + Docker + signing.

## Development Workflow

Three-phase workflow: **Issue → ADR (if needed) → Spec → Implementation**

### Quick Decision Tree

- **Architectural change?** → ADR first (`docs/adr/`)
- **New command?** → Spec first (`docs/commands/`)
- **New internal feature?** → Spec first (`docs/features/`)
- **Bug fix?** → Direct to implementation

### Key Directories

```
docs/adr/       → Architecture Decision Records (why)
docs/features/  → Non-command feature specs (what)
docs/commands/  → CLI command specs (what)
```

### PR Sequence for New Features

1. **PR 1 (ADR)**: `docs(adr): NNNN - title` (if architectural)
2. **PR 2 (Spec)**: `docs(spec): [command|feature] name`
3. **PR 3 (Code)**: `feat(scope): description`

See `docs/workflow.md` for complete guide.

## References

- `docs/workflow.md` - Development workflow (ADR → Spec → Implementation)
- `docs/adr/` - Architecture Decision Records
- `docs/features/` - Non-command feature specifications
- `docs/commands/*.md` - Command specifications
- `docs/contributing.md` - Commit/test/release conventions
- `docs/style/go-style.md` - Go code patterns with examples
- `docs/release.md` - Release automation and supply chain security
- `.goreleaser.yaml` - Build config, ldflags, Docker images
- `SECURITY.md` - Artifact/container verification instructions
