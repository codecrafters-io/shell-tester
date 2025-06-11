#!/bin/bash

set -e

# Ensure we're in the correct directory
cd "$(dirname "$0")/.."

# Build and run
docker build -t shell-tester-ash -f local_testing/ash_shell.Dockerfile .
docker run --rm -it shell-tester-ash make test_ash
