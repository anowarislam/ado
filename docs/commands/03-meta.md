# meta Command Spec

## Command:

```
ado meta
```

### Purpose:

Introspect the ado binary and its environment. This is a diagnostic/debugging command and a reference for how to structure more advanced commands.

### Subcommands:

	1.	ado meta info
	2.	ado meta env
	3.	ado meta features (optional for v0 but specced now)

## ado meta info

### Usage:

	1. ado meta info
	2. ado meta info --output json
	3. ado meta info --output yaml

### Description:
Prints metadata about the current ado binary and its build.

Fields (human-readable mode):

	- Name: ado
	- Version: semantic version (e.g. 0.1.0)
	- Commit: short git SHA, if available
	- BuildTime: ISO8601 string or similar
	- GoVersion: Go compiler version (if available)
	- Platform: os/arch (e.g. darwin/arm64, linux/amd64)

JSON mode:

	- If --output json is provided, output a single JSON object with the fields above, using stable keys:
	- name, version, commit, build_time, go_version, platform

YAML mode:

	- If --output yaml is provided, output a single YAML object with the fields above, using stable keys:
	- name, version, commit, build_time, go_version, platform

Flags:
- --output, -o: text (default), json, yaml

## ado meta env
### Usage:

	1. ado meta env
	2. ado meta env --output json
	3. ado meta env --output yaml

### Description:

Shows relevant environment information that affects ado behavior.

Fields (candidate list, can be narrowed for v0):

	- ConfigPath: resolved config path in use (if any).
	- ConfigSources: the set of locations checked (respects --config when set, otherwise XDG_CONFIG_HOME or HOME).
	- HomeDir: resolved home directory path.
	- CacheDir: resolved cache directory path if used.
	- EnvVariables: selected env variables relevant to ado (currently ADO_CONFIG, ADO_LOG_LEVEL when set).

Behavior:

	- In human-readable mode, present as a sectioned text report.
	- In JSON mode, produce an object containing:
	- config_path, config_sources, home_dir, cache_dir, env.
	- In YAML mode, emit the same keys as YAML.

Flags:
- --output, -o: text (default), json, yaml

## ado meta features (future-friendly)

### Usage:

	1. ado meta features
	2. ado meta features --output json

### Description:
Lists enabled/disabled features or build-time flags (e.g., experimental flags, compiled-in integrations). For v0, returns “No experimental features enabled” in text mode and an empty array in structured modes.

Flags:
- --output, -o: text (default), json, yaml
