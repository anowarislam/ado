# Python Code & Test Tenets

Guidelines for Python lab code (`lab/py`) and its tests. Keep it easy to port to Go and pleasant to maintain.

## Philosophy
- Clarity and explicitness first; favor readable, small functions.
- Minimize cleverness and hidden state; behavior should be obvious from arguments and returns.
- Lab code is disposable but should still be correct and deterministic to ease promotion to Go.
- Type hints where practical; they document intent and catch mistakes early.

## Code Tenets
- Structure
  - Keep prototypes scoped under `lab/py/<area>`; avoid shared mutable globals.
  - Isolate I/O: separate parsing/CLI from pure logic so tests stay simple.
- Naming & style
  - Follow PEP 8: snake_case for functions/vars, CapWords for classes.
  - Docstrings for public functions/classes; inline comments only for non-obvious intent.
  - Use dataclasses for simple data containers instead of ad-hoc dicts.
- Error handling
  - Raise specific exceptions; avoid broad `except Exception` without re-raising context.
  - Fail fast on invalid inputs; keep messages actionable.
- Dependencies
  - Keep deps minimal and pinned in `lab/py/requirements.txt`.
  - Prefer standard library; choose well-maintained libs (e.g., click/typer) when needed.
- I/O & state
  - Avoid implicit globals; pass configuration explicitly.
  - For CLI: parse in the interface layer, keep core functions pure where possible.
- Modern Python patterns
  - Use f-strings for formatting: `f"config at {path}"` over `.format()` or `%` formatting.
  - Prefer `pathlib.Path` over `os.path` for file operations; it's more expressive and chainable.
  - Use `logging` module for diagnostics; reserve `print()` for CLI output only.
  - Type hints: Use `typing` for complex hints like `list[dict[str, Any]]`, `Optional[T]`.
- Type checking
  - Run `mypy lab/py` or `pyright` for static type checking; fix errors before porting to Go.
  - Type hints serve as documentation and catch errors early; use them liberally.
  - Avoid `Any` unless truly dynamic; prefer specific types or unions.
- Import organization
  - Group imports: stdlib, third-party, local (separated by blank lines).
  - Sort alphabetically within each group; use `isort` or `ruff` auto-formatter.
  - Avoid star imports (`from module import *`); be explicit about what you use.

## Test Tenets
- Determinism
  - No time/rand without control; seed/mocks as needed.
  - Avoid external network/files unless using temp dirs/fixtures; keep tests hermetic.
- Structure
  - Use pytest; parametrize for scenario coverage; mark helpers with `@pytest.fixture` when shared setup is needed.
  - Prefer pure function tests; for CLI, use `click.testing.CliRunner` or `subprocess` with temp dirs.
- Assertions
  - Assert specific outputs/structures; include helpful failure messages.
  - One logical expectation per test; use subtests/parametrize to group variants.
- Maintenance
  - Keep tests fast; skip or mark slow/external cases with rationale.
  - Add regression tests alongside bug fixes; fixtures should stay small and readable.
- Pytest plugins
  - `pytest-cov`: Coverage reporting with `pytest --cov=lab/py`.
  - `pytest-xdist`: Parallel test execution with `pytest -n auto` for faster runs.
  - Keep plugin deps in `lab/py/requirements.txt` with versions pinned.

## Examples
- Good
  - Clear error propagation with context:
    ```python
    def load_config(path: Path) -> dict:
        try:
            return yaml.safe_load(path.read_text())
        except FileNotFoundError as exc:
            raise ValueError(f"config not found: {path}") from exc
    ```
  - Pure function with type hints and no globals:
    ```python
    def normalize(text: str, lower: bool = True) -> str:
        text = text.strip()
        return text.lower() if lower else text
    ```
  - Pytest parametrize:
    ```python
    @pytest.mark.parametrize(
        "raw, lower, expected",
        [(" hi ", True, "hi"), ("Hi", False, "Hi")],
    )
    def test_normalize(raw, lower, expected):
        assert normalize(raw, lower) == expected
    ```
  - CLI tested via CliRunner:
    ```python
    def test_echo_cli_upper():
        runner = CliRunner()
        result = runner.invoke(cli, ["echo", "--upper", "hi"])
        assert result.exit_code == 0
        assert result.output.strip() == "HI"
    ```
  - Logging with context:
    ```python
    import logging

    logger = logging.getLogger(__name__)

    def process(path: Path) -> dict:
        logger.debug(f"processing {path}")
        try:
            data = parse(path)
            logger.info(f"parsed {len(data)} entries from {path}")
            return data
        except Exception as exc:
            logger.error(f"failed to process {path}: {exc}")
            raise
    ```
  - Using pathlib.Path:
    ```python
    from pathlib import Path

    def find_configs(root: Path) -> list[Path]:
        return list(root.glob("**/*.yaml"))

    config_path = Path.home() / ".config" / "ado" / "config.yaml"
    if config_path.exists():
        data = config_path.read_text()
    ```
  - Import organization:
    ```python
    # Standard library
    import logging
    from pathlib import Path
    from typing import Optional

    # Third-party
    import click
    import yaml

    # Local
    from ado_cli.config import load_config
    from ado_cli.output import format_table
    ```
