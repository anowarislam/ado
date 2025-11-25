# Documentation Style Guide

Guidelines for writing and maintaining project documentation: command specs, README, and design docs.

## Philosophy
- Documentation is the specification; code implements the spec.
- Examples beat prose; show concrete CLI invocations with expected output.
- Keep docs synchronized with code; outdated docs are worse than no docs.
- Write for future maintainers; assume no prior context.

## Command Specifications (`docs/commands/*.md`)

### Structure

Every command spec follows this template:

```markdown
# Command Name: <command>

Brief one-line description.

## Purpose

Why this command exists; what user problem it solves.

## Usage

```bash
ado <command> [subcommand] [flags] [args]
```

## Flags

### Global Flags (inherited)
- `--config PATH` - Config file path
- `--log-level LEVEL` - Log level
- `--help, -h` - Show help

### Command-Specific Flags
- `--flag-name TYPE` - Description (default: value)

## Examples

### Example 1: Basic usage
```bash
$ ado command arg
Expected output here
```

### Example 2: With flags
```bash
$ ado command --flag value
Expected output here
```

## Edge Cases

### Case 1: Missing required argument
```bash
$ ado command
Error: required argument missing
```

### Case 2: Invalid input
```bash
$ ado command --invalid-flag
Error: unknown flag: --invalid-flag
```

## Implementation

See `cmd/ado/<command>/<command>.go` for implementation.
See `internal/<package>` for shared logic.

## Tests

See `cmd/ado/<command>/<command>_test.go` or `internal/<package>/<package>_test.go`.

## Related Commands

- `ado other-command` - Related functionality
```

### Naming Convention

- Filename: `docs/commands/<NN>-<command>.md` (e.g., `01-echo.md`, `02-help.md`, `03-meta.md`)
- Numbering: Use sequential prefixes for ordering; gaps are okay
- Command name in filename must match actual CLI command

### Content Guidelines

1. **Be specific**: "Show version info" beats "Display information"
2. **Show, don't tell**: Include actual CLI examples with output
3. **Document errors**: Show what happens when things go wrong
4. **Link to code**: Reference implementation files with line numbers if helpful
5. **Keep current**: Update spec when behavior changes; spec is source of truth

### Examples Section

- Start with simplest case (no flags, default behavior)
- Show progressively complex examples
- Include realistic use cases, not toy examples
- Show both success and failure cases
- Format: `$ ado command` followed by expected output

### Edge Cases Section

Document:
- Missing required inputs
- Invalid flag values
- Conflicting flags
- Empty/malformed config files
- Permission errors
- Network failures (if applicable)

## README Maintenance

### Structure

The README should mirror this order:

1. **Project name and tagline** - One sentence: what is this?
2. **Goals** - Why does this exist? What problems does it solve?
3. **Architecture** - High-level: Go binary + Python lab model
4. **Design principles** - Core tenets (see existing README)
5. **Initial commands** - Brief list with examples
6. **Getting started** - How to use (developer perspective)
7. **Roadmap** - Future plans (optional, keep realistic)

### Keep In Sync

When you change behavior:
- Update command examples in README if they're featured
- Update architecture section if structure changes
- Update "Current Status" if milestones are reached
- Ensure `ado meta info` output matches README claims

### Writing Style

- Active voice: "The binary uses Cobra" not "Cobra is used by the binary"
- Present tense: "The CLI provides" not "The CLI will provide"
- Be concise: Remove words that don't add meaning
- Use code blocks for CLI commands and file paths
- Link to detailed docs: "See `docs/architecture.md` for details"

## Markdown Conventions

### Formatting

- **Headers**: `#` for title, `##` for major sections, `###` for subsections
- **Code**: Use `` `backticks` `` for inline code, file paths, commands, flags
- **Blocks**: Use triple backticks with language: `` ```bash ``, `` ```go ``, `` ```yaml ``
- **Emphasis**: Use **bold** for emphasis, *italic* sparingly (prefer bold)
- **Lists**: `-` for unordered, `1.` for ordered; use consistent indentation

### Line Length

- Aim for 80-100 characters for readability
- Break at sentence boundaries, not mid-phrase
- Code blocks can exceed this (don't break commands)
- Tables can exceed this (use judgment)

### Links

- Use relative paths for internal docs: `[structure](../architecture.md)`
- Use descriptive link text: "See the [command spec guide](commands/README.md)" not "click here"
- Link to code with line numbers when helpful: `cmd/ado/root/root.go:14-38`

### Code Blocks

**Bash/Shell commands** - Show prompt and output:

```bash
$ ado meta info
Version: v0.1.0
Commit: abc123
Build Time: 2024-01-15T10:00:00Z
```

**Go code** - Use syntax highlighting:

```go
func main() {
    fmt.Println("Hello")
}
```

**YAML/Config** - Show complete valid examples:

```yaml
version: 1
commands:
  echo:
    enabled: true
