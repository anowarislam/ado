# help Command Spec

Note: The underlying CLI framework will typically provide help behavior (--help and help subcommand). This spec clarifies required customization and content.

### Purpose:

A simple command for help about any command

## Command:

```bash
ado help <command> [subcommand]
```

### Behavior:

For any command X, ado help X should:

- Show:
	- Usage line with flags.
	- Short and long descriptions.
- Examples:
	- At least 2–3 realistic examples.
	- Document:
		- Exit behavior where non-trivial.
		- Output formats if multiple formats are supported.

### Examples:

- ado help
- ado help meta
- ado help meta info
- ado help echo