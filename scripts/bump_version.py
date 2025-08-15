# file: scripts/bump_version.py
# (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
"""Utilities for bumping semantic versions and tagging git commits."""

from __future__ import annotations

import argparse
import re
import subprocess
from pathlib import Path

VERSION_RE = re.compile(r"^v(\d+)\.(\d+)\.(\d+)$")


def get_latest_tag(repo: Path) -> str:
    """Return the latest tag in *repo* or v0.0.0 if none exists."""
    result = subprocess.run(
        ["git", "describe", "--tags", "--abbrev=0"],
        cwd=repo,
        capture_output=True,
        text=True,
    )
    if result.returncode != 0:
        return "v0.0.0"
    return result.stdout.strip()


def increment_version(tag: str, part: str) -> str:
    """Increment *tag* by *part* ("minor" or "major")."""
    match = VERSION_RE.match(tag)
    if not match:
        raise ValueError(f"invalid tag: {tag}")
    major, minor, patch = map(int, match.groups())
    if part == "major":
        major += 1
        minor = 0
        patch = 0
    elif part == "minor":
        minor += 1
        patch = 0
    else:
        raise ValueError("part must be 'major' or 'minor'")
    return f"v{major}.{minor}.{patch}"


def tag_commit(repo: Path, version: str) -> None:
    """Tag HEAD of *repo* with *version* and push to origin if configured."""
    subprocess.run(["git", "tag", version], cwd=repo, check=True)
    remotes = subprocess.run(
        ["git", "remote"], cwd=repo, capture_output=True, text=True, check=True
    ).stdout.split()
    if "origin" in remotes:
        subprocess.run(["git", "push", "origin", version], cwd=repo, check=True)


def bump_version(repo: Path, part: str) -> str:
    """Bump the version in *repo* by *part* and tag the current commit."""
    tag = get_latest_tag(repo)
    version = increment_version(tag, part)
    tag_commit(repo, version)
    return version


def main(argv: list[str] | None = None) -> None:
    """CLI entry point."""
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "part",
        choices=["minor", "major"],
        nargs="?",
        default="minor",
        help="which part of the version to increment",
    )
    args = parser.parse_args(argv)
    repo = Path.cwd()
    version = bump_version(repo, args.part)
    print(version)


if __name__ == "__main__":
    main()
