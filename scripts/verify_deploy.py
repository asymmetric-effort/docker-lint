# file: scripts/verify_deploy.py
# (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
"""Selenium-based post-deployment verification for the deploy-lint site."""

from __future__ import annotations

import os
import sys
import time
from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.common.by import By
from webdriver_manager.chrome import ChromeDriverManager

DEFAULT_URL = "https://deploy-lint.asymmetric-effort.com"


def create_driver() -> webdriver.Chrome:
    """Return a headless Chrome WebDriver instance."""
    options = Options()
    options.add_argument("--headless=new")
    service = Service(ChromeDriverManager().install())
    return webdriver.Chrome(service=service, options=options)


def check_commit(url: str, commit: str, retries: int = 5, delay: int = 5) -> bool:
    """Return True if *url* loads with matching commit hash."""
    driver = create_driver()
    try:
        for _ in range(retries):
            driver.get(url)
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
    url = os.environ.get("VERIFY_URL", DEFAULT_URL)
    if not check_commit(url, commit):
        raise SystemExit("deployment verification failed")


if __name__ == "__main__":
    main()