- Bad
  - Silent failure and broad except:
    ```python
    try:
        data = yaml.safe_load(path.read_text())
    except Exception:
        data = {}
    ```
  - Hidden global state:
    ```python
    cache = {}

    def add_item(key, value):
        cache[key] = value  # implicit mutation
    ```
  - Mutable default arguments:
    ```python
    def append_item(item, items=[]):
        items.append(item)
        return items
    ```
  - Flaky time-based test:
    ```python
    def test_timeout():
        time.sleep(0.1)
        assert done
    ```
  - Late binding in closure (common gotcha):
    ```python
    funcs = [lambda: i for i in range(3)]
    results = [f() for f in funcs]  # Returns [2, 2, 2] not [0, 1, 2]

    # Fix: capture value explicitly
    funcs = [lambda i=i: i for i in range(3)]
    ```
  - Using os.path instead of pathlib:
    ```python
    import os

    # Bad
    config_dir = os.path.join(os.path.expanduser("~"), ".config", "ado")
    if os.path.exists(config_dir):
        files = os.listdir(config_dir)

    # Good
    from pathlib import Path
    config_dir = Path.home() / ".config" / "ado"
    if config_dir.exists():
        files = list(config_dir.iterdir())
    ```

## Real Examples from This Codebase

### CLI structure
See `lab/py/ado_cli/cli.py` for Click-based CLI with subcommands matching Go implementation.

### Type-hinted pure functions
See `lab/py/ado_cli/config.py` for config path resolution logic (mirrors `internal/config/paths.go`).

### Output formatting
See `lab/py/ado_cli/output.py` for structured output helpers (JSON/YAML/text).

### Pytest tests
See `lab/py/tests/test_cli.py` for CLI testing with `CliRunner` and parametrize examples.

### Metadata helpers
See `lab/py/ado_cli/meta.py` for build info prototyping (mirrors `internal/meta/info.go`).

## Promotion to Go

When lab code is ready to port, map Python patterns to Go:

| Python Pattern | Go Equivalent | Notes |
|----------------|---------------|-------|
| `raise ValueError("msg")` | `return fmt.Errorf("msg")` | Exceptions → errors |
| `@dataclass` | `type MyStruct struct{}` | Dataclasses → structs |
| `def func(x: str) -> int:` | `func MyFunc(x string) int` | Type hints → Go types |
| `if __name__ == "__main__":` | `func main() {}` | Entry point |
| `@pytest.mark.parametrize` | Table-driven `t.Run()` | See [go-style.md](go-style.md) |
| `logger.info("msg")` | `slog.Info("msg")` | Logging patterns similar |
| `Path.glob("*.yaml")` | `filepath.Glob("*.yaml")` | File pattern matching |
| `yaml.safe_load(text)` | `yaml.Unmarshal(data, &cfg)` | YAML parsing |

### Promotion checklist:
1. Lock behavior with comprehensive pytest tests and fixtures
2. Document edge cases and error conditions in Python code
3. Implement in Go following [go-style.md](go-style.md) patterns
4. Port tests as table-driven Go tests
5. Verify behavior matches Python prototype with same test fixtures
6. Keep Python code as reference documentation

## Related Guides
- See [go-style.md](go-style.md) for target patterns when porting lab code to production.
- See [ci-style.md](ci-style.md) for `make py.test` and `make py.lint` CI integration.
- See [docs-style.md](docs-style.md) for documenting prototypes and promotion decisions.
- See [README.md](README.md) for style guide index and when to consult each guide.
