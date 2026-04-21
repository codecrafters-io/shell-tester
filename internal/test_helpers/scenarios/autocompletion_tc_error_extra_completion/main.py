#!/usr/bin/env python3
"""
Buggy shell scenario:
- Passes A1 (qp2) and A2 (gm9) autocompletion stages.
- Fails A3 (qm8): on TAB after "xyz" it rings the bell but also rewrites the
  line to "xyz_foo" (should leave the line unchanged).
- Passes programmable-completion stages ne7 (PA1) through pm5 (PA4).
- Fails qf1 (PA5): when the completer returns no output, it rings the bell but
  also returns text + "_foo" as a completion instead of leaving the line unchanged.
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
    "echo":     "echo ",
    "exit":     "exit ",
    "type":     "type ",
    "pwd":      "pwd ",
    "cd":       "cd ",
    "complete": "complete ",
}

# command name -> completer script path (for programmable completion)
_COMPLETIONS: dict[str, str] = {}


# ---------------------------------------------------------------------------
# Programmable completion
# ---------------------------------------------------------------------------

def _programmable_complete(text: str, cmd: str, words: list[str], line: str) -> str | None:
    """Invoke the registered -C completer for cmd and return a completion."""
    cur = text
    prev = words[-1] if words else ""

    comp = _COMPLETIONS[cmd]
    env = os.environ.copy()
    env["COMP_LINE"] = line
    env["COMP_POINT"] = str(len(line))

    proc = subprocess.Popen(
        [comp, cmd, cur, prev],
        env=env,
        stdout=subprocess.PIPE,
        stderr=subprocess.DEVNULL,
        text=True,
    )
    stdout_data = proc.communicate()[0] or ""
    out_lines = [x.strip() for x in stdout_data.splitlines() if x.strip()]

    # Bug: when the completer returns no output, ring the bell (correct) but
    # also return text + "_foo" as a completion instead of leaving the line
    # unchanged (incorrect — should return None).
    if not out_lines:
        sys.stdout.write("\x07")
        sys.stdout.flush()
        return text + "_foo"

    matches = [x for x in out_lines if x.startswith(cur)]
    if len(matches) == 1:
        m = matches[0]
        return m if m.endswith(" ") else m + " "
    if not matches:
        sys.stdout.write("\x07")
        sys.stdout.flush()
        return None

    prefix = os.path.commonprefix(matches)
    if len(prefix) > len(cur):
        return prefix
    sys.stdout.write("\x07")
    sys.stdout.flush()
    return None


# ---------------------------------------------------------------------------
# readline completer
# ---------------------------------------------------------------------------

def _tab_completer(text: str, state: int) -> str | None:
    """readline completer. Called with state=0, 1, 2, … until we return None."""
    if state != 0:
        return None

    # Programmable completion takes priority when the command has a registered -C completer.
    if readline is not None:
        line = readline.get_line_buffer()
        beg = readline.get_begidx()
        before = line[:beg].rstrip()
        if before:
            words = before.split()
            cmd = words[0]
            if cmd in _COMPLETIONS:
                return _programmable_complete(text, cmd, words, line)

    # Bug: ring the bell but still rewrite the line to "xyz_foo" instead of
    # leaving it unchanged.
    if text == "xyz":
        sys.stdout.write("\x07")
        sys.stdout.flush()
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


# ---------------------------------------------------------------------------
# Command execution
# ---------------------------------------------------------------------------

def _find_in_path(name: str) -> str | None:
    for directory in os.environ.get("PATH", "").split(os.pathsep):
        if not directory:
            continue
        full_path = os.path.join(directory, name)
        if os.path.isfile(full_path) and os.access(full_path, os.X_OK):
            return full_path
    return None


def _handle_complete(parts: list[str], raw_line: str) -> int:
    """Handle the complete builtin."""
    if len(parts) >= 3 and parts[1] == "-p":
        cmd = parts[2]
        if cmd in _COMPLETIONS:
            print(f"complete -C '{_COMPLETIONS[cmd]}' {cmd}")
        else:
            print(f"complete: {cmd}: no completion specification")
        return 0
    if len(parts) >= 4 and parts[1] == "-C":
        _COMPLETIONS[parts[3]] = parts[2]
        return 0
    # Flexible parsing for extra whitespace / other orderings.
    try:
        toks = shlex.split(raw_line, posix=True)
    except ValueError:
        return 0
    i = 1
    while i < len(toks):
        if toks[i] == "-C" and i + 2 < len(toks):
            _COMPLETIONS[toks[i + 2]] = toks[i + 1]
            return 0
        i += 1
    return 0


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

    if command_name == "complete":
        return _handle_complete(parts, raw_line)

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
