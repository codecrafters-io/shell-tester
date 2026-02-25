#!/bin/sh
# This buggy implementation of shell will always print the output of background launch with improper formatting
exec python3 $(dirname "$0")/main.py "$@"