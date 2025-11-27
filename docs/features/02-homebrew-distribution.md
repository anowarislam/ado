# Feature: Homebrew Distribution

| Metadata | Value |
|----------|-------|
| **ADR** | N/A |
| **Status** | Draft |
| **Issue** | [#54](https://github.com/anowarislam/ado/issues/54) |
| **Author(s)** | @anowarislam |

## Overview

Automated Homebrew formula distribution via GoReleaser using an in-repository tap, enabling macOS and Linux users to install `ado` with `brew install ado`.

## Motivation

Why is this feature needed? What problem does it solve?

- **Pain point**: macOS users currently must manually download binaries, extract archives, and move files to PATH. This is error-prone and unfamiliar compared to `brew install`.
- **Who benefits**:
  - **macOS developers**: Native package manager integration (Homebrew is de facto standard)
  - **Linux users**: Homebrew on Linux (Linuxbrew) support for consistent experience
  - **DevOps teams**: Automated updates via `brew upgrade` without manual downloads
  - **CI/CD pipelines**: Simple `brew install` step instead of curl/tar/chmod scripts
- **Without this**: Users face manual installation complexity, no automatic updates, no dependency management, and reduced discoverability

### Current State vs. Desired State

| Aspect | Current (Manual Download) | With Homebrew |
|--------|---------------------------|---------------|
| **Installation** | curl → tar → mv → chmod (4 steps) | `brew install ado` (1 step) |
| **Updates** | Check GitHub → download → replace | `brew upgrade ado` (automatic) |
| **Discoverability** | Must know GitHub URL | `brew search ado` finds it |
| **Verification** | Manual checksum comparison | Automatic SHA-256 validation |
| **Uninstall** | Manual deletion | `brew uninstall ado` (clean) |
| **Dependencies** | Manual resolution | Homebrew handles automatically |

## Specification

### Behavior

How the Homebrew distribution works (end-to-end automation):

1. **Developer Workflow** (No changes required)
   - Developer merges PR to `main` branch
   - Release Please creates/updates Release PR with version bump
   - Maintainer merges Release PR
   - GitHub creates new Release tag (e.g., `v1.4.0`)

2. **Automated Build & Release** (GoReleaser)
   - GitHub Actions trigger GoReleaser workflow (`.github/workflows/goreleaser.yml`)
   - GoReleaser builds binaries for all platforms (Linux, macOS, Windows)
   - GoReleaser generates checksums and SLSA provenance attestations
   - GoReleaser publishes release artifacts to GitHub Releases

3. **Homebrew Formula Update** (New - Automated by GoReleaser)
   - GoReleaser reads `brews` configuration from `.goreleaser.yaml`
   - Generates Ruby formula file with version, URL, and SHA-256 checksum
   - Commits formula file to `anowarislam/ado` repository (same repo)
   - Formula file location: `Formula/ado.rb` in the project repository

4. **User Installation**
   ```bash
   # First time setup (add tap) - one-time only
   brew tap anowarislam/ado

   # Install ado
   brew install ado

   # Homebrew downloads binary from GitHub Releases
   # Verifies SHA-256 checksum automatically
   # Installs to /opt/homebrew/bin/ado (Apple Silicon) or /usr/local/bin/ado (Intel)
   ```

5. **User Updates**
   ```bash
   # Check for updates
   brew upgrade ado

   # Homebrew fetches latest cask from tap
   # Downloads new version if available
   # Replaces old binary with new version
   ```

### Configuration

#### GoReleaser Configuration

Add the following section to `.goreleaser.yaml`:

```yaml
# Homebrew Formula (in-repository tap)
# See: https://goreleaser.com/customization/homebrew/
brews:
  - # Repository configuration - same repo, no separate token needed
    repository:
      owner: anowarislam
      name: ado  # Same repository
      branch: main
      # Note: Uses default GITHUB_TOKEN from workflow (no separate token needed)

    # Formula directory in same repository
    folder: Formula

    # Formula metadata
    name: ado
    homepage: https://github.com/anowarislam/ado
    description: "Composable command-line binary to replace ad-hoc shell scripts"
    license: MIT

    # Post-installation message
    caveats: |
      Get started with ado by running:
        ado meta info

      Configuration: ~/.config/ado/config.yaml
      Documentation: https://github.com/anowarislam/ado

    # Git commit configuration
    commit_author:
      name: ado-releaser
      email: noreply@github.com

    commit_msg_template: "chore(brew): update formula to {{ .Tag }}"

    # Skip upload for snapshot builds (local testing only)
    skip_upload: auto

    # Test block (Homebrew requirement)
    test: |
      system "#{bin}/ado", "meta", "info"
      assert_match version.to_s, shell_output("#{bin}/ado meta info")

    # Installation from bottle (precompiled binary)
    install: |
      bin.install "ado"
```

| Option | Type | Value | Purpose |
|--------|------|-------|---------|
| `repository.owner` | string | `anowarislam` | GitHub username/org |
| `repository.name` | string | `ado` | **Same repository** (in-repo tap) |
| `folder` | string | `Formula` | Directory for formula files |
| `name` | string | `ado` | Formula name (used in `brew install`) |
| `homepage` | string | Project URL | Required by Homebrew |
| `description` | string | Brief description | Required by Homebrew (max ~80 chars) |
| `license` | string | `MIT` | SPDX license identifier |
| `caveats` | string | User message | Post-install instructions |
| `skip_upload` | string | `auto` | Skip for snapshot builds |
| `test` | string | Ruby code | Verification tests |
| `install` | string | Ruby code | Installation instructions |

#### GitHub Actions Configuration

**No additional configuration needed** - the existing `GITHUB_TOKEN` in `.github/workflows/goreleaser.yml` has sufficient permissions to push to the same repository.

```yaml
# .github/workflows/goreleaser.yml
# Existing token works (no changes needed)
env:
  GITHUB_TOKEN: ${{ steps.generate-token.outputs.token }}  # Already configured
```

**Key Advantage**: Uses the GitHub App token already configured in your workflow - no separate PAT required.

#### In-Repository Tap Structure

**Repository**: `anowarislam/ado` (same repo as source code)

**Structure**:
```
ado/
├── cmd/                # Go source code
├── internal/           # Internal packages
├── Formula/            # Homebrew formulae (new directory)
│   └── ado.rb          # Auto-generated by GoReleaser
├── docs/               # Documentation
├── .goreleaser.yaml    # Build configuration
└── README.md
```

**Example Generated Formula** (`Formula/ado.rb`):
```ruby
class Ado < Formula
  desc "Composable command-line binary to replace ad-hoc shell scripts"
  homepage "https://github.com/anowarislam/ado"
  version "1.4.0"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/anowarislam/ado/releases/download/v1.4.0/ado_1.4.0_darwin_arm64.tar.gz"
      sha256 "abc123..."
    end
    if Hardware::CPU.intel?
      url "https://github.com/anowarislam/ado/releases/download/v1.4.0/ado_1.4.0_darwin_amd64.tar.gz"
      sha256 "def456..."
    end
  end

  on_linux do
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/anowarislam/ado/releases/download/v1.4.0/ado_1.4.0_linux_arm64.tar.gz"
      sha256 "ghi789..."
    end
    if Hardware::CPU.intel?
      url "https://github.com/anowarislam/ado/releases/download/v1.4.0/ado_1.4.0_linux_amd64.tar.gz"
      sha256 "jkl012..."
    end
  end

  def install
    bin.install "ado"
  end

  test do
    system "#{bin}/ado", "meta", "info"
    assert_match version.to_s, shell_output("#{bin}/ado meta info")
  end

  def caveats
    <<~EOS
      Get started with ado by running:
        ado meta info

      Configuration: ~/.config/ado/config.yaml
      Documentation: https://github.com/anowarislam/ado
    EOS
  end
end
```

### File Locations

| Purpose | Path | Modified By |
|---------|------|-------------|
| GoReleaser config | `.goreleaser.yaml` | Manual (one-time) |
| GitHub Actions workflow | `.github/workflows/goreleaser.yml` | No changes needed |
| Formula directory | `Formula/` | GoReleaser (creates if missing) |
| Generated formula | `Formula/ado.rb` | GoReleaser (automated) |
| Installation docs | `docs/installation.md` | Manual (update) |
| README | `README.md` | Manual (update) |
| Homebrew guide | `docs/homebrew.md` | Manual (create) |

## Examples

### Example 1: First-Time Installation

```bash
# Add the tap (one-time setup)
$ brew tap anowarislam/ado
==> Tapping anowarislam/ado
Cloning into '/opt/homebrew/Library/Taps/anowarislam/ado'...
Tapped 1 formula (15 files, 1.2MB).

# Install ado
$ brew install ado
==> Downloading https://github.com/anowarislam/ado/releases/download/v1.4.0/ado_1.4.0_darwin_arm64.tar.gz
==> Downloading from https://objects.githubusercontent.com/...
######################################################################## 100.0%
==> Installing ado from anowarislam/ado
🍺  ado was successfully installed!

# Verify installation
$ ado meta info
Version:    1.4.0
Commit:     abc1234
Built:      2025-01-15T10:30:00Z
Go version: go1.23.5
Platform:   darwin/arm64

# Check location
$ which ado
/opt/homebrew/bin/ado
```

### Example 2: Upgrading to New Version

```bash
# Check for outdated packages
$ brew outdated
ado (1.4.0) < 1.5.0

# Upgrade ado
$ brew upgrade ado
==> Upgrading anowarislam/ado/ado 1.4.0 -> 1.5.0
==> Downloading https://github.com/anowarislam/ado/releases/download/v1.5.0/...
######################################################################## 100.0%
🍺  ado was successfully upgraded!

# Verify new version
$ ado meta info
Version:    1.5.0
```

### Example 3: Uninstalling

```bash
# Uninstall ado
$ brew uninstall ado
Uninstalling /opt/homebrew/Cellar/ado/1.4.0... (3 files, 8.2MB)

# Remove tap (optional)
$ brew untap anowarislam/ado
Untapping anowarislam/ado...
Untapped 1 formula (15 files, 1.2MB).
```

### Example 4: Local Testing (Development)

```bash
# Build snapshot without publishing
$ goreleaser release --snapshot --clean --skip=publish
  • releasing...
  • loading config file                              file=.goreleaser.yaml
  • building...                                      parallelism=12
  • homebrew formula
    • pushing                                        repository=anowarislam/ado
    • skipped                                        reason=skip_upload is set

# Test formula syntax
$ brew audit --strict --online anowarislam/ado/ado
anowarislam/ado/ado: passed

# Test installation from local tap
$ brew reinstall ado
==> Downloading https://github.com/anowarislam/ado/releases/download/v1.4.0-next/...
🍺  ado was successfully installed!
```

### Example 5: CI/CD Integration

```yaml
# .github/workflows/integration-test.yml
name: Integration Test
on: [push]

jobs:
  test-homebrew-install:
    runs-on: macos-latest
    steps:
      - name: Install ado via Homebrew
        run: |
          brew tap anowarislam/ado
          brew install ado

      - name: Verify installation
        run: |
          ado meta info
          ado --version

      - name: Test basic functionality
        run: |
          ado echo "Hello from CI"
```

## Edge Cases and Error Handling

| Scenario | Expected Behavior | Resolution |
|----------|------------------|------------|
| **Formula commit fails** | GoReleaser fails with "permission denied" | Verify GITHUB_TOKEN has `contents: write` permission |
| **Checksum mismatch** | Homebrew install fails with "checksum mismatch" | User reports issue; regenerate release |
| **Network failure during install** | Homebrew retries download, then fails | User runs `brew install ado` again |
| **Formula syntax error** | `brew audit` fails | Fix `.goreleaser.yaml` config and re-release |
| **Platform mismatch** | Homebrew downloads wrong binary | GoReleaser handles via platform detection |
| **Old Homebrew version** | Formula may not work with legacy Homebrew | Document minimum Homebrew version (4.0.0+) |
| **Tap not added** | `brew install ado` fails with "No available formula" | Error message guides user to add tap first |
| **Multiple architectures (Intel/ARM)** | Homebrew auto-selects correct binary | GoReleaser generates separate URLs per arch |
| **Formula directory missing** | GoReleaser creates `Formula/` automatically | No action needed |
| **Git conflicts on Formula/ado.rb** | GoReleaser commit fails | Rare; manually resolve and re-run release |

## Testing Strategy

### Unit Tests

Not applicable - this feature is build/release configuration, not runtime code.

### Integration Tests

- [ ] **Test GoReleaser snapshot build**: Verify formula generation without publishing
  ```bash
  goreleaser release --snapshot --clean --skip=publish
  ```

- [ ] **Test formula validation**: Run Homebrew audit on generated formula
  ```bash
  brew audit --strict --online anowarislam/ado/ado
  ```

- [ ] **Test formula style**: Verify Ruby syntax and Homebrew style guide
  ```bash
  brew style Formula/ado.rb
  ```

- [ ] **Test installation**: Install from tap and verify functionality
  ```bash
  brew tap anowarislam/ado
  brew install ado
  ado meta info
  ```

- [ ] **Test upgrade**: Simulate version upgrade path
  ```bash
  brew upgrade ado
  ```

- [ ] **Test uninstall**: Verify clean removal
  ```bash
  brew uninstall ado
  which ado  # Should return nothing
  ```

### Manual Testing Checklist

#### Pre-Release Testing

- [ ] Create `Formula/` directory in repository
- [ ] Add `.gitignore` entry to NOT ignore `Formula/` (if accidentally excluded)
- [ ] Add `brews` section to `.goreleaser.yaml`
- [ ] Verify GITHUB_TOKEN has `contents: write` permission (already configured via GitHub App)
- [ ] Run `goreleaser release --snapshot` locally to test formula generation
- [ ] Verify generated formula syntax: `brew audit --strict --online`
- [ ] Create test release tag (e.g., `v1.4.0-rc1`) to trigger GoReleaser
- [ ] Verify formula pushed to same repository automatically
- [ ] Install from tap: `brew tap anowarislam/ado && brew install ado`
- [ ] Verify binary works: `ado meta info`
- [ ] Test upgrade path by creating new test release

#### Post-Release Validation

- [ ] Verify formula updated in repository after real release
- [ ] Test installation on Intel Mac (if available)
- [ ] Test installation on Apple Silicon Mac
- [ ] Test installation on Linux with Homebrew (Linuxbrew)
- [ ] Verify checksums match GitHub Release artifacts
- [ ] Verify download URLs resolve correctly
- [ ] Test `brew upgrade ado` with actual version bump
- [ ] Check caveats message displays correctly
- [ ] Verify binary installed to correct PATH location
- [ ] Test uninstall leaves no artifacts: `brew uninstall ado`
- [ ] Run `brew test ado` to verify formula test block

#### Documentation Validation

- [ ] Update `README.md` with Homebrew installation method
- [ ] Update `docs/installation.md` (make Homebrew primary for macOS)
- [ ] Create `docs/homebrew.md` with detailed guide
- [ ] Update `mkdocs.yml` navigation to include Homebrew docs
- [ ] Verify all installation commands work as documented
- [ ] Test installation flow from user perspective (new user)

## Implementation Checklist

### Phase 1: Repository Setup

- [ ] Create `Formula/` directory in `anowarislam/ado` repository
  ```bash
  mkdir Formula
  git add Formula/.gitkeep  # Keep directory in git
  ```
- [ ] Verify `Formula/` is not in `.gitignore`
- [ ] Commit with message: `build: add Formula directory for Homebrew tap`

### Phase 2: GoReleaser Configuration

- [ ] Add `brews` section to `.goreleaser.yaml`
- [ ] Configure repository (same repo: `anowarislam/ado`)
- [ ] Configure formula metadata (name, description, license)
- [ ] Add caveats message with getting started instructions
- [ ] Add test block for formula verification
- [ ] Set `skip_upload: auto` for snapshot builds
- [ ] **No changes needed** to `.github/workflows/goreleaser.yml` (existing GITHUB_TOKEN works)
- [ ] Commit changes with message: `build(release): add Homebrew formula distribution`

### Phase 3: Testing & Validation

- [ ] Run local snapshot build: `goreleaser release --snapshot --clean --skip=publish`
- [ ] Verify formula generated in `Formula/ado.rb`
- [ ] Verify formula syntax: `brew audit --strict --online Formula/ado.rb`
- [ ] Create test release tag to trigger full workflow (e.g., `v1.4.0-rc1`)
- [ ] Monitor GitHub Actions for successful GoReleaser run
- [ ] Verify formula pushed to `anowarislam/ado` repository
- [ ] Test installation: `brew tap anowarislam/ado && brew install ado`
- [ ] Verify `ado meta info` output
- [ ] Test upgrade with new test release
- [ ] Run `brew test ado` to verify formula test block
- [ ] Document any issues or gotchas in `docs/homebrew.md`

### Phase 4: Documentation Updates

- [ ] Update `README.md`:
  - Add Homebrew section to Installation
  - Show `brew install` as primary method for macOS
  - Example: `brew tap anowarislam/ado && brew install ado`
- [ ] Update `docs/installation.md`:
  - Reorder: Homebrew first for macOS, then binary download
  - Add platform-specific guidance (Intel vs ARM)
  - Document tap addition step
- [ ] Create `docs/homebrew.md`:
  - Detailed Homebrew installation guide
  - In-repo tap explanation
  - Troubleshooting common issues
  - FAQ (tap vs homebrew-core, updates, uninstall)
  - Advanced topics (offline installs, tap removal)
- [ ] Update `mkdocs.yml`:
  - Add "Homebrew Distribution" under "Getting Started" or "Installation"
  - Link to `docs/homebrew.md`
- [ ] Update `CLAUDE.md`:
  - Document Homebrew distribution in Quick Commands
  - Note GoReleaser `brews` integration
  - Document that Formula/ directory is managed by GoReleaser
- [ ] Announce in `CHANGELOG.md` (automated via Release Please)

### Phase 5: Future Enhancements (Optional)

- [ ] Add shell completions to formula:
  ```yaml
  bash_completion.install "completions/ado.bash" => "ado"
  zsh_completion.install "completions/_ado"
  fish_completion.install "completions/ado.fish"
  ```
- [ ] Add man pages to formula:
  ```yaml
  man1.install "man/ado.1"
  ```
- [ ] Consider homebrew-core submission:
  - Requires source builds (not precompiled)
  - Requires dependency vendoring
  - Higher discoverability (`brew install ado` without tap)
  - Review process (1-2 weeks)
- [ ] Add livecheck for automatic version detection:
  ```ruby
  livecheck do
    url :stable
    strategy :github_latest
  end
  ```

## Open Questions

All questions resolved via research in issue #54:

- [x] **Should we use `brews` or `homebrew_casks`?**
  **Answer**: `brews` (formulae) - for CLI tools like `ado`. Casks are for GUI applications. Formulae install precompiled binaries via GoReleaser.

- [x] **Separate tap repository or in-repo tap?**
  **Answer**: In-repo tap (`anowarislam/ado`) - simpler, no separate repository to maintain, no separate token needed. Uses existing GITHUB_TOKEN.

- [x] **How to handle token authentication?**
  **Answer**: Use existing GITHUB_TOKEN (GitHub App token already configured in `.github/workflows/goreleaser.yml`). No separate token needed since formula is committed to same repository.

- [x] **How to test before releasing?**
  **Answer**: Use `goreleaser release --snapshot --clean --skip=publish` for local testing. Validate with `brew audit --strict --online`.

- [x] **What happens when GoReleaser `brews` is removed in v3?**
  **Answer**: Non-issue - we're using `homebrew_casks` from the start, which is the future-proof approach.

## Changelog

| Date | Change | Author |
|------|--------|--------|
| 2025-11-26 | Initial draft | @anowarislam |
