import os
from pathlib import Path
from typing import List, Optional, Tuple


def default_search_paths(home: Optional[Path] = None) -> List[Path]:
    home_dir = Path(home or Path.home())
    paths: List[Path] = []
    xdg = os.getenv("XDG_CONFIG_HOME")

    if xdg:
        paths.append(Path(xdg) / "ado" / "config.yaml")
    else:
        paths.append(home_dir / ".config" / "ado" / "config.yaml")

    paths.append(home_dir / ".ado" / "config.yaml")
    return paths


def resolve_config_path(
    explicit_path: Optional[str],
    env_path: Optional[str],
    home: Optional[Path] = None,
) -> Tuple[Optional[Path], List[Path]]:
    search_paths = default_search_paths(home)

    if explicit_path:
        return Path(explicit_path), [Path(explicit_path), *search_paths]

    if env_path:
        return Path(env_path), [Path(env_path), *search_paths]

    for candidate in search_paths:
        if candidate.exists():
            return candidate, search_paths

    return None, search_paths
