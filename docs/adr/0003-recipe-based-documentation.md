# ADR-0003: Recipe-Based Documentation for CI/CD Patterns

| Metadata | Value |
|----------|-------|
| **Status** | Accepted |
| **Date** | 2025-11-26 |
| **Author(s)** | @anowarislam |
| **Issue** | #46 |
| **Related ADRs** | ADR-0001 (Development Workflow) |

## Context

The `ado` project has evolved a comprehensive CI/CD system with multiple layers of automation, security features, and development practices. However, this knowledge exists only as implementation code and scattered documentation. Teams wanting to replicate this system face several challenges:

- **Fragmented knowledge**: CI/CD practices spread across workflows, hooks, Makefiles, and docs
- **Implicit patterns**: Many decisions and patterns exist only in code, not explicitly documented
- **High barrier to entry**: Understanding the complete system requires reading multiple files and inferring connections
- **Python-specific needs**: Most examples are Go-focused; Python projects need adapted patterns
- **Missing rationale**: Implementation exists, but the "why" behind decisions is unclear
- **No guided path**: Teams don't know where to start or what order to implement features

**What triggered this discussion?**

The realization that:
1. The ado project's CI/CD patterns are production-proven and valuable to share
2. Python startup teams repeatedly face the same CI/CD challenges
3. Existing tutorials teach tools, not complete systems
4. LLM-assisted development benefits from structured, spec-driven workflows
5. Modern supply chain security (SLSA, signing) lacks practical guides

## Decision

We will create a comprehensive **recipe-based documentation system** under `docs/recipes/` that:

1. **Organizes by topic** (not chronological): Foundation → Process → Integration → Implementation
2. **Provides complete examples**: All code is copy-paste ready, not pseudocode
3. **Explains the "why"**: Each pattern includes rationale and trade-offs
4. **Includes visual aids**: 80+ Mermaid diagrams for system understanding
5. **Targets Python projects**: Python-specific adaptations (pyproject.toml, pytest, ruff, PyPI)
6. **Focuses on frameworks**: Principles and patterns, not specific metrics or timelines
7. **Enables self-service**: Teams can implement at their own pace, in their order
8. **Maintains in-repo**: Documentation lives with code, evolves together

**Structure:**
```
docs/recipes/
├── README.md                    # Executive overview
├── 00-overview.md              # Navigation hub
├── 01-philosophy.md            # Foundation: Why & how
├── 02-ci-components.md         # Foundation: System architecture
├── 03-security-features.md     # Foundation: SLSA, signing
├── 04-development-workflow.md  # Process: Issues→ADR→Spec
├── 05-build-automation.md      # Process: Make, hooks
├── 06-release-automation.md    # Process: Zero-touch releases
├── 07-github-integrations.md   # Integration: Apps, Dependabot
├── 08-troubleshooting.md       # Integration: Problem-solving
├── 09-implementation-guide.md  # Implementation: Step-by-step
├── 10-python-adaptation.md     # Implementation: Python patterns
└── 11-principles.md            # Implementation: Standards framework
```

**Integration:**
- Add "CI/CD Recipe" section to MkDocs navigation
- Link from main documentation as "Battle-tested CI/CD patterns"
- Reference from ADRs and specs where relevant
- Keep updated as ado project evolves

## Consequences

### Positive

- **Reproducible patterns**: Teams can replicate ado's CI/CD without reverse-engineering
- **Reduced support burden**: Self-service documentation reduces questions
- **Knowledge preservation**: Decisions and rationale documented for future team members
- **Broader impact**: Helps the wider Python community adopt modern practices
- **Better onboarding**: New contributors understand the complete system faster
- **LLM-friendly**: Structured content enables AI-assisted implementation
- **Living documentation**: Maintained with code, stays current
- **Framework thinking**: Focuses on principles over metrics, more timeless
- **Compliance-ready**: Security patterns meet modern supply chain requirements

### Negative

