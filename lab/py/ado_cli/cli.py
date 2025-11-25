import sys
from typing import Optional

import click

from . import meta
from .output import OUTPUT_CHOICES, format_output


@click.group(
    context_settings={"help_option_names": ["-h", "--help"]},
    invoke_without_command=True,
)
@click.option("--config", type=click.Path(), default=None, help="Path to config file")
@click.option(
    "--log-level",
    type=click.Choice(["debug", "info", "warn", "error"], case_sensitive=False),
    default="info",
    show_default=True,
    help="Log level for output",
)
@click.pass_context
def cli(ctx: click.Context, config: Optional[str], log_level: str):
    """ado lab CLI (Python prototype)."""
    ctx.ensure_object(dict)
    ctx.obj["config"] = config
    ctx.obj["log_level"] = log_level
    if ctx.invoked_subcommand is None:
        click.echo(ctx.get_help())


@cli.command()
@click.argument("message", nargs=-1, required=True)
@click.option("--upper", is_flag=True, default=False, help="Convert message to uppercase")
@click.option("--lower", is_flag=True, default=False, help="Convert message to lowercase")
@click.option(
    "--repeat",
    default=1,
    show_default=True,
    type=click.IntRange(1),
    help="Number of times to repeat the message",
)
@click.option(
    "--output",
    "-o",
    type=click.Choice(OUTPUT_CHOICES, case_sensitive=False),
    default="text",
    show_default=True,
    help="Output format",
)
def echo(message, upper: bool, lower: bool, repeat: int, output: str):
    """Echo input text with optional transformations."""
    if upper and lower:
        raise click.UsageError("cannot use --upper and --lower together")

    joined = " ".join(message)
    if upper:
        joined = joined.upper()
    if lower:
        joined = joined.lower()

    values = [joined for _ in range(repeat)]
    rendered = format_output(output.lower(), values, lambda: "\n".join(values))
    click.echo(rendered)


@cli.group(name="meta")
def meta_group():
    """Introspect the ado binary and its environment."""


@meta_group.command("info")
@click.option(
    "--output",
    "-o",
    type=click.Choice(OUTPUT_CHOICES, case_sensitive=False),
    default="text",
    show_default=True,
    help="Output format",
)
def meta_info(output: str):
    info = meta.BuildInfo()

    payload = {
        "name": info.name,
        "version": info.version,
        "commit": info.commit,
        "build_time": info.build_time,
        "python_version": info.python_version,
        "platform": info.platform,
    }

    def render_text() -> str:
        return (
            f"Name: {info.name}\n"
            f"Version: {info.version}\n"
            f"Commit: {info.commit}\n"
            f"BuildTime: {info.build_time}\n"
            f"PythonVersion: {info.python_version}\n"
            f"Platform: {info.platform}"
        )

    click.echo(format_output(output.lower(), payload, render_text))


@meta_group.command("env")
@click.option(
    "--output",
    "-o",
    type=click.Choice(OUTPUT_CHOICES, case_sensitive=False),
    default="text",
    show_default=True,
    help="Output format",
)
@click.pass_context
def meta_env(ctx: click.Context, output: str):
    info = meta.collect_env_info(ctx.obj.get("config"))

    def render_text() -> str:
        config_path = info["config_path"] or "(none resolved)"
        lines = [f"ConfigPath: {config_path}", "ConfigSources:"]
        if info["config_sources"]:
            lines.extend([f"  - {src}" for src in info["config_sources"]])
        else:
            lines.append("  (none)")
        lines.append(f"HomeDir: {info['home_dir']}")
        lines.append(f"CacheDir: {info['cache_dir']}")
        lines.append("EnvVariables:")
        if info["env"]:
            lines.extend([f"  {k}={v}" for k, v in info["env"].items()])
        else:
            lines.append("  (none set)")
        return "\n".join(lines)

    click.echo(format_output(output.lower(), info, render_text))


@meta_group.command("features")
@click.option(
    "--output",
    "-o",
    type=click.Choice(OUTPUT_CHOICES, case_sensitive=False),
    default="text",
    show_default=True,
    help="Output format",
)
def meta_features(output: str):
    features = meta.collect_features()
    payload = {"features": features}

    def render_text() -> str:
        if not features:
            return "No experimental features enabled"
        return "\n".join(features)

    click.echo(format_output(output.lower(), payload, render_text))


@cli.command(name="help")
@click.argument("topic", nargs=-1)
@click.pass_context
def help_cmd(ctx: click.Context, topic):
    """Show help for a command."""
    root = ctx.find_root()
    root_cmd = root.command

    if not topic:
        click.echo(root_cmd.get_help(root))
        return

    cmd = root_cmd.get_command(root, topic[0])
    if cmd is None:
        raise click.UsageError(f"Unknown command: {' '.join(topic)}")

    sub_ctx = click.Context(cmd, info_name=cmd.name, parent=root)
    click.echo(cmd.get_help(sub_ctx))


def main(argv: Optional[list[str]] = None):
    cli.main(args=argv, prog_name="ado")


if __name__ == "__main__":
    main(sys.argv[1:])
