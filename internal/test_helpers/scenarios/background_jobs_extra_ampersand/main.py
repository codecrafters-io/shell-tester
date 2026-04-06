#!/usr/bin/env python3
"""Buggy shell for ma9 fixtures: `jobs` appends ` &` even for Done entries (bash only does that for Running)."""

import os
import sys
import subprocess
import shlex

PROMPT = "$ "
BUILTINS = {"echo", "exit", "type", "jobs"}

# {"job_id": int, "launch_cmd": str, "process": Popen}
BACKGROUND_JOBS: list[dict] = []


def find_in_path(name: str) -> str | None:
    path = os.environ.get("PATH", "")
    for part in path.split(os.pathsep):
        if not part:
            continue
        full = os.path.join(part, name)
        if os.path.isfile(full) and os.access(full, os.X_OK):
            return full
    return None


def job_status(proc: subprocess.Popen) -> str:
    return "Done" if proc.poll() is not None else "Running"


def format_jobs_line(job_id: int, marker: str, status: str, launch_cmd: str) -> str:
    line = f"[{job_id}]{marker}  {status}                 {launch_cmd}"
    # Intentional bug: always print trailing ` &` like Running, even for Done.
    line += " &"
    return line


def run_builtin(args: list[str], last_exit_code: int) -> tuple[str | None, int | None]:
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
        lines: list[str] = []
        n = len(BACKGROUND_JOBS)
        for i, job in enumerate(BACKGROUND_JOBS):
            jid = job["job_id"]
            launch_cmd = job["launch_cmd"]
            proc = job["process"]
            status = job_status(proc)
            if i == n - 1:
                marker = "+"
            elif n >= 2 and i == n - 2:
                marker = "-"
            else:
                marker = " "
            lines.append(format_jobs_line(jid, marker, status, launch_cmd))
        # Reap terminated jobs after listing (bash-style for this stage).
        BACKGROUND_JOBS[:] = [j for j in BACKGROUND_JOBS if j["process"].poll() is None]
        return "\n".join(lines) if lines else "", None
    return None, None


def run_external(args: list[str], background: bool = False) -> tuple[int | None, subprocess.Popen | None]:
    name = args[0]
    path = find_in_path(name)
    if path is None:
        print(f"{name}: command not found", file=sys.stderr)
        return (127, None)
    try:
        if background:
            proc = subprocess.Popen([path] + args[1:], env=os.environ)
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
                job_id = len(BACKGROUND_JOBS) + 1
                launch_cmd = " ".join([cmd] + args)
                BACKGROUND_JOBS.append(
                    {"job_id": job_id, "launch_cmd": launch_cmd, "process": bg_proc}
                )
                print(f"[{job_id}] {bg_proc.pid}")
                sys.stdout.flush()


if __name__ == "__main__":
    main()
