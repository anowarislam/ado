# Style Guide Index

This directory contains the source of truth for code style, testing, and workflow practices in the `ado` project.

## Purpose

These guides ensure consistency across the two-track development model:
- **Go binary** (`cmd/ado/`, `internal/`) - Production code that ships
- **Python lab** (`lab/py/`) - R&D and prototyping area

## Available Guides

### [go-style.md](go-style.md) - Go Code & Testing
Source of truth for all Go code in `cmd/` and `internal/`:
- Code structure, naming, error handling, concurrency patterns
- Testing tenets: table-driven tests, determinism, coverage
- Logging, benchmarking, and performance considerations
- Real examples from this codebase

**Read this when:**
- Writing new Go commands or internal packages
- Adding tests to existing Go code
- Reviewing Go PRs
- Porting Python lab prototypes to Go

### [python-style.md](python-style.md) - Python Lab Code
Guidelines for experimental code in `lab/py/`:
- Lab-specific philosophy: correct but disposable
- Code tenets: type hints, error handling, minimal globals
- Testing with pytest: parametrize, fixtures, CLI testing
- Promotion path from Python to Go

**Read this when:**
- Prototyping new command logic in Python
- Writing tests for lab code
- Preparing to port Python code to Go
- Reviewing Python lab PRs

### [ci-style.md](ci-style.md) - CI/CD & Releases
Local workflow, GitHub Actions, and release hygiene:
- Make targets that mirror in CI
- Pre-commit hook recommendations
- Concrete GitHub Actions workflow examples
- Versioning, releases, and changelog practices

**Read this when:**
- Setting up or modifying GitHub Actions workflows
- Preparing a release
- Adding pre-commit hooks
- Debugging CI failures

### [docs-style.md](docs-style.md) - Documentation Standards
Writing and maintaining project documentation:
- Command specification format (`docs/commands/*.md`)
- Markdown conventions and structure
- README maintenance practices
- Keeping docs in sync with code

**Read this when:**
- Writing new command specifications
- Updating README or project docs
- Adding examples to documentation
- Ensuring docs match implementation

## When to Consult Multiple Guides

### Before Writing a New Command
1. **Check workflow**: [workflow.md](../workflow.md) → Determine if ADR needed
2. **Spec first**: [docs-style.md](docs-style.md) → Create `docs/commands/<command>.md`
3. **Prototype**: [python-style.md](python-style.md) → Experiment in `lab/py/`
4. **Implement**: [go-style.md](go-style.md) → Port to Go in `cmd/ado/`
5. **CI integration**: [ci-style.md](ci-style.md) → Ensure tests run in workflows

### Before Submitting a PR
1. **Code quality**: [go-style.md](go-style.md) or [python-style.md](python-style.md)
2. **Tests pass**: Both style guides + [ci-style.md](ci-style.md) local checks
3. **Docs updated**: [docs-style.md](docs-style.md) for command specs or README
4. **CI passing**: [ci-style.md](ci-style.md) for workflow requirements

### Before a Release
1. **All tests green**: [ci-style.md](ci-style.md) → Verify workflows
2. **Docs current**: [docs-style.md](docs-style.md) → README, command specs, changelog
3. **Version tagged**: [ci-style.md](ci-style.md) → SemVer conventions

## Cross-References

Each guide references the others where concepts overlap:
- **Go tests** mirror **CI checks** → [go-style.md](go-style.md) ↔ [ci-style.md](ci-style.md)
- **Python lab promotion** targets **Go patterns** → [python-style.md](python-style.md) ↔ [go-style.md](go-style.md)
- **Command specs** inform **implementation** → [docs-style.md](docs-style.md) ↔ [go-style.md](go-style.md)
- **Local workflow** mirrors **CI** → All guides ↔ [ci-style.md](ci-style.md)

## Quick Reference

| Task | Primary Guide | Secondary Guide |
|------|---------------|-----------------|
| Add Go command | go-style.md | docs-style.md (spec first) |
| Prototype in Python | python-style.md | go-style.md (target patterns) |
| Set up CI workflow | ci-style.md | go-style.md, python-style.md |
| Write command spec | docs-style.md | go-style.md (impl patterns) |
| Add Go tests | go-style.md | ci-style.md (CI integration) |
| Add Python tests | python-style.md | ci-style.md (CI integration) |
| Release new version | ci-style.md | docs-style.md (changelog) |

## Philosophy

All guides share these core principles:
- **Clarity over cleverness** - Readable code beats clever tricks
- **Explicit over implicit** - Surface dependencies and behavior
- **Deterministic by default** - Tests should be reliable and fast
- **Documentation-first** - Specs precede implementation
- **CI mirrors local** - Same commands locally and in GitHub Actions

## Maintenance

These style guides are living documents:
- Update guides when patterns evolve
- Add real examples from the codebase as it grows
- Keep cross-references accurate
- Ensure consistency across all four guides
