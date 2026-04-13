#!/bin/sh
# Buggy shell: programmable completion invokes -C completers with correct argv but omits
# COMP_LINE and COMP_POINT from the completer's environment.
exec python3 -u "$(dirname "$0")/main.py" "$@"
