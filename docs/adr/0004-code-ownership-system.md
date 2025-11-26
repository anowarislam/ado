# ADR-0004: Code Ownership System

| Metadata | Value |
|----------|-------|
| **Status** | Accepted |
| **Date** | 2025-11-26 |
| **Author(s)** | @anowarislam |
| **Issue** | #47 |
| **Related ADRs** | ADR-0001 (Development Workflow) |

## Context

The `ado` project follows a three-phase development workflow (Issue → ADR → Spec → Implementation) as documented in ADR-0001. As the project grows beyond a single maintainer, we need a systematic approach to code review that ensures:

- **Expertise-based review**: Changes are reviewed by maintainers with knowledge of the affected areas
- **Clear accountability**: Each component has defined owners responsible for its quality
- **Efficient review process**: Automatic reviewer assignment eliminates manual effort
- **Enforceable standards**: Branch protection ensures required approvals before merging
- **Scalability**: System can grow from 1 maintainer to a team without process changes

**Current State**:
- Manual reviewer assignment (contributors must know who to ask)
- No automated enforcement of review requirements
- Implicit ownership (not documented in code)
- Risk of merging without appropriate review

**Triggering Event**:
Need to establish formal code ownership before adding additional maintainers to the project.

## Decision

We will implement GitHub's native **CODEOWNERS** system with the following approach:

### Core Decision

Use `.github/CODEOWNERS` file to define file-path-based code ownership, integrated with GitHub branch protection to enforce required approvals.

### Why GitHub CODEOWNERS (not Kubernetes-style OWNERS)

After evaluating both approaches:

| Factor | GitHub CODEOWNERS | Kubernetes OWNERS (Prow) |
|--------|-------------------|--------------------------|
| **Setup complexity** | Zero infrastructure (native GitHub) | Requires Prow CI/CD deployment |
| **Maintenance** | Minimal (one file) | Higher (aliases, per-directory files) |
| **Two-phase review** | No (/lgtm + /approve separation not supported) | Yes (reviewers vs approvers) |
| **Team size** | Optimal for 1-20 developers | Designed for 50+ developers |
| **GitHub integration** | Native, automatic | Requires Prow bot |

**Conclusion**: GitHub CODEOWNERS is the right fit for the current and foreseeable team size (1-10 maintainers).

### Implementation Components

**1. `.github/CODEOWNERS` File**

Define ownership using gitignore-style patterns:

```
*                         @anowarislam  # Default
/cmd/ado/                 @anowarislam  # Commands
/internal/                @anowarislam  # Core libraries
/docs/adr/                @anowarislam  # Architecture decisions
/.github/workflows/       @anowarislam  # CI/CD
```

**Key Principles**:
- Last matching pattern wins (more specific overrides general)
- Hierarchical (can subdivide ownership as team grows)
- Explicit (no implicit ownership)

**2. Branch Protection Integration**

Enable "Require review from Code Owners" setting on `main` branch, which:
- Automatically requests reviews from code owners when PR is opened
- Blocks merge until code owner approves
- Dismisses stale approvals when new commits pushed
- Works alongside existing CI checks (tests, coverage, linting)

**3. Documentation**

Create comprehensive documentation:
- `docs/code-ownership.md`: Complete ownership guide with FAQ
- `docs/workflow.md`: Integration with three-phase workflow
- Update contributing guides with review process

**4. Current Ownership Structure**

Start with single code owner (@anowarislam) for all paths, with more specific ownership for critical areas:

- Security-critical files: SECURITY.md, security documentation
- Release automation: .goreleaser.yaml, release-please config
- CI/CD: .github/workflows/
- Architecture decisions: docs/adr/

This provides foundation for distributing ownership as maintainers join.

## Consequences

### Positive

- **Automatic reviewer assignment**: GitHub assigns code owners when PRs opened, eliminating manual effort
- **Enforced review**: Branch protection prevents merge without code owner approval
- **Clear accountability**: CODEOWNERS file documents who owns each component
- **Scalability**: Easy to add more owners (just edit CODEOWNERS file)
- **Audit trail**: Git history tracks ownership changes over time
- **Knowledge transfer**: Review process facilitates learning
- **Consistency**: Ensures architectural decisions maintained
- **Zero infrastructure**: Native GitHub feature, no external tools
- **Compatibility**: Works with existing three-phase workflow (ADR-0001)
- **Future-proof**: Can migrate to Kubernetes-style OWNERS later if team scales to 50+

### Negative

