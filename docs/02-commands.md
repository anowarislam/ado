# Command Specs

## Name

ado - a pluggable command-line toolkit for automation and diagnostics

### Purpose:

Provide consistent, discoverable documentation for all commands, plus a brief high-level description of the philosophy and structure of ado.

## Description

ado as a single binary. It is designed to be extended with new commands. It has a separation of “shipping surface” (Go) and “lab” (Python).

## Usage

ado <command> [subcommand] [flags]

## Global Flags

1. --config string – Config file path
2. --log-level string – Log level (default “info”)
3. --version – Print the version number
4. -h, --help – Help for ado

## Global behavior & conventions

These apply to all commands, including meta, echo, help:

- Binary name: ado
- General syntax: ado <command> [subcommand] [flags] [args]
- Global flags (initial set):
	- --help, -h: show help for the current command or subcommand.
	- --version: print version string plus minimal build info.
	- --config PATH: optional, explicit path to config file.
	- --log-level LEVEL: overrides default log level (info, debug, etc.).
- Exit codes:
	- 0 – success.
	- >0 – failure, command-specific but consistent (later spec).
- Output conventions:
	- Machine-readable modes (e.g. JSON) should be opt-in via --output json.
	- Human-readable default output is structured text, suitable for terminals.
	- All human-readable output goes to stdout; error messages go to stderr.
- Configuration:
	- Default config search order:
		- 1. --config PATH if provided.
		- 2. $XDG_CONFIG_HOME/ado/config.yaml (or $HOME/.config/ado/config.yaml).
		- 3. $HOME/.ado/config.yaml as fallback.
	- Environment variables may override config later, but not required for v0.