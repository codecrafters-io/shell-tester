#!/bin/sh
# This buggy implementation of shell will always print the output of background launch with incorrect PID
exec python3 $(dirname "$0")/main.py "$@"