# ADR-0005: PR-Level Metrics Dashboard

| Metadata | Value |
|----------|-------|
| **Status** | Proposed |
| **Date** | 2025-11-26 |
| **Author(s)** | @anowarislam |
| **Issue** | #50 |
| **Related ADRs** | ADR-0003 (Recipe-Based Documentation) |

## Context

The `ado` project has robust CI/CD infrastructure with comprehensive testing, 80% coverage enforcement, and Codecov integration. However, **PR-level visibility into quality metrics is limited**:

**Current state:**
- Tests run in CI but failures require navigating to Actions tab
- Coverage uploaded to Codecov dashboard (external, requires separate visit)
- No coverage impact shown on PRs (diff coverage, per-file breakdowns)
- No performance regression detection
- No CI cost tracking or trends
- Developers must manually check multiple sources to assess PR quality

**Problems this creates:**
1. **Delayed feedback**: Developers find test failures after context-switching away
2. **Hidden coverage impact**: Don't know if PR lowers coverage until merge
3. **No performance visibility**: Regressions discovered in production
4. **Cost blindness**: No awareness of expensive workflows or optimization opportunities
5. **Poor review experience**: Reviewers lack objective quality signals in PR

**What triggered this discussion?**

Issue #50 proposes a 4-phase enhancement to bring metrics directly into the PR:
1. Enhanced test reporting (rich summaries, failure annotations)
2. Granular coverage enforcement (per-file, per-package, total thresholds)
3. PR coverage comments (diff coverage, uncovered lines visualization)
4. Cost tracking and benchmark monitoring (workflow costs, performance trends)

This is complementary to Issue #44 (broader observability with OpenTelemetry) but focused specifically on **PR-time developer feedback**.

## Decision

**We will implement a GitHub-native PR metrics dashboard** that displays test results, coverage analysis, performance benchmarks, and CI costs directly in Pull Requests.

### Core Decisions

#### 1. Display Location: GitHub-Native Integration

**Primary**: GitHub Status Checks + PR Comments
- Status checks for pass/fail signals (visible in PR header)
- PR comments for detailed metrics (collapsible, updatable)
- No external dashboard required

**Rationale:**
- Developers already use GitHub; no context-switching
- Status checks block merge on failures (enforcement)
- PR comments provide rich formatting (markdown tables, badges)
- No additional hosting or authentication needed

#### 2. Metrics to Track

**Phase 1: Test Reporting**
- Test pass/fail counts by package
- Failure annotations with stack traces
- Test execution time per package

**Phase 2: Coverage Analysis**
- Total coverage percentage (project-wide)
- Per-package coverage thresholds
- Per-file coverage thresholds
- Diff coverage (new/changed code)

**Phase 3: Coverage Visualization**
- Uncovered lines in changed files
- Coverage trend (vs main branch)
- Coverage change per commit

**Phase 4: Performance & Costs**
- CI workflow duration and cost per PR
- Historical cost trends
- Benchmark results with regression detection

#### 3. Tool Stack

**Chosen Tools:**
- **Test reporting**: `robherley/go-test-action` (GitHub Action for rich test summaries)
- **Coverage enforcement**: `vladopajic/go-test-coverage` (YAML-based thresholds)
- **Coverage comments**: `fgrosse/go-coverage-report` (diff coverage PR comments)
- **Cost tracking**: GitHub API (workflow run metadata) + custom reporting
- **Benchmarks**: `benchstat` (Go standard tooling) + custom trend analysis

**Why these tools:**
- All integrate natively with GitHub Actions
- No external services or tokens required
- Leverage existing `go tool cover` and `go test -bench`
- Open source, actively maintained
- Proven in production use

#### 4. Integration with Existing Infrastructure

**Codecov retention:**
- Keep existing Codecov integration (`continue-on-error: true`)
- Use for long-term trend analysis and historical data
- PR dashboard provides immediate feedback; Codecov provides history

**Coverage threshold:**
- Maintain 80% total threshold (existing)
- Add per-package thresholds (80% default, configurable)
- Add per-file thresholds (70% for new files)
- Add diff coverage threshold (85% for changed lines)

**Workflow changes:**
- Enhance `.github/workflows/ci.yml` with new actions
- Add `.github/.testcoverage.yml` for granular thresholds
- Add workflow job for PR comment updates
- Add benchmark workflow (manual trigger + PR automation)

#### 5. Implementation Approach: 4 Sequential Phases

**Phase 1: Enhanced Test Reporting** (Immediate)
- Add `robherley/go-test-action` to CI workflow
- Provides rich test summaries and failure annotations
- Low risk, high value (better test visibility)

**Phase 2: Granular Coverage Enforcement** (Short-term)
- Add `vladopajic/go-test-coverage` for threshold checks
- Create `.github/.testcoverage.yml` config
- Enforce per-package and per-file coverage

**Phase 3: PR Coverage Comments** (Medium-term)
- Add `fgrosse/go-coverage-report` for PR comments
- Display diff coverage and uncovered lines
- Update comment on each push (not spam)

**Phase 4: Cost & Benchmark Tracking** (Long-term)
- Implement GitHub API cost calculation
- Add benchmark workflow with regression detection
- Create historical trend dashboard (GitHub Pages)

