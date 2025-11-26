# Contributing Guide

This document covers conventions for contributing to the `ado` project.

## Development Workflow

This project uses a three-phase development process: **Issue → ADR → Spec → Implementation**

| Change Type | Path |
|-------------|------|
| Architectural change | ADR first (`docs/adr/`) |
| New command | Spec first (`docs/commands/`) |
| New internal feature | Spec first (`docs/features/`) |
| Bug fix | Direct to implementation |

**PR sequence for significant features:**
1. **PR 1 (ADR)**: `docs(adr): NNNN - title` (if architectural)
2. **PR 2 (Spec)**: `docs(spec): [command\|feature] name`
3. **PR 3 (Code)**: `feat(scope): description`

See [workflow.md](workflow.md) for the complete guide with examples.

## Commit Conventions

This project uses [Conventional Commits](https://www.conventionalcommits.org/) enforced by git hooks and CI.

```bash
# Install hooks (required for contributors)
make hooks.install
```

### Commit Format

```
<type>[(scope)][!]: <description>

[optional body]

[optional footer(s)]
```

| Type | Description | Version Bump |
|------|-------------|--------------|
| `feat` | New feature | Minor |
| `fix` | Bug fix | Patch |
| `docs` | Documentation only | None |
| `style` | Code style (formatting) | None |
| `refactor` | Code refactoring | None |
| `perf` | Performance improvement | Patch |
| `test` | Adding/updating tests | None |
| `build` | Build system/dependencies | None |
| `ci` | CI configuration | None |
| `chore` | Other changes | None |

**Breaking changes**: Add `!` after type/scope (e.g., `feat!:` or `feat(api)!:`) - bumps major version.

### Examples

```bash
git commit -m "feat: add export command for metrics"
git commit -m "fix(config): resolve path resolution on Windows"
git commit -m "feat!: change CLI flag syntax for meta command"
```

## Release Workflow

Releases are automated via [release-please](https://github.com/googleapis/release-please):

1. Merge PRs with conventional commit titles to `main`
2. release-please creates/updates a "Release PR" with version bump + CHANGELOG
3. Merge the Release PR → creates GitHub release with tag
4. `goreleaser.yml` workflow triggers on release publish → builds multi-platform binaries

Workflows: `.github/workflows/release-please.yml` (versioning) and `.github/workflows/goreleaser.yml` (builds)

## Testing Guidelines

### Go Tests

- Use table-driven tests with subtests (`t.Run`)
- Keep tests next to code (`*_test.go` files)
- Tests run in repo-local cache (`.gocache`)
- Mark test helpers with `t.Helper()`
- See `docs/style/go-style.md` for detailed tenets

Example pattern:
```go
func TestResolveConfigPath(t *testing.T) {
    tests := []struct {
        name     string
        explicit string
        homeDir  string
        want     string
    }{
        // test cases...
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, _ := ResolveConfigPath(tt.explicit, tt.homeDir)
            if got != tt.want {
                t.Errorf("got %q, want %q", got, tt.want)
            }
        })
    }
}
```

### Python Lab Tests

- Use pytest with parametrization
- Keep tests in `lab/py/tests/`
- See `docs/style/python-style.md` for style guide

## Adding New Commands

This project is **spec-driven** with a structured workflow:

1. **Create issue** using the command proposal template (`.github/ISSUE_TEMPLATE/command_proposal.md`)
2. **Create spec** in `docs/commands/<command>.md` using [TEMPLATE.md](commands/TEMPLATE.md)
3. Prototype in `lab/py/` if logic is complex (optional)
4. Implement in Go under `cmd/ado/<command>/`
5. Export `NewCommand() *cobra.Command`
6. Wire into `cmd/ado/root/root.go`
7. Write table-driven tests

For commands requiring architectural changes (new patterns, new dependencies), create an ADR first. See [workflow.md](workflow.md) for decision criteria.

## Code Constraints

### Error Handling

- Check all errors immediately
- Wrap with context: `fmt.Errorf("context: %w", err)`
- No panics in library code (only `main` may panic on setup failure)

### Code Organization

- Keep packages single-purpose and cohesive
- Avoid circular dependencies and `util` grab-bags
- Accept interfaces, return concrete types
- Define small interfaces at use-sites

### Determinism

- No uncontrolled time/rand in tests
- Avoid real network/file access (use temp dirs)
- Pass dependencies explicitly (no hidden globals)

See `docs/style/go-style.md` for examples.
