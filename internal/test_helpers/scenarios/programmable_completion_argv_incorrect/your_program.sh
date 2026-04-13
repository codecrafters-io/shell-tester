#!/bin/sh
# This buggy implementation of shell will pass incorrect command
# line arguments to the completer script
exec python3 -u "$(dirname "$0")/main.py" "$@"