#!/bin/sh
# This buggy implementation of shell will never reap the background jobs
exec python3 $(dirname "$0")/main.py "$@"