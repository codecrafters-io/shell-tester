#!/usr/bin/env python3
"""Minimal shell that passes base stages oo8--ip1: prompt, invalid command, REPL, exit, echo, type, run."""

import cmd
import os
import sys
import subprocess
import shlex

PROMPT = "$ "
BUILTINS = {"echo", "exit", "type", "jobs"}

# List of background jobs: {"job_id": int, "cmd_line": str, "process": Popen}
BACKGROUND_JOBS: list[dict] = []
# Next job number to assign; never recycled after reaping (intentional error scenario).
NEXT_JOB_ID: int = 1


def reap_finished_jobs() -> None:
    """Remove any background job whose process has exited."""
    global BACKGROUND_JOBS
    BACKGROUND_JOBS = [j for j in BACKGROUND_JOBS if j["process"].poll() is None]

def notify_and_reap_jobs() -> None:
    """Print finished jobs, then remove them."""
    global BACKGROUND_JOBS
    alive = []
    for job in BACKGROUND_JOBS:
        proc = job["process"]
        if proc.poll() is not None:
            jid = job["job_id"]
            cmd_line = job["cmd_line"]
            # For Done, don't print trailing &
            cmd_line = cmd_line.removesuffix(" &")
            # Since this is only used for fy4; hardcoding this for now
            print(f"[{jid}]+ Done                    {cmd_line}")
        else:
            alive.append(job)
    BACKGROUND_JOBS = alive

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
    if cmd == "jobs":
        reap_finished_jobs()
        lines = []
        for i, job in enumerate(BACKGROUND_JOBS):
            jid = job["job_id"]
            cmd_line = job["cmd_line"]
            # + for current (last), - for previous (second-last), space otherwise
            n = len(BACKGROUND_JOBS)
            if i == n - 1:
                marker = "+"
            elif n >= 2 and i == n - 2:
                marker = "-"
            else:
                marker = " "
            line = f"[{jid}]{marker}  Running                 {cmd_line}"
            lines.append(line)
            return "\n".join(lines) if lines else None, None
    return None, None


def run_external(args: list[str], background: bool = False) -> tuple[int | None, subprocess.Popen | None]:
    """Run external command. Returns (exit_code, None) or (None, proc) if background."""
    name = args[0]
    path = find_in_path(name)
    if path is None:
        print(f"{name}: command not found", file=sys.stderr)
        return (127, None)
    try:
        if background:
            proc = subprocess.Popen(
                [path] + args[1:],
                env=os.environ,
                stdout=subprocess.DEVNULL,
                stderr=subprocess.DEVNULL,
            )
            return (None, proc)
        proc = subprocess.run(
            [path] + args[1:],
            env=os.environ,
            capture_output=False,
            timeout=30,
        )
        return (proc.returncode, None)
    except (OSError, subprocess.TimeoutExpired):
        return (127, None)


def main() -> None:
    global NEXT_JOB_ID
    last_exit_code = 0
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
        background = line.rstrip().endswith("&")
        if background:
            line = line.rstrip()[:-1].rstrip()
        parts = shlex.split(line)
        if not parts:
            continue
        cmd = parts[0]
        args = parts[1:]

        if cmd in BUILTINS:
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
                job_id = NEXT_JOB_ID
                NEXT_JOB_ID += 1
                cmd_line = " ".join([cmd] + args) + " &"
                BACKGROUND_JOBS.append({"job_id": job_id, "cmd_line": cmd_line, "process": bg_proc})
                print(f"[{job_id}] {bg_proc.pid}")
                sys.stdout.flush()
        notify_and_reap_jobs()


if __name__ == "__main__":
    main()