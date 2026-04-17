#!/usr/bin/env python3
"""
Buggy shell scenario: passes A1 (qp2) and A2 (gm9) autocompletion stages,
fails A3 (qm8) by appending "_foo" to "xyz" on TAB instead of ringing the bell.

A correct shell receiving TAB on "xyz" (an unrecognised command with no completions)
should ring the bell and leave the line unchanged. This shell instead returns "xyz_foo"
as a completion, so readline rewrites the line to "xyz_foo" and never rings the bell.
"""

from __future__ import annotations

import os
import shlex
import subprocess
import sys

try:
    import readline
except ImportError:
    readline = None  # type: ignore[misc, assignment]

PROMPT = "$ "

BUILTIN_COMPLETIONS: dict[str, str] = {
    "echo": "echo ",
    "exit": "exit ",
    "type": "type ",
    "pwd":  "pwd ",
    "cd":   "cd ",
}


def _tab_completer(text: str, state: int) -> str | None:
    """
    readline completer. Called with state=0, 1, 2, ... until we return None.

    Bug: when text is "xyz" we return "xyz_foo" at state=0 instead of returning
    None (which would make readline ring the bell and leave the line unchanged).
    """
    if state != 0:
        return None

    # Bug: pretend "xyz" has a completion so readline rewrites the line and skips the bell.
    if text == "xyz":
        return "xyz_foo"

    matches = [v for k, v in BUILTIN_COMPLETIONS.items() if k.startswith(text)]

    if len(matches) == 1:
        return matches[0]

    # No unique match: return None so readline rings the bell automatically.
    return None


def _bind_tab_to_completion() -> None:
    """
    Bind TAB to completion for both GNU readline and libedit (macOS system Python).
    libedit ignores `tab: complete`; binding ^I explicitly works for both.
    """
    assert readline is not None
    readline.parse_and_bind("bind ^I rl_complete")
    readline.parse_and_bind("tab: complete")


def _find_in_path(name: str) -> str | None:
    for directory in os.environ.get("PATH", "").split(os.pathsep):
        if not directory:
            continue
        full_path = os.path.join(directory, name)
        if os.path.isfile(full_path) and os.access(full_path, os.X_OK):
            return full_path
    return None


def _run_command(raw_line: str) -> int:
    raw_line = raw_line.strip()
    if not raw_line:
        return 0

    try:
        parts = shlex.split(raw_line, posix=True)
    except ValueError as parse_error:
        print(f"{raw_line}: parse error: {parse_error}", file=sys.stderr)
        return 1

    if not parts:
        return 0

    command_name = parts[0]

    if command_name == "exit":
        exit_code = int(parts[1]) if len(parts) > 1 else 0
        sys.exit(exit_code)

    if command_name == "echo":
        print(" ".join(parts[1:]))
        return 0

    if command_name == "type":
        if len(parts) < 2:
            return 0
        target = parts[1]
        if target in BUILTIN_COMPLETIONS:
            print(f"{target} is a shell builtin")
        else:
            executable_path = _find_in_path(target)
            if executable_path:
                print(f"{target} is {executable_path}")
            else:
                print(f"{target}: not found")
        return 0

    if command_name == "pwd":
        print(os.getcwd())
        return 0

    if command_name == "cd":
        target_dir = parts[1] if len(parts) > 1 else os.environ.get("HOME", "/")
        try:
            os.chdir(target_dir)
        except OSError as cd_error:
            print(f"cd: {cd_error}", file=sys.stderr)
            return 1
        return 0

    executable_path = _find_in_path(command_name)
    if executable_path is None:
        print(f"{command_name}: command not found", file=sys.stderr)
        return 127

    result = subprocess.run(
        [executable_path] + parts[1:],
        env=os.environ,
        stdin=sys.stdin,
        stdout=sys.stdout,
        stderr=sys.stderr,
    )
    return result.returncode


def main() -> None:
    if readline is None:
        print("readline module is required but not available", file=sys.stderr)
        sys.exit(1)

    readline.set_completer(_tab_completer)
    readline.set_completer_delims(" \t\n")
    _bind_tab_to_completion()

    while True:
        try:
            raw_line = input(PROMPT)
        except EOFError:
            break
        except KeyboardInterrupt:
            print("^C")
            continue
        _run_command(raw_line)


if __name__ == "__main__":
    main()
