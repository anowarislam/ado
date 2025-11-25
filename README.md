# ado

`ado` is a single, composable command-line binary designed to replace ad-hoc shell scripts with a more robust, testable, and evolvable tool suite.

The project is built around two ideas:

1. **One binary as the stable surface** – a Go-based CLI that is easy to distribute and safe to depend on.
2. **A lab for fast iteration** – a Python-based R&D area for experiments and prototypes, which can later be promoted into the Go binary once stable.

---

## Goals

- Provide a **single CLI entrypoint** (`ado`) with clear, discoverable commands.
- Offer a **safe alternative to brittle bash scripts**, with:
  - Stronger error handling.
  - Better testability.
  - Consistent UX and logging.
- Maintain a **clean shipping story**:
  - Distribute `ado` as a static or near-static Go binary.
  - Avoid requiring Python at runtime for end users.
- Keep the development loop fast through a **Python lab** for trying ideas.

---

## High-Level Architecture

The repository is organized into three main areas:

- `cmd/ado/`  
  Contains the entrypoint for the `ado` binary. This is where commands are wired together.

- `internal/`  
  Contains internal Go packages with shared logic (metadata, configuration, UI helpers, etc.). These are not intended to be imported by other Go modules.

- `lab/py/`  
  A Python-based lab for experiments, prototype implementations, and one-off analyses. The contents of this directory are for developers, not for end users. The compiled `ado` binary does **not** depend on this directory.

Additional directories:

- `config/` – Example configuration files and templates.
- `docs/` – Detailed specifications and design notes for commands and subsystems.

---

## Design Principles

### 1. Go as the Shipping Surface

The `ado` binary is written in Go. This ensures:

- Straightforward distribution as a single binary for macOS, Linux, and Windows.
- Fast startup and low runtime overhead.
- Strong typing and clear error handling.

End users should only need to download or install the `ado` binary to use it.

### 2. Python as the Lab

Python is used in `lab/py/` as a space to:

- Prototype new command behaviors.
- Explore complex data transformations.
- Develop and refine heuristics before they are ported to Go.

The promotion path is:

1. Prototype in Python under `lab/py/`.
2. Lock in behavior with tests and fixtures.
3. Reimplement the stable logic in Go under `internal/` and expose it through the `ado` CLI.
4. Optionally keep the Python prototype as documentation and for future experimentation.

Python is an implementation detail of the development process, not a hard dependency for end users.

### 3. Opinionated but Extensible CLI UX

The CLI follows a consistent style:

- Binary: `ado`
- Syntax: `ado <command> [subcommand] [flags]`
- Human-readable output by default; machine-readable output via `--output json` where relevant.
- Global flags for:
  - Config file selection (`--config`).
  - Logging verbosity (`--log-level`).
  - Help (`--help`) and version (`--version`).

Each command should:

- Have a clear, single responsibility.
- Provide useful examples via `ado help <command>`.
- Be safe by default (avoid destructive actions unless explicitly requested).

### 4. Config and Environment

Configuration is designed to be predictable and explicit:

- Default config lookup follows:
  1. `--config PATH` if provided.
  2. `$XDG_CONFIG_HOME/ado/config.yaml` or `$HOME/.config/ado/config.yaml`.
  3. `$HOME/.ado/config.yaml` as a fallback.

Environment variables may be used to override specific fields in the future but are not required for the initial version.

The `meta` command allows introspection of the configuration and environment used by `ado`.

### 5. Testing and Stability

As new behaviors mature:

- Add tests around their input/output contracts.
- Keep tests close to the Go implementation (unit tests in `internal/` packages) and optionally integration-level tests in a dedicated directory.

The `ado` binary aims to be stable and predictable. Breaking changes should be deliberate and communicated via versioning.

---

## Initial Commands

The initial set of commands is intentionally small but representative.

### `meta`

Introspects the `ado` binary and its environment.

Examples:

- `ado meta info`  
  Show version, commit, build time, and platform.

- `ado meta env`  
  Show resolved configuration paths and key environment details.

### `echo`

Provides a simple echo function with optional transformations. It is primarily a reference and test command for CLI behavior.

Examples:

- `ado echo hello world`  
  Echo the words as a single line.

- `ado echo --upper hello world`  
  Echo the message in uppercase.

- `ado echo --repeat 3 hello`  
  Print "hello" three times.

### `help`

Displays general help or command-specific details.

Examples:

- `ado help`  
  Show top-level usage and available commands.

- `ado help meta`  
  Show detailed help for the `meta` command.

---

## Roadmap (High-Level)

1. **v0.1.0 – Skeleton**
   - Implement `meta`, `echo`, and `help` as described in the specs.
   - Establish logging, config loading, and basic error handling patterns.

2. **v0.2.x – First Real Utilities**
   - Add the first “real” automation commands (e.g., Git, VMSS, or incident utilities).
   - Establish a pattern for promoting Python lab prototypes to Go commands.

3. **v0.3.x and beyond**
   - Introduce optional extension mechanisms.
   - Expand integration with external systems (Kusto, issue trackers, etc.).
   - Harden behaviors, document guarantees, and refine UX as new commands are added.

---

## Installation

### Binary Download

Download pre-built binaries from the [GitHub Releases](https://github.com/anowarislam/ado/releases) page.

Available platforms:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64, arm64)

```bash
# Example: Download and install on Linux/macOS
curl -LO https://github.com/anowarislam/ado/releases/latest/download/ado_<version>_<os>_<arch>.tar.gz
tar xzf ado_*.tar.gz
sudo mv ado /usr/local/bin/
```

### Docker

Multi-architecture container images are available on GitHub Container Registry:

```bash
# Pull the latest image (auto-selects your architecture)
docker pull ghcr.io/anowarislam/ado:latest

# Run a command
docker run --rm ghcr.io/anowarislam/ado:latest meta info

# Run with a specific version
docker run --rm ghcr.io/anowarislam/ado:1.0.0 echo "Hello"

# Create an alias for convenience
alias ado='docker run --rm -v ~/.config/ado:/root/.config/ado ghcr.io/anowarislam/ado:latest'
```

Available image tags:
- `ghcr.io/anowarislam/ado:latest` - Latest stable release
- `ghcr.io/anowarislam/ado:X.Y.Z` - Specific version

### Build from Source

```bash
git clone https://github.com/anowarislam/ado.git
cd ado
make go.build
./ado meta info
```

---

## Getting Started (Developer Perspective)

1. Clone the repository.
2. Explore `lab/py/` for prototypes and experimental logic.
3. Add or update commands under `cmd/ado/` and `internal/`.
4. Update the docs in `docs/commands/` when new commands or flags are added.
5. Keep the README and `meta` command aligned with the current project state.
