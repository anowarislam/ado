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

Issue #50 proposes comprehensive enhancements to bring metrics directly into the PR:
- Enhanced test reporting (rich summaries, failure annotations)
- Granular coverage enforcement (per-file, per-package, total thresholds)
- PR coverage comments (diff coverage, uncovered lines visualization)
- Cost tracking and benchmark monitoring (workflow costs, performance trends)

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

**Test Reporting**
- Test pass/fail counts by package
- Failure annotations with stack traces
- Test execution time per package

**Coverage Analysis**
- Total coverage percentage (project-wide)
- Per-package coverage thresholds (80% default, configurable)
- Per-file coverage thresholds (70% for new files)
- Diff coverage (85% for changed lines)

**Coverage Visualization**
- Uncovered lines in changed files (PR comment)
- Coverage trend (vs main branch)
- Coverage change per commit

**Performance & Costs**
- CI workflow duration with estimated cost per PR
- Estimated cost calculation (duration × runner rate)
- Benchmark results with regression detection
- Cost estimates clearly labeled (not actual billed minutes)

#### 3. Tool Stack

**Chosen Tools:**
- **Test reporting**: `robherley/go-test-action` (GitHub Action for rich test summaries)
- **Coverage enforcement**: `vladopajic/go-test-coverage` (YAML-based thresholds)
- **Coverage comments**: `fgrosse/go-coverage-report` (diff coverage PR comments)
- **Cost tracking**: Estimated from workflow duration via GitHub Actions Performance Metrics
- **Benchmarks**: `benchstat` (Go standard tooling) + custom trend analysis

**Cost Estimation Methodology:**
- Calculate: `Estimated Cost = Duration (minutes) × Runner Rate`
- Runner rates (November 2025): ubuntu-latest ($0.008/min), macos-latest ($0.016/min)
- Uses workflow timing from GitHub Actions context (no API calls required)
- Note: Previous billing API deprecated April 2025; new consolidated billing API provides actual data with 24-48hr delay

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

#### 5. Implementation Approach

**Single unified implementation** that adds all metrics to CI workflow:

- Enhance `.github/workflows/ci.yml` with new actions and jobs
- Add `.testcoverage.yml` at repo root for granular coverage thresholds
- Add test reporting action (`robherley/go-test-action`, pinned to SHA)
- Add coverage enforcement action (`vladopajic/go-test-coverage`, pinned to SHA)
- Add PR comment action (`fgrosse/go-coverage-report`, pinned to SHA)
- Add workflow job for estimated cost calculation (duration × runner rate)
- Add benchmark workflow with regression detection
- Create/update `.codecov.yml` to disable PR comments (avoid duplication with new comment action)

**Implementation scope:**
- Workflow modifications: ~100-150 lines of YAML
- Configuration file: ~30-50 lines (`.testcoverage.yml` at repo root)
- Codecov config update: 2-3 lines (disable comment feature)
- Documentation updates: CLAUDE.md, recipes, workflow.md
- Tests: Validate workflow syntax, test coverage config parsing

**Codecov coordination:**
- Disable Codecov PR comments to avoid duplication: `comment: false` in `.codecov.yml`
- Keep Codecov upload for historical trends and dashboard
- New PR comment action provides richer diff coverage visualization

The changes are straightforward GitHub Actions additions - no need for multiple PRs.

## Consequences

### Positive

- **Immediate feedback**: Developers see test/coverage results in PR, no navigation required
- **Better enforcement**: Per-file/package thresholds prevent localized coverage drops
- **Improved reviews**: Reviewers have objective quality signals (coverage, performance)
- **Cost awareness**: Teams can identify and optimize expensive workflows
- **Performance safety**: Benchmark regression detection prevents performance bugs
- **GitHub-native**: No external services, tokens, or hosting required
- **Transparent**: All metrics visible to contributors, no hidden dashboards
- **Simple implementation**: Straightforward GitHub Actions additions, single PR

### Negative

- **Maintenance burden**: More GitHub Actions to maintain and debug
- **Noise risk**: Too many metrics could overwhelm PR view (mitigated by collapsible comments)
- **CI complexity**: Additional workflow jobs increase pipeline complexity
- **False positives**: Coverage thresholds may be too strict for some files (mitigated by configuration)
- **Limited history**: Status checks/comments don't provide long-term trends (mitigated by keeping Codecov)
- **CI runtime increase**: Additional jobs may add 1-2 minutes to CI pipeline

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

