# file: tests/test_cname.py
# (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com

"""Tests for GitHub Pages CNAME file."""

from pathlib import Path


def test_cname_file_contains_expected_domain() -> None:
    """CNAME file should exist and contain the custom domain."""
    cname_path = Path("docs") / "CNAME"
    assert cname_path.exists(), "CNAME file missing"
    content = cname_path.read_text(encoding="utf-8").strip()
    assert content == "docker-lint.asymmetric-effort.com"
