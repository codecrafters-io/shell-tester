#!/bin/sh
# This buggy implementation will always print an incorrect job number
exec python3 $(dirname "$0")/main.py "$@"