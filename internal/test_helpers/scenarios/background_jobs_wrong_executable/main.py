#!/usr/bin/env python3
"""Minimal shell: base (oo8–ip1) + af3. Background jobs reparent through this Python so at7 path check always fails (not BusyBox/sh/sleep aliasing)."""

import os
import shlex
import subprocess
import sys

PROMPT = "$ "
BUILTINS = {"echo", "exit", "type", "jobs"}


def find_in_path(name: str) -> str | None:
    path = os.environ.get("PATH", "")
    for part in path.split(os.pathsep):
        if not part:
            continue
        full = os.path.join(part, name)
        if os.path.isfile(full) and os.access(full, os.X_OK):
            return full
    return None


def run_builtin(parts: list[str], last_exit: int) -> int | None:
    cmd = parts[0]
    if cmd == "exit":
        if len(parts) > 1:
            try:
                code = int(parts[1])
            except ValueError:
                code = last_exit
        else:
            code = last_exit
        sys.exit(code)
    if cmd == "echo":
        print(" ".join(parts[1:]) if len(parts) > 1 else "")
        return 0
    if cmd == "type":
        if len(parts) < 2:
            return 0
        name = parts[1]
        if name in BUILTINS:
            print(f"{name} is a shell builtin")
        else:
            p = find_in_path(name)
            print(f"{name} is {p}" if p else f"{name}: not found")
        return 0
    if cmd == "jobs":
        return 0
    return None


def run_external(parts: list[str], background: bool) -> tuple[int | None, subprocess.Popen | None]:
    name = parts[0]
    path = find_in_path(name)
    if path is None:
        print(f"{name}: command not found", file=sys.stderr)
        return 127, None
    cmd_argv = [path] + parts[1:]
    cmd_argv[0] = parts[0]
    if background:
        # Intentional: the PID we print is this interpreter (Popen child), not the real command.
        script = (
            "import os, subprocess as s; s.run("
            + repr(parts)
            + ", env=os.environ, stdin=s.DEVNULL, stdout=s.DEVNULL, stderr=s.DEVNULL)"
        )
        proc = subprocess.Popen(
            [sys.executable, "-c", script],
            stdin=subprocess.DEVNULL,
            stdout=subprocess.DEVNULL,
            stderr=subprocess.DEVNULL,
        )
        return None, proc
    r = subprocess.run(cmd_argv, env=os.environ)
    return (r.returncode if r.returncode is not None else 0), None


def main() -> None:
    last_exit = 0
    job_no = 0
    while True:
        sys.stdout.write(PROMPT)
        sys.stdout.flush()
        line = sys.stdin.readline()
        if not line:
            break
        raw = line.rstrip("\n\r")
        bg = raw.rstrip().endswith("&")
        line = raw.rstrip()[:-1].rstrip() if bg else raw.rstrip()
        if not line:
            continue
        try:
            parts = shlex.split(line)
        except ValueError:
            continue
        if not parts:
            continue
        cmd = parts[0]
        if cmd in BUILTINS:
            ec = run_builtin(parts, last_exit)
            if ec is not None:
                last_exit = ec
            continue
        path = find_in_path(cmd)
        if path is None:
            print(f"{cmd}: command not found", file=sys.stderr)
            last_exit = 127
            continue
        if bg:
            _, proc = run_external(parts, True)
            job_no += 1
            print(f"[{job_no}] {proc.pid}")
            sys.stdout.flush()
        else:
            code, _ = run_external(parts, False)
            last_exit = code if code is not None else 0


if __name__ == "__main__":
    main()
