# Feature Specifications

This directory contains specifications for non-command features in `ado`.

## What is a Feature Spec?

A feature spec documents internal features that are not CLI commands, such as:

- Configuration system enhancements
- Logging infrastructure
- Plugin architecture
- Build/release tooling
- Testing frameworks
- Internal libraries

For CLI command specifications, see [`docs/commands/`](../commands/).

## Feature Index

| ID | Title | Status | ADR |
|----|-------|--------|-----|
| [01](01-structured-logging.md) | Structured Logging | Implemented | [ADR-0002](../adr/0002-structured-logging.md) |
| [02](02-homebrew-distribution.md) | Homebrew Distribution | Draft | N/A |
| [03](03-pr-metrics-dashboard.md) | PR-Level Metrics Dashboard | Draft | [ADR-0005](../adr/0005-pr-metrics-dashboard.md) |
| [04](04-claude-review-optimization.md) | Claude Code Review Optimization | Draft | N/A |

## Creating a Feature Spec

1. Determine if an ADR is required (see [workflow guide](../workflow.md))
2. Create a branch: `spec/feature-name`
3. Copy `TEMPLATE.md` to `NN-feature-name.md`
4. Fill in all sections
5. Submit PR with title: `docs(spec): feature name`
6. After approval, proceed to implementation

## Spec Lifecycle

```
Draft → Approved → Implemented → [Deprecated]
```

| Status | Meaning |
|--------|---------|
| **Draft** | Under discussion, not yet approved |
| **Approved** | Ready for implementation |
| **Implemented** | Implementation complete |
| **Deprecated** | No longer relevant |

## Numbering

Feature specs are numbered sequentially (01, 02, ...). Use the next available number.

## Related Documentation

- [Command Specs](../commands/) - CLI command specifications
- [ADRs](../adr/) - Architecture decisions
- [Workflow Guide](../workflow.md) - Complete development workflow
