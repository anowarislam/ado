# Chapter 7: GitHub Apps, Actions, and Integrations

> From basic GITHUB_TOKEN to production-grade GitHub Apps - mastering GitHub's integration ecosystem for automated workflows.

## Table of Contents

- [Introduction](#introduction)
- [GitHub Authentication Methods](#github-authentication-methods)
- [GitHub Apps Deep Dive](#github-apps-deep-dive)
- [Creating Your GitHub App](#creating-your-github-app)
- [GitHub App Authentication](#github-app-authentication)
- [GitHub Actions Integration](#github-actions-integration)
- [Dependabot Configuration](#dependabot-configuration)
- [Dependabot Auto-Merge](#dependabot-auto-merge)
- [Branch Protection Rules](#branch-protection-rules)
- [GitHub Packages (GHCR)](#github-packages-ghcr)
- [GitHub Pages Deployment](#github-pages-deployment)
- [Codecov Integration](#codecov-integration)
- [Security Integrations](#security-integrations)
- [Permission Matrices](#permission-matrices)
- [Integration Architecture](#integration-architecture)
- [Troubleshooting Guide](#troubleshooting-guide)

---

## Introduction

### Why This Chapter Matters

GitHub provides a rich ecosystem of integrations, but understanding **when to use what** is critical:

```mermaid
graph TB
    A[GitHub Integration Need] --> B{Type of Operation?}

    B -->|Read-only CI checks| C[GITHUB_TOKEN]
    B -->|Trigger other workflows| D[GitHub App]
    B -->|Long-lived automation| D
    B -->|Security-sensitive operations| D

    C --> E[Simple, built-in, limited]
    D --> F[Complex, secure, flexible]

    style C fill:#e1f5e1
    style D fill:#ffe1e1
```

**Key Insight**: The default `GITHUB_TOKEN` is convenient but limited. GitHub Apps unlock advanced automation while maintaining security best practices.

### What You'll Learn

By the end of this chapter, you'll understand:

1. **When to use each authentication method** - PAT vs GITHUB_TOKEN vs GitHub App
2. **How to create and configure GitHub Apps** - step-by-step with screenshots
3. **How to integrate GitHub Apps with Actions** - using `actions/create-github-app-token`
4. **How to configure Dependabot** - for Go, Python, and GitHub Actions
5. **How to implement auto-merge workflows** - safely merging dependency updates
6. **How to configure branch protection** - requiring checks and reviews
7. **How to use GitHub Packages** - publishing container images to GHCR
8. **How to deploy GitHub Pages** - automated documentation deployment
9. **How to integrate Codecov** - coverage reporting and enforcement
10. **How to enable security features** - Dependabot alerts and code scanning

### Real-World Context: The ado Project

Throughout this chapter, we'll reference the `ado` project's actual GitHub integrations:

| Integration | Purpose | Authentication Method |
|-------------|---------|----------------------|
| CI Pipeline | Tests, builds, coverage | `GITHUB_TOKEN` |
| Release Please | Version bumping | GitHub App |
| GoReleaser | Binary releases | GitHub App |
| Dependabot | Dependency updates | Built-in |
| Dependabot Auto-Merge | Automated merging | `GITHUB_TOKEN` |
| GitHub Pages | Documentation deployment | `GITHUB_TOKEN` + OIDC |
| GHCR | Container image storage | `GITHUB_TOKEN` |
| Codecov | Coverage reporting | Codecov token |

---

## GitHub Authentication Methods

### The Three Authentication Methods

GitHub provides three primary ways to authenticate in automated workflows:

#### 1. GITHUB_TOKEN (Built-in)

**What it is**: Automatically created for each workflow run by GitHub Actions.

**Characteristics**:
- Automatically available as `${{ secrets.GITHUB_TOKEN }}`
- Scoped to the current repository
- Expires when the workflow completes
- Cannot trigger other workflows (security limitation)
- Permissions can be configured per workflow

**Example**:
```yaml
name: CI
on: [push]

permissions:
  contents: read
  pull-requests: write

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v6
        with:
          token: ${{ secrets.GITHUB_TOKEN }}
```

**Best for**:
- CI checks (tests, linting, builds)
- Reading repository contents
- Commenting on PRs
- Uploading artifacts
- Basic PR operations

**Limitations**:
- Cannot trigger workflows (prevents recursive triggers)
- Cannot access organization-level resources
- Cannot bypass branch protection rules
- Short-lived (workflow duration only)

#### 2. Personal Access Token (PAT)

**What it is**: A token tied to a specific user account.

**Characteristics**:
- Created manually in user settings
- Can be classic or fine-grained
- Tied to user permissions and rate limits
- Requires manual rotation
- Security risk if user leaves organization

**Example**:
```yaml
jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v6
        with:
          token: ${{ secrets.PERSONAL_ACCESS_TOKEN }}
```

**Best for**:
- **Nothing in production** - use GitHub Apps instead
- Personal projects (with caution)
- Quick prototyping

**Problems**:
- Tied to individual user (breaks if user leaves)
- Hard to audit ("who did this?")
- All permissions tied to user
- Manual rotation required
- Security risk if compromised

#### 3. GitHub App (Recommended for Production)

**What it is**: A first-class integration that authenticates as itself, not a user.

**Characteristics**:
- Identity independent of users
- Fine-grained permissions
- Can trigger workflows
- Can bypass branch protection (if configured)
- Audit trail shows "app" not "user"
- Tokens expire quickly (default: 1 hour)

**Example**:
```yaml
jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - name: Generate token
        id: app-token
        uses: actions/create-github-app-token@v1
        with:
          app-id: ${{ secrets.APP_ID }}
          private-key: ${{ secrets.APP_PRIVATE_KEY }}

      - uses: actions/checkout@v6
        with:
          token: ${{ steps.app-token.outputs.token }}
```

**Best for**:
- Production workflows
- Release automation
- Cross-repository operations
- Triggering workflows
- Long-term automation
- Organization-wide operations

**Advantages**:
- Not tied to individual users
- Fine-grained permissions
- Can trigger workflows
- Better audit trail
- Automatic token expiration
- Can act on behalf of the app

### Decision Matrix

| Need | Use | Why |
|------|-----|-----|
| Run tests in CI | `GITHUB_TOKEN` | Simple, secure, sufficient |
| Trigger release workflow | GitHub App | Can bypass workflow trigger limitations |
| Merge release PR | GitHub App | Can bypass branch protection if needed |
| Comment on PRs | `GITHUB_TOKEN` | Built-in permissions sufficient |
| Cross-repo automation | GitHub App | Organization-level permissions |
| Publish packages | `GITHUB_TOKEN` | Sufficient for GHCR |
| Personal project | `GITHUB_TOKEN` or PAT | Simplicity acceptable |
| Production system | GitHub App | Required for reliability |

### Architecture Comparison

```mermaid
graph TB
    subgraph "GITHUB_TOKEN Flow"
        A1[Workflow Starts] --> A2[Token Auto-Created]
        A2 --> A3[Use in Steps]
        A3 --> A4[Token Expires]
        A4 -.->|Cannot| A5[Trigger Workflows]
    end

    subgraph "GitHub App Flow"
        B1[Workflow Starts] --> B2[Request Token]
        B2 --> B3[App Authenticates]
        B3 --> B4[Token Generated]
        B4 --> B5[Use in Steps]
        B5 --> B6[Token Expires 1hr]
        B6 -.->|Can| B7[Trigger Workflows]
    end

    style A2 fill:#e1f5e1
    style B3 fill:#ffe1e1
```

**Key Insight**: Start with `GITHUB_TOKEN` for everything. Upgrade to GitHub App only when you hit its limitations (typically: triggering workflows or needing bypass capabilities).

---

## GitHub Apps Deep Dive

### What is a GitHub App?

A GitHub App is a first-class GitHub integration that:
- Has its own identity (not tied to a user)
- Requests specific permissions
- Can be installed on organizations or repositories
- Uses JWT + installation tokens for authentication
- Provides better audit trails
- Respects rate limits separately from users

### GitHub App vs OAuth App

| Feature | GitHub App | OAuth App |
|---------|-----------|-----------|
| **Identity** | App itself | User account |
| **Permissions** | Fine-grained | Broad scopes |
| **Installation** | Per org/repo | Per user |
| **Rate limits** | Separate | User's limit |
| **Audit trail** | Shows app name | Shows user |
| **Token lifetime** | 1 hour (default) | Indefinite |
| **Best for** | Automation | User delegation |

### Why GitHub Apps for Release Automation?

**The Problem**: Release automation needs to:
1. Create release PRs (requires write permissions)
2. Merge release PRs (may need to bypass protection)
3. Trigger downstream workflows (blocked by GITHUB_TOKEN)
4. Maintain audit trail (who released what?)

**The Solution**: GitHub App provides:
- Identity independent of team members
- Can trigger workflows (not blocked)
- Can bypass branch protection (if configured)
- Clear audit trail (app name in logs)
- Automatic token expiration (security)

### Real-World Example: ado's Release Flow

```mermaid
sequenceDiagram
    participant Dev as Developer
    participant Main as main branch
    participant RP as Release Please
    participant App as GitHub App
    participant GR as GoReleaser

    Dev->>Main: Push commit (feat:)
    Main->>RP: Trigger release-please.yml
    RP->>App: Request token
    App->>RP: Issue 1-hour token
    RP->>Main: Create/update release PR

    Dev->>Main: Merge release PR
    Main->>GR: Trigger goreleaser.yml
    GR->>App: Request token
    App->>GR: Issue 1-hour token
    GR->>Main: Create GitHub release
    GR->>Main: Upload assets
    GR->>Main: Push container images
```

**Why This Works**:
1. Release Please can trigger on push events (App bypasses limitations)
2. Release PR updates don't require manual intervention
3. GoReleaser can create releases and trigger subsequent workflows
4. Audit trail shows "release-please[bot]" not individual users
5. Tokens expire automatically (no long-lived credentials)

### GitHub App Architecture

```mermaid
graph TB
    subgraph "GitHub App Components"
        A[App Registration] --> B[App ID]
        A --> C[Private Key]
        A --> D[Permissions]
        A --> E[Webhook URL]
    end

    subgraph "Installation"
        F[Install on Org/Repo] --> G[Installation ID]
        G --> H[Installation Token]
    end

    subgraph "Authentication Flow"
        I[Create JWT] --> J[Sign with Private Key]
        J --> K[Request Installation Token]
        K --> L[Use Token for API Calls]
    end

    B --> I
    C --> J
    G --> K
    D --> L

    style A fill:#e1f5e1
    style H fill:#ffe1e1
    style L fill:#e1e5ff
```

**Components Explained**:
- **App ID**: Public identifier for your app
- **Private Key**: Secret used to sign JWTs
- **Installation ID**: Links app to specific org/repo
- **JWT**: Short-lived token proving app identity
- **Installation Token**: API token with app permissions (1 hour default)

---

## Creating Your GitHub App

### Step-by-Step Setup Guide

#### Step 1: Navigate to GitHub App Settings

**UI Path**:
1. Click your profile picture (top right)
2. Settings → Developer settings → GitHub Apps
3. Click "New GitHub App"

**Screenshot Description**:
> A list view showing existing GitHub Apps with a prominent "New GitHub App" button in green at the top right of the page.

#### Step 2: Basic Information

Fill in the basic app details:

**Form Fields**:

```yaml
GitHub App name: "Release Automation for MyProject"
# This name must be unique across GitHub
# Choose something descriptive but not generic

Homepage URL: "https://github.com/myorg/myproject"
# URL users see when viewing the app
# Typically your project's repository

Description: "Automated release management via release-please and GoReleaser"
# What does your app do? Be concise.

Webhook:
  Active: ❌ Unchecked
  # Only needed if you want GitHub to call your server
  # For GitHub Actions integration, you don't need webhooks
```

**Screenshot Description**:
> Form with text fields for App name, Homepage URL, and Description. Webhook section with Active checkbox unchecked and a placeholder URL field (greyed out).

#### Step 3: Permissions Configuration

This is the most critical step. GitHub Apps use **fine-grained permissions**.

**Permission Philosophy**:
- **Principle of least privilege**: Only request what you need
- **Repository vs Organization**: Start with repository-only
- **Read vs Write**: Start with read, upgrade if needed

**For Release Automation (Release Please + GoReleaser)**:

```yaml
Repository permissions:
  Contents: Read and write
    # Why: Create/update release PRs, create releases
  Pull requests: Read and write
    # Why: Create/update/merge release PRs
  Metadata: Read-only
    # Why: Required automatically by GitHub
  Actions: Read and write (optional)
    # Why: If you want to trigger workflow reruns
  Issues: Read and write (optional)
    # Why: If you want to auto-close issues on release

Organization permissions:
  # Leave all as "No access" unless you need org-wide operations
```

**Common Permission Patterns**:

| Use Case | Permissions Needed |
|----------|-------------------|
| Release automation | Contents: write, Pull requests: write |
| CI status checks | Checks: write, Statuses: write |
| Issue management | Issues: write, Pull requests: write |
| Package publishing | Packages: write |
| Workflow dispatch | Actions: write |

**Screenshot Description**:
> Expandable sections for "Repository permissions" and "Organization permissions". Each permission shows a dropdown with options: No access, Read-only, Read and write. Contents and Pull requests set to "Read and write", Metadata automatically set to "Read-only" (greyed out).

#### Step 4: Installation Options

```yaml
Where can this GitHub App be installed?
  ○ Only on this account
  ● Any account
```

**Choose based on**:
- **Only on this account**: Private app, just for your org
- **Any account**: If you plan to share the app (e.g., open source tool)

**For most projects**: Choose "Only on this account" for security.

**Screenshot Description**:
> Radio button selection with two options. "Only on this account" selected. Helper text below explains this app can only be installed on the account that created it.

#### Step 5: Create the App

Click **"Create GitHub App"** button at bottom.

**What happens next**:
1. App is created with a unique App ID
2. You're redirected to the app's settings page
3. You can now generate a private key

**Screenshot Description**:
> Success page showing "GitHub App created successfully" message with the app's name and a large App ID displayed (e.g., "App ID: 123456"). Below are tabs for General, Permissions, Install App, Advanced.

#### Step 6: Generate Private Key

**This is critical for authentication**:

1. Scroll down to "Private keys" section
2. Click "Generate a private key"
3. A `.pem` file downloads automatically
4. **SAVE THIS FILE SECURELY** - you can't download it again

**Screenshot Description**:
> Section titled "Private keys" with explanatory text: "Private keys are used to sign access token requests." A button labeled "Generate a private key" in green. Below, a list showing "No private keys" if none exist yet, or a list of keys with creation date and fingerprint if any exist.

**Security Warning**:
```
⚠️  CRITICAL: The private key file that downloads is the ONLY copy.
- Store it securely (password manager, secrets vault)
- Never commit it to git
- Never share it
- If lost, you must generate a new key
```

**File contents** (example):
```pem
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA2Z3qX+0FAKEC0dGToEXAMPLEKEYiLCXnSq7vL9...
(many more lines of base64-encoded key data)
...
-----END RSA PRIVATE KEY-----
```

#### Step 7: Note Your App ID

Your App ID is visible at the top of the settings page:

```
App ID: 123456
```

**You'll need two secrets for GitHub Actions**:
1. `APP_ID`: The numeric ID (e.g., 123456)
2. `APP_PRIVATE_KEY`: The entire contents of the `.pem` file

#### Step 8: Install the App

**Installation grants the app access to specific repositories**:

1. Click "Install App" tab (top of page)
2. Click "Install" next to your organization/account
3. Choose installation scope:
   - **All repositories**: App can access everything
   - **Only select repositories**: Choose specific repos

**For security**: Start with "Only select repositories" and choose your project.

**Screenshot Description**:
> Page showing "Install AppName" with your organization name. Two radio buttons: "All repositories" and "Only select repositories". Below the second option, a searchable dropdown list of repositories with checkboxes. A green "Install" button at bottom.

4. Click "Install"
5. Authorize the installation (may require password re-entry)

**Post-installation**: You'll see a confirmation screen with:
- Installation ID (visible in URL)
- List of repositories with access
- Ability to configure installation later

**Screenshot Description**:
> Confirmation page showing "AppName was installed successfully" with a list of repositories it can access. Each repo shows its name and a gear icon for settings. At top, a note: "This installation gives AppName access to: [list of permissions]".

---

## GitHub App Authentication

### Authentication Flow Overview

GitHub Apps use a **two-step authentication process**:

```mermaid
sequenceDiagram
    participant W as Workflow
    participant A as actions/create-github-app-token
    participant GH as GitHub API

    W->>A: Provide App ID + Private Key
    A->>A: Generate JWT (signed with private key)
    A->>GH: Request installation token (with JWT)
    GH->>GH: Verify JWT signature
    GH->>GH: Check app permissions
    GH->>A: Return installation token
    A->>W: Output token (expires in 1 hour)
    W->>GH: Use token for API calls
```

**Two token types**:
1. **JWT (JSON Web Token)**: Proves the app's identity (short-lived, seconds)
2. **Installation Token**: Used for actual API calls (default: 1 hour)

### Storing Secrets in GitHub

Before using the app in workflows, you must store two secrets:

#### Adding Repository Secrets

**UI Path**:
1. Go to your repository on GitHub
2. Settings → Secrets and variables → Actions
3. Click "New repository secret"

**Add these two secrets**:

**Secret 1: APP_ID**
```
Name: APP_ID
Value: 123456
```
(Use your actual App ID from step 7)

**Secret 2: APP_PRIVATE_KEY**
```
Name: APP_PRIVATE_KEY
Value: (paste entire contents of .pem file)
```

**Example**:
```
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA2Z3qX+0FAKEC0dGToEXAMPLEKEYiLCXnSq7vL9...
(paste all lines exactly as they appear in the .pem file)
...
-----END RSA PRIVATE KEY-----
```

**Screenshot Description**:
> Secrets page showing a list of secrets with masked values. A green "New repository secret" button at top right. Each secret row shows name, last updated date, and an "Update" button.

**Security Notes**:
- Secrets are encrypted at rest by GitHub
- Secrets are never exposed in logs (automatically masked)
- Only workflows can access repository secrets
- Consider using environment secrets for stricter control

### Using the GitHub App in Workflows

#### Method 1: Using actions/create-github-app-token

**This is the recommended approach** - it handles all the JWT complexity:

```yaml
name: Release Please

on:
  push:
    branches:
      - main

jobs:
  release-please:
    runs-on: ubuntu-latest
    steps:
      # Step 1: Generate token from GitHub App
      - name: Generate token
        id: generate-token
        uses: actions/create-github-app-token@v1
        with:
          app-id: ${{ secrets.APP_ID }}
          private-key: ${{ secrets.APP_PRIVATE_KEY }}

      # Step 2: Use the token in subsequent steps
      - name: Run release-please
        uses: googleapis/release-please-action@v4
        with:
          token: ${{ steps.generate-token.outputs.token }}
```

**How it works**:
1. The action reads `APP_ID` and `APP_PRIVATE_KEY` from secrets
2. Generates a JWT signed with the private key
3. Requests an installation token from GitHub
4. Outputs the token for use in later steps

**Token properties**:
- Expires in 1 hour (default)
- Has permissions defined in app configuration
- Can bypass `GITHUB_TOKEN` limitations

#### Method 2: Scoping to Specific Repositories

If your app has access to multiple repos but you only want the token for specific ones:

```yaml
- name: Generate token
  id: generate-token
  uses: actions/create-github-app-token@v1
  with:
    app-id: ${{ secrets.APP_ID }}
    private-key: ${{ secrets.APP_PRIVATE_KEY }}
    repositories: |
      myproject
      myproject-docs
```

#### Method 3: Organization-Level Tokens

For organization-wide operations:

```yaml
- name: Generate org token
  id: generate-token
  uses: actions/create-github-app-token@v1
  with:
    app-id: ${{ secrets.APP_ID }}
    private-key: ${{ secrets.APP_PRIVATE_KEY }}
    owner: myorg
```

### Real-World Example: ado's Release Flow

#### release-please.yml

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
      # Generate GitHub App token
      - name: Generate token
        id: generate-token
        uses: actions/create-github-app-token@v1
        with:
          app-id: ${{ secrets.APP_ID }}
          private-key: ${{ secrets.APP_PRIVATE_KEY }}

      # Use token to run release-please
      - name: Run release-please
        id: release
        uses: googleapis/release-please-action@v4
        with:
          token: ${{ steps.generate-token.outputs.token }}
          config-file: release-please-config.json
          manifest-file: .release-please-manifest.json
```

**Why GitHub App here?**
- Release Please needs to create/update PRs (requires write permissions)
- PR updates should trigger CI workflows (GITHUB_TOKEN can't do this)
- Audit trail shows "release-please[bot]" not individual user

#### goreleaser.yml

```yaml
name: GoReleaser

on:
  release:
    types: [published]

permissions:
  contents: write
  packages: write
  id-token: write
  attestations: write

jobs:
  goreleaser:
    runs-on: ubuntu-latest
    steps:
      # Generate GitHub App token for releases
      - name: Generate token
        id: generate-token
        uses: actions/create-github-app-token@v1
        with:
          app-id: ${{ secrets.APP_ID }}
          private-key: ${{ secrets.APP_PRIVATE_KEY }}

      - uses: actions/checkout@v6
        with:
          fetch-depth: 0

      # Use GITHUB_TOKEN for package registry (app tokens can't push packages)
      - name: Login to GHCR
        uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      # Use GitHub App token for release creation
      - name: Run goreleaser
        uses: goreleaser/goreleaser-action@v6
        with:
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ steps.generate-token.outputs.token }}
```

**Why mixed tokens?**
- GitHub App token: For creating releases and uploading assets
- GITHUB_TOKEN: For pushing container images to GHCR (app tokens don't work for packages)

### Token Lifetime and Refresh

**Default behavior**:
- Installation tokens expire after 1 hour
- If your workflow runs longer, the token expires mid-workflow

**Solution for long workflows**:
Generate a new token when needed:

```yaml
jobs:
  long-running:
    runs-on: ubuntu-latest
    steps:
      - name: Generate initial token
        id: token1
        uses: actions/create-github-app-token@v1
        with:
          app-id: ${{ secrets.APP_ID }}
          private-key: ${{ secrets.APP_PRIVATE_KEY }}

      # ... use token1 for first hour of work ...

      - name: Generate fresh token
        id: token2
        uses: actions/create-github-app-token@v1
        with:
          app-id: ${{ secrets.APP_ID }}
          private-key: ${{ secrets.APP_PRIVATE_KEY }}

      # ... use token2 for next hour ...
```

### Debugging Authentication Issues

**Common problems and solutions**:

| Problem | Symptom | Solution |
|---------|---------|----------|
| Invalid App ID | "App not found" error | Verify `APP_ID` secret matches app settings |
| Invalid Private Key | "Authentication failed" | Re-download private key, ensure entire PEM file copied |
| Missing Installation | "Installation not found" | Install app on repository (see Step 8) |
| Insufficient Permissions | "Resource not accessible" | Update app permissions, re-accept installation |
| Token Expired | "Unauthorized" mid-workflow | Generate fresh token (see above) |

**Debugging workflow**:

```yaml
- name: Generate token
  id: generate-token
  uses: actions/create-github-app-token@v1
  with:
    app-id: ${{ secrets.APP_ID }}
    private-key: ${{ secrets.APP_PRIVATE_KEY }}

# Debug: Verify token was generated
- name: Verify token
  run: |
    if [ -z "${{ steps.generate-token.outputs.token }}" ]; then
      echo "ERROR: Token generation failed"
      exit 1
    fi
    echo "Token generated successfully (length: ${#TOKEN})"
  env:
    TOKEN: ${{ steps.generate-token.outputs.token }}

# Debug: Test token with simple API call
- name: Test token
  run: |
    RESPONSE=$(curl -s -H "Authorization: Bearer $TOKEN" \
      https://api.github.com/repos/${{ github.repository }})
    echo "API test: $(echo $RESPONSE | jq -r .name)"
  env:
    TOKEN: ${{ steps.generate-token.outputs.token }}
```

---

## GitHub Actions Integration

### Workflow Triggers

GitHub Actions workflows can be triggered by many events. Understanding **when to use GitHub App tokens** is key:

#### GITHUB_TOKEN Sufficient

```yaml
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v6
        # Default: uses GITHUB_TOKEN

      - name: Run tests
        run: go test ./...

      - name: Comment on PR
        if: github.event_name == 'pull_request'
        uses: actions/github-script@v7
        with:
          script: |
            github.rest.issues.createComment({
              owner: context.repo.owner,
              repo: context.repo.repo,
              issue_number: context.issue.number,
              body: '✅ Tests passed!'
            })
```

**When GITHUB_TOKEN is enough**:
- Running tests
- Building artifacts
- Commenting on PRs
- Uploading artifacts
- Updating check runs

#### GitHub App Required

```yaml
name: Release Please

on:
  push:
    branches: [main]

jobs:
  release-please:
    runs-on: ubuntu-latest
    steps:
      - name: Generate token
        id: app-token
        uses: actions/create-github-app-token@v1
        with:
          app-id: ${{ secrets.APP_ID }}
          private-key: ${{ secrets.APP_PRIVATE_KEY }}

      - name: Create release PR
        uses: googleapis/release-please-action@v4
        with:
          token: ${{ steps.app-token.outputs.token }}
```

**When GitHub App is required**:
- Creating PRs that trigger CI workflows
- Bypassing branch protection rules
- Triggering workflow_dispatch events
- Cross-repository operations
- Long-lived automation

### Workflow Permissions

**Configure permissions at workflow or job level**:

```yaml
name: CI

# Workflow-level permissions (default for all jobs)
permissions:
  contents: read
  pull-requests: write

jobs:
  test:
    # Inherits workflow-level permissions
    runs-on: ubuntu-latest
    steps:
      - run: echo "Has contents:read and pull-requests:write"

  deploy:
    # Override for specific job
    permissions:
      contents: write
      packages: write
    runs-on: ubuntu-latest
    steps:
      - run: echo "Has contents:write and packages:write"
```

**Available permissions**:
- `actions: read|write`
- `checks: read|write`
- `contents: read|write`
- `deployments: read|write`
- `id-token: write` (for OIDC)
- `issues: read|write`
- `packages: read|write`
- `pages: write`
- `pull-requests: read|write`
- `repository-projects: read|write`
- `security-events: read|write`
- `statuses: read|write`

**Best practice**: Use principle of least privilege:

```yaml
# BAD: Overly permissive
permissions: write-all

# GOOD: Specific permissions
permissions:
  contents: read
  pull-requests: write
```

### Conditional Workflows

**Run workflows based on conditions**:

```yaml
name: Dependabot Auto-Merge

on:
  pull_request:
    types: [opened, synchronize, reopened]

jobs:
  auto-merge:
    # Only run for Dependabot PRs
    if: github.actor == 'dependabot[bot]'
    runs-on: ubuntu-latest
    steps:
      - name: Fetch metadata
        id: metadata
        uses: dependabot/fetch-metadata@v2

      # Only auto-merge patch/minor updates
      - name: Auto-merge
        if: |
          steps.metadata.outputs.update-type == 'version-update:semver-patch' ||
          steps.metadata.outputs.update-type == 'version-update:semver-minor'
        run: gh pr merge --auto --squash "$PR_URL"
        env:
          PR_URL: ${{ github.event.pull_request.html_url }}
          GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

**Common conditions**:
- `github.event_name == 'push'` - Specific event type
- `github.ref == 'refs/heads/main'` - Specific branch
- `github.actor == 'dependabot[bot]'` - Specific actor
- `startsWith(github.ref, 'refs/tags/')` - Tag push
- `github.event.pull_request.draft == false` - Ready PRs only

### Reusable Workflows

**Create shared workflows** for common patterns:

**.github/workflows/reusable-test.yml**:
```yaml
name: Reusable Test Workflow

on:
  workflow_call:
    inputs:
      go-version:
        required: false
        type: string
        default: '1.23'
    secrets:
      codecov-token:
        required: false

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v6

      - name: Set up Go
        uses: actions/setup-go@v6
        with:
          go-version: ${{ inputs.go-version }}

      - name: Run tests
        run: go test -coverprofile=coverage.out ./...

      - name: Upload coverage
        if: secrets.codecov-token != ''
        uses: codecov/codecov-action@v5
        with:
          token: ${{ secrets.codecov-token }}
```

**Use in other workflows**:
```yaml
name: CI

on: [push, pull_request]

jobs:
  test:
    uses: ./.github/workflows/reusable-test.yml
    with:
      go-version: '1.23'
    secrets:
      codecov-token: ${{ secrets.CODECOV_TOKEN }}
```

### Matrix Strategies

**Test across multiple versions/platforms**:

```yaml
name: Cross-Platform Tests

on: [push, pull_request]

jobs:
  test:
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
        go-version: ['1.22', '1.23']
    runs-on: ${{ matrix.os }}
    steps:
      - uses: actions/checkout@v6

      - name: Set up Go
        uses: actions/setup-go@v6
        with:
          go-version: ${{ matrix.go-version }}

      - name: Run tests
        run: go test ./...
```

**Result**: 6 jobs run in parallel (3 OS × 2 Go versions)

### Caching Dependencies

**Speed up workflows with caching**:

```yaml
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v6

      # Go: Cache automatically with go-version
      - name: Set up Go
        uses: actions/setup-go@v6
        with:
          go-version: '1.23'
          cache: true

      # Python: Cache pip dependencies
      - name: Set up Python
        uses: actions/setup-python@v6
        with:
          python-version: '3.12'
          cache: pip
          cache-dependency-path: requirements.txt

      # Node: Cache npm dependencies
      - name: Set up Node
        uses: actions/setup-node@v4
        with:
          node-version: '20'
          cache: npm
```

**Manual caching**:
```yaml
- name: Cache Go modules
  uses: actions/cache@v4
  with:
    path: |
      ~/.cache/go-build
      ~/go/pkg/mod
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
    restore-keys: |
      ${{ runner.os }}-go-
```

### Artifact Management

**Upload artifacts for later use**:

```yaml
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v6
      - run: go build -o myapp

      - name: Upload binary
        uses: actions/upload-artifact@v4
        with:
          name: myapp-linux-amd64
          path: myapp
          retention-days: 7

  test-binary:
    needs: build
    runs-on: ubuntu-latest
    steps:
      - name: Download binary
        uses: actions/download-artifact@v4
        with:
          name: myapp-linux-amd64

      - name: Test binary
        run: ./myapp --version
```

### Environment Protection

**Add manual approval gates**:

```yaml
jobs:
  deploy:
    runs-on: ubuntu-latest
    environment:
      name: production
      url: https://myapp.com
    steps:
      - name: Deploy to production
        run: ./deploy.sh
```

**Configure in repository settings**:
- Settings → Environments → New environment
- Add protection rules:
  - Required reviewers
  - Wait timer
  - Environment secrets

---

## Dependabot Configuration

### What is Dependabot?

Dependabot automatically:
1. Monitors your dependencies for updates
2. Creates PRs for security vulnerabilities
3. Creates PRs for version updates
4. Groups related updates together

**Benefits**:
- Stay up-to-date with minimal effort
- Security patches applied quickly
- Consistent PR format for automation
- Works across multiple ecosystems

### Basic Configuration

Create `.github/dependabot.yml`:

```yaml
version: 2

updates:
  # Monitor Go dependencies
  - package-ecosystem: gomod
    directory: "/"
    schedule:
      interval: weekly
      day: monday
    labels:
      - dependencies
      - go
```

**Key fields**:
- `package-ecosystem`: Which dependency system (gomod, pip, npm, etc.)
- `directory`: Where the manifest file lives
- `schedule`: How often to check for updates
- `labels`: Applied to PRs for filtering/automation

### Multi-Ecosystem Configuration: ado Example

The `ado` project monitors three ecosystems:

```yaml
version: 2

updates:
  # Go dependencies
  - package-ecosystem: gomod
    directory: "/"
    schedule:
      interval: weekly
      day: monday
    labels:
      - dependencies
      - go
    commit-message:
      prefix: "build(deps)"
    groups:
      go-dependencies:
        patterns:
          - "*"

  # Python dependencies
  - package-ecosystem: pip
    directory: "/lab/py"
    schedule:
      interval: weekly
      day: monday
    labels:
      - dependencies
      - python
    commit-message:
      prefix: "build(deps)"
    groups:
      python-dependencies:
        patterns:
          - "*"

  # GitHub Actions
  - package-ecosystem: github-actions
    directory: "/"
    schedule:
      interval: monthly
    labels:
      - dependencies
      - ci
    commit-message:
      prefix: "build(deps)"
    groups:
      github-actions:
        patterns:
          - "*"
```

**Why different schedules?**
- Go/Python: Weekly - frequent updates, want to catch issues quickly
- GitHub Actions: Monthly - less frequent updates, more stability preferred

### Grouping Updates

**Problem**: Individual PRs for each dependency creates noise.

**Solution**: Group related updates into single PRs:

```yaml
updates:
  - package-ecosystem: gomod
    directory: "/"
    schedule:
      interval: weekly
    groups:
      # Group all Go dependencies together
      go-dependencies:
        patterns:
          - "*"

      # Or: Group by pattern
      cobra-ecosystem:
        patterns:
          - "github.com/spf13/cobra"
          - "github.com/spf13/viper"
```

**Result**: One PR titled "Bump go-dependencies group" instead of multiple PRs.

### Commit Message Customization

**Match your conventional commit format**:

```yaml
updates:
  - package-ecosystem: gomod
    directory: "/"
    commit-message:
      prefix: "build(deps)"
      prefix-development: "chore(deps-dev)"
      include: "scope"
```

**Resulting commit**:
```
build(deps): bump github.com/spf13/cobra from 1.8.0 to 1.8.1
```

### Ignoring Dependencies

**Prevent updates for specific packages**:

```yaml
updates:
  - package-ecosystem: gomod
    directory: "/"
    ignore:
      # Ignore all updates to this package
      - dependency-name: "github.com/legacy-package/*"

      # Ignore major version updates only
      - dependency-name: "github.com/some-package"
        update-types: ["version-update:semver-major"]
```

### Security vs Version Updates

Dependabot creates two types of PRs:

#### 1. Security Updates

**Triggered by**: GitHub Security Advisories

**Characteristics**:
- Created immediately (ignores schedule)
- Labeled with `security` automatically
- Includes vulnerability details
- Higher priority

**Example**:
```
🔒 Bump golang.org/x/crypto from 0.14.0 to 0.15.0

Fixes: CVE-2023-48795 (High severity)
```

#### 2. Version Updates

**Triggered by**: Regular schedule checks

**Characteristics**:
- Created per schedule (e.g., weekly)
- Labeled per configuration
- Follows grouping rules
- Can be safely delayed

**Example**:
```
Bump github.com/spf13/cobra from 1.8.0 to 1.8.1
```

### Supported Ecosystems

| Ecosystem | package-ecosystem value | Manifest file |
|-----------|------------------------|---------------|
| Go | `gomod` | `go.mod` |
| Python | `pip` | `requirements.txt`, `Pipfile` |
| JavaScript | `npm` | `package.json` |
| Docker | `docker` | `Dockerfile` |
| GitHub Actions | `github-actions` | `.github/workflows/*.yml` |
| Terraform | `terraform` | `*.tf` |
| Cargo (Rust) | `cargo` | `Cargo.toml` |
| Composer (PHP) | `composer` | `composer.json` |
| Maven | `maven` | `pom.xml` |
| Gradle | `gradle` | `build.gradle` |
| NuGet | `nuget` | `*.csproj` |

### Advanced Configuration

#### Version Constraints

```yaml
updates:
  - package-ecosystem: gomod
    directory: "/"
    # Only update to versions matching this pattern
    allow:
      - dependency-type: "production"
    ignore:
      - dependency-name: "*"
        update-types: ["version-update:semver-major"]
```

#### Pull Request Limits

```yaml
updates:
  - package-ecosystem: gomod
    directory: "/"
    # Limit concurrent open PRs
    open-pull-requests-limit: 5
```

**Default**: 5 PRs per ecosystem

**Why limit**: Avoid overwhelming maintainers

#### Custom PR Labels

```yaml
updates:
  - package-ecosystem: gomod
    directory: "/"
    labels:
      - dependencies
      - automated
      - go
      - "needs-review"  # Multi-word labels need quotes
```

#### Assignees and Reviewers

```yaml
updates:
  - package-ecosystem: gomod
    directory: "/"
    assignees:
      - "myusername"
    reviewers:
      - "team-lead"
      - "security-team"
```

#### Schedule Customization

```yaml
updates:
  - package-ecosystem: gomod
    directory: "/"
    schedule:
      interval: weekly
      day: monday
      time: "05:00"        # 5 AM UTC
      timezone: "America/New_York"
```

**Intervals**: daily, weekly, monthly

#### Target Branch

```yaml
updates:
  - package-ecosystem: gomod
    directory: "/"
    # Create PRs against develop instead of main
    target-branch: "develop"
```

### Real-World Example: Multi-Language Monorepo

```yaml
version: 2

updates:
  # Backend (Go)
  - package-ecosystem: gomod
    directory: "/backend"
    schedule:
      interval: weekly
    groups:
      backend-dependencies:
        patterns: ["*"]
    labels:
      - dependencies
      - backend
      - go

  # Frontend (npm)
  - package-ecosystem: npm
    directory: "/frontend"
    schedule:
      interval: weekly
    groups:
      frontend-dependencies:
        patterns: ["*"]
    labels:
      - dependencies
      - frontend
      - javascript

  # Infrastructure (Terraform)
  - package-ecosystem: terraform
    directory: "/terraform"
    schedule:
      interval: monthly
    labels:
      - dependencies
      - infrastructure
    # Don't auto-update infrastructure too aggressively
    open-pull-requests-limit: 2

  # GitHub Actions
  - package-ecosystem: github-actions
    directory: "/"
    schedule:
      interval: monthly
    groups:
      actions:
        patterns: ["*"]
    labels:
      - dependencies
      - ci
```

---

## Dependabot Auto-Merge

### Why Auto-Merge?

**Problem**: Dependabot creates many PRs, each requiring:
1. Wait for CI to pass
2. Manual review
3. Manual merge
4. Repeat for next update

**Solution**: Automatically merge PRs that meet safety criteria:
- ✅ Patch/minor updates (no breaking changes)
- ✅ All CI checks pass
- ✅ No merge conflicts
- ❌ Major updates (require manual review)

### Safety Principles

**What to auto-merge**:
- Patch updates (1.2.3 → 1.2.4) - bug fixes only
- Minor updates (1.2.0 → 1.3.0) - new features, backward compatible

**What NOT to auto-merge**:
- Major updates (1.0.0 → 2.0.0) - breaking changes possible
- Security updates - may need investigation
- Failed CI checks - obvious

### Basic Auto-Merge Workflow

**.github/workflows/dependabot-auto-merge.yml**:

```yaml
name: Dependabot Auto-Merge

on:
  pull_request:
    types: [opened, synchronize, reopened]

permissions:
  contents: write
  pull-requests: write

jobs:
  auto-merge:
    name: Auto-merge Dependabot PRs
    runs-on: ubuntu-latest
    # Only run for Dependabot PRs
    if: github.actor == 'dependabot[bot]'
    steps:
      # Fetch PR metadata (update type, dependency name, etc.)
      - name: Fetch Dependabot metadata
        id: metadata
        uses: dependabot/fetch-metadata@v2
        with:
          github-token: ${{ secrets.GITHUB_TOKEN }}

      # Auto-merge patch and minor updates only
      - name: Auto-merge patch and minor updates
        if: |
          steps.metadata.outputs.update-type == 'version-update:semver-patch' ||
          steps.metadata.outputs.update-type == 'version-update:semver-minor'
        run: gh pr merge --auto --squash "$PR_URL"
        env:
          PR_URL: ${{ github.event.pull_request.html_url }}
          GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}

      # Comment on major updates instead of merging
      - name: Comment on major updates
        if: steps.metadata.outputs.update-type == 'version-update:semver-major'
        run: |
          gh pr comment "$PR_URL" --body "⚠️ This is a **major version update**. Please review manually before merging."
        env:
          PR_URL: ${{ github.event.pull_request.html_url }}
          GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

**How it works**:
1. Trigger on any PR activity
2. Check if PR is from Dependabot
3. Fetch metadata about the update
4. If patch/minor: enable auto-merge
5. If major: add comment requesting manual review
6. PR merges automatically once CI passes

### Understanding Auto-Merge

**Key concept**: `gh pr merge --auto` doesn't merge immediately. It:
1. Enables auto-merge on the PR
2. Waits for all required checks to pass
3. Waits for required reviews (if any)
4. Merges automatically when conditions met

**Merge strategies**:
- `--squash`: All commits squashed (recommended for dependency updates)
- `--merge`: Standard merge commit
- `--rebase`: Rebase and merge

### Metadata Available

The `dependabot/fetch-metadata` action provides:

```yaml
steps.metadata.outputs:
  dependency-names: "github.com/spf13/cobra"
  dependency-type: "direct:production"
  update-type: "version-update:semver-minor"
  previous-version: "1.8.0"
  new-version: "1.8.1"
  package-ecosystem: "gomod"
  target-branch: "main"
```

**Use cases**:
- Filter by dependency name
- Treat production vs dev dependencies differently
- Special handling for specific ecosystems

### Advanced Auto-Merge Patterns

#### Pattern 1: Different Rules per Ecosystem

```yaml
jobs:
  auto-merge:
    runs-on: ubuntu-latest
    if: github.actor == 'dependabot[bot]'
    steps:
      - name: Fetch metadata
        id: metadata
        uses: dependabot/fetch-metadata@v2

      # Go: Auto-merge all non-major updates
      - name: Auto-merge Go updates
        if: |
          steps.metadata.outputs.package-ecosystem == 'gomod' &&
          steps.metadata.outputs.update-type != 'version-update:semver-major'
        run: gh pr merge --auto --squash "$PR_URL"
        env:
          PR_URL: ${{ github.event.pull_request.html_url }}
          GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}

      # GitHub Actions: Auto-merge all updates (low risk)
      - name: Auto-merge GitHub Actions
        if: steps.metadata.outputs.package-ecosystem == 'github-actions'
        run: gh pr merge --auto --squash "$PR_URL"
        env:
          PR_URL: ${{ github.event.pull_request.html_url }}
          GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}

      # Python: Only auto-merge patches (more cautious)
      - name: Auto-merge Python patches only
        if: |
          steps.metadata.outputs.package-ecosystem == 'pip' &&
          steps.metadata.outputs.update-type == 'version-update:semver-patch'
        run: gh pr merge --auto --squash "$PR_URL"
        env:
          PR_URL: ${{ github.event.pull_request.html_url }}
          GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

#### Pattern 2: Production vs Development Dependencies

```yaml
- name: Auto-merge dev dependencies
  if: |
    steps.metadata.outputs.dependency-type == 'direct:development' ||
    steps.metadata.outputs.dependency-type == 'indirect'
  run: gh pr merge --auto --squash "$PR_URL"
  env:
    PR_URL: ${{ github.event.pull_request.html_url }}
    GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}

- name: Comment on production dependencies
  if: steps.metadata.outputs.dependency-type == 'direct:production'
  run: |
    gh pr comment "$PR_URL" --body "⚠️ Production dependency update. Requires manual review."
  env:
    PR_URL: ${{ github.event.pull_request.html_url }}
    GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

#### Pattern 3: Approve Then Merge

Some repositories require PR approval even for bots:

```yaml
jobs:
  auto-merge:
    runs-on: ubuntu-latest
    if: github.actor == 'dependabot[bot]'
    steps:
      - name: Fetch metadata
        id: metadata
        uses: dependabot/fetch-metadata@v2

      # Step 1: Approve the PR
      - name: Approve PR
        if: |
          steps.metadata.outputs.update-type == 'version-update:semver-patch' ||
          steps.metadata.outputs.update-type == 'version-update:semver-minor'
        run: gh pr review --approve "$PR_URL"
        env:
          PR_URL: ${{ github.event.pull_request.html_url }}
          GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}

      # Step 2: Enable auto-merge
      - name: Enable auto-merge
        if: |
          steps.metadata.outputs.update-type == 'version-update:semver-patch' ||
          steps.metadata.outputs.update-type == 'version-update:semver-minor'
        run: gh pr merge --auto --squash "$PR_URL"
        env:
          PR_URL: ${{ github.event.pull_request.html_url }}
          GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### Security Considerations

**Question**: Is auto-merging dependency updates safe?

**Answer**: Yes, with proper safeguards:

✅ **Safeguards in place**:
1. CI must pass (tests, linting, coverage)
2. Only patch/minor updates (semantic versioning guarantee)
3. Branch protection rules enforced
4. Dependabot PRs are isolated (one dependency per PR)
5. Can be reverted immediately if issues found

❌ **Don't auto-merge**:
- Major version updates (breaking changes likely)
- Updates from untrusted sources
- Critical infrastructure dependencies without review
- If your test coverage is insufficient

**Risk mitigation**:
```yaml
# Require high test coverage
- name: Check coverage
  run: |
    COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | tr -d '%')
    if [ "${COVERAGE%.*}" -lt 80 ]; then
      echo "Coverage too low for auto-merge"
      exit 1
    fi
```

### Monitoring Auto-Merged PRs

**Add notifications**:

```yaml
- name: Notify on auto-merge
  if: success()
  run: |
    gh pr comment "$PR_URL" --body "✅ Auto-merged after CI passed."
  env:
    PR_URL: ${{ github.event.pull_request.html_url }}
    GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

**Track metrics**:
- How many PRs auto-merged vs manual review?
- Any reverted auto-merges?
- Time saved vs manual merging

### Disabling Auto-Merge Temporarily

**When you need manual control**:

```yaml
# Add to workflow to disable temporarily
- name: Check for auto-merge disable flag
  run: |
    if gh pr view "$PR_URL" --json labels --jq '.labels[].name' | grep -q "no-auto-merge"; then
      echo "Auto-merge disabled by label"
      exit 78  # Neutral exit code (doesn't fail workflow)
    fi
  env:
    PR_URL: ${{ github.event.pull_request.html_url }}
    GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

**Usage**: Add `no-auto-merge` label to Dependabot PR to skip auto-merge.

---

## Branch Protection Rules

### Why Branch Protection?

Branch protection ensures:
1. No direct pushes to main branch
2. All changes go through PRs
3. CI must pass before merge
4. Required reviewers approve changes
5. Commit history stays clean

**Without protection**:
```bash
git push origin main  # Danger: Bypasses CI, no review
```

**With protection**:
```bash
git push origin main  # Error: Protected branch
# Must: Create PR → CI passes → Review → Merge
```

### Basic Protection Setup

**UI Path**:
1. Repository → Settings → Branches
2. Add branch protection rule
3. Branch name pattern: `main` (or `main|develop` for multiple)
4. Configure protection settings

**Screenshot Description**:
> Branch protection rules page showing a form with "Branch name pattern" field at top (with "main" entered), followed by multiple checkboxes for protection settings organized in sections.

### Essential Protection Rules

#### Rule 1: Require Pull Request

```
☑️ Require a pull request before merging
  ☑️ Require approvals: 1
  ☐ Dismiss stale pull request approvals when new commits are pushed
  ☐ Require review from Code Owners
  ☐ Require approval of the most recent reviewable push
```

**Effect**:
- Direct pushes blocked
- At least 1 approval required
- Maintains review quality

**For solo projects**: Set to 0 approvals (just require PR, no review)

#### Rule 2: Require Status Checks

```
☑️ Require status checks to pass before merging
  ☑️ Require branches to be up to date before merging

  Search for status checks:
    ☑️ Conventional Commits
    ☑️ Go
    ☑️ Python Lab
    ☑️ Documentation
    ☑️ Docker
```

**Effect**:
- All CI checks must pass
- Branch must be up-to-date with main
- Prevents merging broken code

**Critical**: Check "Require branches to be up to date" to prevent integration issues.

#### Rule 3: Additional Settings

```
☑️ Require conversation resolution before merging
  # All PR comments must be resolved

☑️ Require linear history
  # No merge commits, only squash or rebase

☑️ Require signed commits (optional)
  # All commits must be GPG-signed

☐ Allow force pushes
  # Leave unchecked for main branch

☐ Allow deletions
  # Leave unchecked for main branch
```

### Required Status Checks

**How to identify check names**:
1. Run a workflow in a PR
2. Check the "Checks" tab on PR
3. Note the exact name of each check

**Example from ado's CI workflow**:

```yaml
# .github/workflows/ci.yml
jobs:
  commits:
    name: Conventional Commits  # ← This is the check name
    # ...

  go:
    name: Go  # ← This is the check name
    # ...

  python:
    name: Python Lab  # ← This is the check name
    # ...
```

**Add to branch protection**:
- ✅ Conventional Commits
- ✅ Go
- ✅ Python Lab
- ✅ Documentation
- ✅ Docker

**Screenshot Description**:
> Search box with "Search for status checks" placeholder. Below, a list of checkboxes with job names from CI workflows. Each checked item shows a green checkmark and the job name (e.g., "Go", "Python Lab").

### Protection Levels

#### Level 1: Basic (Solo Projects)

```
☑️ Require a pull request before merging
  Approvals: 0

☑️ Require status checks to pass before merging
  ☑️ Go
  ☑️ Python Lab
```

**Why**: Enforces CI checks, allows self-merge.

#### Level 2: Small Team

```
☑️ Require a pull request before merging
  Approvals: 1
  ☑️ Dismiss stale pull request approvals when new commits are pushed

☑️ Require status checks to pass before merging
  ☑️ Require branches to be up to date before merging
  ☑️ All CI checks

☑️ Require conversation resolution before merging
```

**Why**: Requires peer review, ensures fresh approval after changes.

#### Level 3: Production (Large Teams)

```
☑️ Require a pull request before merging
  Approvals: 2
  ☑️ Dismiss stale pull request approvals when new commits are pushed
  ☑️ Require review from Code Owners

☑️ Require status checks to pass before merging
  ☑️ Require branches to be up to date before merging
  ☑️ All CI checks (comprehensive list)

☑️ Require signed commits
☑️ Require linear history
☑️ Require conversation resolution before merging
☑️ Include administrators (no bypass)
```

**Why**: Maximum safety for critical systems.

### Bypass Permissions

**Allow specific users/teams to bypass rules**:

```
☐ Allow specified actors to bypass required pull requests
  Add: myusername, security-team
```

**Use cases**:
- Emergency hotfixes
- Release automation (GitHub App)
- Security patches

**Best practice**: Minimize bypass permissions. Use GitHub App for automation instead.

### GitHub App and Branch Protection

**Problem**: GitHub App needs to merge release PRs, but branch protection blocks merges.

**Solution 1**: Allow GitHub App to bypass (recommended)

```
☑️ Allow specified actors to bypass required pull requests
  Add: release-please[bot]
```

**Solution 2**: Auto-merge with status checks

Don't bypass, just auto-merge once checks pass:

```yaml
# release-please workflow doesn't merge
# Separate workflow merges after CI passes
```

**ado's approach**: Release Please creates PRs, CI passes, maintainer manually merges. No bypass needed.

### Rulesets (Beta)

GitHub is migrating to "Rulesets" - a more flexible system:

**UI Path**: Settings → Rules → Rulesets

**Benefits**:
- Apply rules across multiple branches with patterns
- More granular bypass permissions
- Better performance
- Import/export rulesets

**Example ruleset**:
```yaml
name: "Protect release branches"
target: branch
enforcement: active

conditions:
  ref_name:
    include:
      - "refs/heads/main"
      - "refs/heads/release/**"

rules:
  - type: pull_request
    parameters:
      required_approving_review_count: 1

  - type: required_status_checks
    parameters:
      strict_required_status_checks_policy: true
      required_status_checks:
        - context: "Go"
        - context: "Python Lab"
```

**Migration**: Existing branch protection rules auto-migrate to rulesets when enabled.

### Testing Branch Protection

**Verify protection is working**:

```bash
# Test 1: Direct push should fail
git checkout main
git commit --allow-empty -m "test"
git push origin main
# Expected: Error: Protected branch

# Test 2: PR without CI should block merge
git checkout -b test-branch
git commit --allow-empty -m "test"
git push origin test-branch
# Create PR, try to merge without waiting for CI
# Expected: Merge button disabled

# Test 3: PR with CI should allow merge
# Wait for CI to pass
# Expected: Merge button enabled
```

### Protection Rules Anti-Patterns

❌ **Don't**:
- Allow force pushes to main
- Allow deletions of main
- Disable status checks "temporarily" (you'll forget)
- Add too many bypass actors
- Require signed commits without team buy-in

✅ **Do**:
- Start with basic rules, add more over time
- Document why each rule exists
- Review rules quarterly
- Use GitHub App for automation instead of bypass
- Enable "Require branches to be up to date"

---

## GitHub Packages (GHCR)

### What is GHCR?

GitHub Container Registry (GHCR) is GitHub's Docker image hosting service:

**Benefits**:
- Integrated with GitHub (same permissions)
- Free for public repositories
- Fast (globally distributed CDN)
- OCI-compliant (Docker, OCI images)
- Supports container signing
- Works with GitHub Actions

**Image naming**:
```
ghcr.io/<owner>/<image>:<tag>
```

**Example**:
```
ghcr.io/anowarislam/ado:1.0.0
ghcr.io/anowarislam/ado:latest
```

### Enabling GHCR

**Repository Settings**:
1. Settings → Actions → General
2. Scroll to "Workflow permissions"
3. Select "Read and write permissions"
4. ✅ Check "Allow GitHub Actions to create and approve pull requests"

**Why**: Allows workflows to push images to GHCR.

### Publishing Images from GitHub Actions

#### Basic Example

```yaml
name: Build and Push Image

on:
  push:
    tags:
      - 'v*'

permissions:
  contents: read
  packages: write

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v6

      # Login to GHCR
      - name: Login to GHCR
        uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      # Build and push
      - name: Build and push
        uses: docker/build-push-action@v5
        with:
          context: .
          push: true
          tags: |
            ghcr.io/${{ github.repository }}:${{ github.ref_name }}
            ghcr.io/${{ github.repository }}:latest
```

**Key points**:
- Use `GITHUB_TOKEN` (not GitHub App) for packages
- Requires `packages: write` permission
- `github.actor` is the user/bot triggering the workflow

#### Multi-Architecture Images

**Support both AMD64 and ARM64**:

```yaml
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v6

      # Enable multi-platform builds
      - name: Set up QEMU
        uses: docker/setup-qemu-action@v3

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Login to GHCR
        uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      # Build for multiple platforms
      - name: Build and push multi-arch
        uses: docker/build-push-action@v5
        with:
          context: .
          platforms: linux/amd64,linux/arm64
          push: true
          tags: |
            ghcr.io/${{ github.repository }}:${{ github.ref_name }}
            ghcr.io/${{ github.repository }}:latest
```

### GoReleaser Integration

**ado's approach**: GoReleaser builds separate images per architecture, then creates a manifest:

**.goreleaser.yaml** (excerpt):
```yaml
dockers:
  # AMD64 image
  - id: ado-amd64
    goos: linux
    goarch: amd64
    dockerfile: goreleaser.Dockerfile
    use: buildx
    image_templates:
      - "ghcr.io/anowarislam/ado:{{ .Version }}-amd64"
      - "ghcr.io/anowarislam/ado:latest-amd64"
    build_flag_templates:
      - "--platform=linux/amd64"
      - "--label=org.opencontainers.image.source=https://github.com/anowarislam/ado"
      - "--label=org.opencontainers.image.version={{ .Version }}"

  # ARM64 image
  - id: ado-arm64
    goos: linux
    goarch: arm64
    dockerfile: goreleaser.Dockerfile
    use: buildx
    image_templates:
      - "ghcr.io/anowarislam/ado:{{ .Version }}-arm64"
      - "ghcr.io/anowarislam/ado:latest-arm64"
    build_flag_templates:
      - "--platform=linux/arm64"

# Multi-arch manifests
docker_manifests:
  - name_template: "ghcr.io/anowarislam/ado:{{ .Version }}"
    image_templates:
      - "ghcr.io/anowarislam/ado:{{ .Version }}-amd64"
      - "ghcr.io/anowarislam/ado:{{ .Version }}-arm64"

  - name_template: "ghcr.io/anowarislam/ado:latest"
    image_templates:
      - "ghcr.io/anowarislam/ado:latest-amd64"
      - "ghcr.io/anowarislam/ado:latest-arm64"
```

**Result**: Users can pull `ghcr.io/anowarislam/ado:latest` and automatically get the correct architecture.

### Image Metadata Labels

**OCI labels provide metadata**:

```dockerfile
LABEL org.opencontainers.image.title="ado"
LABEL org.opencontainers.image.description="CLI tool for X"
LABEL org.opencontainers.image.url="https://github.com/anowarislam/ado"
LABEL org.opencontainers.image.source="https://github.com/anowarislam/ado"
LABEL org.opencontainers.image.version="1.0.0"
LABEL org.opencontainers.image.revision="abc123"
LABEL org.opencontainers.image.licenses="MIT"
```

**Benefits**:
- Visible in GHCR UI
- Machine-readable metadata
- Links back to source code
- License compliance

### Image Visibility

**By default**: Images are private (even for public repos).

**Make images public**:
1. Go to package page: `https://github.com/users/<username>/packages/container/<package-name>`
2. Package settings → Change visibility
3. Select "Public"
4. Confirm

**Screenshot Description**:
> Package settings page showing "Danger Zone" section at bottom. Red button labeled "Change visibility" with dropdown showing "Public" and "Private" options. Warning text explains consequences of making package public.

### Pulling Images

**Public images** (no authentication):
```bash
docker pull ghcr.io/anowarislam/ado:latest
```

**Private images** (requires authentication):
```bash
# Create personal access token with read:packages scope
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin

docker pull ghcr.io/myorg/private-image:latest
```

### Image Lifecycle

**Retention policies**:
- Untagged images: Delete after 14 days (default)
- Can configure per package

**Tagging strategy**:
```
ghcr.io/owner/image:1.0.0       # Specific version (immutable)
ghcr.io/owner/image:1.0         # Minor version (updates)
ghcr.io/owner/image:1           # Major version (updates)
ghcr.io/owner/image:latest      # Latest release (updates)
ghcr.io/owner/image:main        # Latest commit on main (updates)
```

**ado's strategy**:
- `1.0.0`: Specific release (never changes)
- `latest`: Most recent release (updates on each release)

### Container Signing

**Sign images with cosign** (see Security Integrations section):

```yaml
- name: Install cosign
  uses: sigstore/cosign-installer@v3

- name: Sign container images
  run: |
    cosign sign --yes ghcr.io/${{ github.repository }}:${{ github.ref_name }}
```

**Users verify signatures**:
```bash
cosign verify ghcr.io/anowarislam/ado:1.0.0 \
  --certificate-identity-regexp=https://github.com/anowarislam/ado \
  --certificate-oidc-issuer=https://token.actions.githubusercontent.com
```

---

## GitHub Pages Deployment

### What is GitHub Pages?

GitHub Pages hosts static websites directly from repositories:

**Use cases**:
- Documentation sites
- Project landing pages
- API references
- Coverage reports
- User guides

**ado's use**: Documentation site built with MkDocs Material.

**URL**: `https://<username>.github.io/<repo>/`

**Example**: https://anowarislam.github.io/ado/

### Enabling GitHub Pages

**Repository Settings**:
1. Settings → Pages
2. Source: "GitHub Actions"
3. (Old method: Deploy from branch - don't use this)

**Screenshot Description**:
> GitHub Pages settings showing "Source" dropdown with "GitHub Actions" selected (recommended). Below, a notice: "Your site is live at https://username.github.io/repo/". Branch and folder dropdowns are greyed out when using GitHub Actions.

**Why GitHub Actions source?**:
- Full control over build process
- Can use any static site generator
- Better caching and performance
- Integrates with CI/CD

### Documentation Workflow

**ado's documentation workflow** (.github/workflows/docs.yml):

```yaml
name: Documentation

on:
  # Deploy on releases
  release:
    types: [published]

  # Deploy on docs changes to main
  push:
    branches: [main]
    paths:
      - 'docs/**'
      - 'mkdocs.yml'
      - 'CHANGELOG.md'

  # Build only (don't deploy) on PRs
  pull_request:
    paths:
      - 'docs/**'
      - 'mkdocs.yml'

# Required permissions for Pages
permissions:
  contents: read
  pages: write
  id-token: write

# Prevent concurrent deploys
concurrency:
  group: "pages"
  cancel-in-progress: false

jobs:
  # Build job
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0  # Full history for git dates

      - name: Set up Python
        uses: actions/setup-python@v5
        with:
          python-version: '3.12'

      # Cache pip dependencies
      - uses: actions/cache@v4
        with:
          path: ~/.cache/pip
          key: ${{ runner.os }}-pip-${{ hashFiles('mkdocs.yml') }}

      # Install documentation tools
      - name: Install MkDocs
        run: |
          pip install \
            mkdocs-material \
            mkdocs-minify-plugin \
            pillow \
            cairosvg

      # Copy CHANGELOG to docs
      - name: Copy CHANGELOG
        run: cp CHANGELOG.md docs/changelog.md

      # Build site
      - name: Build documentation
        run: mkdocs build --strict

      # Upload artifact (only if not PR)
      - name: Upload artifact
        if: github.event_name != 'pull_request'
        uses: actions/upload-pages-artifact@v3
        with:
          path: site

  # Deploy job (only on main/release)
  deploy:
    if: github.event_name != 'pull_request'
    needs: build
    runs-on: ubuntu-latest
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    steps:
      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v4
```

**Key features**:
1. **Build on PRs** - validates docs before merge
2. **Deploy on main** - keeps docs up-to-date
3. **Deploy on release** - major doc updates
4. **Caching** - speeds up builds
5. **Concurrency control** - prevents race conditions

### MkDocs Configuration

**mkdocs.yml** (excerpt):

```yaml
site_name: ado
site_url: https://anowarislam.github.io/ado/
repo_url: https://github.com/anowarislam/ado

theme:
  name: material
  palette:
    - media: "(prefers-color-scheme: light)"
      scheme: default
      primary: indigo
      toggle:
        icon: material/brightness-7
        name: Switch to dark mode
    - media: "(prefers-color-scheme: dark)"
      scheme: slate
      primary: indigo
      toggle:
        icon: material/brightness-4
        name: Switch to light mode
  features:
    - navigation.instant
    - navigation.tracking
    - navigation.tabs
    - search.suggest
    - content.code.copy

markdown_extensions:
  - pymdownx.highlight
  - pymdownx.superfences:
      custom_fences:
        - name: mermaid
          class: mermaid
          format: !!python/name:pymdownx.superfences.fence_code_format

plugins:
  - search
  - minify:
      minify_html: true

nav:
  - Home: index.md
  - Getting Started:
      - Installation: installation.md
  - Development:
      - Workflow: workflow.md
  - Changelog: changelog.md
```

**Why MkDocs Material?**:
- Beautiful out-of-the-box
- Fast search
- Mobile-responsive
- Dark mode
- Mermaid diagrams
- Code highlighting
- Easy navigation

### Custom Domain

**Use custom domain** (e.g., docs.myproject.com):

1. Add CNAME record in DNS:
   ```
   docs.myproject.com CNAME username.github.io.
   ```

2. Add custom domain in Settings → Pages

3. ✅ Enforce HTTPS (automatic with GitHub)

4. Add CNAME file to docs:
   ```
   echo "docs.myproject.com" > docs/CNAME
   ```

### Versioning Documentation

**Use mike for versioned docs**:

```bash
pip install mike

# Deploy version
mike deploy 1.0 latest --push --update-aliases

# Set default version
mike set-default latest

# List versions
mike list
```

**mkdocs.yml**:
```yaml
extra:
  version:
    provider: mike
```

**Result**: Docs available at:
- `/1.0/` - Version 1.0
- `/latest/` - Latest version
- Version switcher in UI

### Deployment Environments

**Use GitHub environments for approval gates**:

```yaml
deploy:
  needs: build
  runs-on: ubuntu-latest
  environment:
    name: github-pages
    url: ${{ steps.deployment.outputs.page_url }}
  steps:
    - uses: actions/deploy-pages@v4
```

**Benefits**:
- Deployment history
- Protection rules (require approval)
- Environment secrets
- Deployment logs

### Link Validation

**Validate internal links** before deployment:

```yaml
validate-links:
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4

    - name: Set up Python
      uses: actions/setup-python@v5
      with:
        python-version: '3.12'

    - name: Install dependencies
      run: |
        pip install mkdocs-material linkchecker

    - name: Build site
      run: mkdocs build

    - name: Check links
      run: |
        cd site && python -m http.server 8000 &
        sleep 2
        linkchecker http://localhost:8000 \
          --ignore-url='https?://.*' \
          --check-extern=false
        pkill -f "python -m http.server"
```

**Why**: Catch broken links before deploying.

---

## Codecov Integration

### What is Codecov?

Codecov is a coverage reporting and visualization service:

**Features**:
- Visual coverage reports
- Line-by-line coverage in PR comments
- Coverage trends over time
- Threshold enforcement
- Branch comparison
- Badge generation

**ado's usage**: Upload Go test coverage to Codecov.

### Setup Process

#### Step 1: Sign Up for Codecov

1. Go to https://codecov.io/
2. Sign in with GitHub
3. Authorize Codecov app
4. Select repositories to enable

**Free tier**: Unlimited public repos, 5 private repos.

#### Step 2: Get Repository Token

**For private repositories** (public repos don't need token):

1. Go to repository settings on Codecov
2. Copy the "Repository Upload Token"
3. Add to GitHub Secrets as `CODECOV_TOKEN`

**Screenshot Description**:
> Codecov repository settings showing a section titled "Repository Upload Token". A masked token value with a "Copy" button next to it. Below, explanatory text: "This token is used to upload coverage reports to Codecov."

#### Step 3: Upload Coverage in CI

**Add to CI workflow**:

```yaml
name: CI

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v6

      - name: Set up Go
        uses: actions/setup-go@v6
        with:
          go-version: '1.23'

      # Generate coverage report
      - name: Run tests with coverage
        run: go test -coverprofile=coverage.out -covermode=atomic ./...

      # Upload to Codecov
      - name: Upload coverage
        uses: codecov/codecov-action@v5
        with:
          token: ${{ secrets.CODECOV_TOKEN }}  # Required for private repos
          files: ./coverage.out
          flags: unittests
          name: codecov-ado
          fail_ci_if_error: false  # Don't fail CI if Codecov is down
```

**ado's implementation** (from ci.yml):

```yaml
- name: Test
  run: go test -coverprofile=coverage.out -covermode=atomic ./cmd/ado/... ./internal/...

- name: Upload coverage
  uses: codecov/codecov-action@v5
  with:
    files: ./coverage.out
    flags: unittests
    name: codecov-ado
  continue-on-error: true  # Don't fail CI if Codecov has issues
```

**Why `continue-on-error: true`?**
- Codecov is external service (can be down)
- Coverage is important but not CI-blocking
- Local coverage check already enforced (80% threshold)

### Coverage Reports

**Codecov provides**:

#### 1. PR Comments

Automatically comments on PRs with:
- Overall coverage change
- Files with changed coverage
- Uncovered lines

**Example**:
```
Codecov Report
Coverage: 82.45% (Δ +0.23%)

Files with coverage changes:
  internal/config/paths.go: 85.71% (Δ +5.00%)
  cmd/ado/meta/meta.go: 100.00% (no change)

Additional details:
  - 10 lines covered
  - 2 lines missed
```

#### 2. Web Dashboard

**View at**: `https://app.codecov.io/gh/<username>/<repo>`

**Features**:
- Coverage trends over time
- File-by-file coverage
- Sunburst visualization
- Commit history
- Branch comparison

#### 3. Coverage Badge

**Add to README.md**:
```markdown
[![codecov](https://codecov.io/gh/anowarislam/ado/branch/main/graph/badge.svg)](https://codecov.io/gh/anowarislam/ado)
```

**Result**: ![codecov](https://img.shields.io/badge/coverage-82%25-green)

### Configuration File

**Create codecov.yml** in repository root:

```yaml
# codecov.yml
coverage:
  status:
    project:
      default:
        target: 80%        # Project must have 80% coverage
        threshold: 2%      # Allow up to 2% decrease
    patch:
      default:
        target: 80%        # New code must have 80% coverage
        threshold: 2%

comment:
  layout: "header, diff, files"
  behavior: default
  require_changes: true  # Only comment if coverage changed

ignore:
  - "cmd/ado/root/root.go"  # CLI entry point (hard to test)
  - "**/*_test.go"           # Don't include test files in coverage
  - "lab/**"                 # Exclude experimental code
```

**Status checks**:
- `codecov/project` - Overall project coverage
- `codecov/patch` - Coverage of changed lines in PR

### Flags for Multi-Language Projects

**Separate coverage by language**:

```yaml
# Upload Go coverage
- name: Upload Go coverage
  uses: codecov/codecov-action@v5
  with:
    files: ./coverage-go.out
    flags: go,unittests

# Upload Python coverage
- name: Upload Python coverage
  uses: codecov/codecov-action@v5
  with:
    files: ./coverage-python.xml
    flags: python,unittests
```

**codecov.yml**:
```yaml
flags:
  go:
    paths:
      - cmd/
      - internal/
  python:
    paths:
      - lab/py/
```

**Result**: Separate coverage metrics for each language.

### Coverage Enforcement

**Two-layer enforcement** (recommended approach):

#### Layer 1: Local CI Check (ado's approach)

```yaml
- name: Check coverage threshold (80%)
  run: |
    COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | tr -d '%')
    echo "Total coverage: ${COVERAGE}%"
    COVERAGE_INT=$(echo "$COVERAGE" | awk '{printf "%d", $1}')
    if [ "$COVERAGE_INT" -lt 80 ]; then
      echo "::error::Coverage ${COVERAGE}% is below 80% threshold"
      exit 1
    fi
```

**Advantages**:
- Immediate feedback
- Doesn't depend on external service
- Blocks merge if threshold not met
- Fast (no network calls)

#### Layer 2: Codecov Status Check

```yaml
# codecov.yml
coverage:
  status:
    project:
      default:
        target: 80%
        informational: false  # Block PR if not met
```

**Advantages**:
- Historical trends
- Visual reports
- Line-by-line view
- PR comments

**ado's strategy**: Enforce locally (Layer 1), report to Codecov (Layer 2 informational).

### Troubleshooting

**Common issues**:

| Issue | Cause | Solution |
|-------|-------|----------|
| No coverage reported | Wrong file path | Verify `files` parameter |
| "Unable to locate" error | Coverage file not generated | Check test command |
| Token error (private repo) | Missing CODECOV_TOKEN | Add secret |
| Report not showing | Wrong branch | Check branch filter in Codecov settings |
| Old reports showing | Caching issue | Clear Codecov cache |

**Debug upload**:
```yaml
- name: Upload coverage (debug)
  uses: codecov/codecov-action@v5
  with:
    files: ./coverage.out
    verbose: true  # Enable debug logging
```

---

## Security Integrations

### Overview

GitHub provides multiple security features to protect your repository:

```mermaid
graph TB
    subgraph "Security Features"
        A[Dependabot Alerts]
        B[Code Scanning]
        C[Secret Scanning]
        D[Dependency Review]
        E[Security Advisories]
    end

    subgraph "Supply Chain Security"
        F[SLSA Attestations]
        G[Container Signing]
        H[SBOM Generation]
    end

    A --> I[Automated PRs]
    B --> J[PR Comments]
    C --> K[Alert + Block]
    F --> L[Verification]
    G --> L

    style A fill:#ffe1e1
    style F fill:#e1f5e1
```

### 1. Dependabot Security Alerts

**What it is**: Automated vulnerability detection for dependencies.

**Enable**:
1. Settings → Security → Code security and analysis
2. ✅ Enable Dependabot alerts
3. ✅ Enable Dependabot security updates

**Screenshot Description**:
> Security settings page showing "Code security and analysis" section. Multiple toggles for "Dependabot alerts" (On), "Dependabot security updates" (On), "Grouped security updates" (On). Each with "Enable" or "Disable" buttons and explanatory text.

**How it works**:
1. GitHub scans your dependencies (go.mod, requirements.txt, etc.)
2. Compares against GitHub Advisory Database
3. Creates alert when vulnerability found
4. Optionally creates PR with fix

**Example alert**:
```
🔒 High severity vulnerability in golang.org/x/crypto

CVE-2023-48795: Potential for SSH terrapin attack

Recommendation: Update to golang.org/x/crypto@0.15.0 or later

Dependabot will create a PR to fix this vulnerability.
```

**Alert severity levels**:
- **Critical**: Immediate action required
- **High**: Prompt attention needed
- **Medium**: Review and plan fix
- **Low**: Fix when convenient

### 2. Code Scanning (CodeQL)

**What it is**: Automated code analysis for security vulnerabilities.

**Setup**:

**.github/workflows/codeql.yml**:
```yaml
name: CodeQL

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]
  schedule:
    - cron: '0 0 * * 1'  # Weekly on Monday

jobs:
  analyze:
    name: Analyze
    runs-on: ubuntu-latest
    permissions:
      security-events: write
      contents: read

    strategy:
      matrix:
        language: ['go', 'python']

    steps:
      - uses: actions/checkout@v6

      # Initialize CodeQL
      - name: Initialize CodeQL
        uses: github/codeql-action/init@v3
        with:
          languages: ${{ matrix.language }}

      # Autobuild (or run custom build)
      - name: Autobuild
        uses: github/codeql-action/autobuild@v3

      # Analyze code
      - name: Perform CodeQL Analysis
        uses: github/codeql-action/analyze@v3
```

**What it detects**:
- SQL injection
- XSS vulnerabilities
- Path traversal
- Command injection
- Authentication issues
- Cryptographic weaknesses

**Results appear**:
- Security tab → Code scanning alerts
- PR comments (if issues in changed code)

**Screenshot Description**:
> Security tab showing "Code scanning" section with list of alerts. Each alert shows severity (Critical/High/Medium/Low), rule name, affected file, and line number. Filters on left for severity, status, and branch.

### 3. Secret Scanning

**What it is**: Detects committed secrets (API keys, tokens, passwords).

**Enable**:
1. Settings → Security → Code security and analysis
2. ✅ Enable secret scanning
3. ✅ Enable push protection (blocks commits with secrets)

**Detected secrets**:
- GitHub tokens
- AWS access keys
- Azure credentials
- Google Cloud keys
- Slack tokens
- Stripe keys
- Private keys (SSH, GPG)
- Database connection strings

**Push protection**:
```bash
git push origin main
# Error: Push blocked - secret detected

   GitHub Token detected in commit abc123
   File: config/token.txt
   Pattern: ghp_xxxxx

   To proceed:
   1. Remove the secret
   2. Commit the fix
   3. Force push (if already pushed)
```

**False positives**: Mark as "Close as false positive" in Security tab.

### 4. Dependency Review

**What it is**: Review dependency changes in PRs before merging.

**.github/workflows/dependency-review.yml**:
```yaml
name: Dependency Review

on:
  pull_request:
    branches: [main]

permissions:
  contents: read
  pull-requests: write

jobs:
  review:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v6

      - name: Dependency Review
        uses: actions/dependency-review-action@v4
        with:
          fail-on-severity: high
          comment-summary-in-pr: always
```

**What it checks**:
- New dependencies added
- Dependencies with known vulnerabilities
- License compatibility
- Deprecated packages

**PR comment example**:
```
Dependency Review

Changes detected:
  ✅ Added: github.com/spf13/cobra@1.8.1 (MIT license)
  ⚠️  Added: some-package@2.0.0 (High severity vulnerability)
  ❌ Added: bad-package@1.0.0 (GPL license - incompatible)

Recommendation: Review high-severity vulnerability before merging.
```

### 5. SLSA Provenance (Build Attestations)

**What it is**: Cryptographic proof of how artifacts were built.

**Enable in workflow** (ado's approach):

```yaml
name: GoReleaser

permissions:
  attestations: write  # Required for attestations
  id-token: write      # Required for signing

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      # ... build artifacts ...

      # Attest binaries
      - name: Attest build provenance
        uses: actions/attest-build-provenance@v2
        with:
          subject-path: |
            dist/*.tar.gz
            dist/*.zip
            dist/checksums.txt
```

**What it proves**:
- Built in GitHub Actions (not local machine)
- Exact commit SHA
- Workflow file used
- Build environment details
- Timestamp

**Users verify**:
```bash
gh attestation verify ado_1.0.0_linux_amd64.tar.gz --owner anowarislam

✓ Verification succeeded!

Attestation:
  - Source: github.com/anowarislam/ado
  - Commit: abc123...
  - Workflow: .github/workflows/goreleaser.yml
  - Built on: 2024-11-26T10:30:00Z
```

**See Chapter 3** for detailed SLSA implementation guide.

### 6. Container Signing (Sigstore/Cosign)

**What it is**: Cryptographically sign container images.

**Enable in workflow** (ado's approach):

```yaml
jobs:
  release:
    runs-on: ubuntu-latest
    permissions:
      packages: write
      id-token: write  # For keyless signing
    steps:
      # ... build and push images ...

      # Install cosign
      - name: Install cosign
        uses: sigstore/cosign-installer@v3

      # Sign images (keyless with OIDC)
      - name: Sign container images
        run: |
          cosign sign --yes ghcr.io/${{ github.repository }}:${{ github.ref_name }}
          cosign sign --yes ghcr.io/${{ github.repository }}:latest
```

**Users verify**:
```bash
cosign verify ghcr.io/anowarislam/ado:1.0.0 \
  --certificate-identity-regexp=https://github.com/anowarislam/ado \
  --certificate-oidc-issuer=https://token.actions.githubusercontent.com

Verification successful!
```

**Why keyless signing?**:
- No private keys to manage
- Uses GitHub's OIDC identity
- Signatures stored in public transparency log
- Can't be forged

**See Chapter 3** for detailed signing implementation.

### 7. Security Policy

**Create SECURITY.md**:

```markdown
# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 1.x     | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

Please report security vulnerabilities to security@example.com.

Do NOT open public issues for security vulnerabilities.

Expected response time: 48 hours

## Security Advisories

View published advisories at:
https://github.com/username/repo/security/advisories

## Verification

All releases are signed and attested:

### Verify binary attestation
\`\`\`bash
gh attestation verify ado_1.0.0_linux_amd64.tar.gz --owner username
\`\`\`

### Verify container signature
\`\`\`bash
cosign verify ghcr.io/username/repo:1.0.0 \
  --certificate-identity-regexp=https://github.com/username/repo \
  --certificate-oidc-issuer=https://token.actions.githubusercontent.com
\`\`\`
```

**Benefits**:
- Shows users you take security seriously
- Provides responsible disclosure process
- Documents verification procedures
- Appears in Security tab

### 8. Security Advisories

**Create private advisory**:
1. Security tab → Advisories
2. "New draft security advisory"
3. Fill in details (vulnerability, severity, affected versions)
4. Optionally create private fork to fix
5. Request CVE (optional)
6. Publish when ready

**Screenshot Description**:
> Form for creating security advisory with fields: Title, Severity (dropdown: Low/Medium/High/Critical), CWE (Common Weakness Enumeration dropdown), Description (markdown editor), Affected versions, Patched versions, Workarounds. Buttons: "Create advisory" and "Cancel".

**Why use advisories?**:
- Coordinate disclosure with fix
- Get CVE assigned automatically
- Notify users via Dependabot
- Work privately before public disclosure

---

## Permission Matrices

### GitHub Actions Permissions

**GITHUB_TOKEN permissions** by operation:

| Operation | Permission Required | Read | Write |
|-----------|-------------------|------|-------|
| Clone repository | `contents` | ✅ | |
| Create release | `contents` | | ✅ |
| Upload release assets | `contents` | | ✅ |
| Comment on PR | `pull-requests` | | ✅ |
| Create PR | `pull-requests` | | ✅ |
| Merge PR | `pull-requests` | | ✅ |
| Push to GHCR | `packages` | | ✅ |
| Upload artifacts | `actions` | | ✅ |
| Update check runs | `checks` | | ✅ |
| Deploy Pages | `pages` | | ✅ |
| Create attestation | `attestations` | | ✅ |
| Sign with OIDC | `id-token` | | ✅ |

**Workflow examples**:

```yaml
# Minimal CI (tests only)
permissions:
  contents: read

# CI with PR comments
permissions:
  contents: read
  pull-requests: write

# Release workflow
permissions:
  contents: write
  packages: write
  id-token: write
  attestations: write

# Pages deployment
permissions:
  contents: read
  pages: write
  id-token: write
```

### GitHub App Permissions

**Required permissions** by use case:

| Use Case | Repository Permissions | Organization Permissions |
|----------|----------------------|-------------------------|
| Release automation | Contents: write<br>Pull requests: write | None |
| Issue management | Issues: write<br>Pull requests: write | None |
| Status checks | Checks: write<br>Statuses: write | None |
| Workflow dispatch | Actions: write | None |
| Org-wide automation | (per use case) | Members: read (min) |
| Package publishing | Packages: write | None |

**ado's GitHub App permissions**:
```yaml
Repository permissions:
  Contents: Read and write          # Create releases, tags
  Pull requests: Read and write     # Create/update release PRs
  Metadata: Read-only               # Required by GitHub

Organization permissions:
  None                              # Only operates on single repo
```

### Integration Permissions Checklist

**Setting up a new integration**:

| Integration | GitHub Permissions | External Credentials | Branch Protection |
|-------------|-------------------|---------------------|-------------------|
| Basic CI | `contents: read` | None | ✅ Must pass |
| Release Please | GitHub App: write | App ID + Private Key | ❌ Can bypass |
| GoReleaser | GitHub App: write | App ID + Private Key | N/A (post-merge) |
| Dependabot | Built-in | None | ✅ Must pass |
| Dependabot Auto-Merge | `contents: write`<br>`pull-requests: write` | None | ✅ Must pass |
| GHCR | `packages: write` | None | N/A |
| GitHub Pages | `pages: write`<br>`id-token: write` | None | N/A |
| Codecov | None (upload only) | CODECOV_TOKEN (private repos) | N/A |
| CodeQL | `security-events: write` | None | ✅ Optional |

---

## Integration Architecture

### Complete Integration Flow

```mermaid
sequenceDiagram
    participant Dev as Developer
    participant GH as GitHub
    participant CI as CI Workflow
    participant RP as Release Please
    participant GR as GoReleaser
    participant GHCR as GitHub Packages
    participant Pages as GitHub Pages
    participant Codecov as Codecov

    Dev->>GH: Push feature branch
    GH->>CI: Trigger CI workflow
    CI->>CI: Run tests
    CI->>Codecov: Upload coverage
    CI->>Dev: PR checks pass ✅

    Dev->>GH: Merge PR to main
    GH->>RP: Trigger release-please
    RP->>GH: Create/update release PR

    Dev->>GH: Merge release PR
    GH->>GH: Create GitHub release
    GH->>GR: Trigger goreleaser
    GR->>GHCR: Push container images
    GR->>GH: Upload binaries + attestations
    GR->>GH: Sign containers
    GH->>Pages: Trigger docs deployment
    Pages->>Pages: Build and deploy docs
```

### Integration Dependency Graph

```mermaid
graph TB
    subgraph "Core Infrastructure"
        A[GitHub Repository]
        B[GitHub Actions]
    end

    subgraph "Quality Gates"
        C[CI Workflow]
        D[Branch Protection]
        E[CodeQL]
    end

    subgraph "Dependency Management"
        F[Dependabot]
        G[Dependency Review]
        H[Auto-Merge]
    end

    subgraph "Release Automation"
        I[Release Please]
        J[GoReleaser]
        K[GitHub App]
    end

    subgraph "Distribution"
        L[GitHub Releases]
        M[GHCR]
        N[GitHub Pages]
    end

    subgraph "Security"
        O[Attestations]
        P[Container Signing]
        Q[Secret Scanning]
    end

    subgraph "Reporting"
        R[Codecov]
        S[Security Advisories]
    end

    A --> B
    B --> C
    B --> E
    C --> D
    A --> F
    F --> G
    F --> H
    K --> I
    K --> J
    I --> L
    J --> L
    J --> M
    J --> N
    J --> O
    J --> P
    C --> R
    A --> Q
    A --> S

    style A fill:#e1f5e1
    style K fill:#ffe1e1
    style O fill:#e1e5ff
```

### Network Diagram

**Integration communication paths**:

```
┌─────────────────┐
│  GitHub Repo    │
└────────┬────────┘
         │
    ┌────┴────┐
    │         │
┌───▼──┐  ┌──▼────┐
│  CI  │  │ Pages │
└──┬───┘  └───────┘
   │
   ├──────► Codecov (external)
   │
   └──────► GHCR (internal)

┌──────────────────┐
│  GitHub App      │
└────────┬─────────┘
         │
    ┌────┴─────┐
    │          │
┌───▼────┐  ┌──▼────┐
│Release │  │GoRel  │
│Please  │  │easer  │
└────────┘  └───┬───┘
                │
           ┌────┴────┐
           │         │
        ┌──▼───┐  ┌──▼────┐
        │Attest│  │ Sign  │
        └──────┘  └───────┘
```

**External services**:
- Codecov: Outbound HTTPS (coverage upload)
- Sigstore: Outbound HTTPS (signing, transparency log)

**Internal services**:
- GHCR: GitHub-hosted, same authentication
- GitHub Pages: GitHub-hosted, same authentication

### Data Flow

**Coverage data flow**:
```
Test execution (CI)
  → Generate coverage.out
    → Upload to Codecov API
      → Codecov processes
        → Comment on PR
        → Update dashboard
```

**Release artifact flow**:
```
Code push → main
  → Release Please detects conventional commits
    → Create release PR
      → Human merges PR
        → GitHub creates release
          → GoReleaser builds artifacts
            → Upload binaries to GitHub Releases
            → Push images to GHCR
            → Generate attestations
            → Sign containers
              → Verify in SECURITY.md
```

### Integration Costs

**GitHub-hosted (free for public repos)**:
- CI minutes: 2,000 min/month (private repos)
- GHCR storage: 500 MB (private repos)
- GitHub Pages: Unlimited bandwidth (public repos)
- Actions: Unlimited (public repos)

**External services**:
- Codecov: Free for public repos, paid for private
- Sigstore: Free (transparency log)

**Cost optimization**:
- Cache dependencies (faster builds, fewer minutes)
- Use concurrency limits (prevent runaway builds)
- Schedule heavy workflows (e.g., CodeQL weekly, not every push)

---

## Troubleshooting Guide

### Common Issues and Solutions

#### Issue 1: GitHub App Token Generation Fails

**Symptom**:
```
Error: Failed to generate token: App not found
```

**Causes & solutions**:

| Cause | Solution |
|-------|----------|
| Wrong App ID | Verify `APP_ID` secret matches app settings |
| Wrong Private Key | Re-download private key, ensure entire PEM copied |
| App not installed | Install app on repository (Settings → GitHub Apps) |
| Private key expired | Generate new private key, update secret |

**Debug**:
```yaml
- name: Debug secrets
  run: |
    echo "App ID length: ${#APP_ID}"
    echo "Private key starts with: $(echo "$PRIVATE_KEY" | head -c 50)"
  env:
    APP_ID: ${{ secrets.APP_ID }}
    PRIVATE_KEY: ${{ secrets.APP_PRIVATE_KEY }}
```

#### Issue 2: GITHUB_TOKEN Permissions Denied

**Symptom**:
```
Error: Resource not accessible by integration
```

**Cause**: Missing permissions in workflow.

**Solution**: Add required permissions:
```yaml
permissions:
  contents: write      # If creating releases
  pull-requests: write # If commenting on PRs
```

**Verify current permissions**:
```yaml
- name: Check permissions
  run: |
    gh api repos/${{ github.repository }} --jq '.permissions'
  env:
    GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

#### Issue 3: Dependabot PRs Not Auto-Merging

**Symptom**: PR stays open despite `--auto` flag.

**Possible causes**:

| Cause | Check | Solution |
|-------|-------|----------|
| CI not passing | PR checks tab | Fix failing tests |
| Branch protection requires review | Branch settings | Add approval or change settings |
| Merge conflicts | PR status | Rebase or recreate PR |
| Auto-merge not enabled | PR sidebar | Verify "Enable auto-merge" checked |

**Verify auto-merge status**:
```yaml
- name: Check auto-merge status
  run: |
    gh pr view "$PR_URL" --json autoMergeRequest
  env:
    PR_URL: ${{ github.event.pull_request.html_url }}
    GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

#### Issue 4: Codecov Upload Fails

**Symptom**:
```
Error: Unable to locate coverage file
```

**Solutions**:

```yaml
# Fix 1: Verify file exists
- name: Check coverage file
  run: ls -la coverage.out

# Fix 2: Use correct path
- name: Upload coverage
  uses: codecov/codecov-action@v5
  with:
    files: ./coverage.out  # Not coverage.out (missing ./)

# Fix 3: Generate coverage correctly
- name: Run tests
  run: go test -coverprofile=coverage.out ./...  # Not just go test
```

#### Issue 5: GitHub Pages Not Deploying

**Symptom**: Site not updating after push.

**Checklist**:

```
☑️ Source set to "GitHub Actions" (not "Deploy from branch")
☑️ Workflow has correct permissions (pages: write, id-token: write)
☑️ Upload artifact step included
☑️ Deploy step included
☑️ No errors in workflow logs
```

**Common fix**:
```yaml
# Make sure you have BOTH steps

# Step 1: Upload artifact
- uses: actions/upload-pages-artifact@v3
  with:
    path: site

# Step 2: Deploy artifact
- uses: actions/deploy-pages@v4
```

#### Issue 6: Branch Protection Blocking Everything

**Symptom**: Can't merge anything, even with passing checks.

**Debug steps**:

1. Check required status checks match exactly:
   ```
   Branch protection: "Go"
   Workflow job name: "go"  ❌ Case mismatch!
   ```

2. Verify all checks appear in PR:
   - PR → Checks tab → Note exact names
   - Branch protection → Edit → Search for those exact names

3. Check "Require branches to be up to date":
   - If checked, branch must be rebased before merge
   - If failing, try rebasing: `git rebase origin/main`

#### Issue 7: Container Image Push Fails

**Symptom**:
```
Error: denied: permission_denied: write_package
```

**Solutions**:

```yaml
# Fix 1: Verify permissions
permissions:
  packages: write  # Must be present

# Fix 2: Use GITHUB_TOKEN (not GitHub App token)
- name: Login to GHCR
  uses: docker/login-action@v3
  with:
    registry: ghcr.io
    username: ${{ github.actor }}
    password: ${{ secrets.GITHUB_TOKEN }}  # NOT app token

# Fix 3: Verify package name
# Must be: ghcr.io/<owner>/<repo>:<tag>
# Not: ghcr.io/<repo>:<tag>
```

#### Issue 8: Attestation Verification Fails

**Symptom**:
```
gh attestation verify: no attestations found
```

**Causes**:

1. **Attestations not generated**: Add attestation step to workflow
2. **Wrong artifact name**: Match exactly with what was built
3. **Private repo**: Use `gh auth login` first
4. **Artifact from fork**: Attestations only for source repo

**Verify attestations exist**:
```bash
gh api \
  /repos/OWNER/REPO/attestations/SHA256 \
  --jq '.attestations[] | {bundle_sha256: .bundle.dsseEnvelope.payload}'
```

#### Issue 9: Release Please Not Creating PRs

**Symptom**: No release PR after merge to main.

**Checklist**:

```
☑️ Commits follow conventional format (feat:, fix:, etc.)
☑️ Commits on main branch (not feature branch)
☑️ release-please-config.json exists
☑️ .release-please-manifest.json exists
☑️ Workflow triggered (check Actions tab)
☑️ GitHub App token generated successfully
```

**Common issue**: Commit merged via squash with non-conventional title.

**Fix**: Ensure PR title is conventional (it becomes the commit message).

#### Issue 10: Secrets Not Available in Workflow

**Symptom**:
```
Error: Secret CODECOV_TOKEN not found
```

**Causes & solutions**:

| Cause | Solution |
|-------|----------|
| Secret not added | Add in Settings → Secrets → Actions |
| Wrong secret name | Match exactly (case-sensitive) |
| Fork PR | Secrets not available to fork PRs (security) |
| Environment secret | Specify `environment` in job |

**For fork PRs**: Use `pull_request_target` (carefully!) or skip secret steps:
```yaml
- name: Upload coverage
  if: github.event.pull_request.head.repo.full_name == github.repository
  uses: codecov/codecov-action@v5
```

### Debug Strategies

#### Strategy 1: Enable Debug Logging

**Enable in repository**:
1. Settings → Secrets → Actions
2. Add secret: `ACTIONS_STEP_DEBUG` = `true`
3. Add secret: `ACTIONS_RUNNER_DEBUG` = `true`

**Result**: Verbose logs for all workflow steps.

#### Strategy 2: Add Debug Steps

```yaml
- name: Debug context
  run: |
    echo "Event: ${{ github.event_name }}"
    echo "Ref: ${{ github.ref }}"
    echo "SHA: ${{ github.sha }}"
    echo "Actor: ${{ github.actor }}"
    echo "Repo: ${{ github.repository }}"

- name: Debug environment
  run: |
    echo "Working directory: $(pwd)"
    echo "Files:"
    ls -la
    echo "Environment variables:"
    env | sort
```

#### Strategy 3: Test Locally

**Use act to run workflows locally**:
```bash
# Install act
brew install act  # macOS
# or
curl https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash

# Run workflow
act -j test  # Run 'test' job

# With secrets
act -j test -s GITHUB_TOKEN=ghp_xxx
```

**Limitations**: Not all features work (e.g., OIDC, attestations).

---

## Summary

This chapter covered GitHub's integration ecosystem:

### Key Takeaways

1. **Authentication**:
   - Use `GITHUB_TOKEN` for most CI operations
   - Use GitHub App for triggering workflows and release automation
   - Avoid PATs in production

2. **GitHub Apps**:
   - Create with minimal permissions
   - Store App ID and Private Key as secrets
   - Use `actions/create-github-app-token` in workflows

3. **Dependabot**:
   - Configure for all ecosystems (Go, Python, GitHub Actions)
   - Group updates to reduce PR noise
   - Auto-merge patch/minor updates safely

4. **Branch Protection**:
   - Require PR + status checks
   - Start minimal, add rules gradually
   - Use GitHub App bypass for automation

5. **GitHub Packages (GHCR)**:
   - Use for container images
   - Multi-arch images for compatibility
   - Sign images with cosign

6. **GitHub Pages**:
   - Deploy via GitHub Actions (not branch)
   - Use MkDocs Material for documentation
   - Validate links before deployment

7. **Security**:
   - Enable Dependabot alerts
   - Use CodeQL for code scanning
   - Generate SLSA attestations
   - Sign containers with Sigstore

### Integration Checklist

**For a new project, enable**:

```
Essential:
☑️ Branch protection (PR required, status checks)
☑️ Dependabot (gomod, pip, github-actions)
☑️ Basic CI (tests, linting)
☑️ GitHub Actions permissions (minimal)

Important:
☑️ Dependabot auto-merge
☑️ Codecov integration
☑️ GitHub Pages (documentation)
☑️ Secret scanning + push protection

Advanced:
☑️ GitHub App (release automation)
☑️ GHCR (container images)
☑️ CodeQL (code scanning)
☑️ SLSA attestations
☑️ Container signing
```

### Next Steps

- **Chapter 8**: Troubleshooting common CI/CD issues
- **Chapter 9**: Step-by-step implementation guide
- **Chapter 10**: Python-specific adaptations

---

**You now have a complete understanding of GitHub's integration ecosystem. Ready to implement a production-grade CI/CD system with proper security and automation.**
