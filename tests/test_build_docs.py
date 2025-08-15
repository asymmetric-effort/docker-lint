# file: tests/test_build_docs.py
# (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com

"""Tests for documentation site generation."""

from pathlib import Path

from scripts.build_docs import build_site, collect_markdown_files


def test_collect_markdown_files(tmp_path: Path) -> None:
    """collect_markdown_files should find README.md and docs/*.md"""
    (tmp_path / "docs").mkdir()
    (tmp_path / "docs" / "a.md").write_text("# A", encoding="utf-8")
    (tmp_path / "README.md").write_text("# R", encoding="utf-8")
    files = collect_markdown_files(tmp_path)
    expected = {tmp_path / "README.md", tmp_path / "docs" / "a.md"}
    assert set(files) == expected


def test_build_site_generates_html(tmp_path: Path) -> None:
    """build_site should emit html pages, assets, and an index."""
    src = tmp_path
    (src / "docs" / "img").mkdir(parents=True)
    (src / "docs" / "img" / "docker-linter.png").write_bytes(b"icon")
    (src / "docs" / "img" / "favicon.ico").write_bytes(b"icon")
    (src / "docs" / "a.md").write_text("# A", encoding="utf-8")
    (src / "README.md").write_text("![icon](docs/img/docker-linter.png)\n# R", encoding="utf-8")
    (src / "LICENSE").write_text("MIT", encoding="utf-8")
    out = src / "site"
    build_site(src, out)
    assert (out / "README.html").exists()
    assert (out / "docs" / "a.html").exists()
    assert (out / "index.html").exists()
    assert (out / "license.html").exists()


def test_build_site_adds_seo_metadata(tmp_path: Path) -> None:
    """Generated pages should include basic SEO meta tags."""
    src = tmp_path
    (src / "docs").mkdir()
    (src / "docs" / "a.md").write_text("# A\n\nAbout A", encoding="utf-8")
    (src / "README.md").write_text("# R\n\nAbout R", encoding="utf-8")
    (src / "LICENSE").write_text("MIT", encoding="utf-8")
    out = src / "site"
    build_site(src, out)
    html = (out / "README.html").read_text(encoding="utf-8")
    assert "<title>R | docker-lint</title>" in html
    assert '<meta name="description" content="About R"' in html

