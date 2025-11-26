---
name: Command Proposal
about: Propose a new CLI command for ado
title: 'feat(cmd): '
labels: enhancement, command, needs-triage
assignees: ''
---

## Command Summary

**Proposed command:**

```bash
ado [command] [subcommand] [flags] [args]
```

**One-sentence description:**

## Problem Statement

**What user problem does this command solve?**

- What can't users do today?
- What workaround are they using?

## Proposed Behavior

### Basic Usage

```bash
# Example 1: Most common use case
ado command arg
# Expected output:
# ...
```

### With Flags

```bash
# Example 2: With options
ado command --flag value arg
# Expected output:
# ...
```

### Structured Output

```bash
# Example 3: JSON output for scripting
ado command --output json arg
# Expected output:
# {"field": "value"}
```

## Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--flag` | `-f` | string | `""` | Description |
| `--output` | `-o` | enum | `text` | Output format: text, json, yaml |

## Error Cases

| Condition | Expected Behavior |
|-----------|------------------|
| No arguments | Error: "at least one argument required" |
| Invalid flag | Error: "invalid value for --flag" |

## Scope Assessment

**Does this require an ADR first?**

- [ ] **Yes** → This introduces new command patterns or architectural changes
- [ ] **No** → This follows existing command patterns (like `echo`, `meta`)
- [ ] **Unsure** → Need maintainer input

**Does this need Python prototyping first?**

- [ ] **Yes** → Complex logic that benefits from rapid iteration in `lab/py/`
- [ ] **No** → Straightforward implementation

## Related Commands

How does this relate to existing commands?

- `ado existing-command` - Relationship/difference

## Acceptance Criteria

- [ ] Command spec approved in `docs/commands/NN-name.md`
- [ ] Implementation follows spec exactly
- [ ] All examples work as documented
- [ ] All error cases behave as specified
- [ ] Tests pass with >80% coverage
- [ ] Help text (`--help`) is accurate

## Checklist

- [ ] I have searched existing issues and `docs/commands/` for similar commands
- [ ] I have read the [workflow guide](../../docs/workflow.md)
- [ ] I have determined whether an ADR is needed (see above)
- [ ] I am willing to write the spec if this proposal is accepted
