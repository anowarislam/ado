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

# Python Lab
make py.install          # Install lab package
make py.test             # Run pytest
make py.lint             # Run ruff

# All
make test                # Go + Python tests
make hooks.install       # Install commit hooks
```

## Architecture

```
cmd/ado/<command>/    → Cobra commands, each exports NewCommand()
internal/             → Shared packages (meta, config, ui)
lab/py/               → Python prototypes (promote to Go when stable)
```

**Command wiring**: `cmd/ado/root/root.go` registers all subcommands.

**Build metadata**: `internal/meta/info.go` - Version/Commit/BuildTime set via ldflags.

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

## References

- `docs/contributing.md` - Full commit/test/release conventions
- `docs/style/go-style.md` - Go code patterns with examples
- `docs/commands/*.md` - Command specifications
- `.goreleaser.yaml` - Build config and ldflags
