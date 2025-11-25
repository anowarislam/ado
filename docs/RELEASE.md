# Release Automation Guide

This document describes the fully automated release pipeline for this project. The system automatically creates releases, generates changelogs, and builds cross-platform binaries with zero manual intervention.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Components](#components)
- [Prerequisites](#prerequisites)
- [Setup Guide](#setup-guide)
  - [Step 1: Create GitHub App](#step-1-create-github-app)
  - [Step 2: Configure Repository Secrets](#step-2-configure-repository-secrets)
  - [Step 3: Create Workflow Files](#step-3-create-workflow-files)
  - [Step 4: Create Configuration Files](#step-4-create-configuration-files)
- [How It Works](#how-it-works)
- [Conventional Commits](#conventional-commits)
- [Version Bumping Rules](#version-bumping-rules)
- [Changelog Sections](#changelog-sections)
- [Manual Operations](#manual-operations)
- [Troubleshooting](#troubleshooting)
- [Security Considerations](#security-considerations)

---

## Overview

This release automation system provides:

- **Automated version bumping** based on conventional commits
- **Automated changelog generation** organized by commit type
- **Automated GitHub Release creation** with release notes
- **Cross-platform binary builds** (Linux, macOS, Windows × amd64, arm64)
- **Checksum generation** for release verification
- **Zero manual intervention** after initial setup

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         RELEASE AUTOMATION FLOW                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   Developer pushes    →   Release PR    →   Merge PR   →   GitHub Release   │
│   to main branch          created           approved       + Binaries        │
│                                                                              │
│   feat: add feature       v1.2.0            Automatic      ado_1.2.0_*      │
│   fix: bug fix            CHANGELOG.md      tag + release  checksums.txt    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Architecture

### High-Level Flow

```mermaid
flowchart TB
    subgraph Developer["Developer Workflow"]
        A[Push commits to main] --> B{Commit types?}
        B -->|feat/fix/etc| C[Triggers Release Please]
        B -->|chore without scope| D[No release action]
    end

    subgraph ReleasePlease["Release Please Workflow"]
        C --> E[Analyze commits since last release]
        E --> F[Calculate next version]
        F --> G[Generate changelog]
        G --> H{Release PR exists?}
        H -->|No| I[Create Release PR]
        H -->|Yes| J[Update Release PR]
        I --> K[Wait for merge]
        J --> K
    end

    subgraph Release["Release Creation"]
        K --> L[PR merged]
        L --> M[Create Git tag]
        M --> N[Create GitHub Release]
        N --> O[Trigger GoReleaser]
    end

    subgraph GoReleaser["GoReleaser Workflow"]
        O --> P[Checkout tag]
        P --> Q[Run tests]
        Q --> R[Build binaries]
        R --> S[Create archives]
        S --> T[Generate checksums]
        T --> U[Upload to GitHub Release]
    end

    U --> V[Release Complete]

    style A fill:#e1f5fe
    style V fill:#c8e6c9
    style N fill:#fff9c4
    style U fill:#fff9c4
```

### Detailed Workflow Sequence

```mermaid
sequenceDiagram
    autonumber
    participant Dev as Developer
    participant GH as GitHub
    participant RP as Release Please
    participant GR as GoReleaser
    participant App as GitHub App

    Dev->>GH: Push commit to main
    GH->>RP: Trigger workflow (push event)
    RP->>App: Request token
    App->>RP: Return short-lived token
    RP->>GH: Analyze commits since last tag

    alt New releasable commits exist
        RP->>GH: Create/Update Release PR
        Note over GH: PR includes version bump<br/>and changelog updates
    end

    Dev->>GH: Merge Release PR
    GH->>RP: Trigger workflow (push event)
    RP->>App: Request token
    App->>RP: Return short-lived token
    RP->>GH: Create tag (v1.x.x)
    RP->>GH: Create GitHub Release

    GH->>GR: Trigger workflow (release event)
    GR->>App: Request token
    App->>GR: Return short-lived token
    GR->>GR: Build binaries (6 platforms)
    GR->>GR: Create archives + checksums
    GR->>GH: Upload assets to release

    Note over GH: Release complete with<br/>binaries and checksums
```

### Component Interaction

```mermaid
flowchart LR
    subgraph Triggers["Trigger Events"]
        T1[push to main]
        T2[release published]
        T3[workflow_dispatch]
    end

    subgraph Workflows["GitHub Actions Workflows"]
        W1[release-please.yml]
        W2[goreleaser.yml]
    end

    subgraph Config["Configuration Files"]
        C1[release-please-config.json]
        C2[.release-please-manifest.json]
        C3[.goreleaser.yaml]
    end

    subgraph Auth["Authentication"]
        A1[GitHub App]
        A2[APP_ID secret]
        A3[APP_PRIVATE_KEY secret]
    end

    subgraph Outputs["Release Artifacts"]
        O1[Git Tag]
        O2[GitHub Release]
        O3[CHANGELOG.md]
        O4[Binary Archives]
        O5[checksums.txt]
    end

    T1 --> W1
    T2 --> W2
    T3 --> W2

    W1 --> C1
    W1 --> C2
    W2 --> C3

    A1 --> A2
    A1 --> A3
    A2 --> W1
    A3 --> W1
    A2 --> W2
    A3 --> W2

    W1 --> O1
    W1 --> O2
    W1 --> O3
    W2 --> O4
    W2 --> O5

    style A1 fill:#ffcdd2
    style O2 fill:#c8e6c9
    style O4 fill:#c8e6c9
```

---

## Components

### 1. Release Please

[Release Please](https://github.com/googleapis/release-please) automates CHANGELOG generation, version bumps, and GitHub Releases based on conventional commits.

**Responsibilities:**
- Parse commit messages following conventional commits spec
- Determine next semantic version
- Generate/update CHANGELOG.md
- Create and maintain Release PRs
- Create GitHub Releases and tags when PRs are merged

### 2. GoReleaser

[GoReleaser](https://goreleaser.com/) builds Go binaries for multiple platforms and uploads them to GitHub Releases.

**Responsibilities:**
- Cross-compile binaries for multiple OS/architecture combinations
- Create distributable archives (tar.gz, zip)
- Generate checksums for verification
- Upload release assets to GitHub

### 3. GitHub App

A GitHub App provides secure, scoped authentication for the release workflows.

**Why GitHub App instead of PAT or GITHUB_TOKEN:**

| Method | Token Lifetime | Maintenance | Security | Release Creation |
|--------|---------------|-------------|----------|------------------|
| `GITHUB_TOKEN` | Workflow run | None | Good | **Broken on new repos** |
| Personal Access Token | Up to 1 year | Annual rotation | Moderate | Works |
| **GitHub App** | 1 hour (auto) | **None** | **Best** | **Works** |

> **Note:** There is a [known GitHub platform bug](https://github.com/orgs/community/discussions/180369) where `GITHUB_TOKEN` cannot create releases on newer repositories. Using a GitHub App is the recommended workaround.

---

## Prerequisites

Before setting up the release automation, ensure you have:

- [ ] A GitHub repository with admin access
- [ ] Go project with a valid `go.mod` file
- [ ] Commits following [Conventional Commits](https://www.conventionalcommits.org/) specification
- [ ] GitHub CLI (`gh`) installed locally (optional, for testing)

---

## Setup Guide

### Step 1: Create GitHub App

#### 1.1 Navigate to GitHub App Creation

Go to: **https://github.com/settings/apps/new**

#### 1.2 Fill in Basic Information

| Field | Value | Notes |
|-------|-------|-------|
| **GitHub App name** | `your-project-release-bot` | Must be unique across GitHub |
| **Homepage URL** | `https://github.com/your-username/your-repo` | Your repository URL |
| **Webhook** | ❌ Uncheck "Active" | Not needed for this use case |

#### 1.3 Configure Permissions

Scroll to **"Permissions"** section and set **Repository permissions**:

| Permission | Access Level | Purpose |
|------------|--------------|---------|
| **Contents** | Read and write | Create tags, releases, update files |
| **Pull requests** | Read and write | Create and update release PRs |
| **Metadata** | Read-only | Auto-selected, required |

Leave all other permissions as "No access".

#### 1.4 Create the App

- Scroll to bottom
- Under "Where can this GitHub App be installed?", select **"Only on this account"**
- Click **"Create GitHub App"**

#### 1.5 Note the App ID

After creation, you'll be on the app's settings page. Note the **App ID** (a number like `123456`) displayed near the top.

#### 1.6 Generate Private Key

- Scroll down to **"Private keys"** section
- Click **"Generate a private key"**
- A `.pem` file will automatically download
- **Store this file securely** - you'll need it for the next step

#### 1.7 Install the App

- In the left sidebar, click **"Install App"**
- Click **"Install"** next to your account
- Select **"Only select repositories"**
- Choose your repository from the dropdown
- Click **"Install"**

### Step 2: Configure Repository Secrets

Navigate to your repository's secrets: `https://github.com/YOUR-USERNAME/YOUR-REPO/settings/secrets/actions`

Add two repository secrets:

#### Secret 1: APP_ID

| Field | Value |
|-------|-------|
| **Name** | `APP_ID` |
| **Secret** | The App ID number from Step 1.5 |

#### Secret 2: APP_PRIVATE_KEY

| Field | Value |
|-------|-------|
| **Name** | `APP_PRIVATE_KEY` |
| **Secret** | Entire contents of the `.pem` file from Step 1.6 |

**To get the .pem contents:**
```bash
cat ~/Downloads/your-app-name.*.private-key.pem
```

Copy everything including the `-----BEGIN RSA PRIVATE KEY-----` and `-----END RSA PRIVATE KEY-----` lines.

### Step 3: Create Workflow Files

Create the following workflow files in `.github/workflows/`:

#### 3.1 Release Please Workflow

**File:** `.github/workflows/release-please.yml`

```yaml
name: Release Please

on:
  push:
    branches:
      - main

jobs:
  release-please:
    name: Release Please
    runs-on: ubuntu-latest
    steps:
      - name: Generate token
        id: generate-token
        uses: actions/create-github-app-token@v1
        with:
          app-id: ${{ secrets.APP_ID }}
          private-key: ${{ secrets.APP_PRIVATE_KEY }}

      - name: Run release-please
        id: release
        uses: googleapis/release-please-action@v4
        with:
          token: ${{ steps.generate-token.outputs.token }}
          config-file: release-please-config.json
          manifest-file: .release-please-manifest.json
```

#### 3.2 GoReleaser Workflow

**File:** `.github/workflows/goreleaser.yml`

```yaml
name: GoReleaser

on:
  release:
    types: [published]
  workflow_dispatch:
    inputs:
      tag:
        description: 'Tag to release (e.g., v1.0.0)'
        required: true
        type: string

jobs:
  goreleaser:
    name: Build and Release
    runs-on: ubuntu-latest
    steps:
      - name: Generate token
        id: generate-token
        uses: actions/create-github-app-token@v1
        with:
          app-id: ${{ secrets.APP_ID }}
          private-key: ${{ secrets.APP_PRIVATE_KEY }}

      - name: Checkout code
        uses: actions/checkout@v4
        with:
          fetch-depth: 0
          ref: ${{ github.event.inputs.tag || github.ref }}

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
          cache: true

      - name: Run tests
        run: go test ./...

      - name: Run goreleaser
        uses: goreleaser/goreleaser-action@v6
        with:
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ steps.generate-token.outputs.token }}
```

### Step 4: Create Configuration Files

#### 4.1 Release Please Config

**File:** `release-please-config.json`

```json
{
  "$schema": "https://raw.githubusercontent.com/googleapis/release-please/main/schemas/config.json",
  "packages": {
    ".": {
      "release-type": "go",
      "changelog-path": "CHANGELOG.md",
      "bump-minor-pre-major": true,
      "bump-patch-for-minor-pre-major": true,
      "extra-files": [],
      "changelog-sections": [
        {"type": "feat", "section": "Features"},
        {"type": "fix", "section": "Bug Fixes"},
        {"type": "perf", "section": "Performance Improvements"},
        {"type": "revert", "section": "Reverts"},
        {"type": "build", "section": "Dependencies"},
        {"type": "ci", "section": "CI/CD"},
        {"type": "refactor", "section": "Code Refactoring"},
        {"type": "docs", "section": "Documentation"}
      ]
    }
  },
  "include-component-in-tag": false,
  "include-v-in-tag": true
}
```

**Configuration explained:**

| Option | Value | Description |
|--------|-------|-------------|
| `release-type` | `"go"` | Configures for Go projects |
| `changelog-path` | `"CHANGELOG.md"` | Where changelog is written |
| `bump-minor-pre-major` | `true` | Before v1.0.0, `feat` bumps minor |
| `bump-patch-for-minor-pre-major` | `true` | Before v1.0.0, `fix` bumps patch |
| `include-v-in-tag` | `true` | Tags are `v1.0.0` not `1.0.0` |
| `changelog-sections` | `[...]` | Maps commit types to sections |

#### 4.2 Release Please Manifest

**File:** `.release-please-manifest.json`

```json
{
  ".": "0.0.0"
}
```

> **Note:** Set this to your current version, or `"0.0.0"` for new projects. Release Please will update this automatically.

#### 4.3 GoReleaser Config

**File:** `.goreleaser.yaml`

```yaml
version: 2

project_name: your-project

builds:
  - id: your-project
    binary: your-project
    main: ./cmd/your-project
    env:
      - CGO_ENABLED=0
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64
    ldflags:
      - -s -w
      - -X main.Version={{.Version}}
      - -X main.Commit={{.Commit}}
      - -X main.BuildTime={{.Date}}

archives:
  - format: tar.gz
    name_template: "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}"
    format_overrides:
      - goos: windows
        format: zip

checksum:
  name_template: checksums.txt

changelog:
  sort: asc
  filters:
    exclude:
      - "^docs:"
      - "^test:"
      - "^chore:"
      - Merge pull request
      - Merge branch

release:
  github:
    owner: your-username
    name: your-repo
  draft: false
  prerelease: auto
```

**Customize these values:**
- `project_name`: Your project name
- `main`: Path to your main package
- `ldflags`: Adjust to match your version variables
- `release.github.owner`: Your GitHub username
- `release.github.name`: Your repository name

---

## How It Works

### The Release Cycle

```mermaid
stateDiagram-v2
    [*] --> Development: Start
    Development --> CommitPushed: Push to main
    CommitPushed --> AnalyzeCommits: Release Please triggers

    AnalyzeCommits --> NoRelease: No releasable commits
    AnalyzeCommits --> CreatePR: Has feat/fix/etc

    NoRelease --> Development: Continue development

    CreatePR --> PROpen: Release PR created/updated
    PROpen --> PRMerged: Developer merges PR
    PROpen --> PROpen: More commits pushed

    PRMerged --> TagCreated: Git tag created
    TagCreated --> ReleaseCreated: GitHub Release created
    ReleaseCreated --> GoReleaserTriggers: release event fires

    GoReleaserTriggers --> BuildBinaries: Build for all platforms
    BuildBinaries --> UploadAssets: Upload to release
    UploadAssets --> [*]: Release complete
```

### Step-by-Step Breakdown

#### Phase 1: Commit Analysis

When you push to `main`:

1. Release Please workflow triggers
2. GitHub App token is generated (valid for 1 hour)
3. Release Please fetches all commits since the last release tag
4. Commits are parsed for conventional commit prefixes

#### Phase 2: Release PR Management

If releasable commits are found:

1. Next version is calculated based on commit types
2. CHANGELOG.md entries are generated
3. A Release PR is created or updated
4. PR includes:
   - Updated `.release-please-manifest.json`
   - Updated `CHANGELOG.md`
   - Any `extra-files` configured

#### Phase 3: Release Creation

When the Release PR is merged:

1. Release Please workflow triggers again
2. Detects the merged release PR
3. Creates a Git tag (e.g., `v1.2.0`)
4. Creates a GitHub Release with changelog as release notes

#### Phase 4: Binary Building

When the GitHub Release is published:

1. GoReleaser workflow triggers via `release: [published]` event
2. Checks out the tagged commit
3. Runs tests to ensure release quality
4. Builds binaries for all configured platforms
5. Creates archives (tar.gz for Unix, zip for Windows)
6. Generates checksums
7. Uploads all assets to the GitHub Release

---

## Conventional Commits

This system relies on [Conventional Commits](https://www.conventionalcommits.org/) for automatic versioning.

### Commit Message Format

```
<type>[(scope)][!]: <description>

[optional body]

[optional footer(s)]
```

### Examples

```bash
# Feature (bumps MINOR version)
feat: add user authentication

# Feature with scope
feat(api): add new /users endpoint

# Bug fix (bumps PATCH version)
fix: resolve null pointer in config parser

# Bug fix with scope
fix(cli): handle spaces in file paths

# Breaking change (bumps MAJOR version)
feat!: redesign configuration format

# Breaking change with scope
fix(api)!: change response format for /users

# Other types (included in changelog, no version bump alone)
docs: update API documentation
build(deps): upgrade cobra to v1.8.0
ci: add code coverage reporting
refactor: simplify error handling
perf: optimize database queries
```

### Commit Types Reference

| Type | Description | Version Bump | Changelog Section |
|------|-------------|--------------|-------------------|
| `feat` | New feature | MINOR | Features |
| `fix` | Bug fix | PATCH | Bug Fixes |
| `perf` | Performance improvement | PATCH | Performance Improvements |
| `revert` | Reverts a previous commit | PATCH | Reverts |
| `build` | Build system or dependencies | None* | Dependencies |
| `ci` | CI/CD changes | None* | CI/CD |
| `refactor` | Code refactoring | None* | Code Refactoring |
| `docs` | Documentation | None* | Documentation |
| `test` | Tests | None* | Not included |
| `chore` | Maintenance | None* | Not included |
| `style` | Code style | None* | Not included |

> *These types don't trigger a version bump on their own, but are included in the release if other releasable commits exist.

---

## Version Bumping Rules

```mermaid
flowchart TD
    A[Analyze Commits] --> B{Any BREAKING CHANGE?}
    B -->|Yes| C[MAJOR bump]
    B -->|No| D{Any feat commits?}
    D -->|Yes| E[MINOR bump]
    D -->|No| F{Any fix/perf commits?}
    F -->|Yes| G[PATCH bump]
    F -->|No| H[No version bump]

    C --> I[v1.0.0 → v2.0.0]
    E --> J[v1.0.0 → v1.1.0]
    G --> K[v1.0.0 → v1.0.1]
    H --> L[No release PR created]

    style C fill:#ffcdd2
    style E fill:#fff9c4
    style G fill:#c8e6c9
    style H fill:#e0e0e0
```

### Pre-1.0.0 Behavior

When the version is below 1.0.0 (configured via `bump-minor-pre-major` and `bump-patch-for-minor-pre-major`):

| Commit Type | Normal Behavior | Pre-1.0.0 Behavior |
|-------------|-----------------|---------------------|
| Breaking change | MAJOR bump | MINOR bump |
| `feat` | MINOR bump | MINOR bump |
| `fix` | PATCH bump | PATCH bump |

---

## Changelog Sections

The changelog is organized by commit type. Each section corresponds to a commit type prefix.

### Example CHANGELOG.md

```markdown
# Changelog

## [1.2.0](https://github.com/user/repo/compare/v1.1.0...v1.2.0) (2024-01-15)

### Features

* **api:** add batch processing endpoint ([abc1234](https://github.com/user/repo/commit/abc1234))
* **cli:** support configuration via environment variables ([def5678](https://github.com/user/repo/commit/def5678))

### Bug Fixes

* **parser:** handle UTF-8 BOM in input files ([111aaaa](https://github.com/user/repo/commit/111aaaa))

### Dependencies

* **deps:** bump golang.org/x/text from 0.13.0 to 0.14.0 ([222bbbb](https://github.com/user/repo/commit/222bbbb))

### CI/CD

* add automated security scanning ([333cccc](https://github.com/user/repo/commit/333cccc))
```

### Customizing Sections

To modify which commit types appear in the changelog, edit `changelog-sections` in `release-please-config.json`:

```json
"changelog-sections": [
  {"type": "feat", "section": "Features"},
  {"type": "fix", "section": "Bug Fixes"},
  {"type": "perf", "section": "Performance Improvements"},
  {"type": "revert", "section": "Reverts"},
  {"type": "build", "section": "Dependencies"},
  {"type": "ci", "section": "CI/CD"},
  {"type": "refactor", "section": "Code Refactoring"},
  {"type": "docs", "section": "Documentation"}
]
```

To **hide** a type from the changelog while still tracking it:

```json
{"type": "chore", "section": "Miscellaneous", "hidden": true}
```

---

## Manual Operations

### Manually Triggering GoReleaser

If you need to rebuild release assets:

```bash
# Via GitHub CLI
gh workflow run goreleaser.yml -f tag=v1.2.0

# Or via GitHub UI
# Go to Actions → GoReleaser → Run workflow → Enter tag
```

### Creating a Release Without New Commits

If you need to create a release PR manually:

1. Create an empty commit with a releasable type:
   ```bash
   git commit --allow-empty -m "fix: trigger release"
   git push origin main
   ```

2. Or manually create the release via GitHub UI

### Skipping CI for a Commit

Add `[skip ci]` to your commit message:

```bash
git commit -m "docs: update README [skip ci]"
```

---

## Troubleshooting

### Common Issues

#### 1. "author_id does not have push access"

**Cause:** GitHub platform bug with `GITHUB_TOKEN` on newer repositories.

**Solution:** Use GitHub App authentication as described in this guide.

#### 2. Release PR Not Created

**Possible causes:**
- No releasable commits (only `chore`, `test`, `style` commits)
- Commits don't follow conventional commit format
- Release Please manifest has wrong version

**Debug steps:**
```bash
# Check recent commits
git log --oneline -10

# Verify commit format
git log -1 --format="%B"

# Check manifest version
cat .release-please-manifest.json
```

#### 3. GoReleaser Fails

**Possible causes:**
- Tests failing
- Invalid `.goreleaser.yaml` configuration
- Missing Go dependencies

**Debug steps:**
```bash
# Test locally
goreleaser check
goreleaser build --snapshot --clean

# Run tests
go test ./...
```

#### 4. "Duplicate release tag"

**Cause:** Tag or release already exists.

**Solution:** This is usually safe to ignore. It means release-please is trying to process an already-released version.

#### 5. GitHub App Token Issues

**Debug steps:**
```bash
# Verify secrets are set
gh secret list

# Check workflow run logs
gh run view <run-id> --log
```

### Viewing Workflow Logs

```bash
# List recent workflow runs
gh run list --limit 10

# View specific run
gh run view <run-id>

# View failed job logs
gh run view <run-id> --log-failed
```

---

## Security Considerations

### GitHub App Security

1. **Minimal Permissions:** The app only has `contents` and `pull-requests` write access
2. **Short-lived Tokens:** App tokens expire after 1 hour
3. **Repository Scoped:** App is installed only on specific repositories
4. **Audit Trail:** All actions are logged under the app's identity

### Secret Management

1. **Never commit secrets:** The `.pem` file and App ID should never be in version control
2. **Rotate if compromised:** If the private key is exposed, regenerate it immediately in the app settings
3. **Use environment-specific secrets:** For organizations, consider using environment-level secrets

### Binary Verification

Users can verify downloaded binaries using the checksums:

```bash
# Download checksum file
curl -LO https://github.com/user/repo/releases/download/v1.0.0/checksums.txt

# Verify a specific file
sha256sum -c checksums.txt --ignore-missing
```

---

## File Structure Summary

After setup, your repository should have:

```
.
├── .github/
│   └── workflows/
│       ├── release-please.yml    # Creates releases
│       └── goreleaser.yml        # Builds binaries
├── .goreleaser.yaml              # GoReleaser configuration
├── .release-please-manifest.json # Current version tracking
├── release-please-config.json    # Release Please configuration
├── CHANGELOG.md                  # Auto-generated changelog
└── ...
```

---

## Quick Reference

### Commit → Version Mapping

```
feat: ...           → v1.0.0 → v1.1.0  (MINOR)
fix: ...            → v1.0.0 → v1.0.1  (PATCH)
feat!: ...          → v1.0.0 → v2.0.0  (MAJOR)
build(deps): ...    → (no bump, included in next release)
```

### Release Flow Commands

```bash
# Check for open release PR
gh pr list --label "autorelease: pending"

# View release PR
gh pr view <number>

# Merge release PR (triggers release)
gh pr merge <number> --squash

# Check release status
gh release list

# View specific release
gh release view v1.0.0
```

---

## References

- [Release Please Documentation](https://github.com/googleapis/release-please)
- [GoReleaser Documentation](https://goreleaser.com/intro/)
- [Conventional Commits Specification](https://www.conventionalcommits.org/)
- [GitHub Apps Documentation](https://docs.github.com/en/apps)
- [actions/create-github-app-token](https://github.com/actions/create-github-app-token)
