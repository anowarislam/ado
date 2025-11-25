# Project Structure Spec

Project name: ado

## Goal

Provide a single Go-based binary (ado) that exposes subcommands. Python lives in a separate lab directory for prototyping logic, but is not required at runtime.

## Repository layout (spec)

Top-level layout:

- cmd/ado/
    - Contains the CLI entrypoint (main package) and wiring for subcommands.
- internal/
    - Internal packages for shared logic, not part of the public API.
    - Example: internal/meta, internal/ui, internal/config.
- pkg/ (optional, future)
    - Publicly reusable pieces, if you later want third parties to import.
- lab/py/
- Python lab for R&D, experiments, and prototypes.
    - Never required for the compiled ado binary to work.
- config/
    - Default config templates, example YAML files.
- docs/
    - Additional documents, examples, design notes.
- tests/ (if you want high-level integration tests separated from unit tests)
- README.md
- LICENSE
- Makefile or equivalent build tooling (optional at this spec stage).


### Concrete tree (descriptive, not binding):
```
ado/
  cmd/
    ado/
      main.go                # Entry for the CLI (spec only, not implemented yet)
  internal/
    meta/                    # Build info, environment, self-introspection logic
    ui/                      # Text UI conventions: colors, error formatting, tables
    config/                  # Config loading & merging logic
  lab/
    py/
      README.md              # How to run prototypes
      autoscale/
      vmss/
      kusto/
  config/
    ado.example.yaml     # Example config, used by README/docs
  docs/
    01-design-cli-structure.md
    02-commands.md
    commands/
      01-echo.md
      02-help.md
      03-meta.md
  README.md
  LICENSE
```