```

### Tables

Use for structured comparisons:

| Command | Purpose | Example |
|---------|---------|---------|
| `meta info` | Show version | `ado meta info` |
| `meta env` | Show config paths | `ado meta env` |

Align pipes for readability in source; rendering will normalize.

## Design Documents (`docs/*.md`)

### Purpose

Design docs capture architectural decisions and specifications before implementation:
- `docs/architecture.md` - Directory layout and organization
- `docs/commands-overview.md` - Global CLI conventions and behavior
- `docs/style/*.md` - Code and workflow style guides

### Structure

1. **Goal** - What is being specified?
2. **Context** - Why is this needed? What problem does it solve?
3. **Specification** - Detailed technical spec (structure, behavior, contracts)
4. **Examples** - Concrete examples (code, commands, configs)
5. **Alternatives considered** (optional) - What else was considered and why not?
6. **Open questions** (optional) - What's still undecided?

### When to Create

- Before implementing a new subsystem (config, logging, plugins)
- When establishing a convention (CLI patterns, testing style)
- For complex features requiring upfront design
- To document architectural decisions

### Keep Updated

- Mark sections as `[IMPLEMENTED]` or `[PROPOSED]` when helpful
- Update specs when implementation deviates (spec stays source of truth)
- Link from code back to design docs in comments

## Synchronization Checklist

When making changes, update docs in this order:

### Adding a New Command

1. ✅ Create `docs/commands/<NN>-<command>.md` with full spec
2. ✅ Update `README.md` "Initial Commands" section if it's user-facing
3. ✅ Update `CLAUDE.md` "Current Command Set" if implementation is complete
4. ✅ Implement in `cmd/ado/<command>/<command>.go`
5. ✅ Add tests; ensure examples in spec work
6. ✅ Update `AGENTS.md` if workflow changes (new Make targets, etc.)

### Changing Command Behavior

1. ✅ Update `docs/commands/<command>.md` spec first (examples, edge cases)
2. ✅ Update code to match new spec
3. ✅ Update/add tests to match new spec
4. ✅ Update `README.md` if examples are featured there
5. ✅ Add changelog entry or update release notes

### Changing Architecture

1. ✅ Update `docs/architecture.md` or create new design doc
2. ✅ Update `README.md` "Architecture" section if user-visible
3. ✅ Update `CLAUDE.md` if structure changes affect future development
4. ✅ Update code
5. ✅ Update `AGENTS.md` if workflow changes

### Release Preparation

1. ✅ Ensure all command specs match implementation
2. ✅ Update `README.md` "Current Status" and roadmap
3. ✅ Update `CHANGELOG.md` with user-visible changes
4. ✅ Update version in `internal/meta/info.go` (or rely on ldflags)
5. ✅ Tag release with `vX.Y.Z`

## Real Examples from This Codebase

### Command Spec
See `docs/commands/03-meta.md` for `ado meta` command specification with examples.

### Design Doc
See `docs/architecture.md` for directory layout specification.

### Style Guide
See `docs/style/go-style.md` for code style documentation format.

### README
See `README.md` for top-level architecture and goals documentation.

## Bad Practices to Avoid

### In Command Specs
- ❌ Vague descriptions: "This command does stuff" → ✅ "Show build version and commit hash"
- ❌ Missing examples → ✅ Always include at least 2-3 concrete examples
- ❌ No error cases → ✅ Document what happens when things fail
- ❌ Outdated flags → ✅ Keep flags synchronized with implementation

### In README
- ❌ Overpromising: "Will support plugins" → ✅ "Designed for future plugin support"
- ❌ Stale examples → ✅ Test examples actually work
- ❌ Too much detail → ✅ Link to detailed docs; README is overview
- ❌ No getting started → ✅ Show how to build/run/test

### In Design Docs
- ❌ Implementation details only → ✅ Explain the "why" first
- ❌ Never updating after implementation → ✅ Mark sections as [IMPLEMENTED] or update spec
- ❌ No alternatives discussed → ✅ Show what was considered and why not
- ❌ Burying decisions in prose → ✅ Use tables, lists, clear structure

### In Markdown
- ❌ No syntax highlighting on code blocks → ✅ Always specify language: `` ```go ``
- ❌ Broken relative links → ✅ Test links work from docs directory
- ❌ "Click here" link text → ✅ Descriptive text: "See the [style guide](README.md)"
- ❌ Walls of text → ✅ Break into sections with headers, lists, examples

## Related Guides
- See [go-style.md](go-style.md) for code documentation (package comments, godoc).
- See [python-style.md](python-style.md) for Python docstrings and type hints.
- See [ci-style.md](ci-style.md) for changelog and release note practices.
- See [README.md](README.md) for style guide index and cross-references.
