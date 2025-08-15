# file: scripts/check_gpg_keys.py
# (c) 2025 Asymmetric Effort, LLC. scaldwell@asymmetric-effort.com
"""Utilities to verify CI GPG signing keys."""

from __future__ import annotations

import os
import subprocess
import sys
from datetime import datetime, timezone

SIGNER_KEY_ID = "8528A7AE7B308461"


def _run_gpg(args: list[str], data: str | None = None) -> subprocess.CompletedProcess[str]:
    """Execute gpg with *args* and optional input *data*."""
    return subprocess.run(
        ["gpg", "--batch", "--yes", *args],
        input=data,
        text=True,
        capture_output=True,
        check=True,
    )


def _parse_fingerprint(data: str) -> str:
    """Extract the first fingerprint from colon-formatted *data*."""
    for line in data.splitlines():
        if line.startswith("fpr:"):
            return line.split(":")[9]
    raise ValueError("No fingerprint found")


def import_key(key_data: str) -> str:
    """Import *key_data* and return its fingerprint."""
    info = _run_gpg([
        "--import-options",
        "show-only",
        "--with-colons",
        "--import",
    ], key_data).stdout
    fpr = _parse_fingerprint(info)
    _run_gpg(["--import"], key_data)
    return fpr


def ensure_signer_key(key_id: str) -> None:
    """Ensure that *key_id* is present in the keyring."""
    try:
        _run_gpg(["--list-keys", key_id])
    except subprocess.CalledProcessError:
        _run_gpg(["--keyserver", "hkps://keys.openpgp.org", "--recv-keys", key_id])


def is_key_unexpired(fpr: str) -> bool:
    """Return True if key *fpr* is unexpired."""
    out = _run_gpg(["--with-colons", "--list-keys", fpr]).stdout
    for line in out.splitlines():
        if line.startswith("pub:"):
            expire = line.split(":")[6]
            if not expire:
                return True
            return int(expire) > int(datetime.now(timezone.utc).timestamp())
    return False


def is_key_signed_by(fpr: str, signer_fpr: str) -> bool:
    """Return True if key *fpr* is signed by *signer_fpr*."""
    out = _run_gpg(["--with-colons", "--check-sigs", fpr]).stdout
    for line in out.splitlines():
        if line.startswith("sig:"):
            parts = line.split(":")
            if len(parts) > 12 and parts[12].upper() == signer_fpr.upper():
                return True
    return False


def check_gpg_keys(signer_key_id: str = SIGNER_KEY_ID) -> None:
    """Validate CI signing keys against *signer_key_id*."""
    pub_key = os.environ.get("GPG_SIGNING_KEY_CI_PUBLIC")
    priv_key = os.environ.get("GPG_SIGNING_KEY_CI_PRIVATE")
    if not pub_key or not priv_key:
        raise SystemExit("Missing key data")

    pub_fpr = import_key(pub_key)
    priv_fpr = import_key(priv_key)
    if pub_fpr != priv_fpr:
        raise SystemExit("Public and private key mismatch")

    ensure_signer_key(signer_key_id)
    if not is_key_unexpired(pub_fpr):
        raise SystemExit("Key expired")
    if not is_key_signed_by(pub_fpr, signer_key_id):
        raise SystemExit("Key not signed by required signer")


if __name__ == "__main__":
    try:
        check_gpg_keys()
    except subprocess.CalledProcessError as exc:
        sys.stderr.write(exc.stderr)
        raise SystemExit(1) from exc
