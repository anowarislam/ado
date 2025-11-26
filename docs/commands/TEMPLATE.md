# [command] Command Spec

<!--
Instructions:
1. Copy this file to `NN-command.md` (use next available number)
2. Fill in all sections below
3. Submit PR with title: `docs(spec): command name`
4. Update status after approval/implementation
-->

| Metadata | Value |
|----------|-------|
| **ADR** | ADR-NNNN or N/A |
| **Status** | Draft |
| **Issue** | #NNN |

## Command

```bash
ado [command] [subcommand] [flags] [args]
```

## Purpose

Why this command exists; what user problem it solves. (1-2 sentences)

## Usage Examples

```bash
# Example 1: Basic usage
ado command arg
# Expected output

# Example 2: With flags
ado command --flag value
# Expected output

# Example 3: Common workflow
ado command --output json | jq '.field'
# Expected output
```

## Arguments

| Argument | Type | Required | Description |
|----------|------|----------|-------------|
| `arg1` | string | Yes | Description |
| `arg2...` | string[] | No | Description (variadic) |

## Flags

### Command-Specific Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--flag` | `-f` | string | `""` | Description |
| `--count` | `-n` | int | `1` | Description |
| `--output` | `-o` | enum | `text` | Output format: text, json, yaml |
| `--verbose` | `-v` | bool | `false` | Enable verbose output |

### Inherited Global Flags

All commands inherit these flags from the root command:

- `--config PATH` - Config file path (default: auto-detected)
- `--log-level LEVEL` - Log level: debug, info, warn, error (default: info)
- `--help, -h` - Show help for command

## Behavior

Describe the command's logic step by step:

1. Parse and validate arguments
2. Load configuration (if applicable)
3. Perform main operation
4. Format and output result

### Output Formats

**Text (default):**
```
Human-readable output here
```

**JSON (`--output json`):**
```json
{
  "field": "value",
  "nested": {
    "key": "value"
  }
}
```

**YAML (`--output yaml`):**
```yaml
field: value
nested:
  key: value
```

## Error Cases

| Condition | Exit Code | Error Message |
|-----------|-----------|---------------|
| No arguments provided | 1 | `Error: at least one argument required` |
| Invalid flag value | 1 | `Error: invalid value for --flag: "X"` |
| Conflicting flags | 1 | `Error: --flag1 and --flag2 are mutually exclusive` |
| File not found | 1 | `Error: file not found: "path"` |
| Permission denied | 1 | `Error: permission denied: "path"` |

## Implementation

| Purpose | Path |
|---------|------|
| Command | `cmd/ado/[command]/[command].go` |
| Tests | `cmd/ado/[command]/[command]_test.go` |
| Shared logic | `internal/[package]/` |

### Implementation Notes

- Export `NewCommand() *cobra.Command`
- Wire into `cmd/ado/root/root.go` via `AddCommand()`
- Use table-driven tests with `t.Run()` subtests
- Wrap errors with context: `fmt.Errorf("command: %w", err)`

## Testing Checklist

- [ ] Basic usage works as documented
- [ ] All flags work correctly
- [ ] All output formats work (text, json, yaml)
- [ ] Error cases return correct exit codes
- [ ] Error messages are clear and actionable
- [ ] Help text (`--help`) is accurate
- [ ] Edge cases handled (empty input, large input, etc.)

## Related Commands

- `ado related-command` - Description of relationship
- `ado other-command` - Description of relationship
