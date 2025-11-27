# Development Workflow

This document describes the three-phase workflow for proposing and implementing changes in `ado`.

## Overview

All significant changes follow this workflow:

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│    Issue     │ ──▶ │  ADR (PR 1)  │ ──▶ │  Spec (PR 2) │ ──▶ │  Code (PR 3) │
│  (Proposal)  │     │  (if needed) │     │              │     │              │
└──────────────┘     └──────────────┘     └──────────────┘     └──────────────┘
```

Each phase produces a separate PR. Phase transitions require approval.

## Quick Reference

| Change Type | ADR? | Spec Location | Branch Pattern |
|-------------|------|---------------|----------------|
| New command | Usually no | `docs/commands/` | `spec/cmd-name` |
| New internal feature | Maybe | `docs/features/` | `spec/feature-name` |
| Architectural change | **Yes** | After ADR | `adr/NNNN-title` |
| Bug fix | No | N/A | `fix/description` |
| Documentation | No | N/A | `docs/description` |

## Phase 0: Issue Creation

All work begins with a GitHub Issue using the appropriate template:

| Template | Use When |
|----------|----------|
| [**ADR Proposal**](https://github.com/anowarislam/ado/issues/new?template=adr_proposal.md) | Architectural decisions, new patterns, breaking changes |
| [**Feature Proposal**](https://github.com/anowarislam/ado/issues/new?template=feature_proposal.md) | Non-command features (config, internal packages) |
| [**Command Proposal**](https://github.com/anowarislam/ado/issues/new?template=command_proposal.md) | New CLI commands |
| [**Bug Report**](https://github.com/anowarislam/ado/issues/new?template=bug_report.md) | Fixing broken behavior |

### Decision Tree: Do I Need an ADR?

```
Is this a new architectural pattern?           → YES → ADR required
Does this change existing architecture?        → YES → ADR required
Does this add new dependencies/tools?          → YES → ADR required
Does this affect multiple components?          → YES → ADR required
Is this security-related?                      → YES → ADR required
Is this a single command following patterns?   → NO  → Skip to Spec
Is this a bug fix or documentation?            → NO  → Skip to Implementation
```

## Phase 1: Decision (ADR)

**When required:** See decision tree above.

**Location:** `docs/adr/NNNN-title.md`

**Template:** [docs/adr/TEMPLATE.md](adr/TEMPLATE.md)

### Process

1. Create branch: `git checkout -b adr/NNNN-short-title`
2. Copy template: `docs/adr/TEMPLATE.md` → `docs/adr/NNNN-title.md`
3. Fill in all sections:
   - **Context**: Why is this decision needed?
   - **Decision**: What are we deciding?
   - **Consequences**: What are the trade-offs?
   - **Alternatives**: What else did we consider?
4. Set Status: `Proposed`
5. Create PR with title: `docs(adr): NNNN - title`
6. Request review from maintainers
7. Address feedback
8. Once approved: Update Status to `Accepted`, merge

### PR Checklist for ADR

- [ ] ADR follows template structure
- [ ] Context clearly explains the problem
- [ ] Decision is clearly stated
- [ ] Alternatives are documented with rationale
- [ ] Consequences (positive/negative) are realistic
- [ ] Links to originating Issue
- [ ] Added to `docs/adr/README.md` index

### After ADR Approval

1. Update the ADR index (`docs/adr/README.md`)
2. Update `mkdocs.yml` navigation if needed
3. Proceed to Phase 2 (Specification)

## Phase 2: Specification

**Location:** Depends on change type.

| Type | Location | Template |
|------|----------|----------|
| Command | `docs/commands/NN-name.md` | [docs/commands/TEMPLATE.md](commands/TEMPLATE.md) |
| Feature | `docs/features/NN-name.md` | [docs/features/TEMPLATE.md](features/TEMPLATE.md) |

### Process

1. Create branch: `git checkout -b spec/feature-name`
2. Copy appropriate template
3. Fill in all sections:
   - **Purpose/Overview**: What does this do?
   - **Examples**: Concrete usage examples
   - **Behavior**: Step-by-step logic
   - **Error Cases**: What can go wrong?
   - **Testing**: How to verify correctness?
4. Link to ADR (if applicable)
5. Create PR with title: `docs(spec): [command|feature] name`
6. Request review
7. Address feedback
8. Once approved: Update Status to `Approved`, merge

### PR Checklist for Spec

- [ ] Spec follows template structure
- [ ] Examples are concrete and testable
- [ ] Error cases documented with expected behavior
- [ ] Implementation locations specified
- [ ] Testing strategy defined
- [ ] Links to ADR (if applicable)
- [ ] Links to originating Issue
- [ ] Added to appropriate index (commands or features README)
- [ ] Added to `mkdocs.yml` navigation

### Tests-First Guidance

Before implementation, consider writing test cases based on the spec:

```go
func TestCommandName(t *testing.T) {
    tests := []struct {
        name    string
        args    []string
        want    string
        wantErr bool
    }{
        // From spec: Example 1 - Basic usage
        {"basic", []string{"arg"}, "expected output", false},
        // From spec: Error case - No arguments
        {"no args", []string{}, "", true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

This creates a clear contract between spec and implementation.

## Phase 3: Implementation

**Location:** Code in `cmd/ado/` or `internal/`

### Process

1. Create branch: `git checkout -b feat/feature-name` (or `fix/`)
2. Implement according to spec:
   - Follow the spec exactly
   - All examples must work as documented
   - All error cases must behave as specified
3. Write/complete tests
4. Update documentation:
   - CLAUDE.md if architecture changed
   - README.md if user-facing
5. Run validation: `make validate`
6. Create PR with conventional commit title
7. Request review
8. Address feedback
9. Once approved: Merge

### PR Checklist for Implementation

- [ ] Code follows [Go style guide](style/go-style.md)
- [ ] All spec requirements implemented
- [ ] All examples from spec work correctly
- [ ] All error cases behave as specified
- [ ] Tests pass (`make test`)
- [ ] Linting passes (`make lint`)
- [ ] CLAUDE.md updated if architecture changed
- [ ] README.md updated if user-facing
- [ ] Links to Spec and Issue in PR description

### Implementation Patterns

**For commands:**
```
1. Create cmd/ado/<command>/<command>.go
2. Export NewCommand() *cobra.Command
3. Wire into cmd/ado/root/root.go via AddCommand()
4. Create cmd/ado/<command>/<command>_test.go
5. Add shared logic to internal/ if needed
```

**For features:**
```
1. Create internal/<feature>/ package
2. Define interfaces at use-sites
3. Implement core logic
4. Integrate with existing code
5. Write comprehensive tests
```

## Branch Naming

| Phase | Pattern | Example |
|-------|---------|---------|
| ADR | `adr/NNNN-short-title` | `adr/0002-plugin-architecture` |
| Spec | `spec/feature-name` | `spec/config-validation` |
| Implementation | `feat/feature-name` | `feat/config-validation` |
| Bug fix | `fix/issue-description` | `fix/config-path-windows` |
| Documentation | `docs/description` | `docs/api-reference` |

## Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

| Phase | Pattern | Example |
|-------|---------|---------|
| ADR | `docs(adr): NNNN - title` | `docs(adr): 0002 - plugin architecture` |
| Spec | `docs(spec): type name` | `docs(spec): command config` |
| Implementation | `feat(scope): description` | `feat(config): add validation` |
| Bug fix | `fix(scope): description` | `fix(config): handle missing file` |

## PR Labels

| Label | Meaning |
|-------|---------|
| `adr` | Architecture Decision Record |
| `spec` | Specification document |
| `enhancement` | Feature implementation |
| `command` | New CLI command |
| `bug` | Bug fix |
| `documentation` | Documentation only |
| `needs-discussion` | Requires team input |
| `needs-adr` | ADR required before proceeding |
| `needs-spec` | Spec required before proceeding |

## Code Ownership & Reviewer Assignment

This project uses GitHub's CODEOWNERS file (`.github/CODEOWNERS`) to define who reviews and approves changes to different parts of the codebase.

### How It Works

When you open a PR, GitHub automatically:
1. Analyzes which files you changed
2. Looks up code owners from `.github/CODEOWNERS`
3. Requests reviews from matching code owners
4. Marks them as **required reviewers** (PR cannot merge without their approval)

### Current Code Owners

| Component | Owner | What They Review |
|-----------|-------|------------------|
| **Default (all files)** | @anowarislam | All changes unless more specific owner exists |
| **Commands** (`/cmd/ado/`) | @anowarislam | CLI command implementations |
| **Internal Packages** (`/internal/`) | @anowarislam | Core libraries and utilities |
| **Documentation** (`/docs/`) | @anowarislam | All documentation |
| **ADRs** (`/docs/adr/`) | @anowarislam | Architecture Decision Records |
| **CI/CD** (`/.github/workflows/`) | @anowarislam | GitHub Actions workflows |
| **Build System** (`/Makefile`, `/make/`) | @anowarislam | Build automation |
| **Python Lab** (`/lab/py/`) | @anowarislam | Python prototyping |

**Complete mapping**: See [`.github/CODEOWNERS`](../.github/CODEOWNERS)

### Review Requirements by Phase

Each phase has specific review requirements:

**Phase 1 - ADR Review**:
- **Reviewer**: Tech lead (@anowarislam)
- **Approval Criteria**:
  - Context clearly explains problem
  - Decision is well-reasoned
  - Alternatives documented with rationale
  - Consequences realistic (no unverified claims)
- **Timeline**: Within 48 business hours

**Phase 2 - Spec Review**:
- **Reviewer**: Package owner (from CODEOWNERS)
- **Approval Criteria**:
  - Examples are concrete and testable
  - Error cases well-defined
  - Implementation locations specified
  - Testing strategy clear
- **Timeline**: Within 48 business hours

**Phase 3 - Implementation Review**:
- **Reviewer**: Automatic based on modified files (CODEOWNERS)
- **Approval Criteria**:
  - Follows spec exactly
  - All CI checks pass (tests, coverage, lint, build)
  - 80%+ test coverage maintained
  - No security vulnerabilities
  - Code follows style guide
- **Timeline**: Within 48 business hours

### Branch Protection Rules

The `main` branch requires:
- ✅ Code owner approval
- ✅ All CI checks pass
- ✅ All conversations resolved
- ✅ Branch up to date with main
- ✅ Commits signed
- ✅ Conventional commit format

## CODEOWNERS & Automated Reviewer Assignment

This project uses GitHub's **CODEOWNERS** file to automate code review assignments and enforce approval requirements. This ensures changes are reviewed by maintainers with expertise in the affected areas.

### What is CODEOWNERS?

The CODEOWNERS file (`.github/CODEOWNERS`) defines which GitHub users or teams are responsible for reviewing changes to specific parts of the codebase. It uses gitignore-style patterns to map file paths to owners.

**Location**: `.github/CODEOWNERS`

**Format**:
```
# Pattern         Owner(s)
*                 @anowarislam           # Default owner
/cmd/ado/         @anowarislam           # CLI commands
/internal/config/ @anowarislam           # Config package
/docs/adr/        @anowarislam           # Architecture decisions
```

**Key Principles**:
- **Pattern-based**: Uses glob patterns like `.gitignore`
- **Hierarchical**: More specific patterns override general ones
- **Last match wins**: If multiple patterns match, the last one takes precedence
- **Automatic**: GitHub requests reviews automatically when PRs are opened

### How CODEOWNERS Integrates with Branch Protection

CODEOWNERS works together with branch protection rules to enforce code review:

**Without CODEOWNERS**:
- PR can be approved by any collaborator
- No automatic reviewer assignment
- Manual effort to find right reviewers

**With CODEOWNERS + Branch Protection**:
- GitHub automatically assigns code owners as reviewers
- Branch protection requires approval from code owners
- PR cannot merge until code owner approves
- Dismisses stale approvals when new commits are pushed

**Branch Protection Setting**: "Require review from Code Owners" is **enabled** on the `main` branch.

### Automatic Reviewer Assignment

When you open a PR, GitHub automatically:

1. **Analyzes Changed Files**: Scans all files modified in the PR
2. **Matches Patterns**: Finds matching patterns in `.github/CODEOWNERS`
3. **Requests Reviews**: Adds code owners as **required reviewers**
4. **Blocks Merge**: PR cannot merge until code owners approve

**Example**:

You modify these files:
```
internal/config/loader.go
internal/config/loader_test.go
docs/features/01-config-validation.md
```

GitHub automatically requests review from:
- `@anowarislam` (owns `/internal/config/`)
- `@anowarislam` (owns `/docs/features/`)

Result: One code owner covers all changes, so **1 approval required**.

### Current CODEOWNERS Structure

The `ado` project uses a comprehensive ownership model covering all major components:

| Component | Path Pattern | Owner | Purpose |
|-----------|--------------|-------|---------|
| **Default (Catch-All)** | `*` | @anowarislam | All files unless more specific pattern matches |
| **CLI Commands** | `/cmd/ado/` | @anowarislam | Command implementations and tests |
| **Internal Packages** | `/internal/` | @anowarislam | Core libraries (config, logging, meta, ui) |
| **Config System** | `/internal/config/` | @anowarislam | Configuration loading and validation |
| **Logging System** | `/internal/logging/` | @anowarislam | Structured logging implementation |
| **Python Lab** | `/lab/py/` | @anowarislam | Python prototyping environment |
| **Documentation** | `/docs/` | @anowarislam | All documentation and guides |
| **ADRs** | `/docs/adr/` | @anowarislam | Architecture Decision Records |
| **Command Specs** | `/docs/commands/` | @anowarislam | CLI command specifications |
| **Feature Specs** | `/docs/features/` | @anowarislam | Non-command feature specifications |
| **CI/CD Workflows** | `/.github/workflows/` | @anowarislam | GitHub Actions automation |
| **Build System** | `/Makefile`, `/make/` | @anowarislam | Build and task automation |
| **Git Hooks** | `/.githooks/` | @anowarislam | Pre-commit and pre-push hooks |
| **Release Config** | `/.goreleaser.yaml`, `/release-please-config.json` | @anowarislam | Release automation |
| **Dependencies** | `go.mod`, `go.sum`, `lab/py/pyproject.toml` | @anowarislam | Dependency management |
| **Security** | `/SECURITY.md`, `/docs/recipes/04-security-features.md` | @anowarislam | Security policies and features |

**Complete mapping**: See [`.github/CODEOWNERS`](../.github/CODEOWNERS)

### Configuring Branch Protection with CODEOWNERS

To enforce code owner approvals, the following branch protection settings are enabled on `main`:

**Required Settings**:
1. **Require a pull request before merging** ✅
   - Ensures all changes go through PR workflow
2. **Require approvals** ✅
   - Minimum: 1 approval
3. **Require review from Code Owners** ✅
   - **This is the critical setting that enforces CODEOWNERS**
4. **Dismiss stale pull request approvals when new commits are pushed** ✅
   - Ensures code owner re-reviews changes
5. **Require approval of the most recent reviewable push** ✅
   - Prevents merging with outdated approvals

**GitHub Settings Path**:
```
Settings → Branches → main → Edit Branch Protection Rule
  ↓
[✓] Require a pull request before merging
  [✓] Require approvals: 1
  [✓] Dismiss stale pull request approvals when new commits are pushed
  [✓] Require review from Code Owners  ← KEY SETTING
  [✓] Require approval of the most recent reviewable push
```

### Benefits of Using CODEOWNERS

**For Contributors**:
- ✅ No guessing who should review your PR
- ✅ Automatic reviewer assignment saves time
- ✅ Clear ownership reduces back-and-forth
- ✅ Faster review turnaround (right person notified immediately)

**For Maintainers**:
- ✅ Enforces review by subject-matter experts
- ✅ Distributes review load across team (when multiple owners)
- ✅ Prevents accidental merges without proper review
- ✅ Tracks ownership as codebase grows
- ✅ Reduces review bottlenecks

**For the Project**:
- ✅ Maintains architectural consistency
- ✅ Ensures security-critical changes get extra scrutiny
- ✅ Documents responsibility boundaries
- ✅ Facilitates knowledge transfer through reviews
- ✅ Scales as team grows (add more code owners easily)

### Example: CODEOWNERS File for `ado` Project

Here's a simplified example showing the key patterns used in this project:

```bash
# .github/CODEOWNERS
#
# Code owners are automatically requested for review when someone
# opens a pull request that modifies code that they own.
#
# Syntax: [file-pattern] @username @org/team-name
# Order matters: Last matching pattern takes precedence.

# =============================================================================
# Default Owner - catches all files unless overridden
# =============================================================================
* @anowarislam

# =============================================================================
# Build & Release Configuration - requires build system expertise
# =============================================================================
/.github/workflows/     @anowarislam
/.goreleaser.yaml       @anowarislam
/Makefile               @anowarislam
/make/                  @anowarislam

# =============================================================================
# Go Source Code - requires Go expertise
# =============================================================================
/cmd/ado/               @anowarislam
/internal/              @anowarislam

# More specific package ownership
/internal/config/       @anowarislam
/internal/logging/      @anowarislam

# =============================================================================
# Documentation - requires technical writing skills
# =============================================================================
/docs/                  @anowarislam
/docs/adr/              @anowarislam  # Architecture decisions need tech lead
/docs/commands/         @anowarislam
/docs/features/         @anowarislam

# =============================================================================
# Security-Critical Files - requires security review
# =============================================================================
SECURITY.md                           @anowarislam
/docs/recipes/04-security-features.md @anowarislam
```

**Pattern Precedence Example**:

If you modify `/internal/config/loader.go`:
1. Matches `*` → `@anowarislam`
2. Matches `/internal/` → `@anowarislam`
3. Matches `/internal/config/` → `@anowarislam` (last match wins)

Result: `@anowarislam` is the code owner.

### Integration Architecture

CODEOWNERS integrates with multiple GitHub features:

```
┌─────────────────────────────────────────────────────────────────┐
│                         GitHub PR Workflow                       │
└─────────────────────────────────────────────────────────────────┘
                                  │
                    1. PR Opened (files modified)
                                  │
                                  ▼
        ┌──────────────────────────────────────────────┐
        │          .github/CODEOWNERS File             │
        │  Pattern matching: files → owners            │
        └──────────────────┬───────────────────────────┘
                           │
         2. Automatic reviewer assignment
                           │
                           ▼
        ┌──────────────────────────────────────────────┐
        │        Required Reviewers Added              │
        │  (Code owners marked as required)            │
        └──────────────────┬───────────────────────────┘
                           │
            3. Review process begins
                           │
          ┌────────────────┴────────────────┐
          │                                 │
          ▼                                 ▼
┌────────────────────┐          ┌────────────────────┐
│   CI/CD Checks     │          │  Code Owner Review │
│  - Tests (80%)     │          │  - Architecture    │
│  - Linting         │          │  - Code quality    │
│  - Build           │          │  - Security        │
│  - Coverage        │          │  - Design          │
└────────┬───────────┘          └────────┬───────────┘
         │                               │
         └───────────┬───────────────────┘
                     │
      4. Both CI and code owner approval required
                     │
                     ▼
        ┌──────────────────────────────────────────────┐
        │       Branch Protection Rules Check          │
        │  ✓ CI passed                                 │
        │  ✓ Code owner approved                       │
        │  ✓ Conversations resolved                    │
        │  ✓ Branch up to date                         │
        │  ✓ Commits signed                            │
        └──────────────────┬───────────────────────────┘
                           │
              5. All requirements met
                           │
                           ▼
        ┌──────────────────────────────────────────────┐
        │           Merge Button Enabled               │
        │      (Squash & Merge to main)                │
        └──────────────────────────────────────────────┘
```

**Key Integration Points**:

1. **CODEOWNERS → Reviewer Assignment**: Automatic when PR is opened
2. **Branch Protection → CODEOWNERS**: "Require review from Code Owners" setting
3. **CI/CD → Code Owner Review**: Both must pass (AND logic, not OR)
4. **Stale Approval Dismissal**: New commits trigger re-review by code owners
5. **Merge Blocking**: Cannot merge without code owner approval, even if CI passes

### Re-requesting Review

If you push new changes after review:
1. Previous approvals are automatically dismissed
2. Code owner must re-review
3. Use the circular arrow icon next to reviewer name to re-request

### Getting Help

If you're unsure who should review your changes:
1. Check `.github/CODEOWNERS` for the files you modified
2. GitHub will auto-assign when you open the PR
3. Ask in PR comments if unclear

### Complete Documentation

See [Code Ownership Guide](code-ownership.md) for:
- Detailed review responsibilities
- How to become a code owner
- Emergency procedures
- FAQ

## PR Quality Standards

Every PR is automatically validated against quality standards via the PR Metrics Dashboard (see [ADR-0005](adr/0005-pr-metrics-dashboard.md)).

### Automated Quality Checks

When you open or update a PR, GitHub Actions automatically runs comprehensive checks:

| Check | Threshold | Enforcement |
|-------|-----------|-------------|
| **Total Coverage** | 80% minimum | Blocks merge if below threshold |
| **Package Coverage** | 80% per package (configurable) | Blocks merge if any package below threshold |
| **File Coverage** | 70% per file | Blocks merge if any file below threshold |
| **Diff Coverage** | 85% of changed lines | Blocks merge if new/modified code insufficiently tested |
| **Test Results** | All tests passing | Blocks merge on any test failure |
| **Build** | Binary compiles | Blocks merge on build failure |

### What You'll See in Your PR

**1. Status Checks** (PR header):
```
✅ Go / Test Coverage (82.5% - threshold: 80%)
✅ Go / Package Coverage Enforcement
✅ Go / Test Results (142 tests passing)
❌ Go / Diff Coverage (75% - need 85%)
```

**2. PR Comment** (automatically posted/updated):
```markdown
## Test Coverage Report

**Total Coverage:** 82.5% (+1.2% vs main) ✅
**Diff Coverage:** 87.5% (35/40 changed lines) ✅
**Estimated CI Cost:** $0.027 (3m 24s @ $0.008/min)

### Changed Files
| File | Coverage | Status |
|------|----------|--------|
| internal/config/loader.go | 85.0% (34/40 lines) | ✅ |
| internal/logging/handler.go | 72.0% (18/25 lines) | ⚠️ Below threshold |
```

**3. Test Failure Annotations** (Files changed tab):
- Failed tests appear as inline annotations
- Click annotation to see full stack trace
- Line numbers link directly to failure location

### Coverage Thresholds by Package

Coverage thresholds are configurable in `.testcoverage.yml`:

```yaml
override:
  github.com/anowarislam/ado/internal/meta: 90       # Critical packages
  github.com/anowarislam/ado/internal/config: 85     # Important logic
  github.com/anowarislam/ado/cmd/ado/version: 60     # Simple commands
```

**Philosophy**: Higher thresholds for critical code, lower for trivial wiring.

### Handling Coverage Failures

If your PR is blocked by coverage:

1. **Check the PR comment** to see which files/packages need improvement
2. **Add tests** for uncovered lines (preferred)
3. **Adjust thresholds** in `.testcoverage.yml` if justified (requires explanation in PR)
4. **Document why** certain lines are hard to test (e.g., error handling for rare conditions)

**Do not bypass** coverage checks without team discussion.

### Performance Benchmarks

Optional benchmark workflow tracks performance regressions:

- Automatically runs on PRs that change `.go` files
- Compares PR performance vs main branch
- Posts warning if >5% slower
- Does not block merge (informational only)

### Cost Awareness

Each PR shows estimated CI cost to help teams:

- Identify expensive workflows
- Optimize long-running tests
- Track CI budget over time

**Note**: Costs are estimates based on workflow duration, not actual billing.

### Quality Philosophy

**Our approach**:
- **Prevent regressions**: Coverage can only increase, never decrease
- **Granular enforcement**: Package-level thresholds prevent localized quality drops
- **Diff coverage**: New code must be well-tested (85% minimum)
- **Transparency**: All metrics visible in PR, no hidden dashboards

**Why this matters**:
- Catches bugs before production
- Maintains codebase health long-term
- Provides objective review criteria
- Encourages thorough testing

See [Feature Spec](features/03-pr-metrics-dashboard.md) for complete implementation details.

## Examples

### Example 1: Adding a New Command (No ADR Needed)

A new command that follows existing patterns:

1. **Issue**: Create Command Proposal issue
2. **Spec PR**: `docs(spec): command export`
   - Branch: `spec/export`
   - File: `docs/commands/04-export.md`
3. **Implementation PR**: `feat(cmd): add export command`
   - Branch: `feat/export`
   - Files: `cmd/ado/export/export.go`, `cmd/ado/export/export_test.go`

### Example 2: Adding Plugin System (ADR Required)

A new architectural pattern affecting multiple components:

1. **Issue**: Create ADR Proposal issue
2. **ADR PR**: `docs(adr): 0002 - plugin architecture`
   - Branch: `adr/0002-plugin-architecture`
   - File: `docs/adr/0002-plugin-architecture.md`
3. **Spec PR**: `docs(spec): feature plugin-system`
   - Branch: `spec/plugin-system`
   - File: `docs/features/01-plugin-system.md`
4. **Implementation PR**: `feat(plugins): implement plugin loader`
   - Branch: `feat/plugin-system`
   - Files: `internal/plugins/...`

### Example 3: Bug Fix (Skip to Implementation)

Fixing broken behavior:

1. **Issue**: Create Bug Report issue
2. **Implementation PR**: `fix(config): handle missing file gracefully`
   - Branch: `fix/config-missing-file`
   - Include test that reproduces the bug

## Workflow Diagram

```
                    ┌─────────────────────────────────────────┐
                    │            GitHub Issue                 │
                    │  (ADR/Feature/Command/Bug Proposal)     │
                    └─────────────────┬───────────────────────┘
                                      │
                    ┌─────────────────▼───────────────────┐
                    │         Triage & Discussion         │
                    │   - Assign labels                   │
                    │   - Determine if ADR needed         │
                    └─────────────────┬───────────────────┘
                                      │
              ┌───────────────────────┼───────────────────────┐
              │                       │                       │
              ▼                       ▼                       ▼
    ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
    │   ADR Required  │    │  Spec Required  │    │  Direct to Code │
    │                 │    │  (No ADR)       │    │  (Bug fix)      │
    └────────┬────────┘    └────────┬────────┘    └────────┬────────┘
             │                      │                      │
             ▼                      │                      │
    ┌─────────────────┐             │                      │
    │   PR #1: ADR    │             │                      │
    │   docs/adr/     │             │                      │
    └────────┬────────┘             │                      │
             │                      │                      │
             ▼                      ▼                      │
    ┌─────────────────────────────────────────┐            │
    │              PR #2: Spec                │            │
    │   docs/commands/ or docs/features/      │            │
    └─────────────────┬───────────────────────┘            │
                      │                                    │
                      ▼                                    ▼
    ┌─────────────────────────────────────────────────────────────┐
    │                    PR #3: Implementation                    │
    │   cmd/ado/<command>/ or internal/<feature>/                 │
    └─────────────────────────────────────────────────────────────┘
```

## Tips for Success

### Writing Good Specs

1. **Be concrete**: Use actual examples, not abstract descriptions
2. **Define errors first**: Know what can go wrong before coding
3. **Think in tests**: Each spec section should map to test cases
4. **Keep it focused**: One spec = one feature/command

### Working with LLMs

This workflow is designed for LLM-assisted development:

1. **Specs as prompts**: A well-written spec gives Claude/Copilot clear instructions
2. **Tests as validation**: Pre-written tests verify LLM output
3. **ADRs as context**: Architectural decisions provide background for code generation
4. **Async development**: Submit spec, let LLM implement, review results

### Common Pitfalls

| Pitfall | Solution |
|---------|----------|
| Skipping ADR for "simple" arch changes | Use the decision tree honestly |
| Vague specs | Add more concrete examples |
| Implementing before spec approval | Wait for approval to avoid rework |
| Not updating indexes | Always update README.md and mkdocs.yml |

## Related Documentation

- [Code Ownership](code-ownership.md) - Code review and approval process
- [ADRs](adr/) - Architecture Decision Records
- [Command Specs](commands/) - CLI command specifications
- [Feature Specs](features/) - Non-command feature specifications
- [Contributing Guide](contributing.md) - General contribution guidelines
- [Go Style Guide](style/go-style.md) - Code patterns and conventions