- **Maintenance burden**: 84,596 words to keep updated as project evolves
- **Potential divergence**: Recipe examples may drift from actual implementation
- **Initial effort**: Significant upfront time to create comprehensive content
- **Scope creep risk**: Temptation to document everything vs. essential patterns
- **Python-specific**: Less directly applicable to non-Python projects (though patterns translate)

### Neutral

- **Documentation architecture**: Establishes pattern for other recipe-style docs
- **Size consideration**: Large documentation (84KB total) but organized and navigable
- **External vs internal**: Primarily targets external teams, but useful for internal onboarding too
- **Opinionated approach**: Reflects ado's specific choices, not universal truth

## Alternatives Considered

### Alternative 1: Single Tutorial Page

**Description:** Create one comprehensive tutorial page covering all CI/CD aspects.

**Why not chosen:**
- **Too overwhelming**: 84,000 words in one page is unusable
- **Hard to navigate**: No way to jump to specific topics
- **Poor for reference**: Can't bookmark specific sections easily
- **Unmaintainable**: Single file becomes merge conflict nightmare

### Alternative 2: External Wiki

**Description:** Use GitHub Wiki or separate documentation site for CI/CD patterns.

**Why not chosen:**
- **Disconnected from code**: Wiki doesn't evolve with codebase
- **No version control**: Can't track changes or review edits
- **No CI validation**: Links break, examples drift, no automated checks
- **Extra maintenance**: Separate place to update

### Alternative 3: Blog Post Series

**Description:** Write series of blog posts about CI/CD implementation.

**Why not chosen:**
- **Temporal structure**: Hard to keep updated over time
- **No organization**: Posts scattered chronologically, not by topic
- **External platform**: Not maintained with code
- **No search**: Harder to find specific information

### Alternative 4: Separate Repository

**Description:** Create dedicated `ado-cicd-guide` repository.

**Why not chosen:**
- **Divergence risk**: Examples drift from actual implementation
- **Duplication**: Two repos to maintain
- **Discoverability**: Users must find separate repo
- **Reference complexity**: Harder to link between docs and code

### Alternative 5: Just Code Comments

**Description:** Document patterns only in code comments and existing docs.

**Why not chosen:**
- **Scattered knowledge**: No central place to understand complete system
- **Implementation-first**: Hard to learn architecture from code
- **No guided path**: Users don't know where to start
- **Implicit rationale**: "Why" decisions remain hidden

## Implementation Notes

**Completed in PR #45:**
- [x] All 13 recipe chapters created (84,596 words)
- [x] Added to MkDocs navigation under "CI/CD Recipe"
- [x] 80+ Mermaid diagrams for visual learning
- [x] 300+ copy-paste ready code examples
- [x] Enhanced pre-push hook validates documentation build
- [x] CI validates all internal links and anchors
- [x] Python-specific adaptations throughout
- [x] Framework-focused (not metrics-driven)

**Ongoing maintenance:**
- Update recipes when ado CI/CD patterns change
- Add new chapters as new patterns emerge
- Keep code examples synchronized with actual implementation
- Review and update on each major release
- Incorporate feedback from teams using the recipes

**Success criteria:**
- Teams successfully implement CI/CD using recipe guidance
- Internal team members reference recipes for understanding
- External contributions reference recipes for context
- Documentation stays current with codebase (checked in CI)

## References

- [Issue #46](https://github.com/anowarislam/ado/issues/46): Feature proposal
- [PR #45](https://github.com/anowarislam/ado/pull/45): Implementation
- [ADR-0001](0001-development-workflow.md): Development workflow (referenced in recipes)
- [CNCF Project Documentation](https://www.cncf.io/): Recipe pattern inspiration
- [AWS Solutions Library](https://aws.amazon.com/solutions/): Solution-focused documentation pattern
- [Conventional Commits](https://www.conventionalcommits.org/): Referenced throughout recipes
- [SLSA Framework](https://slsa.dev/): Security patterns documented in recipe
- [Sigstore](https://www.sigstore.dev/): Signing patterns documented in recipe
