#!/bin/sh
# This buggy implementation will always print a trailing ampersand for finished jobs
exec python3 $(dirname "$0")/main.py "$@"