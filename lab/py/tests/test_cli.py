import json
import sys
from pathlib import Path

from click.testing import CliRunner

# Ensure lab package is importable when running from repo root
LAB_ROOT = Path(__file__).resolve().parents[1]
if str(LAB_ROOT) not in sys.path:
    sys.path.insert(0, str(LAB_ROOT))

from ado_cli.cli import cli  # noqa: E402


def test_echo_upper():
    runner = CliRunner()
    result = runner.invoke(cli, ["echo", "--upper", "hello", "world"])
    assert result.exit_code == 0
    assert result.output.strip() == "HELLO WORLD"


def test_echo_json_repeat():
    runner = CliRunner()
    result = runner.invoke(
        cli, ["echo", "--repeat", "2", "--output", "json", "hi"], catch_exceptions=False
    )
    assert result.exit_code == 0
    data = json.loads(result.output)
    assert data == ["hi", "hi"]


def test_meta_info_json():
    runner = CliRunner()
    result = runner.invoke(cli, ["meta", "info", "--output", "json"])
    assert result.exit_code == 0
    payload = json.loads(result.output)
    assert payload["name"] == "ado"
    assert "platform" in payload


def test_meta_env_prefers_ado_config(tmp_path: Path):
    env_config = tmp_path / "env-config.yaml"
    env_config.write_text("content")

    runner = CliRunner()
    env = {"ADO_CONFIG": str(env_config)}
    result = runner.invoke(cli, ["meta", "env", "--output", "json"], env=env)
    assert result.exit_code == 0

    payload = json.loads(result.output)
    assert payload["config_path"] == str(env_config)
    assert payload["config_sources"][0] == str(env_config)
    assert payload["env"]["ADO_CONFIG"] == str(env_config)


def test_help_command_shows_echo_help():
    runner = CliRunner()
    result = runner.invoke(cli, ["help", "echo"])
    assert result.exit_code == 0
    assert "Echo input text with optional transformations." in result.output
