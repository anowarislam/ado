# ado

A single, composable command-line binary designed to replace ad-hoc shell scripts with a more robust, testable, and evolvable tool suite.

[![CI](https://github.com/anowarislam/ado/actions/workflows/ci.yml/badge.svg)](https://github.com/anowarislam/ado/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/anowarislam/ado)](https://go.dev/)
[![License](https://img.shields.io/github/license/anowarislam/ado)](https://github.com/anowarislam/ado/blob/main/LICENSE)
[![Release](https://img.shields.io/github/v/release/anowarislam/ado)](https://github.com/anowarislam/ado/releases)

## Overview

The project is built around two ideas:

1. **One binary as the stable surface** - a Go-based CLI that is easy to distribute and safe to depend on.
2. **A lab for fast iteration** - a Python-based R&D area for experiments and prototypes, which can later be promoted into the Go binary once stable.

## Features

- :rocket: **Single binary** - Easy to distribute, no runtime dependencies
- :test_tube: **Python lab** - Fast prototyping before Go implementation
- :whale: **Docker support** - Multi-arch container images
- :gear: **Automated releases** - Semantic versioning with conventional commits
- :shield: **Well tested** - Comprehensive test coverage

## Quick Install

=== "Binary"

    ```bash
    curl -LO https://github.com/anowarislam/ado/releases/latest/download/ado_VERSION_OS_ARCH.tar.gz
    tar xzf ado_*.tar.gz
    sudo mv ado /usr/local/bin/
    ```

=== "Docker"

    ```bash
    docker pull ghcr.io/anowarislam/ado:latest
    docker run --rm ghcr.io/anowarislam/ado:latest meta info
    ```

=== "Source"

    ```bash
    git clone https://github.com/anowarislam/ado.git
    cd ado
    make go.build
    ```

## Quick Usage

```bash
# Show version info
ado meta info

# Echo with transformations
ado echo --upper hello world

# Get help
ado help
```

## Next Steps

- [Installation](installation.md) - Detailed installation instructions
- [Quick Start](quickstart.md) - Get up and running quickly
- [Commands](commands-overview.md) - Available commands reference
- [Makefile Reference](make.md) - Build, test, and development targets
- [Contributing](contributing.md) - How to contribute

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/anowarislam/ado/blob/main/LICENSE) file for details.
