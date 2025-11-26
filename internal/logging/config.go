package logging

// Config holds logging configuration.
type Config struct {
	// Level is the minimum log level: debug, info, warn, error.
	// Default: "info"
	Level string

	// Format is the output format: auto, text, json.
	// "auto" detects TTY and uses text for terminals, JSON otherwise.
	// Default: "auto"
	Format string

	// Output is the output destination: stderr, stdout.
	// Default: "stderr"
	Output string
}

// DefaultConfig returns the default logging configuration.
func DefaultConfig() Config {
	return Config{
		Level:  "info",
		Format: "auto",
		Output: "stderr",
	}
}

// Validate validates the configuration and returns a normalized copy.
// Invalid values are replaced with defaults.
func (c Config) Validate() Config {
	result := c

	// Validate level
	switch c.Level {
	case "debug", "info", "warn", "error":
		// Valid
	default:
		result.Level = "info"
	}

	// Validate format
	switch c.Format {
	case "auto", "text", "json":
		// Valid
	default:
		result.Format = "auto"
	}

	// Validate output
	switch c.Output {
	case "stderr", "stdout":
		// Valid
	default:
		result.Output = "stderr"
	}

	return result
}

// IsValidLevel checks if the given level string is valid.
func IsValidLevel(level string) bool {
	switch level {
	case "debug", "info", "warn", "error":
		return true
	default:
		return false
	}
}
