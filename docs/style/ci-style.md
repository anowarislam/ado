# CI & Release Style Guide

Guidelines for local checks, GitHub Actions, and release hygiene.

## Philosophy
- Fast feedback locally; CI should mirror local commands and stay deterministic.
- Keep pipelines simple and explicit; prefer small, composable steps.
- Fail loudly on lint/test breaks; avoid silent ignores.
- Releases follow SemVer; changelog and docs should reflect user-visible changes.

## Local Workflow
- Use the Makefile entrypoints:
  - `make test` for the unified test run (Go + Python lab).
  - `make go.vet`, `make go.fmt`, `make go.tidy` for Go hygiene.
  - `make py.lint`, `make py.test` for Python lab.
- Keep build artifacts out of the repo; clean with `make go.clean` / `make py.clean` when needed.
- Run tests before commits; add regression tests with fixes.

## Git Hooks / Pre-commit
- Recommended hooks (not enforced yet):
  - gofmt/goimports on staged Go files.
  - go test ./... with repo-local `GOCACHE=$(pwd)/.gocache`.
  - Ruff/pytest for `lab/py` when Python files change.
  - Markdown lints optional; keep docs readable.
- Keep hooks fast; skip heavy network operations. Provide overrides for emergency commits.

## GitHub Actions

### Principles
- Mirror local Make targets; avoid bespoke scripts.
- Split jobs for Go and Python; run in parallel when possible.
- Cache Go build/test cache and pip wheels to speed runs; fall back gracefully.
- Pin action versions; avoid `@main` or `@master` (use `@v4`, `@v5`, etc.).
- Minimal secrets; prefer GitHub-provided tokens.

### Format Check Pattern

For Go, use this pattern to ensure gofmt compliance:

```bash
test -z "$(gofmt -l .)"
```

This fails if any files need formatting. Alternative verbose version:

```bash
if [ -n "$(gofmt -l .)" ]; then
  echo "The following files need gofmt:"
  gofmt -l .
  exit 1
fi
```

### Concrete Workflow Examples

#### Go Workflow

```yaml
name: Go
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - name: Check out code
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'
          cache: true

      - name: Format check
        run: test -z "$(gofmt -l .)"

      - name: Tidy check
        run: |
          go mod tidy
          git diff --exit-code go.mod go.sum

      - name: Vet
        run: make go.vet

      - name: Test
        run: make go.test

      - name: Build
        run: make go.build
```

#### Python Lab Workflow

```yaml
name: Python Lab
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - name: Check out code
        uses: actions/checkout@v4

      - name: Set up Python
        uses: actions/setup-python@v5
        with:
          python-version: '3.12'
          cache: 'pip'
          cache-dependency-path: lab/py/requirements.txt

      - name: Install dependencies
        run: make py.install

      - name: Lint
        run: make py.lint

      - name: Test
        run: make py.test
```

#### Combined Workflow (Parallel Jobs)

```yaml
name: CI
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  go:
    name: Go Tests
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
          cache: true
      - run: test -z "$(gofmt -l .)"
      - run: make go.vet
      - run: make go.test
      - run: make go.build

  python:
    name: Python Lab Tests
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-python@v5
        with:
          python-version: '3.12'
          cache: 'pip'
          cache-dependency-path: lab/py/requirements.txt
      - run: make py.install
      - run: make py.lint
      - run: make py.test
```

### Multi-Platform Testing

For binaries targeting multiple OS/architectures:

```yaml
name: Multi-Platform
on: [push, pull_request]

jobs:
  test:
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
        go-version: ['1.22']
    runs-on: ${{ matrix.os }}
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go-version }}
          cache: true
      - run: make go.test
      - run: make go.build
```

### Coverage Reporting

Integrate with codecov or coveralls:

```yaml
- name: Test with coverage
  run: go test -coverprofile=coverage.out -covermode=atomic ./...

- name: Upload coverage to Codecov
  uses: codecov/codecov-action@v4
  with:
    files: ./coverage.out
    flags: unittests
    name: codecov-ado
```

## Versioning & Releases

### Semantic Versioning
Follow SemVer strictly:
- **MAJOR** (v1.0.0 → v2.0.0): Incompatible changes (breaking API, CLI flags, config format).
- **MINOR** (v1.0.0 → v1.1.0): Backward-compatible features (new commands, new flags).
- **PATCH** (v1.0.0 → v1.0.1): Backward-compatible fixes (bug fixes, docs, performance).

Tag releases with `vX.Y.Z`; keep changelog entries and docs aligned with behavior.