- **Initial overhead**: Need to document ownership model and train contributors
- **Single owner bottleneck**: With one maintainer, all PRs require @anowarislam approval (can't merge own PRs without workaround)
- **No two-phase review**: Can't separate detailed code review from architectural approval
- **Manual re-assignment**: Can't automatically load-balance among multiple owners of same path
- **Pattern maintenance**: CODEOWNERS file can become complex if overly granular

### Neutral

- **Ownership model**: Reflects current reality (single maintainer) but provides structure for growth
- **GitHub dependency**: Relies on GitHub-specific feature (but already committed to GitHub)
- **Documentation burden**: Adds code-ownership.md to maintain (but provides value)

## Alternatives Considered

### Alternative 1: Manual Review Assignment

**Description**: Continue current practice of manually assigning reviewers in GitHub UI.

**Why not chosen**:
- **No enforcement**: Easy to forget to assign reviewer
- **No documentation**: Ownership not visible in codebase
- **Not scalable**: Becomes painful with multiple maintainers
- **Inconsistent**: Different contributors might assign different reviewers

### Alternative 2: Kubernetes-Style OWNERS Files

**Description**: Use Prow CI/CD with OWNERS files in each directory, supporting /lgtm and /approve commands.

**Why not chosen**:
- **Infrastructure overhead**: Requires deploying and maintaining Prow bot
- **Over-engineered**: Two-phase review (reviewer vs approver) not needed for small team
- **Complexity**: Multiple OWNERS files across directories harder to maintain
- **Cost**: Prow deployment and operation requires effort
- **Team size**: Designed for 50-500+ developer teams (overkill for 1-10)

However, this remains a **valid future option** if:
- Team grows to 20+ active developers
- Need emerges for two-phase review (implementation vs architecture approval)
- Multiple SIGs (Special Interest Groups) form around different areas

### Alternative 3: No Formal Ownership

**Description**: Keep implicit ownership based on commit history and expertise.

**Why not chosen**:
- **No automation**: Manual effort to find reviewers
- **No enforcement**: Nothing prevents merge without review
- **Not documented**: New contributors don't know who to ask
- **Doesn't scale**: Breaks down as team grows

### Alternative 4: Review-by-Anyone Model

**Description**: Any maintainer can approve any change (no specialized ownership).

**Why not chosen**:
- **Loss of expertise**: Changes reviewed by people unfamiliar with the area
- **Inconsistent architecture**: Different reviewers may have different standards
- **Security risk**: Critical files (workflows, security docs) should have specific reviewers
- **Knowledge silos**: Doesn't encourage cross-component understanding

## Implementation Notes

**Completed**:
- [x] Created `.github/CODEOWNERS` with comprehensive path mappings
- [x] Created `docs/code-ownership.md` (comprehensive guide with FAQ)
- [x] Updated `docs/workflow.md` (added CODEOWNERS integration section)
- [x] Updated `docs/contributing.md` (added code review process)
- [x] Updated `CONTRIBUTING.md` (added ownership info)
- [x] Updated `README.md` (added ownership references)
- [x] Updated `mkdocs.yml` navigation
- [x] Created GitHub Issue #47
- [x] Created this ADR

**Next Steps**:
- [ ] Enable branch protection setting: "Require review from Code Owners" on main branch
- [ ] Merge PR with CODEOWNERS implementation
- [ ] Monitor and adjust ownership patterns as needed

**Migration Path (if Needed)**:

If team grows to 20+ developers and Kubernetes-style OWNERS becomes beneficial:

**Phase 1**: Add OWNERS_ALIASES at repository root (defines groups)
**Phase 2**: Create per-directory OWNERS files referencing aliases
**Phase 3**: Deploy Prow bot for /lgtm and /approve commands
**Phase 4**: Migrate patterns from CODEOWNERS to OWNERS files
**Phase 5**: Enable Tide for automatic merge queue

Estimated effort: 2-4 weeks depending on team coordination.

**Success Criteria**:
- All PRs have automatic code owner assignment
- Code owner approval required before merge (enforced by GitHub)
- Contributors understand ownership model (documented in code-ownership.md)
- Ownership boundaries clear and documented (CODEOWNERS file)

## References

- [Issue #47](https://github.com/anowarislam/ado/issues/47): Feature proposal and implementation tracking
- [ADR-0001](0001-development-workflow.md): Three-phase development workflow
- [GitHub CODEOWNERS Documentation](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-code-owners)
- [GitHub Branch Protection](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-protected-branches/about-protected-branches)
- [Kubernetes OWNERS Guide](https://www.kubernetes.dev/docs/guide/owners/) - Alternative approach for larger teams
- [Prow Documentation](https://docs.prow.k8s.io/) - Kubernetes CI/CD with OWNERS support
- [CNCF Security Best Practices](https://contribute.cncf.io/projects/best-practices/security/security-hygiene/)
