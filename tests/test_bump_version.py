# file: tests/test_bump_version.py
# (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
"""Tests for the bump_version utility."""

from __future__ import annotations

from pathlib import Path
import subprocess

from scripts.bump_version import bump_version, increment_version


def init_repo(path: Path) -> None:
    """Initialize a git repo with a single commit."""
    subprocess.run(["git", "init"], cwd=path, check=True)
    subprocess.run(["git", "config", "user.email", "test@example.com"], cwd=path, check=True)
    subprocess.run(["git", "config", "user.name", "Test"], cwd=path, check=True)
    (path / "file.txt").write_text("data", encoding="utf-8")
    subprocess.run(["git", "add", "file.txt"], cwd=path, check=True)
    subprocess.run(["git", "commit", "-m", "init"], cwd=path, check=True)


def git_tags(path: Path) -> list[str]:
    """Return a list of tags in *path*."""
    result = subprocess.run(["git", "tag"], cwd=path, capture_output=True, text=True, check=True)
    return result.stdout.split()


def test_increment_version_minor() -> None:
    """increment_version should bump the minor part."""
    assert increment_version("v1.2.3", "minor") == "v1.3.0"


def test_increment_version_major() -> None:
    """increment_version should bump the major part."""
    assert increment_version("v1.2.3", "major") == "v2.0.0"


def test_bump_version_creates_minor_tag(tmp_path: Path) -> None:
    """bump_version should tag v0.1.0 when no tags exist."""
    init_repo(tmp_path)
    version = bump_version(tmp_path, "minor")
    assert version == "v0.1.0"
    assert "v0.1.0" in git_tags(tmp_path)


def test_bump_version_creates_major_tag(tmp_path: Path) -> None:
    """bump_version should bump to v1.0.0 from v0.1.0."""
    init_repo(tmp_path)
    subprocess.run(["git", "tag", "v0.1.0"], cwd=tmp_path, check=True)
    version = bump_version(tmp_path, "major")
    assert version == "v1.0.0"
    assert "v1.0.0" in git_tags(tmp_path)
