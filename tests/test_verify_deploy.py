# file: tests/test_verify_deploy.py
# (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
"""Tests for the post-deployment verification script."""

from __future__ import annotations

from selenium.webdriver.common.by import By

from scripts.verify_deploy import check_commit


class FakeMeta:
    """Simple stub for selenium WebElement representing a meta tag."""

    def __init__(self, content: str) -> None:
        self._content = content

    def get_attribute(self, name: str) -> str:
        assert name == "content"
        return self._content


class FakeDriver:
    """Stub WebDriver that returns a predetermined meta tag."""

    def __init__(self, content: str | None) -> None:
        self.content = content
        self.quit_called = False
        self.visits: list[str] = []

    def get(self, url: str) -> None:  # pragma: no cover - trivial
        self.visits.append(url)

    def find_element(self, by: str, selector: str) -> FakeMeta:
        assert by == By.CSS_SELECTOR
        assert selector == "meta[name='docker-lint:commit']"
        if self.content is None:
            raise Exception("meta not found")
        return FakeMeta(self.content)

    def quit(self) -> None:  # pragma: no cover - trivial
        self.quit_called = True


def test_check_commit_matches(monkeypatch) -> None:
    """check_commit should return True when commit matches."""
    driver = FakeDriver("abc123")
    monkeypatch.setattr("scripts.verify_deploy.create_driver", lambda: driver)
    assert check_commit("http://example", "abc123", retries=1, delay=0)
    assert driver.quit_called
    assert driver.visits == ["http://example"]


def test_check_commit_mismatch(monkeypatch) -> None:
    """check_commit should return False when commit mismatches."""
    driver = FakeDriver("xyz")
    monkeypatch.setattr("scripts.verify_deploy.create_driver", lambda: driver)
    assert not check_commit("http://example", "abc123", retries=1, delay=0)