**Why sequential:**
- Each phase delivers value independently
- Allows testing and iteration before next phase
- Reduces risk of breaking existing CI
- Easier code review and rollback if needed

## Consequences

### Positive

- **Immediate feedback**: Developers see test/coverage results in PR, no navigation required
- **Better enforcement**: Per-file/package thresholds prevent localized coverage drops
- **Improved reviews**: Reviewers have objective quality signals (coverage, performance)
- **Cost awareness**: Teams can identify and optimize expensive workflows
- **Performance safety**: Benchmark regression detection prevents performance bugs
- **GitHub-native**: No external services, tokens, or hosting required
- **Transparent**: All metrics visible to contributors, no hidden dashboards
- **Incremental adoption**: 4-phase approach allows gradual rollout

### Negative

- **Maintenance burden**: More GitHub Actions to maintain and debug
- **Noise risk**: Too many metrics could overwhelm PR view (mitigated by collapsible comments)
- **CI complexity**: Additional workflow jobs increase pipeline complexity
- **False positives**: Coverage thresholds may be too strict for some files (mitigated by configuration)
- **Limited history**: Status checks/comments don't provide long-term trends (mitigated by keeping Codecov)
- **Implementation time**: 4 phases could take 15-20 hours total

### Neutral

- **Codecov complementary**: Keep for historical trends; PR dashboard for immediate feedback
- **Configuration overhead**: Need to maintain `.testcoverage.yml` in addition to workflow
- **GitHub coupling**: Tied to GitHub platform (but already committed to GitHub)
- **Action dependencies**: Rely on third-party actions (but all are popular, well-maintained)

## Alternatives Considered

### Alternative 1: Codecov Comments Only

**Description:** Use only Codecov's built-in PR comments for coverage visibility.

**Why not chosen:**
- No test result summaries or annotations
- No per-package or per-file thresholds
- No cost or benchmark tracking
- Codecov comments can be delayed (upload lag)
- Requires Codecov Pro for advanced features (cost)

### Alternative 2: External Dashboard (e.g., Vercel, Custom Build)

**Description:** Build separate web dashboard hosted on Vercel or GitHub Pages with detailed metrics.

**Why not chosen:**
- Requires context-switching away from PR
- Needs authentication setup
- Additional hosting/maintenance burden
- Developers less likely to check external dashboard
- Status checks still needed for enforcement (no elimination of GitHub integration)

### Alternative 3: Custom GitHub App

**Description:** Build custom GitHub App that posts rich metrics via Checks API.

**Why not chosen:**
- High development and maintenance cost
- Requires hosting for app server
- Authentication and token management complexity
- Existing GitHub Actions provide same functionality
- Overkill for current needs

### Alternative 4: GitHub Actions Summary Only

**Description:** Use GitHub Actions job summaries (markdown in Actions tab) for metrics.

**Why not chosen:**
- Hidden in Actions tab (not visible in PR)
- Requires navigation away from PR
- No status check integration (can't block merge)
- Does not solve the "lack of PR visibility" problem

### Alternative 5: Do Nothing (Keep Current State)

**Description:** Maintain current setup (Codecov + basic CI).

**Why not chosen:**
- Does not address developer feedback delays
- Coverage impact remains hidden until merge
- No performance regression detection
- No cost visibility or optimization opportunities
- Reviewers lack objective quality signals

## Implementation Notes

**Completed after ADR approval:**

1. **Feature Spec**: `docs/features/03-pr-metrics-dashboard.md`
   - Define exact workflow changes
   - Specify comment format and content
   - Document configuration schema
   - Provide example outputs

2. **Implementation (4 PRs)**:
   - PR 1: Phase 1 - Enhanced test reporting
   - PR 2: Phase 2 - Granular coverage enforcement
   - PR 3: Phase 3 - PR coverage comments
   - PR 4: Phase 4 - Cost and benchmark tracking

3. **Documentation Updates**:
   - Update `docs/recipes/03-ci-components.md` with new patterns
   - Update `CLAUDE.md` with new workflow details
   - Update `docs/workflow.md` with PR quality standards

**Success Criteria:**
- Developers see test failures in PR without navigating to Actions
- Coverage impact visible in PR comment within 2 minutes of push
- Per-package coverage violations block merge
- Performance regressions detected and reported in PR
- CI workflow costs visible and trending downward

## References

- [Issue #50: PR-level metrics dashboard](https://github.com/anowarislam/ado/issues/50)
- [Issue #44: Observability strategy with OpenTelemetry](https://github.com/anowarislam/ado/issues/44) (complementary)
- [ADR-0003: Recipe-Based Documentation](0003-recipe-based-documentation.md) (CI/CD patterns)
- [robherley/go-test-action](https://github.com/robherley/go-test-action) - Test reporting
- [vladopajic/go-test-coverage](https://github.com/vladopajic/go-test-coverage) - Coverage thresholds
- [fgrosse/go-coverage-report](https://github.com/fgrosse/go-coverage-report) - PR comments
- [GitHub Actions Checks API](https://docs.github.com/en/rest/checks) - Status checks
- [Codecov GitHub Action](https://github.com/codecov/codecov-action) - Existing integration
- [GitHub Actions Usage API](https://docs.github.com/en/rest/actions/workflows#get-workflow-usage) - Cost tracking
