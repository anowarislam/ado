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

- [ADRs](adr/) - Architecture Decision Records
- [Command Specs](commands/) - CLI command specifications
- [Feature Specs](features/) - Non-command feature specifications
- [Contributing Guide](contributing.md) - General contribution guidelines
- [Go Style Guide](style/go-style.md) - Code patterns and conventions
