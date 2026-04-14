#!/bin/sh
# This buggy implementation of shell will launch an incorrect executable in background
exec python3 $(dirname "$0")/main.py "$@"