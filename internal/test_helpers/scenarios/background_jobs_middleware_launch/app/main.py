#!/usr/bin/env python3
"""Minimal shell: base stages (oo8–ip1) plus af3 (jobs) and at7 (background + PID)."""

import os
import sys
import subprocess
import shlex

PROMPT = "$ "
BUILTINS = frozenset({"echo", "exit", "type", "jobs"})


def find_in_path(name: str) -> str | None:
    for part in os.environ.get("PATH", "").split(os.pathsep):
        if not part:
            continue
        full = os.path.join(part, name)
        if os.path.isfile(full) and os.access(full, os.X_OK):
            return full
    return None


def run_builtin(argv: list[str], last_exit: int) -> tuple[str | None, int | None]:
    if not argv:
        return None, None
    cmd = argv[0]
    if cmd == "exit":
        code = int(argv[1]) if len(argv) > 1 else last_exit
        return None, code
    if cmd == "echo":
        return (" ".join(argv[1:]) if len(argv) > 1 else ""), None
    if cmd == "type":
        if len(argv) < 2:
            return None, None
        name = argv[1]
        if name in BUILTINS:
            return f"{name} is a shell builtin", None
        p = find_in_path(name)
        if p:
            return f"{name} is {p}", None
        return f"{name}: not found", None
    if cmd == "jobs":
        return None, None
    return None, None


def run_external(argv: list[str], background: bool) -> tuple[int | None, subprocess.Popen | None]:
    name = argv[0]
    path = find_in_path(name)
    if path is None:
        print(f"{name}: command not found", file=sys.stderr)
        return 127, None
    try:
        kw: dict = {"env": os.environ, "executable": path}
        if background:
            kw["stdout"] = subprocess.DEVNULL
            kw["stderr"] = subprocess.DEVNULL
            proc = subprocess.Popen([name] + argv[1:], **kw)
            return None, proc
        r = subprocess.run([name] + argv[1:], **kw, timeout=120)
        return r.returncode, None
    except (OSError, subprocess.TimeoutExpired):
        return 127, None


def main() -> None:
    last_exit = 0
    job_no = 0
    while True:
        sys.stdout.write(PROMPT)
        sys.stdout.flush()
        try:
            line = sys.stdin.readline()
        except (EOFError, KeyboardInterrupt):
            break
        if not line:
            break
        raw = line.rstrip("\n\r")
        bg = raw.rstrip().endswith("&")
        if bg:
            raw = raw.rstrip()[:-1].rstrip()
        parts = shlex.split(raw)
        if not parts:
            continue
        if parts[0] in BUILTINS:
            out, code = run_builtin(parts, last_exit)
            if code is not None:
                sys.exit(code if code >= 0 else 0)
            if out is not None:
                print(out)
            continue
        code, proc = run_external(parts, background=bg)
        if code is not None:
            last_exit = code
        else:
            job_no += 1
            print(f"[{job_no}] {proc.pid}")
            sys.stdout.flush()


if __name__ == "__main__":
    main()
