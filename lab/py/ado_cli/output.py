import json
from typing import Any, Callable, Iterable

import yaml

OutputFormat = str

OUTPUT_CHOICES: Iterable[str] = ("text", "json", "yaml")


def format_output(format: OutputFormat, payload: Any, render_text: Callable[[], str]) -> str:
    if format == "json":
        return json.dumps(payload, indent=2)
    if format == "yaml":
        return yaml.safe_dump(payload, sort_keys=False)
    return render_text()
