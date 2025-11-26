# ADR-0001: ADR + Spec Development Workflow

| Metadata | Value |
|----------|-------|
| **Status** | Accepted |
| **Date** | 2025-11-25 |
| **Author(s)** | @anowarislam |
| **Issue** | N/A (bootstrap) |
| **Related ADRs** | N/A (first ADR) |

## Context

The `ado` project follows a spec-driven development approach where command specifications in `docs/commands/` are written before implementation. However, the project lacks:

1. **A formal process for architectural decisions** - No documentation of why choices were made
2. **Feature specifications for non-command features** - Only CLI commands have specs
3. **A documented workflow connecting Issues to Implementation** - No clear path from idea to code
4. **Clear guidance on when different documentation types are required** - Contributors unsure what to write

As the project grows, these gaps create:

- Inconsistent decision-making processes
- Lost context on why decisions were made
- Unclear contribution paths for different types of changes
- Difficulty onboarding new contributors

## Decision

We will implement a **three-phase development workflow**:

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│    Issue     │ ──▶ │  ADR (PR 1)  │ ──▶ │  Spec (PR 2) │ ──▶ │  Code (PR 3) │
│  (Proposal)  │     │  (if needed) │     │              │     │              │
└──────────────┘     └──────────────┘     └──────────────┘     └──────────────┘
```

### Phase 1: Decision (ADR)

- Location: `docs/adr/`
- Format: `NNNN-title.md` (4-digit numbering)
- Required for architectural decisions
- Results in 1 PR

### Phase 2: Specification

- Commands: `docs/commands/` (existing)
- Features: `docs/features/` (new)
- Details the "what" to implement
- Results in 1 PR

### Phase 3: Implementation

- Code in `cmd/ado/` or `internal/`
- Tests alongside implementation
- Results in 1 PR

### ADRs are required for:

- New architectural patterns
- Breaking changes to existing patterns
- New dependencies or tools
- Security-related decisions
- Changes affecting multiple components

### ADRs are optional for:

- Single command additions (use command spec)
- Bug fixes
- Documentation improvements
- Minor refactoring

## Consequences

### Positive

- **Clear documentation trail** for architectural decisions
- **Consistent process** for all contributors
- **Explicit checkpoints** prevent wasted implementation effort
- **Historical context** preserved for future maintainers
- **Specs serve as acceptance criteria** for implementation PRs
- **Enables async collaboration** with LLM-assisted development

### Negative

- **Additional process overhead** for small changes
- **Learning curve** for new contributors
- **Risk of over-documentation** for simple changes

### Neutral

- Requires discipline to maintain workflow
- May evolve as project needs change
- Three PRs per feature may feel slow initially

## Alternatives Considered

### Alternative 1: Single PR with Inline Documentation

Keep all decisions in commit messages and PR descriptions.

**Why not chosen:** Context is scattered and hard to find. GitHub PRs are not optimized for long-term reference. Decisions get buried in closed PR history.

### Alternative 2: Wiki-based Documentation

Use GitHub Wiki for decisions and specs.

**Why not chosen:** Wiki is disconnected from code review process. Changes not versioned with code. No PR-based approval workflow.

### Alternative 3: RFCs Instead of ADRs

Use a Request for Comments (RFC) process for all changes.

**Why not chosen:** RFCs are heavier weight and better suited for open-source projects with many external contributors. ADRs are specifically designed for architectural decisions and are more lightweight for a focused project.

### Alternative 4: All-in-One Design Documents

Combine ADR, spec, and implementation plan in single documents.

**Why not chosen:** Separating phases allows for checkpoints and prevents wasted effort. ADRs focus on "why", specs on "what", keeping concerns separated.

## Implementation Notes

This ADR establishes the workflow. Implementation includes:

1. **PR 1 (this PR):** ADR framework
   - `docs/adr/README.md` - Index and process guide
   - `docs/adr/TEMPLATE.md` - ADR template
   - `docs/adr/0001-development-workflow.md` - This ADR
   - `.github/ISSUE_TEMPLATE/adr_proposal.md` - Issue template
   - Updates to `CLAUDE.md` and `README.md`

2. **PR 2:** Specification framework
   - `docs/features/` directory with template
   - `docs/commands/TEMPLATE.md`
   - `docs/workflow.md` - Complete workflow guide
   - Issue templates for feature and command proposals

3. **PR 3:** Documentation integration
   - Updates to `docs/contributing.md`
   - Updates to `docs/style/docs-style.md`
   - Updates to `CONTRIBUTING.md`

## References

- [ADR GitHub Organization](https://adr.github.io/)
- [Documenting Architecture Decisions](https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions) - Michael Nygard
- [Lightweight Architecture Decision Records](https://www.thoughtworks.com/radar/techniques/lightweight-architecture-decision-records) - ThoughtWorks
- [GitHub Spec Kit](https://github.com/github/spec-kit) - Spec-driven development patterns
