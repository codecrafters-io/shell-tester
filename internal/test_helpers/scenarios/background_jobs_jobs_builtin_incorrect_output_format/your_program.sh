#!/bin/sh
# This buggy implementation of shell will always print the output of jobs builtin in a wrong format
exec python3 $(dirname "$0")/main.py "$@"