# Architecture Decision Records

This directory contains Architecture Decision Records (ADRs) for the `ado` project.

## What is an ADR?

An ADR is a document that captures an important architectural decision along with its context and consequences. ADRs help us:

- Understand **why** decisions were made
- Onboard new contributors with historical context
- Avoid revisiting decisions without new information
- Document trade-offs and alternatives considered

## ADR Index

| ID | Title | Status | Date |
|----|-------|--------|------|
| [0001](0001-development-workflow.md) | ADR + Spec Development Workflow | Accepted | 2025-11-25 |
| [0002](0002-structured-logging.md) | Structured Logging | Accepted | 2025-11-26 |
| [0003](0003-recipe-based-documentation.md) | Recipe-Based Documentation for CI/CD Patterns | Accepted | 2025-11-26 |

## ADR Lifecycle

```
Proposed → Accepted → [Deprecated | Superseded]
```

| Status | Meaning |
|--------|---------|
| **Proposed** | Under discussion, not yet approved |
| **Accepted** | Approved and in effect |
| **Deprecated** | No longer relevant (context changed) |
| **Superseded** | Replaced by a newer ADR |

## When to Write an ADR

**Required for:**

- New architectural patterns
- Breaking changes to existing patterns
- Adding new dependencies or tools
- Security-related decisions
- Changes affecting multiple components

**Not required for:**

- Bug fixes
- Single command additions (use [command spec](../commands/) instead)
- Documentation improvements
- Minor refactoring

## Decision Tree

```
Is this a new architectural pattern?           → YES → ADR required
Does this change existing architecture?        → YES → ADR required
Does this add new dependencies/tools?          → YES → ADR required
Does this affect multiple components?          → YES → ADR required
Is this security-related?                      → YES → ADR required
Is this a single command following patterns?   → NO  → Use command spec
Is this a bug fix or documentation?            → NO  → Direct to implementation
```

## Creating an ADR

1. Check if an ADR already exists for your topic
2. Create an Issue using the "ADR Proposal" template
3. Once discussion concludes, create a branch: `adr/NNNN-short-title`
4. Copy `TEMPLATE.md` to `NNNN-title.md`
5. Fill in all sections
6. Submit PR with title: `docs(adr): NNNN - title`
7. After approval, update status to "Accepted"

## Numbering

ADRs are numbered sequentially (0001, 0002, ...). Numbers are never reused even if an ADR is deprecated or superseded.

## Related Documentation

- [Development Workflow](../workflow.md) - Complete workflow guide
- [Command Specs](../commands/) - CLI command specifications
- [Feature Specs](../features/) - Non-command feature specifications
- [Contributing Guide](../contributing.md) - Contribution guidelines

## References

- [ADR GitHub Organization](https://adr.github.io/)
- [Documenting Architecture Decisions](https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions) - Michael Nygard
