# ADR-0002: Structured Logging

| Metadata | Value |
|----------|-------|
| **Status** | Accepted |
| **Date** | 2025-11-26 |
| **Author(s)** | @anowarislam |
| **Issue** | #37 |
| **Related ADRs** | N/A |

## Context

The `ado` CLI currently uses ad-hoc `fmt.Println` and `fmt.Fprintf` calls for output. This approach has several limitations:

- **No log levels**: Cannot distinguish debug output from errors
- **Not structured**: Cannot parse output in automated pipelines
- **No timestamps**: Difficult to correlate events
- **Inconsistent format**: Different commands format output differently
- **No configuration**: Cannot control verbosity at runtime

As `ado` grows, we need a consistent logging infrastructure that:

1. Supports log levels (Debug, Info, Warn, Error)
2. Provides structured output (JSON) for machine consumption
3. Defaults to human-readable output in terminals
4. Is configurable via config file and CLI flags
5. Maintains minimal dependencies

## Decision

**We will use Go's standard library `log/slog` package** for structured logging in `ado`.

Key implementation details:

1. **Logger Interface**: Define a thin interface in `internal/logging` wrapping `slog`
2. **Default Handler**: Text format for TTY, JSON for non-TTY (auto-detect)
3. **Configuration**: Support via `~/.config/ado/config.yaml` and `--log-level` flag
4. **Levels**: Debug, Info, Warn, Error (matching slog's built-in levels)
5. **Context Fields**: Support structured fields for machine-readable logs

```go
// Usage example
log := logging.New(cfg)
log.Info("validating config", "path", configPath)
log.Debug("parsed yaml", "keys", len(rawMap))
log.Error("validation failed", "error", err)
```

## Consequences

### Positive

- **Zero dependencies**: `log/slog` is in Go's stdlib (1.21+), no external packages needed
- **Official solution**: Maintained by the Go team, well-documented, stable API
- **Structured by design**: Key-value attributes are first-class citizens
- **Flexible handlers**: Easy to switch between text/JSON output
- **Future-proof**: As the official Go logging solution, it will continue to evolve
- **Performance**: Designed to be allocation-efficient

### Negative

- **Go 1.21+ required**: Already satisfied (ado requires Go 1.23+)
- **Less features than third-party**: No built-in log rotation, file output requires custom handler
- **Learning curve**: Developers familiar with zerolog/zap may need adjustment

### Neutral

- **API style**: Uses `slog.Info("msg", "key", value)` vs `zerolog`'s fluent `.Str("key", value).Msg("msg")`
- **Handler system**: Custom handlers can be added for advanced use cases

## Alternatives Considered

### Alternative 1: zerolog

**Description:** Zero-allocation JSON logger, extremely fast, fluent API.

```go
log.Info().Str("path", configPath).Msg("validating config")
```

**Why not chosen:**
- Adds external dependency
- Performance gains not meaningful for CLI (not high-throughput server)
- Fluent API is more verbose than slog's variadic approach
- slog is now the official Go solution

### Alternative 2: zap

**Description:** Uber's structured, leveled logging library. Battle-tested in production.

```go
logger.Info("validating config", zap.String("path", configPath))
```

**Why not chosen:**
- Adds external dependency
- More complex API (need to call `zap.String()`, `zap.Int()`, etc.)
- Designed for server-side high-performance logging
- Overkill for a CLI tool

### Alternative 3: Keep using fmt.Println

**Description:** Continue with current approach, no structured logging.

**Why not chosen:**
- Cannot implement log levels
- Cannot provide JSON output for pipelines
- Cannot configure verbosity
- Does not solve the original problem

## Implementation Notes

After approval, implement in two phases:

1. **Spec**: `docs/features/01-structured-logging.md`
   - Define Logger interface
   - Specify configuration options
   - Document output formats

2. **Implementation**: `internal/logging/`
   - `logger.go`: Logger interface and implementation
   - `handler.go`: Custom handlers (auto-detect TTY)
   - `config.go`: Configuration integration
   - Tests for all components

Integration with existing code:
- Add `logging.Logger` to command context
- Gradually replace `fmt.Println` with structured log calls
- Preserve user-facing output (results) separate from logs

## References

- [Go slog package documentation](https://pkg.go.dev/log/slog)
- [Go blog: Structured Logging with slog](https://go.dev/blog/slog)
- [Issue #37: Feature: Structured Logging](https://github.com/anowarislam/ado/issues/37)
- [zerolog](https://github.com/rs/zerolog) - Alternative considered
- [zap](https://github.com/uber-go/zap) - Alternative considered
