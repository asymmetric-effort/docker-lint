# file: scripts/build_docs.py
# (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com

"""Generate static HTML documentation from Markdown sources."""

from __future__ import annotations

import argparse
from pathlib import Path
from typing import List, Tuple

import markdown

HTML_TEMPLATE = """<!doctype html>
<html lang=\"en\">
<head>
<meta charset=\"utf-8\">
<title>{title}</title>
<meta name=\"description\" content=\"{description}\">
</head>
<body>
{content}
</body>
</html>
"""


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


def extract_title(text: str, fallback: str) -> str:
    """Return the first Markdown heading or fallback."""
    for line in text.splitlines():
        if line.startswith("#"):
            return line.lstrip("#").strip()
    return fallback


def extract_description(text: str) -> str:
    """Return the first non-heading line for description."""
    for line in text.splitlines():
        stripped = line.strip()
        if stripped and not stripped.startswith("#"):
            return stripped
    return ""


def convert_file(md_path: Path) -> Tuple[str, str, str]:
    """Convert a Markdown file to HTML and extract metadata."""
    text = md_path.read_text(encoding="utf-8")
    html = markdown.markdown(text)
    title = extract_title(text, md_path.stem)
    description = extract_description(text)
    return html, title, description


def render_page(title: str, description: str, body: str) -> str:
    """Render a full HTML page with SEO tags."""
    return HTML_TEMPLATE.format(title=title, description=description, content=body)


def write_html(html: str, dest: Path) -> None:
    """Write HTML content to the destination path."""
    dest.parent.mkdir(parents=True, exist_ok=True)
    dest.write_text(html, encoding="utf-8")


def build_site(src: Path, out: Path) -> None:
    """Generate HTML pages and an index from source Markdown files."""
    files = collect_markdown_files(src)
    for md_file in files:
        body, title, description = convert_file(md_file)
        page_html = render_page(f"{title} | docker-lint", description, body)
        dest = out / md_file.relative_to(src).with_suffix(".html")
        write_html(page_html, dest)

    license_text = (src / "LICENSE").read_text(encoding="utf-8")
    license_html = markdown.markdown(license_text)
    license_page = render_page("License | docker-lint", "", license_html)
    write_html(license_page, out / "license.html")

    links = [
        f'<li><a href="{md.relative_to(src).with_suffix(".html").as_posix()}">{md.stem}</a></li>'
        for md in files
    ]
    links.append('<li><a href="license.html">License</a></li>')
    index_body = "<ul>\n" + "\n".join(links) + "\n</ul>"
    index_html = render_page("Index | docker-lint", "Documentation index", index_body)
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
