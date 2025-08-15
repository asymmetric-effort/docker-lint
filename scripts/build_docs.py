# file: scripts/build_docs.py
# (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com

"""Generate static HTML documentation from Markdown sources."""

from __future__ import annotations

import argparse
from pathlib import Path
from typing import List

import markdown


def collect_markdown_files(src: Path) -> List[Path]:
    """Return README.md and all Markdown files under docs/."""
    files: List[Path] = []
    readme = src / "README.md"
    if readme.exists():
        files.append(readme)
    docs_dir = src / "docs"
    if docs_dir.exists():
        files.extend(sorted(docs_dir.rglob("*.md")))
    return files


def convert_file(md_path: Path) -> str:
    """Convert a Markdown file to HTML."""
    text = md_path.read_text(encoding="utf-8")
    return markdown.markdown(text)


def write_html(html: str, dest: Path) -> None:
    """Write HTML content to the destination path."""
    dest.parent.mkdir(parents=True, exist_ok=True)
    dest.write_text(html, encoding="utf-8")


def build_site(src: Path, out: Path) -> None:
    """Generate HTML pages and an index from source Markdown files."""
    files = collect_markdown_files(src)
    for md_file in files:
        html = convert_file(md_file)
        dest = out / md_file.relative_to(src).with_suffix(".html")
        write_html(html, dest)

    license_text = (src / "LICENSE").read_text(encoding="utf-8")
    license_html = markdown.markdown(license_text)
    write_html(license_html, out / "license.html")

    links = [
        f'<li><a href="{md.relative_to(src).with_suffix(".html").as_posix()}">{md.stem}</a></li>'
        for md in files
    ]
    links.append('<li><a href="license.html">License</a></li>')
    index_html = "<ul>\n" + "\n".join(links) + "\n</ul>"
    write_html(index_html, out / "index.html")


def main() -> None:
    """Entry point for command-line execution."""
    parser = argparse.ArgumentParser(description="Build documentation site.")
    parser.add_argument("--src", type=Path, default=Path("."), help="Source directory")
    parser.add_argument("--output", type=Path, required=True, help="Output directory")
    args = parser.parse_args()
    build_site(args.src, args.output)


if __name__ == "__main__":
    main()
