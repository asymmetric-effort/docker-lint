# file: tests/test_website_structure.py
# (c) 2025 Asymmetric Effort, LLC. MIT License.

"""Tests for static Docker Lint website structure."""

from pathlib import Path

import pytest
from bs4 import BeautifulSoup


def test_index_structure() -> None:
    """index.html should include banner, navigation, and footer."""
    html = Path("index.html").read_text(encoding="utf-8")
    soup = BeautifulSoup(html, "html.parser")
    banner = soup.find("header", id="banner")
    assert banner is not None
    assert "Docker Lint" in banner.get_text()
    top_nav = soup.find("nav", id="top-nav")
    assert top_nav is not None
    texts = [a.get_text() for a in top_nav.find_all("a")]
    assert "HOME" in texts and "CODE REPO" in texts
    footer = soup.find("footer")
    assert footer is not None
    assert "(c) 2025 Asymmetric Effort, LLC. MIT License." in footer.get_text()


def test_banner_height_var() -> None:
    """styles.css defines 100px banner height variable."""
    css = Path("css/styles.css").read_text(encoding="utf-8")
    assert "--banner-h: 100px" in css


def test_app_aria_busy() -> None:
    """app.js references aria-busy state."""
    js = Path("js/app.js").read_text(encoding="utf-8")
    assert "aria-busy" in js

