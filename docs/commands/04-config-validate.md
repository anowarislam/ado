# config validate Command Spec

| Metadata | Value |
|----------|-------|
| **ADR** | N/A |
| **Status** | Draft |
| **Issue** | #36 |

## Command

```bash
ado config validate [--file PATH] [--strict] [--output FORMAT]
```

## Purpose

Validate configuration files against the expected schema and report errors. Useful for:
- CI pipelines that need to verify config before deployment
- Debugging configuration issues during development
- Pre-commit hooks to catch config errors early

## Usage Examples

```bash
# Example 1: Validate auto-detected config
$ ado config validate
✓ Config valid: /Users/user/.config/ado/config.yaml

# Example 2: Validate specific file
$ ado config validate --file ./custom-config.yaml
✓ Config valid: ./custom-config.yaml

# Example 3: Strict mode (warnings become errors)
$ ado config validate --strict
✗ Config invalid: /Users/user/.config/ado/config.yaml
  Error: unknown key "foobar" at line 5

# Example 4: JSON output for CI pipelines
$ ado config validate --output json
{
  "valid": true,
  "path": "/Users/user/.config/ado/config.yaml",
  "errors": [],
  "warnings": []
}

# Example 5: Validation failure
$ ado config validate --file ./broken.yaml
✗ Config invalid: ./broken.yaml
  Error: invalid YAML syntax at line 3: mapping values are not allowed here
```

## Arguments

This command takes no positional arguments.

## Flags

### Command-Specific Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--file` | `-f` | string | auto-detect | Path to config file to validate |
| `--strict` | `-s` | bool | `false` | Treat warnings as errors (exit 1) |
| `--output` | `-o` | enum | `text` | Output format: `text`, `json` |

### Inherited Global Flags

All commands inherit these flags from the root command:

- `--config PATH` - Config file path (default: auto-detected)
- `--log-level LEVEL` - Log level: debug, info, warn, error (default: info)
- `--help, -h` - Show help for command

**Note**: `--file` takes precedence over `--config` for this command.

## Behavior

### Validation Steps

1. **Resolve config path**
   - If `--file` provided: use that path
   - Else if `--config` provided: use that path
   - Else: auto-detect using `ResolveConfigPath()` from `internal/config`

2. **Check file exists**
   - If file not found: report error, exit 1

3. **Parse YAML**
   - If invalid YAML syntax: report error with line number, exit 1

4. **Validate structure** (when schema is defined)
   - Check for unknown keys → warning (or error in strict mode)
   - Check value types match expected types → error
   - Check required keys present → error

5. **Report results**
   - Success: print confirmation, exit 0
   - Warnings only (non-strict): print warnings, exit 0
   - Warnings (strict) or errors: print issues, exit 1

### Output Formats

**Text (default):**

Success:
```
✓ Config valid: /path/to/config.yaml
```

Success with warnings (non-strict):
```
✓ Config valid: /path/to/config.yaml
  Warning: unknown key "deprecated_option" at line 12
```

Failure:
```
✗ Config invalid: /path/to/config.yaml
  Error: invalid YAML syntax at line 3: mapping values are not allowed here
```

**JSON (`--output json`):**

```json
{
  "valid": true,
  "path": "/path/to/config.yaml",
  "errors": [],
  "warnings": [
    {
      "message": "unknown key \"deprecated_option\"",
      "line": 12,
      "severity": "warning"
    }
  ]
}
```

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Valid config (no errors, warnings allowed in non-strict) |
| 1 | Invalid config or errors encountered |

## Error Cases

| Condition | Exit Code | Output |
|-----------|-----------|--------|
| No config found | 1 | `Error: no config file found. Searched: [paths...]` |
| File not found | 1 | `Error: config file not found: "/path/to/file"` |
| Permission denied | 1 | `Error: permission denied: "/path/to/file"` |
| Invalid YAML syntax | 1 | `Error: invalid YAML at line N: <parser message>` |
| Unknown keys (non-strict) | 0 | `Warning: unknown key "foo" at line N` |
| Unknown keys (strict) | 1 | `Error: unknown key "foo" at line N` |
| Invalid value type | 1 | `Error: invalid type for "key": expected string, got int` |

## Config Schema

The current minimal config schema:

```yaml
# ~/.config/ado/config.yaml
version: 1                    # Config schema version (required)

# Future: command defaults, aliases, plugins, etc.
```

### Schema Rules

| Key | Type | Required | Description |
|-----|------|----------|-------------|
| `version` | int | Yes | Config schema version (currently: 1) |

**Note**: The schema will expand as features are added. Unknown keys generate warnings to support forward compatibility.

## Implementation

| Purpose | Path |
|---------|------|
| Parent command | `cmd/ado/config/config.go` |
| Validate subcommand | `cmd/ado/config/validate.go` |
| Tests | `cmd/ado/config/validate_test.go` |
| Validation logic | `internal/config/validate.go` |
| Validation tests | `internal/config/validate_test.go` |

### Implementation Notes

- `cmd/ado/config/config.go` exports `NewCommand() *cobra.Command` with `validate` subcommand
- Wire into `cmd/ado/root/root.go` via `AddCommand(config.NewCommand())`
- Validation logic in `internal/config/` for reuse by other commands
- Use `gopkg.in/yaml.v3` for YAML parsing with line number support
- Table-driven tests covering all error cases

### Code Structure

```go
// internal/config/validate.go
package config

type ValidationResult struct {
    Valid    bool
    Path     string
    Errors   []ValidationIssue
    Warnings []ValidationIssue
}

type ValidationIssue struct {
    Message  string
    Line     int
    Severity string // "error" or "warning"
}

// Validate validates a config file and returns the result.
func Validate(path string) (*ValidationResult, error)
```

```go
// cmd/ado/config/validate.go
package config

func newValidateCmd() *cobra.Command {
    // Implementation
}
```

## Testing Checklist

- [ ] No config found returns error with search paths
- [ ] Missing file returns clear error
- [ ] Permission denied handled gracefully
- [ ] Invalid YAML reports line number
- [ ] Valid minimal config passes
- [ ] Unknown keys produce warnings
- [ ] `--strict` makes warnings into errors
- [ ] `--file` overrides auto-detection
- [ ] `--output json` produces valid JSON
- [ ] Exit codes match specification
- [ ] Help text (`--help`) is accurate

### Test Cases

```go
func TestConfigValidate(t *testing.T) {
    tests := []struct {
        name     string
        content  string
        strict   bool
        wantErr  bool
        wantWarn bool
    }{
        {"valid minimal", "version: 1\n", false, false, false},
        {"missing version", "foo: bar\n", false, true, false},
        {"unknown key non-strict", "version: 1\nunknown: value\n", false, false, true},
        {"unknown key strict", "version: 1\nunknown: value\n", true, true, false},
        {"invalid yaml", "version: 1\n  bad indent", false, true, false},
    }
    // ...
}
```

## Related Commands

- `ado meta env` - Shows config search paths (useful for debugging which config is loaded)
- `ado config init` - (Future) Initialize a new config file
- `ado config show` - (Future) Display current config with sources

## Open Questions

- [ ] Should we support YAML anchors and aliases?
- [ ] Should `--fix` flag auto-correct simple issues (remove unknown keys)?
- [ ] Should validation warn about deprecated keys from older schema versions?

## Changelog

| Date | Change | Author |
|------|--------|--------|
| 2025-11-26 | Initial draft | @anowarislam |
