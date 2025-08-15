# file: tests/test_check_gpg_keys.py
# (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
"""Tests for the check_gpg_keys utility."""

from __future__ import annotations

import os
from pathlib import Path
import subprocess

import pytest

from scripts.check_gpg_keys import check_gpg_keys


def _run(cmd: list[str], homedir: Path) -> subprocess.CompletedProcess[str]:
    """Run a gpg command within *homedir*."""
    env = os.environ | {"GNUPGHOME": str(homedir)}
    return subprocess.run(cmd, env=env, check=True, capture_output=True, text=True)


def _generate_key(email: str, homedir: Path) -> str:
    """Generate a new key and return its fingerprint."""
    _run([
        "gpg",
        "--batch",
        "--pinentry-mode=loopback",
        "--passphrase",
        "",
        "--quick-generate-key",
        email,
        "default",
        "default",
        "never",
    ], homedir)
    out = _run(["gpg", "--with-colons", "--list-keys", email], homedir).stdout
    for line in out.splitlines():
        if line.startswith("fpr:"):
            return line.split(":")[9]
    raise RuntimeError("Fingerprint not found")


def _export_keys(email: str, homedir: Path) -> tuple[str, str]:
    """Return (public, private) ASCII armored keys for *email*."""
    pub = _run(["gpg", "--armor", "--export", email], homedir).stdout
    priv = _run(["gpg", "--armor", "--export-secret-keys", email], homedir).stdout
    return pub, priv


def _sign(target_fpr: str, signer_fpr: str, homedir: Path) -> None:
    """Sign *target_fpr* with *signer_fpr*."""
    _run([
        "gpg",
        "--batch",
        "--pinentry-mode=loopback",
        "--passphrase",
        "",
        "--quick-sign-key",
        "--local-user",
        signer_fpr,
        target_fpr,
    ], homedir)


@pytest.mark.parametrize("signed", [True, False])
def test_check_gpg_keys(tmp_path: Path, monkeypatch: pytest.MonkeyPatch, signed: bool) -> None:
    """check_gpg_keys should validate signature presence."""
    homedir = tmp_path / "gnupg"
    homedir.mkdir()

    signer_fpr = _generate_key("signer@example.com", homedir)
    target_fpr = _generate_key("ci@example.com", homedir)
    if signed:
        _sign(target_fpr, signer_fpr, homedir)
    pub, priv = _export_keys("ci@example.com", homedir)

    monkeypatch.setenv("GPG_SIGNING_KEY_CI_PUBLIC", pub)
    monkeypatch.setenv("GPG_SIGNING_KEY_CI_PRIVATE", priv)
    monkeypatch.setenv("GNUPGHOME", str(homedir))

    if signed:
        check_gpg_keys(signer_fpr)
    else:
        with pytest.raises(SystemExit):
            check_gpg_keys(signer_fpr)
