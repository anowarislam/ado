# Go Style & Testing Guide

Single source of truth for Go code and tests in this repo.

## Philosophy
- Clarity over cleverness; obvious control flow beats micro-optimizations.
- Simplicity and readability first; future readers matter more than brevity.
- Maintainability: small, cohesive units that are easy to change.
- Explicit beats implicit: surface dependencies; avoid magic globals.
- Deterministic by default; hidden state and time/rand coupling are bugs.
- Tests are executable specs of behavior and contracts.

## Code Tenets
- Structure
  - Command wiring under `cmd/ado/<command>`; shared logic under `internal/<domain>`.
  - Keep packages cohesive; avoid circular deps and `util` grab-bags.
  - Prefer dependency injection (constructors/functions) over globals.
- Naming & style
  - MixedCaps, concise, descriptive. Exported symbols need doc comments.
  - Accept interfaces, return concrete types; keep interfaces tiny and defined at use sites.
- Error handling
  - Check errors immediately; wrap with context (`fmt.Errorf("context: %w", err)`).
  - Do not panic in libraries; only `main` may panic on unrecoverable setup failures.
  - Distinguish user errors from system failures; keep messages actionable.
- Concurrency
  - Respect `context.Context` for cancellation/timeouts; avoid leaking goroutines.
  - Share memory by communicating; keep synchronization simple and explicit.
- I/O & state
  - Avoid hidden configuration; pass config/deps explicitly.
  - Be conservative with third-party deps; standard library first, then well-maintained libs (e.g., Cobra).
- Logging
  - Use structured logging (`log/slog` in stdlib, or logrus/zap for richer features).
  - Log errors at boundaries (handlers, main loops), not deep in call stacks.
  - Reserve `log.Fatal` for main package only; libraries return errors.
  - Use appropriate levels: DEBUG for dev tracing, INFO for user-visible events, ERROR for failures.
- Package documentation
  - Every package needs a package comment: `// Package meta provides build introspection.`
  - Place package comment in a file matching the package name or in `doc.go`.
  - Export godocs for all public APIs; unexported symbols can have inline comments.
  - Use complete sentences; start with the symbol name: `// BuildInfo contains version metadata.`
- Performance
  - Optimize only with benchmarks proving the need (avoid premature optimization).
  - Use `go test -bench=. -benchmem` for measurement; compare before/after with benchstat.
  - Profile with `go tool pprof` when bottlenecks matter; document tradeoffs in comments.
  - Prefer clarity unless profiling shows measurable impact on user-facing latency.

## Test Tenets
- Contracts first
  - Table-driven tests for scenarios and edge cases; use subtests (`t.Run`) to group.
  - Keep fixtures/goldens small and readable; prefer unit tests near code.
- Determinism
  - No uncontrolled time/rand; seed or mock. Avoid real networks/files unless scoped to temp dirs.
- Assertions
  - Plain `testing` package; clear fail messages with expected vs. actual.
  - One logical expectation per test case; mark helpers with `t.Helper()`.
- Maintenance
  - Tests should be fast and reliable; gate slow/flaky cases behind build tags or skips with rationale.
  - Add regression tests alongside fixes.
- Coverage
  - Aim for high coverage of public APIs and critical paths; 100% is not required.
  - Use `go test -cover` and `go tool cover -html=coverage.out` for reports.
  - Focus on meaningful assertions, not coverage percentage gaming.
  - Uncovered code should be unreachable or explicitly documented as out-of-scope.
- Test organization
  - Unit tests live next to code (`*_test.go` in same package or `*_test.go` in `<pkg>_test` package).
  - Integration tests can live in `tests/` at repo root if they import multiple packages.
  - Use build tags for slow/external tests: `//go:build integration` and run with `go test -tags=integration`.
  - Keep fast unit tests as default; CI should run them on every commit.

