#!/bin/sh
# This buggy implementation of shell will always print incorrect output of the jobs builtin
exec python3 $(dirname "$0")/main.py "$@"