## Security Considerations

**Supply Chain Security:**
- **Pin all third-party actions to SHA** (not tags) for immutability
- Aligns with SLSA Level 3 requirements (ADR-0003)
- Example: `uses: robherley/go-test-action@v1.2.3` → `uses: robherley/go-test-action@abc123...`

**Permissions:**
- Test reporting: `contents: read` (default)
- Coverage comments: `pull-requests: write` (minimal required)
- Cost tracking: `actions: read` (workflow metadata only)
- Follows principle of least privilege

**Dependency Verification:**
All chosen actions are:
- Actively maintained (verified 2024 commits)
- Widely adopted in community
- GitHub Actions-native or well-established

## Rollback Strategy

If issues arise post-implementation:

1. **Immediate rollback**: Disable actions via workflow conditions
   ```yaml
   if: false  # Temporary disable while investigating
   ```

2. **Partial rollback**: Comment out specific jobs (test reporting, coverage comments, cost tracking independently)

3. **Full rollback**: Revert workflow changes via git
   - Codecov remains primary coverage source (no data loss)
   - All changes are additive (no removal of existing functionality)

4. **Configuration tuning**: Adjust thresholds in `.testcoverage.yml` without workflow changes

**Risk mitigation**: Changes are low-risk, easily reversible, and Codecov provides continuity.

## Implementation Notes

**Completed after ADR approval:**

1. **Feature Spec**: `docs/features/03-pr-metrics-dashboard.md`
   - Define exact workflow changes
   - Specify comment format and content
   - Document configuration schema
   - Provide example outputs

2. **Implementation (Single PR)**:
   - Add all GitHub Actions to `.github/workflows/ci.yml`
   - Create `.github/.testcoverage.yml` configuration
   - Add benchmark workflow
   - Update documentation (CLAUDE.md, recipes, workflow.md)
   - Test workflow changes

3. **Files Modified**:
   - `.github/workflows/ci.yml` (enhanced with metrics actions)
   - `.testcoverage.yml` (new configuration file at repo root)
   - `.github/workflows/benchmark.yml` (new workflow)
   - `.codecov.yml` (disable PR comments to avoid duplication)
   - `docs/recipes/03-ci-components.md` (updated patterns)
   - `CLAUDE.md` (updated workflow details)
   - `docs/workflow.md` (PR quality standards)

**Success Criteria (Measurable):**
- **Test visibility**: Developers see test failures in PR without navigating to Actions (100% of PRs)
- **Coverage timeliness**: Coverage impact visible in PR comment within 2 minutes of push (95% of runs)
- **Enforcement effectiveness**: Per-package coverage violations block merge with clear error messages
- **Performance safety**: Performance regressions >5% detected and reported in PR
- **Cost visibility**: Estimated CI costs visible in PR comments with <5% overhead
- **Cost accuracy**: Estimates within 10% of actual billed minutes (verified quarterly)
- **Developer satisfaction**: Zero complaints about metric noise in first 3 months

## References

- [Issue #50: PR-level metrics dashboard](https://github.com/anowarislam/ado/issues/50)
- [Issue #44: Observability strategy with OpenTelemetry](https://github.com/anowarislam/ado/issues/44) (complementary)
- [ADR-0003: Recipe-Based Documentation](0003-recipe-based-documentation.md) (CI/CD patterns)
- [robherley/go-test-action](https://github.com/robherley/go-test-action) - Test reporting
- [vladopajic/go-test-coverage](https://github.com/vladopajic/go-test-coverage) - Coverage thresholds
- [fgrosse/go-coverage-report](https://github.com/fgrosse/go-coverage-report) - PR comments
- [GitHub Actions Checks API](https://docs.github.com/en/rest/checks) - Status checks
- [Codecov GitHub Action](https://github.com/codecov/codecov-action) - Existing integration
- [GitHub Actions Performance Metrics](https://docs.github.com/en/actions/concepts/metrics) - Runtime data for cost estimation
- [New Billing Platform Usage API](https://docs.github.com/en/billing/using-the-new-billing-platform/automating-usage-reporting) - Actual billing data (24-48hr delay, optional)
