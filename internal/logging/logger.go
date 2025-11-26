// Package logging provides structured logging for ado using Go's log/slog.
package logging

import (
	"context"
	"io"
	"log/slog"
	"os"
)

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

// logger wraps slog.Logger to implement the Logger interface.
type logger struct {
	slog *slog.Logger
}

// New creates a new Logger from the given configuration.
func New(cfg Config) Logger {
	level := parseLevel(cfg.Level)
	output := resolveOutput(cfg.Output)
	handler := createHandler(cfg.Format, output, level)

	return &logger{
		slog: slog.New(handler),
	}
}

// Default returns the default logger (info level, auto format, stderr).
func Default() Logger {
	return New(Config{
		Level:  "info",
		Format: "auto",
		Output: "stderr",
	})
}

// Debug logs at debug level.
func (l *logger) Debug(msg string, args ...any) {
	l.slog.Debug(msg, args...)
}

// Info logs at info level.
func (l *logger) Info(msg string, args ...any) {
	l.slog.Info(msg, args...)
}

// Warn logs at warn level.
func (l *logger) Warn(msg string, args ...any) {
	l.slog.Warn(msg, args...)
}

// Error logs at error level.
func (l *logger) Error(msg string, args ...any) {
	l.slog.Error(msg, args...)
}

// With returns a new Logger with the given attributes.
func (l *logger) With(args ...any) Logger {
	return &logger{
		slog: l.slog.With(args...),
	}
}

// Handler returns the underlying slog.Handler.
func (l *logger) Handler() slog.Handler {
	return l.slog.Handler()
}

// parseLevel converts a string level to slog.Level.
func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// resolveOutput returns the writer for the given output name.
func resolveOutput(output string) io.Writer {
	switch output {
	case "stdout":
		return os.Stdout
	case "stderr":
		return os.Stderr
	default:
		return os.Stderr
	}
}

// createHandler creates the appropriate slog.Handler based on format.
func createHandler(format string, output io.Writer, level slog.Level) slog.Handler {
	opts := &slog.HandlerOptions{
		Level: level,
	}

	switch resolveFormat(format, output) {
	case "json":
		return slog.NewJSONHandler(output, opts)
	default:
		return slog.NewTextHandler(output, opts)
	}
}

// resolveFormat determines the actual format to use.
// If format is "auto", it detects based on TTY.
func resolveFormat(format string, output io.Writer) string {
	switch format {
	case "json":
		return "json"
	case "text":
		return "text"
	case "auto":
		if isTTY(output) {
			return "text"
		}
		return "json"
	default:
		return "text"
	}
}

// isTTY checks if the writer is a terminal.
func isTTY(w io.Writer) bool {
	if f, ok := w.(*os.File); ok {
		stat, err := f.Stat()
		if err != nil {
			return false
		}
		return (stat.Mode() & os.ModeCharDevice) != 0
	}
	return false
}

// NopLogger returns a logger that discards all output.
// Useful for testing or when logging should be disabled.
func NopLogger() Logger {
	return &logger{
		slog: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
}

// FromContext returns a logger from context, or the default logger if none.
func FromContext(ctx context.Context) Logger {
	if l, ok := ctx.Value(loggerKey{}).(Logger); ok {
		return l
	}
	return Default()
}

// WithContext returns a new context with the logger attached.
func WithContext(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, l)
}

// loggerKey is the context key for storing the logger.
type loggerKey struct{}
