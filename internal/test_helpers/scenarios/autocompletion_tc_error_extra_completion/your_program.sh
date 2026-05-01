#!/bin/sh
# Buggy shell: passes A1 (qp2) and A2 (gm9) but fails A3 (qm8).
# When TAB is pressed on "xyz", it appends "_foo" instead of ringing the bell.
exec python3 -u "$(dirname "$0")/main.py" "$@"
