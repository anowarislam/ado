# Feature: [Name]

<!--
Instructions:
1. Copy this file to `NN-feature-name.md` (use next available number)
2. Fill in all sections below
3. Submit PR with title: `docs(spec): feature name`
4. Update status after approval/implementation
-->

| Metadata | Value |
|----------|-------|
| **ADR** | ADR-NNNN or N/A |
| **Status** | Draft |
| **Issue** | #NNN |
| **Author(s)** | @username |

## Overview

Brief description of the feature (1-2 sentences).

## Motivation

Why is this feature needed? What problem does it solve?

- What pain point are we addressing?
- Who benefits from this feature?
- What happens if we don't implement this?

## Specification

### Behavior

Detailed description of how the feature works.

1. Step 1
2. Step 2
3. Step 3

### Configuration

If applicable, configuration options:

```yaml
# Example config in ~/.config/ado/config.yaml
feature_name:
  enabled: true
  option: value
```

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `enabled` | bool | `false` | Enable this feature |
| `option` | string | `""` | Description of option |

### API/Interface

If applicable, public interfaces:

```go
// Package feature provides...
package feature

// FeatureName does...
type FeatureName interface {
    // Method does...
    Method(arg Type) (Result, error)
}
```

### File Locations

| Purpose | Path |
|---------|------|
| Main implementation | `internal/feature/` |
| Tests | `internal/feature/*_test.go` |
| Config integration | `internal/config/feature.go` |

## Examples

### Example 1: Basic Usage

```bash
# Setup
...

# Usage
...

# Expected output
...
```

### Example 2: Advanced Usage

```bash
# More complex scenario
...
```

## Edge Cases and Error Handling

| Scenario | Expected Behavior |
|----------|------------------|
| Missing config | Default to X |
| Invalid input | Return error: "message" |
| Network failure | Retry N times, then fail with error |
| Concurrent access | Thread-safe via mutex |

## Testing Strategy

### Unit Tests

- [ ] Test case 1: description
- [ ] Test case 2: description
- [ ] Test case 3: description

### Integration Tests

- [ ] End-to-end scenario 1
- [ ] End-to-end scenario 2

### Manual Testing Checklist

- [ ] Step 1: Do X, verify Y
- [ ] Step 2: Do A, verify B

## Implementation Checklist

- [ ] Create `internal/feature/` package
- [ ] Implement core functionality
- [ ] Add configuration support (if applicable)
- [ ] Write unit tests (target: 80% coverage)
- [ ] Write integration tests
- [ ] Update CLAUDE.md if architecture changes
- [ ] Add to README if user-facing

## Open Questions

Questions to resolve before/during implementation:

- [ ] Question 1?
- [ ] Question 2?

## Changelog

| Date | Change | Author |
|------|--------|--------|
| YYYY-MM-DD | Initial draft | @username |