### Release Automation with goreleaser

Use [goreleaser](https://goreleaser.com/) for building multi-platform binaries:

**.goreleaser.yaml example:**

```yaml
project_name: ado

builds:
  - id: ado
    binary: ado
    main: ./cmd/ado
    env:
      - CGO_ENABLED=0
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64
    ldflags:
      - -s -w
      - -X github.com/anowarislam/ado/internal/meta.Version={{.Version}}
      - -X github.com/anowarislam/ado/internal/meta.Commit={{.Commit}}
      - -X github.com/anowarislam/ado/internal/meta.BuildTime={{.Date}}

archives:
  - format: tar.gz
    name_template: "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}"
    format_overrides:
      - goos: windows
        format: zip

checksum:
  name_template: "checksums.txt"

release:
  github:
    owner: anowarislam
    name: ado
  draft: false
  prerelease: auto
```

**GitHub Actions release workflow:**

```yaml
name: Release
on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - name: Run goreleaser
        uses: goreleaser/goreleaser-action@v6
        with:
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### Changelog Practices

- Use conventional commits: `feat:`, `fix:`, `docs:`, `chore:` prefixes.
- Generate changelogs from git history or PR titles.
- Include representative CLI examples in release notes.
- Mention which tests were added or fixed.
- Document migration path for breaking changes.

**Example release notes format:**

```markdown
## v0.2.0 - 2024-01-15

### Features
- Added `ado meta env` to show resolved config paths (#23)
- Support for XDG_CONFIG_HOME in config resolution (#24)

### Fixes
- Fixed panic when config file is missing (#25)

### Breaking Changes
- None

### Examples
$ ado meta env
Config path: /Users/user/.config/ado/config.yaml
XDG_CONFIG_HOME: /Users/user/.config
```

### Artifact Handling

- Attach built binaries to GitHub releases.
- Include checksums (SHA256) for verification.
- Sign releases if distributing outside GitHub (GPG/cosign).
- Document installation instructions in README.

## Repository Settings

### Branch Protection for `main`
- Require status checks to pass (Go tests, Python tests).
- Require at least 1 approving review for PRs.
- Disallow force push and branch deletion.
- Require branches to be up to date before merging.

### PR and Issue Templates

**`.github/pull_request_template.md`:**

```markdown
## Summary
Brief description of changes.

## Type
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Checklist
- [ ] Tests added/updated and passing (`make test`)
- [ ] Docs updated (README, command specs, CLAUDE.md)
- [ ] Changelog updated or issue linked
- [ ] Ran `make go.fmt` and `make go.vet`
- [ ] Ran `make py.lint` if Python changed

## Testing
How was this tested? Include CLI examples.
```

**`.github/ISSUE_TEMPLATE/bug_report.md`:**

```markdown
---
name: Bug report
about: Report a bug in ado
---

## Bug Description
Clear description of the bug.

## Steps to Reproduce
1. Run `ado <command>`
2. Observe error

## Expected Behavior
What should happen.

## Actual Behavior
What actually happened (error message, output).

## Environment
Run `ado meta info` and paste output:
```

### Dependency Management

- Use Dependabot for automated dependency updates:

**`.github/dependabot.yml`:**

```yaml
version: 2
updates:
  - package-ecosystem: gomod
    directory: "/"
    schedule:
      interval: weekly
    labels:
      - dependencies
      - go

  - package-ecosystem: pip
    directory: "/lab/py"
    schedule:
      interval: weekly
    labels:
      - dependencies
      - python

  - package-ecosystem: github-actions
    directory: "/"
    schedule:
      interval: monthly
    labels:
      - dependencies
      - ci
```

## Bad Practices to Avoid
- Overloaded CI jobs that do everything; prefer focused steps.
- Ignoring lint failures or marking steps as non-blocking without rationale.
- Pushing untested changes; relying on CI to catch obvious failures.
- Unpinned action versions; reliance on external network or flaky resources without retries.
- Skipping local checks before pushing (`make test` should be routine).
- Silent failures in CI (always use `set -e` in scripts or check exit codes).
- Not caching dependencies (slow CI wastes developer time).
- Mixing multiple concerns in a single commit (makes revert/bisect harder).

## Related Guides
- See [go-style.md](go-style.md) for Go code patterns that CI validates.
- See [python-style.md](python-style.md) for Python lab patterns that CI lints/tests.
- See [docs-style.md](docs-style.md) for keeping release notes and changelogs in sync.
- See [README.md](README.md) for style guide index and cross-references.
