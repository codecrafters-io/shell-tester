#!/usr/bin/env python3
"""
Interactive shell for fixture recording: passes programmable-completion stages ne7–zi0 (PA1–PA6),
fails nr7 (PA7) by invoking -C completers without COMP_LINE and COMP_POINT in the environment.
"""

from __future__ import annotations

import os
import shlex
import subprocess
import sys
import threading

try:
    import readline
except ImportError:
    readline = None  # type: ignore[misc, assignment]

PROMPT = "$ "
BUILTINS = {"echo", "exit", "type", "complete"}

# command name -> completer path
_COMPLETIONS: dict[str, str] = {}


def _try_register_complete(line: str) -> bool:
    """Parse `complete ... -C <path> <name>` (flexible whitespace). Returns True if handled."""
    try:
        parts = shlex.split(line, posix=True)
    except ValueError:
        return False
    if not parts or parts[0] != "complete":
        return False
    i = 1
    while i < len(parts):
        if parts[i] == "-C" and i + 2 < len(parts):
            path = parts[i + 1]
            cmd = parts[i + 2]
            _COMPLETIONS[cmd] = path
            return True
        i += 1
    return False


def _list_completions() -> None:
    for cmd in sorted(_COMPLETIONS.keys()):
        path = _COMPLETIONS[cmd]
        print(f"complete -C '{path}' {cmd}")


def _run_builtin(parts: list[str], last_exit: int) -> int | None:
    if not parts:
        return None
    name = parts[0]
    if name == "exit":
        code = int(parts[1]) if len(parts) > 1 else last_exit
        sys.exit(max(code, 0))
    if name == "echo":
        print(" ".join(parts[1:]) if len(parts) > 1 else "")
        return None
    if name == "type":
        if len(parts) < 2:
            return None
        t = parts[1]
        if t in BUILTINS:
            print(f"{t} is a shell builtin")
        else:
            print(f"{t}: not found")
        return None
    if name == "complete":
        if len(parts) == 1:
            _list_completions()
            return None
        return None
    return None


def _find_in_path(name: str) -> str | None:
    for d in os.environ.get("PATH", "").split(os.pathsep):
        if not d:
            continue
        full = os.path.join(d, name)
        if os.path.isfile(full) and os.access(full, os.X_OK):
            return full
    return None


def _run_external(parts: list[str]) -> int:
    exe = _find_in_path(parts[0])
    if exe is None:
        print(f"{parts[0]}: command not found", file=sys.stderr)
        return 127
    argv = [exe] + parts[1:]
    argv[0] = parts[0]
    p = subprocess.run(
        argv,
        env=os.environ,
        stdin=sys.stdin,
        stdout=sys.stdout,
        stderr=sys.stderr,
    )
    return p.returncode


def _process_line(line: str) -> int:
    last = 0
    line = line.strip("\n\r")
    if not line:
        return last
    if _try_register_complete(line):
        return last
    try:
        parts = shlex.split(line, posix=True)
    except ValueError:
        print(f"{line}: parse error", file=sys.stderr)
        return 1
    if not parts:
        return last
    if parts[0] in BUILTINS:
        r = _run_builtin(parts, last)
        return last if r is None else r
    return _run_external(parts)


def _programmable_completer(text: str, state: int) -> str | None:
    if state != 0 or readline is None:
        return None

    line = readline.get_line_buffer()
    beg = readline.get_begidx()
    point = len(line)

    before = line[:beg].rstrip()
    if not before:
        return None
    words = before.split()
    cmd = words[0]
    if cmd not in _COMPLETIONS:
        return None

    prev = words[-1] if words else ""
    cur = text

    comp = _COMPLETIONS[cmd]
    # Intentional bug: do not pass COMP_LINE / COMP_POINT (bash sets these for -C completers).
    env = os.environ.copy()
    env.pop("COMP_LINE", None)
    env.pop("COMP_POINT", None)

    argv = [comp, cmd, cur, prev]

    # Stream completer stderr to the shell's stderr while capturing stdout for candidates.
    # PA5 completer may sleep a long time after writing stderr; bash waits for exit.
    proc = subprocess.Popen(
        argv,
        env=env,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True,
    )
    err_chunks: list[str] = []

    def _forward_stderr() -> None:
        if proc.stderr is None:
            return
        for line in proc.stderr:
            sys.stderr.write(line)
            sys.stderr.flush()
            err_chunks.append(line)

    thr = threading.Thread(target=_forward_stderr, daemon=True)
    thr.start()
    stdout_data = proc.stdout.read() if proc.stdout else ""
    rc = proc.wait()
    thr.join()

    out_lines = [x.strip() for x in stdout_data.splitlines() if x.strip() != ""]
    err_text = "".join(err_chunks)
    err_lines = [x for x in err_text.splitlines() if x.strip() != ""]

    if rc != 0:
        return None

    if not out_lines and err_lines:
        return err_lines[0]

    if not out_lines:
        sys.stdout.write("\x07")
        sys.stdout.flush()
        return None

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


def _bind_tab_triggers_completion() -> None:
    """Make Tab run the completer.

    macOS system Python links libedit, which ignores GNU's `tab: complete` — Tab then inserts a
    literal \\t. Binding ^I (TAB) to rl_complete works for libedit and GNU readline.
    """
    assert readline is not None
    readline.parse_and_bind("bind ^I rl_complete")
    readline.parse_and_bind("tab: complete")


def main() -> None:
    if readline is None:
        print("readline is required for this scenario.", file=sys.stderr)
        sys.exit(1)

    readline.set_completer(_programmable_completer)
    readline.set_completer_delims(" \t\n")
    _bind_tab_triggers_completion()

    last_exit = 0
    while True:
        try:
            line = input(PROMPT)
        except EOFError:
            break
        except KeyboardInterrupt:
            print("^C")
            continue
        last_exit = _process_line(line)


if __name__ == "__main__":
    main()
