# Contributing Guide

This document covers conventions for contributing to the `ado` project.

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

This project is **spec-driven**:

1. Create spec in `docs/commands/<command>.md` first
2. Prototype in `lab/py/` if logic is complex
3. Implement in Go under `cmd/ado/<command>/`
4. Export `NewCommand() *cobra.Command`
5. Wire into `cmd/ado/root/root.go`
6. Write table-driven tests

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
