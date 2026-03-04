#!/usr/bin/env python3
"""Minimal shell that passes base stages oo8--ip1: prompt, invalid command, REPL, exit, echo, type, run."""

import os
import sys
import subprocess
import shlex

PROMPT = "$ "
BUILTINS = {"echo", "exit", "type"}


def find_in_path(name: str) -> str | None:
    """First executable in PATH with this name (skip non-executable)."""
    path = os.environ.get("PATH", "")
    for part in path.split(os.pathsep):
        if not part:
            continue
        full = os.path.join(part, name)
        if os.path.isfile(full) and os.access(full, os.X_OK):
            return full
    return None


def run_builtin(args: list[str], last_exit_code: int) -> tuple[str | None, int | None]:
    """Run builtin. Returns (output_line_or_None, exit_code_or_None). None means continue (no exit)."""
    if not args:
        return None, None
    cmd = args[0].lower()
    if cmd == "exit":
        code = int(args[1]) if len(args) > 1 else last_exit_code
        if code < 0:
            code = 0
        return None, code
    if cmd == "echo":
        out = " ".join(args[1:]) if len(args) > 1 else ""
        return out, None
    if cmd == "type":
        if len(args) < 2:
            return None, None
        name = args[1]
        if name in BUILTINS:
            return f"{name} is a shell builtin", None
        path = find_in_path(name)
        if path:
            return f"{name} is {path}", None
        return f"{name}: not found", None
    return None, None


def run_external(args: list[str], background: bool = False) -> tuple[int | None, subprocess.Popen | None]:
    """Run external command. Returns (exit_code, None) or (None, popen) if background."""
    name = args[0]
    path = find_in_path(name)
    if path is None:
        print(f"{name}: command not found", file=sys.stderr)
        return (127, None)
    try:
        # Build command so that arg0 is args[0] (the original, like "sleep"), not the full path.
        cmd_argv = [path] + args[1:]
        cmd_argv[0] = args[0]  # Set arg0 properly

        if background:
            proc = subprocess.Popen(
                cmd_argv,
                env=os.environ,
                stdout=subprocess.DEVNULL,
                stderr=subprocess.DEVNULL,
            )
            return (None, proc)
        proc = subprocess.run(
            cmd_argv,
            env=os.environ,
            capture_output=False,
            timeout=30,
        )
        return (proc.returncode, None)
    except (OSError, subprocess.TimeoutExpired):
        return (127, None)


def main() -> None:
    last_exit_code = 0
    job_number = 0
    while True:
        sys.stdout.write(PROMPT)
        sys.stdout.flush()
        try:
            line = sys.stdin.readline()
        except (EOFError, KeyboardInterrupt):
            break
        if not line:
            break
        line = line.rstrip("\n\r")
        # Command reflection: full line as "$ <command>"
        background = line.rstrip().endswith("&")
        if background:
            line = line.rstrip()[:-1].rstrip()  # strip trailing & and spaces
        parts = shlex.split(line)
        if not parts:
            continue
        cmd = parts[0]
        args = parts[1:]

        if cmd in BUILTINS:
            # Builtins always run in foreground (ignore &)
            out, exit_code = run_builtin([cmd] + args, last_exit_code)
            if exit_code is not None:
                sys.exit(exit_code)
            if out is not None:
                print(out)
                sys.stdout.flush()
        else:
            exit_code, bg_proc = run_external([cmd] + args, background=background)
            if exit_code is not None:
                last_exit_code = exit_code
            else:
                job_number += 1
                ## Intentional bug: Print 1234 instead of the actual PID
                print(f"[{job_number}] 1234")
                sys.stdout.flush()


if __name__ == "__main__":
    main()
