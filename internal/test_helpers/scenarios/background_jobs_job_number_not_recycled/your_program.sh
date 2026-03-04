#!/bin/sh
# This buggy implementation of shell will never recycle job numbers and keep on increasing
exec python3 $(dirname "$0")/main.py "$@"