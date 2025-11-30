# Installation

Multiple installation methods are available for `ado`.

## Homebrew (macOS & Linux)

**Recommended installation method for macOS users.**

Homebrew provides automatic updates, dependency management, and clean uninstallation.

### Installation

```bash
# First-time setup: add the ado tap
brew tap anowarislam/ado

# Install ado
brew install ado

# Verify installation
ado meta info
```

### Supported Platforms

- macOS (Intel and Apple Silicon)
- Linux with Homebrew (Linuxbrew)

### Updating

```bash
# Check for updates
brew outdated

# Upgrade ado
brew upgrade ado
```

### Uninstalling

```bash
# Uninstall ado
brew uninstall ado

# Remove tap (optional)
brew untap anowarislam/ado
```

---

## Binary Download

Download pre-built binaries from the [GitHub Releases](https://github.com/anowarislam/ado/releases) page.

### Supported Platforms

| OS | Architecture | File |
|----|--------------|------|
| Linux | amd64 | `ado_VERSION_linux_amd64.tar.gz` |
| Linux | arm64 | `ado_VERSION_linux_arm64.tar.gz` |
| macOS | amd64 | `ado_VERSION_darwin_amd64.tar.gz` |
| macOS | arm64 (Apple Silicon) | `ado_VERSION_darwin_arm64.tar.gz` |
| Windows | amd64 | `ado_VERSION_windows_amd64.zip` |
| Windows | arm64 | `ado_VERSION_windows_arm64.zip` |

### Linux / macOS

```bash
# Download latest release (replace VERSION, OS, ARCH)
curl -LO https://github.com/anowarislam/ado/releases/latest/download/ado_VERSION_OS_ARCH.tar.gz

# Extract
tar xzf ado_*.tar.gz

# Move to PATH
sudo mv ado /usr/local/bin/

# Verify installation
ado meta info
```

### Windows

1. Download the `.zip` file from [Releases](https://github.com/anowarislam/ado/releases)
2. Extract the archive
3. Add the directory to your `PATH`
4. Open a new terminal and run `ado meta info`

### Verifying Downloads

All releases include SHA256 checksums and cryptographic attestations.

#### Checksum Verification

```bash
# Download checksums
curl -LO https://github.com/anowarislam/ado/releases/download/vX.Y.Z/checksums.txt

# Verify
sha256sum -c checksums.txt --ignore-missing
```

#### Attestation Verification (Recommended)

Verify the binary was built by our CI pipeline:

```bash
# Requires GitHub CLI
gh attestation verify ado_X.Y.Z_linux_amd64.tar.gz --owner anowarislam
```

See [Security Policy](https://github.com/anowarislam/ado/blob/main/SECURITY.md) for more details.

## Docker

Multi-architecture container images are available on GitHub Container Registry.

### Pull and Run

```bash
# Pull latest
docker pull ghcr.io/anowarislam/ado:latest

# Run a command
docker run --rm ghcr.io/anowarislam/ado:latest meta info

# Run with specific version
docker run --rm ghcr.io/anowarislam/ado:1.0.0 echo "Hello"
```

### Available Tags

| Tag | Description |
|-----|-------------|
| `latest` | Latest stable release |
| `X.Y.Z` | Specific version (e.g., `1.0.0`) |
| `X.Y.Z-amd64` | AMD64-specific image |
| `X.Y.Z-arm64` | ARM64-specific image |

### Shell Alias

For convenience, create an alias:

```bash
# Add to ~/.bashrc or ~/.zshrc
alias ado='docker run --rm ghcr.io/anowarislam/ado:latest'
```

> **Note:** The container runs as a non-root user (UID 65534) for security.

### Verifying Container Images

Container images are signed and can be verified with [cosign](https://github.com/sigstore/cosign):

```bash
cosign verify ghcr.io/anowarislam/ado:latest \
  --certificate-identity-regexp="https://github.com/anowarislam/ado/" \
  --certificate-oidc-issuer="https://token.actions.githubusercontent.com"
```

## Build from Source

### Prerequisites

- Go 1.23 or later
- Make
- Git

### Steps

```bash
# Clone repository
git clone https://github.com/anowarislam/ado.git
cd ado

# Build
make go.build

# Verify
./ado meta info

# Optional: Install to PATH
sudo mv ./ado /usr/local/bin/
```

### Development Build

For development with debug symbols:

```bash
# Build without stripping
go build -o ado ./cmd/ado

# Run tests
make test
```

## Updating

### Binary

Download and replace the existing binary with the new version.

### Docker

```bash
docker pull ghcr.io/anowarislam/ado:latest
```

### From Source

```bash
git pull origin main
make go.build
```

## Uninstalling

### Binary

```bash
sudo rm /usr/local/bin/ado
rm -rf ~/.config/ado  # Optional: remove config
```

### Docker

```bash
docker rmi ghcr.io/anowarislam/ado:latest
```

## Next Steps

- [Quick Start](quickstart.md) - Get started with basic usage
- [Commands](commands-overview.md) - See all available commands
