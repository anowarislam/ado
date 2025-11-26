# Contributing to ado

Thank you for your interest in contributing to `ado`! This document provides guidelines and information for contributors.

## Quick Links

- [Development Workflow](docs/workflow.md) - Issue → ADR → Spec → Implementation
- [Code Ownership](docs/code-ownership.md) - Review process and approvals
- [Code of Conduct](CODE_OF_CONDUCT.md)
- [Security Policy](SECURITY.md)
- [Detailed Contributing Guide](docs/contributing.md)
- [Style Guides](docs/style/)

## Getting Started

### Prerequisites

- **Go 1.23+** for CLI development
- **Python 3.12+** for lab experiments (optional)
- **Make** for build automation
- **Git** with hooks support

### Setup

```bash
# Clone the repository
git clone https://github.com/anowarislam/ado.git
cd ado

# Install git hooks (recommended)
make hooks.install

# Build and test
make go.build
make test
```

### Git Hooks

Run `make hooks.install` to enable local validation:

| Hook | When | What it Checks |
|------|------|----------------|
| `commit-msg` | Every commit | Conventional commit format |
| `pre-push` | Before push | Tests, coverage (80%), build |

**Skip hooks** (use sparingly): `git push --no-verify`

## Code Review Process

All PRs require review and approval from code owners before merging.

**How it works**:
1. When you open a PR, GitHub automatically requests reviews from code owners based on files changed
2. Code owners are defined in `.github/CODEOWNERS`
3. Branch protection requires approval from code owners before merge
4. Code owners aim to review within 48 business hours

**Current code owner**: @anowarislam reviews all changes

**Learn more**: See [code-ownership.md](docs/code-ownership.md) for detailed review process and FAQ.

## Development Workflow

This project uses a **three-phase workflow** for significant changes:

1. **Decision (ADR)** - For architectural changes, create an ADR first
2. **Specification** - Write spec before implementation (command or feature)
3. **Implementation** - Code, tests, docs

See [workflow.md](docs/workflow.md) for the complete guide with examples.

### 1. Create a Branch

```bash
# For ADRs
git checkout -b adr/NNNN-short-title

# For specs
git checkout -b spec/feature-name

# For implementation
git checkout -b feat/your-feature-name
# or
git checkout -b fix/your-bug-fix
```

### 2. Make Changes

- Follow the [Go Style Guide](docs/style/go-style.md)
- Write tests for new functionality
- Update documentation if needed
- For new commands: create spec first ([template](docs/commands/TEMPLATE.md))

### 3. Validate Locally

```bash
# Run all checks (mirrors CI)
make test          # Run Go + Python tests
make go.vet        # Lint Go code
make go.fmt        # Check formatting
make py.lint       # Lint Python (if changed)

# Check test coverage (80% minimum required)
make go.test.cover.check
```

#### Test Coverage Policy

- **Minimum threshold: 80%** - CI enforces this for all PRs
- Check coverage locally: `make go.test.cover`
- View HTML report: `make go.test.cover.html`
- New code should include tests that maintain or improve coverage

### 4. Commit with Conventional Commits

```bash
# Format: <type>[(scope)][!]: <description>
git commit -m "feat(cli): add new command"
git commit -m "fix(config): handle missing file"
git commit -m "docs: update README"
```

**Types:** `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`

### 5. Create Pull Request

```bash
git push origin your-branch
gh pr create --title "feat(cli): add new command" --body "Description"
```

## What We're Looking For

### Good First Issues

Look for issues labeled [`good first issue`](https://github.com/anowarislam/ado/labels/good%20first%20issue).

### Types of Contributions

- **Bug fixes** - Fix issues, improve error handling
- **Features** - New commands, flags, or functionality
- **Documentation** - Improve guides, add examples
- **Tests** - Increase coverage, add edge cases
- **Performance** - Optimize code, reduce binary size

## Pull Request Guidelines

1. **One concern per PR** - Keep PRs focused
2. **Write tests** - New features need tests
3. **Update docs** - Keep documentation current
4. **Follow style** - Use existing patterns
5. **Sign commits** - DCO sign-off appreciated
6. **Wait for code owner approval** - PRs require approval from code owners before merging

**Note**: All PRs must pass CI checks (tests, linting, coverage) AND receive code owner approval before merging to main.

## Questions?

- Open a [Discussion](https://github.com/anowarislam/ado/discussions)
- Check [existing issues](https://github.com/anowarislam/ado/issues)

---

For more detailed information, see the [full contributing guide](docs/contributing.md).
