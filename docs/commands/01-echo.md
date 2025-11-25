# echo Command Spec

## Command:

```bash
ado echo
```

### Purpose:

A simple, deterministic command for validating CLI parsing, piping, and formatting conventions. Serves as a baseline example for developers adding new commands.

### Usage examples:

    1. ado echo hello world
    2. ado echo --upper hello world
    3. ado echo --output yaml hello world
    4. ado echo --repeat 3 hello

### Arguments:

Positional message...: one or more tokens (strings) to be echoed.

### Flags:

- --upper:
    - Description: convert the joined message to uppercase before output.
    - Default: false.

- --lower:
    - Description: convert the joined message to lowercase before output.
    - Default: false.
    - If both --upper and --lower are provided, this is considered invalid; the command should fail with a clear error.

- --repeat N:
    - Description: number of times to repeat the message.
    - Default: 1.
    - Constraints: N >= 1. Out-of-range values cause a validation error.

- --output FORMAT, -o FORMAT:
    - Allowed values: text, json, yaml
    - Default: text

### Behavior:

- The message is formed by joining positional arguments with a single space.
- Transformations (--upper, --lower) are applied to the entire joined string.
- If --repeat N > 1:
	- In text mode: print the transformed message N times, each on its own line.
	- In json/yaml mode: print an array of strings, length N.

### Error cases:

- No arguments:
  - Behavior: exit code >0, and print an error explaining that at least one word is required.

- Conflicting flags (--upper and --lower):  
  - Behavior: exit code >0, descriptive error.

- Invalid repeat (--repeat < 1):
  - Behavior: exit code >0, descriptive error.

- Invalid output format (not in text|json|yaml):
  - Behavior: exit code >0, descriptive error.
