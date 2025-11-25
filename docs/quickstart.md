# Quick Start

Get up and running with `ado` in minutes.

## Verify Installation

After [installing](installation.md) `ado`, verify it's working:

```bash
$ ado meta info
Version: v1.0.0
Commit: abc1234
Build Time: 2025-01-01T00:00:00Z
Platform: darwin/arm64
```

## Basic Commands

### Get Help

```bash
# Show all commands
$ ado help

# Get help for specific command
$ ado help meta
$ ado help echo
```

### Meta Information

The `meta` command provides information about the `ado` binary and environment.

```bash
# Show version and build info
$ ado meta info
Version: v1.0.0
Commit: abc1234
Build Time: 2025-01-01T00:00:00Z
Platform: darwin/arm64

# Show environment and config paths
$ ado meta env
Config Paths:
  XDG_CONFIG_HOME: /home/user/.config
  Config Dir: /home/user/.config/ado
  Config File: /home/user/.config/ado/config.yaml
```

### Echo Command

The `echo` command provides text transformation capabilities.

```bash
# Basic echo
$ ado echo hello world
hello world

# Uppercase transformation
$ ado echo --upper hello world
HELLO WORLD

# Repeat text
$ ado echo --repeat 3 hello
hello
hello
hello

# Combine flags
$ ado echo --upper --repeat 2 hello
HELLO
HELLO
```

## Configuration

`ado` looks for configuration in these locations (in order):

1. `--config PATH` flag (if provided)
2. `$XDG_CONFIG_HOME/ado/config.yaml`
3. `$HOME/.config/ado/config.yaml`
4. `$HOME/.ado/config.yaml`

### Example Config

```yaml
# ~/.config/ado/config.yaml
version: 1
log_level: info
```

## Common Workflows

### Check Version Before Update

```bash
# Current version
ado meta info

# After update, verify new version
ado meta info
```

### Debug Configuration Issues

```bash
# Check which config paths are being used
ado meta env
```

## Next Steps

- [Commands Reference](commands-overview.md) - Detailed command documentation
- [Architecture](architecture.md) - Understand the project structure
- [Contributing](contributing.md) - Help improve `ado`
