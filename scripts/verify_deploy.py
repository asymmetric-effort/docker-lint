# file: scripts/verify_deploy.py
# (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
"""Selenium-based post-deployment verification for docker-lint site."""

from __future__ import annotations

import sys
import time
from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.common.by import By
from webdriver_manager.chrome import ChromeDriverManager

URL = "https://docker-lint.asymmetric-effort.com"


def check_commit(commit: str, retries: int = 5, delay: int = 5) -> bool:
    """Return True if site loads with matching commit hash."""
    options = Options()
    options.add_argument("--headless=new")
    driver = webdriver.Chrome(ChromeDriverManager().install(), options=options)
    try:
        for _ in range(retries):
            driver.get(URL)
            try:
                meta = driver.find_element(By.CSS_SELECTOR, "meta[name='docker-lint:commit']")
                if meta.get_attribute("content") == commit:
                    return True
            except Exception:
                pass
            time.sleep(delay)
        return False
    finally:
        driver.quit()


def main() -> None:
    """CLI entry point."""
    if len(sys.argv) != 2:
        raise SystemExit("usage: verify_deploy.py <commit>")
    commit = sys.argv[1]
    if not check_commit(commit):
        raise SystemExit("deployment verification failed")


if __name__ == "__main__":
    main()
