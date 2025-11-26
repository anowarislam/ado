# Feature: Structured Logging

| Metadata | Value |
|----------|-------|
| **ADR** | [ADR-0002](../adr/0002-structured-logging.md) |
| **Status** | Draft |
| **Issue** | #37 |
| **Author(s)** | @anowarislam |

## Overview

Structured logging infrastructure using Go's `log/slog` package, providing leveled logging with configurable output formats for both human readability and machine consumption.

## Motivation

Why is this feature needed? What problem does it solve?

- **Pain point**: Current ad-hoc `fmt.Println` calls provide no log levels, timestamps, or structure
- **Who benefits**:
  - Developers debugging issues (log levels)
  - Operators monitoring in production (structured JSON)
  - CI pipelines parsing output (machine-readable)
- **Without this**: Inconsistent output, difficult debugging, unusable in automated systems

## Specification

### Behavior

1. **Logger Initialization**
   - Create logger from config at startup
   - Auto-detect TTY for format selection (text vs JSON)
   - Respect `--log-level` flag override

2. **Log Levels** (ascending severity)
   - `Debug`: Verbose debugging information
   - `Info`: General operational messages (default)
   - `Warn`: Potential issues that don't stop execution
   - `Error`: Failures that affect functionality

3. **Output Formats**
   - **Text** (default for TTY): Human-readable, colored output
   - **JSON** (default for non-TTY): Machine-parseable structured logs

4. **Structured Fields**
   - All log calls support key-value pairs
   - Keys are strings, values can be any type
   - Fields are formatted based on output mode

### Configuration

```yaml
# Example config in ~/.config/ado/config.yaml
version: 1
logging:
  level: info      # debug, info, warn, error
  format: auto     # auto, text, json
  output: stderr   # stderr, stdout
```

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `logging.level` | string | `info` | Minimum log level to output |
| `logging.format` | string | `auto` | Output format (auto detects TTY) |
| `logging.output` | string | `stderr` | Output destination |

### CLI Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--log-level` | string | from config | Override log level |

### API/Interface

```go
// Package logging provides structured logging for ado.
package logging

import "log/slog"

// Logger provides structured, leveled logging.
type Logger interface {
    // Debug logs at debug level with optional structured fields.
    Debug(msg string, args ...any)

    // Info logs at info level with optional structured fields.
    Info(msg string, args ...any)

    // Warn logs at warn level with optional structured fields.
    Warn(msg string, args ...any)

    // Error logs at error level with optional structured fields.
    Error(msg string, args ...any)

    // With returns a new Logger with the given attributes.
    With(args ...any) Logger

    // Handler returns the underlying slog.Handler.
    Handler() slog.Handler
}

// Config holds logging configuration.
type Config struct {
    Level  string // debug, info, warn, error
    Format string // auto, text, json
    Output string // stderr, stdout
}

// New creates a new Logger from the given configuration.
func New(cfg Config) Logger

// NewFromConfig creates a Logger from the application config.
func NewFromConfig(appCfg *config.Config) Logger

// Default returns the default logger (info level, auto format, stderr).
func Default() Logger
```

### File Locations

| Purpose | Path |
|---------|------|
| Main implementation | `internal/logging/logger.go` |
| Configuration | `internal/logging/config.go` |
| Handlers | `internal/logging/handler.go` |
| Tests | `internal/logging/*_test.go` |

## Examples

### Example 1: Basic Usage

```go
// In a command
log := logging.New(logging.Config{Level: "info", Format: "auto"})

log.Info("validating config", "path", "/home/user/.config/ado/config.yaml")
// Text output (TTY):
// INFO validating config path=/home/user/.config/ado/config.yaml
//
// JSON output (non-TTY):
// {"level":"INFO","msg":"validating config","path":"/home/user/.config/ado/config.yaml"}

log.Debug("parsed yaml", "keys", 3)
// (not shown at info level)

log.Error("validation failed", "error", "missing version")
// Text output:
// ERROR validation failed error="missing version"
```

### Example 2: Using With() for Context

```go
// Add persistent context
cmdLog := log.With("command", "config", "subcommand", "validate")

cmdLog.Info("starting validation")
// {"level":"INFO","msg":"starting validation","command":"config","subcommand":"validate"}

cmdLog.Debug("checking file", "path", path)
// {"level":"DEBUG","msg":"checking file","command":"config","subcommand":"validate","path":"..."}
```

### Example 3: CLI Flag Override

```bash
# Run with debug logging
ado --log-level debug config validate

# Output includes debug messages:
# DEBUG resolving config path sources=[XDG_CONFIG_HOME, ~/.config/ado, ~/.ado]
# DEBUG found config file path=/home/user/.config/ado/config.yaml
# INFO validating config path=/home/user/.config/ado/config.yaml
```

### Example 4: Configuration File

```yaml
# ~/.config/ado/config.yaml
version: 1
logging:
  level: debug
  format: json
  output: stderr
```

```bash
ado config validate
# All output in JSON format, including debug messages
```

## Edge Cases and Error Handling

| Scenario | Expected Behavior |
|----------|------------------|
| Invalid log level in config | Default to "info", log warning |
| Invalid format in config | Default to "auto", log warning |
| No config file | Use defaults (info, auto, stderr) |
| `--log-level` flag | Overrides config file setting |
| TTY detection fails | Default to text format |
| Nil logger | Panic with clear error (programming error) |

## Testing Strategy

### Unit Tests

- [ ] Test level filtering (debug messages hidden at info level)
- [ ] Test text output format
- [ ] Test JSON output format
- [ ] Test TTY auto-detection
- [ ] Test With() creates new logger with context
- [ ] Test config parsing and defaults
- [ ] Test invalid config values fall back to defaults

### Integration Tests

- [ ] Test end-to-end with config file
- [ ] Test `--log-level` flag override
- [ ] Test JSON output is parseable

### Manual Testing Checklist

- [ ] Run `ado config validate` and verify output format
- [ ] Run with `--log-level debug` and verify debug messages appear
- [ ] Pipe output to a file and verify JSON format
- [ ] Run in terminal and verify text format with colors

## Implementation Checklist

- [ ] Create `internal/logging/` package
- [ ] Implement `Logger` interface wrapping `slog.Logger`
- [ ] Implement `Config` struct and parsing
- [ ] Implement TTY detection for auto format
- [ ] Create text handler with optional colors
- [ ] Create JSON handler
- [ ] Add `logging` section to config schema
- [ ] Add `--log-level` flag to root command
- [ ] Write unit tests (target: 80% coverage)
- [ ] Write integration tests
- [ ] Update CLAUDE.md with new package

## Open Questions

Resolved by ADR-0002:

- [x] Which logging library? **Answer: log/slog (stdlib)**
- [x] External dependencies? **Answer: None needed**

## Changelog

| Date | Change | Author |
|------|--------|--------|
| 2025-11-26 | Initial draft | @anowarislam |
