import os
import platform
from dataclasses import dataclass
from pathlib import Path
from typing import Dict, List, Optional

from .config import resolve_config_path


@dataclass
class BuildInfo:
    name: str = "ado"
    version: str = "0.0.0-dev"
    commit: str = "none"
    build_time: str = "unknown"
    python_version: str = platform.python_version()
    platform: str = f"{platform.system().lower()}/{platform.machine().lower()}"


def collect_env_info(explicit_config: Optional[str]) -> Dict[str, object]:
    home_dir = Path.home()
    env_config = os.getenv("ADO_CONFIG")
    resolved, sources = resolve_config_path(explicit_config, env_config, home=home_dir)

    env_vars: Dict[str, str] = {}
    for key in ("ADO_CONFIG", "ADO_LOG_LEVEL"):
        val = os.getenv(key)
        if val is not None:
            env_vars[key] = val

    return {
        "config_path": str(resolved) if resolved else None,
        "config_sources": [str(p) for p in sources],
        "home_dir": str(home_dir),
        "cache_dir": str(Path(os.path.expanduser("~")).joinpath("Library", "Caches"))
        if platform.system().lower() == "darwin"
        else str(Path.home() / ".cache"),
        "env": env_vars,
    }


def collect_features() -> List[str]:
    return []
