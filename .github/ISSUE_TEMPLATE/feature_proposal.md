---
name: Feature Proposal
about: Propose a new internal feature for ado (non-command features)
title: 'feat: '
labels: enhancement, needs-triage
assignees: ''
---

## Feature Summary

**One-sentence description of the feature.**

## Problem Statement

**What problem does this solve?**

- What pain point are we addressing?
- Who is affected by this problem?

## Proposed Solution

**How should this work?**

High-level description of the solution.

## Use Cases

**Who benefits and how?**

1. As a [user type], I want to [action] so that [benefit].
2. ...

## Scope Assessment

**What type of change is this?**

- [ ] **New CLI command** → Use [Command Proposal](command_proposal.md) template instead
- [ ] **New internal feature** → Config, logging, plugins, etc.
- [ ] **Enhancement to existing feature**
- [ ] **Developer tooling** → Build, test, CI improvements

**Does this require an ADR first?**

An ADR is required if this:
- [ ] Introduces a new architectural pattern
- [ ] Changes existing architecture
- [ ] Adds new dependencies or tools
- [ ] Affects multiple components
- [ ] Is security-related

If any are checked, create an [ADR Proposal](adr_proposal.md) first.

## Example Usage

**How would this feature be used?**

```yaml
# Example config
new_feature:
  enabled: true
  option: value
```

Or for code:

```go
// Example API usage
feature.DoSomething(ctx, options)
```

## Acceptance Criteria

- [ ] Criterion 1
- [ ] Criterion 2
- [ ] Tests added with >80% coverage
- [ ] Documentation updated

## Alternatives Considered

**Other approaches you've considered:**

1. Alternative 1: ...
2. Alternative 2: ...

## Additional Context

Screenshots, mockups, or additional details.

## Checklist

- [ ] I have searched existing issues for duplicates
- [ ] I have read the [workflow guide](../../docs/workflow.md)
- [ ] I have determined whether an ADR is needed (see above)
- [ ] I am willing to write the spec if this proposal is accepted
