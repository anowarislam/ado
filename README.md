# ado

[![CI](https://github.com/anowarislam/ado/actions/workflows/ci.yml/badge.svg)](https://github.com/anowarislam/ado/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/anowarislam/ado)](https://go.dev/)
[![License](https://img.shields.io/github/license/anowarislam/ado)](LICENSE)
[![Release](https://img.shields.io/github/v/release/anowarislam/ado)](https://github.com/anowarislam/ado/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/anowarislam/ado)](https://goreportcard.com/report/github.com/anowarislam/ado)

A single, composable command-line binary designed to replace ad-hoc shell scripts with a more robust, testable, and evolvable tool suite.

## Quick Links

| | |
|---|---|
| [Installation](#installation) | Get started with ado |
| [Documentation](docs/) | Detailed guides and specs |
| [Contributing](CONTRIBUTING.md) | How to contribute |
| [Changelog](CHANGELOG.md) | Version history |
| [Security](SECURITY.md) | Security policy |

## Overview

The project is built around two ideas:

1. **One binary as the stable surface** - a Go-based CLI that is easy to distribute and safe to depend on.
2. **A lab for fast iteration** - a Python-based R&D area for experiments and prototypes, which can later be promoted into the Go binary once stable.

## Installation

### Binary Download

Download pre-built binaries from the [GitHub Releases](https://github.com/anowarislam/ado/releases) page.

**Platforms:** Linux, macOS, Windows (amd64, arm64)

```bash
# Linux/macOS
curl -LO https://github.com/anowarislam/ado/releases/latest/download/ado_VERSION_OS_ARCH.tar.gz
tar xzf ado_*.tar.gz
sudo mv ado /usr/local/bin/

# Verify installation
ado meta info
```

### Docker

```bash
# Pull and run
docker pull ghcr.io/anowarislam/ado:latest
docker run --rm ghcr.io/anowarislam/ado:latest meta info

# Create an alias for convenience
alias ado='docker run --rm ghcr.io/anowarislam/ado:latest'
```

> **Note:** The container runs as a non-root user (UID 65534) for security.
> Config file mounting is not supported in the containerized version.
> Container images are signed—see [Security Policy](SECURITY.md#container-image-verification) for verification.

### Build from Source

```bash
git clone https://github.com/anowarislam/ado.git
cd ado
make go.build
./ado meta info
```

## Usage

```bash
# Show version and build info
ado meta info

# Show environment and config paths
ado meta env

# Echo with transformations
ado echo hello world
ado echo --upper hello world
ado echo --repeat 3 hello

# Get help
ado help
ado help meta
```

## Architecture

```
ado/
├── cmd/ado/           # CLI entrypoint and commands
├── internal/          # Shared Go packages (not importable)
├── lab/py/            # Python lab for prototypes (dev only)
├── docs/              # Documentation and specs
│   ├── commands/      # Command specifications
│   └── style/         # Code style guides
└── make/              # Makefile includes
```

| Directory | Purpose | Shipped |
|-----------|---------|---------|
| `cmd/ado/` | CLI commands (Cobra) | Yes |
| `internal/` | Shared logic (config, meta, ui) | Yes |
| `lab/py/` | Python experiments | No |

## Development

### Prerequisites

- Go 1.23+
- Python 3.12+ (optional, for lab)
- Make

### Quick Start

```bash
# Clone and setup
git clone https://github.com/anowarislam/ado.git
cd ado
make hooks.install     # Install commit hooks

# Development workflow
make go.build          # Build binary
make test              # Run all tests
make lint              # Run linters
make validate          # Full validation (lint + test)
make ci                # Full CI pipeline locally
```

### Available Make Targets

| Target | Description |
|--------|-------------|
| `make go.build` | Build the ado binary |
| `make go.test` | Run Go tests |
| `make go.vet` | Run Go linter |
| `make go.fmt` | Format Go code |
| `make py.test` | Run Python tests |
| `make py.lint` | Run Python linter |
| `make test` | Run all tests |
| `make lint` | Run all linters |
| `make validate` | Lint + test |
| `make ci` | Full CI pipeline |
| `make hooks.install` | Install git hooks |

## Contributing

We welcome contributions! Please see:

- [Contributing Guide](CONTRIBUTING.md) - How to contribute
- [Code of Conduct](CODE_OF_CONDUCT.md) - Community guidelines
- [Security Policy](SECURITY.md) - Reporting vulnerabilities

### Quick Contribution Steps

1. Fork the repository
2. Create a feature branch (`git checkout -b feat/amazing-feature`)
3. Make changes following our [style guides](docs/style/)
4. Commit with [conventional commits](https://www.conventionalcommits.org/) (`git commit -m "feat: add amazing feature"`)
5. Push and open a Pull Request

## Documentation

| Document | Description |
|----------|-------------|
| [Development Workflow](docs/workflow.md) | Issue → ADR → Spec → Implementation |
| [ADRs](docs/adr/) | Architecture Decision Records |
| [Feature Specs](docs/features/) | Non-command feature specifications |
| [Architecture](docs/architecture.md) | Repository layout |
| [Commands](docs/commands-overview.md) | Command specifications |
| [Contributing](docs/contributing.md) | Detailed contribution guide |
| [Release Guide](docs/release.md) | Release automation |
| [Go Style](docs/style/go-style.md) | Go code patterns |
| [Python Style](docs/style/python-style.md) | Python code patterns |

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

Built with:
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [GoReleaser](https://goreleaser.com/) - Release automation
- [Release Please](https://github.com/googleapis/release-please) - Version management
