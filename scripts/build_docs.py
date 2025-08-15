# file: scripts/build_docs.py
# (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com

"""Generate static HTML documentation from Markdown sources."""

from __future__ import annotations

import argparse
import os
import shutil
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


def wrap_html(content: str, out: Path, dest: Path) -> str:
    """Wrap body content with basic HTML including icon and favicon."""
    img_path = Path(os.path.relpath(out / "img" / "docker-linter.png", dest.parent)).as_posix()
    favicon_path = Path(os.path.relpath(out / "img" / "favicon.ico", dest.parent)).as_posix()
    return (
        "<!DOCTYPE html>\n<html>\n<head>\n"
        f'<link rel="icon" href="{favicon_path}">\n'
        "</head>\n<body>\n"
        f'<img src="{img_path}" alt="docker-lint icon"/>\n'
        f"{content}\n"
        "</body>\n</html>\n"
    )


def copy_static_assets(src: Path, out: Path) -> None:
    """Copy docs/img assets to the output directory."""
    img_src = src / "docs" / "img"
    if not img_src.exists():
        return
    targets = [out / "img", out / "docs" / "img"]
    for dest in targets:
        dest.mkdir(parents=True, exist_ok=True)
        for file in img_src.iterdir():
            if file.is_file():
                shutil.copy2(file, dest / file.name)


def build_site(src: Path, out: Path) -> None:
    """Generate HTML pages and an index from source Markdown files."""
    files = collect_markdown_files(src)
    copy_static_assets(src, out)
    for md_file in files:
        html = convert_file(md_file)
        dest = out / md_file.relative_to(src).with_suffix(".html")
        wrapped = wrap_html(html, out, dest)
        write_html(wrapped, dest)

    license_text = (src / "LICENSE").read_text(encoding="utf-8")
    license_html = markdown.markdown(license_text)
    license_dest = out / "license.html"
    write_html(wrap_html(license_html, out, license_dest), license_dest)

    links = [
        f'<li><a href="{md.relative_to(src).with_suffix(".html").as_posix()}">{md.stem}</a></li>'
        for md in files
    ]
    links.append('<li><a href="license.html">License</a></li>')
    index_html = "<ul>\n" + "\n".join(links) + "\n</ul>"
    index_dest = out / "index.html"
    write_html(wrap_html(index_html, out, index_dest), index_dest)


def main() -> None:
    """Entry point for command-line execution."""
    parser = argparse.ArgumentParser(description="Build documentation site.")
    parser.add_argument("--src", type=Path, default=Path("."), help="Source directory")
    parser.add_argument("--output", type=Path, required=True, help="Output directory")
    args = parser.parse_args()
    build_site(args.src, args.output)


if __name__ == "__main__":
    main()
