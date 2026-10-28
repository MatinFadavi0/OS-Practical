#!/usr/bin/env python3
"""Maintainer-only cross-build helper. Students use the resulting executable."""

import argparse
import hashlib
import os
from pathlib import Path
import subprocess


parser = argparse.ArgumentParser(description=__doc__)
parser.add_argument("--target", action="append", help="e.g. linux/arm64; repeatable")
args = parser.parse_args()
root = Path(__file__).resolve().parents[1]
targets = args.target or [f"{system}/{arch}" for system in ("darwin", "linux", "windows")
                          for arch in ("amd64", "arm64")]
dist = root / "dist"
dist.mkdir(exist_ok=True)
for target in targets:
    system, arch = target.split("/")
    if system not in ("darwin", "linux", "windows") or arch not in ("amd64", "arm64"):
        parser.error(f"unsupported target {target}")
    suffix = ".exe" if system == "windows" else ""
    output = dist / f"langar-lab-{system}-{arch}{suffix}"
    print(f"Building {target}", flush=True)
    subprocess.run(["go", "build", "-trimpath", "-buildvcs=false", "-o", str(output), "./cmd/langar-lab"],
                   cwd=root, env=dict(os.environ, GOOS=system, GOARCH=arch, CGO_ENABLED="0"),
                   check=True)
files = sorted(p for p in dist.glob("langar-lab-*") if p.is_file())
(dist / "SHA256SUMS").write_text("".join(
    hashlib.sha256(p.read_bytes()).hexdigest() + "  " + p.name + "\n" for p in files))