## Examples
- Good
  - Clear error propagation:
    ```go
    data, err := loader.Load(ctx, path)
    if err != nil {
    	return fmt.Errorf("load config %q: %w", path, err)
    }
    ```
  - Small interface defined at use-site:
    ```go
    type Clock interface {
    	Now() time.Time
    }
    ```
  - Table-driven test:
    ```go
    func TestNormalize(t *testing.T) {
    	tests := []struct {
    		name string
    		in   string
    		want string
    	}{
    		{"trim", "  hi  ", "hi"},
    		{"lower", "HI", "hi"},
    	}
    	for _, tt := range tests {
    		t.Run(tt.name, func(t *testing.T) {
    			if got := normalize(tt.in); got != tt.want {
    				t.Fatalf("normalize(%q) = %q, want %q", tt.in, got, tt.want)
    			}
    		})
    	}
    }
    ```
  - Concurrency with cancellation:
    ```go
    go func(ctx context.Context) {
    	defer close(done)
    	for {
    		select {
    		case <-ctx.Done():
    			return
    		default:
    			doWork()
    		}
    	}
    }(ctx)
    ```
  - Channel pattern (worker pool):
    ```go
    func processJobs(ctx context.Context, jobs <-chan Job, results chan<- Result) {
    	for {
    		select {
    		case <-ctx.Done():
    			return
    		case job, ok := <-jobs:
    			if !ok {
    				return // Channel closed
    			}
    			results <- process(job)
    		}
    	}
    }
    ```
  - Structured logging:
    ```go
    import "log/slog"

    func processConfig(path string) error {
    	slog.Debug("processing config", "path", path)
    	data, err := os.ReadFile(path)
    	if err != nil {
    		slog.Error("failed to read config", "path", path, "error", err)
    		return fmt.Errorf("read config: %w", err)
    	}
    	slog.Info("config loaded", "path", path, "size", len(data))
    	return nil
    }
    ```
  - Package documentation:
    ```go
    // Package config handles configuration file discovery and loading.
    //
    // Config files are searched in XDG_CONFIG_HOME, then fallback locations.
    // See ResolveConfigPath for the complete search order.
    package config
    ```
  - Build constraints for platform-specific code:
    ```go
    //go:build linux || darwin

    package platform

    func getConfigDir() string {
    	return "/etc/ado"
    }
    ```
- Bad
  - Ignoring errors:
    ```go
    data, _ := loader.Load(path) // swallowed error
    process(data)
    ```
  - God interface:
    ```go
    type Service interface {
    	Start() error
    	Stop() error
    	Restart() error
    	Status() (string, error)
    	Logs() ([]byte, error)
    }
    ```
  - Panicking on user input:
    ```go
    func parse(path string) Config {
    	return mustLoad(path) // panics on error
    }
    ```
  - Flaky time-based test:
    ```go
    func TestTimeout(t *testing.T) {
    	time.Sleep(100 * time.Millisecond)
    	if !done {
    		t.Fatal("not done")
    	}
    }
    ```
  - Logging in libraries instead of returning errors:
    ```go
    func loadConfig(path string) (*Config, error) {
    	data, err := os.ReadFile(path)
    	if err != nil {
    		log.Printf("failed to read: %v", err) // Bad: library shouldn't log
    		return nil, err
    	}
    	return parse(data)
    }
    ```

## Real Examples from This Codebase

### Table-driven tests
See `internal/config/paths_test.go:8-42` for `TestDefaultSearchPaths` and `TestResolveConfigPath`.

### Build metadata pattern
See `internal/meta/info.go:8-35` for `BuildInfo` struct and `CurrentBuildInfo()` function.
Variables like `Version`, `Commit`, `BuildTime` are set via ldflags at build time.

### Command wiring
See `cmd/ado/root/root.go:14-38` for root command setup with Cobra.
Each command (`echo`, `meta`) exports `NewCommand() *cobra.Command` and registers in `AddCommand()`.

### Error wrapping
See `internal/config/paths.go:26-38` for `ResolveConfigPath` error handling and path resolution logic.

### Config resolution
See `internal/config/paths.go:8-23` for `DefaultSearchPaths` implementing XDG Base Directory pattern.

## Related Guides
- See [ci-style.md](ci-style.md) for GitHub Actions workflows that validate these patterns.
- See [python-style.md](python-style.md) for lab prototyping; promote to Go using these patterns.
- See [docs-style.md](docs-style.md) for command spec format; implementation should match specs.
- See [README.md](README.md) for style guide index and cross-references.
