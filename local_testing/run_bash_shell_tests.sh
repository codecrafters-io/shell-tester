#!/bin/bash

set -e

# Ensure we're in the correct directory
cd "$(dirname "$0")/.."

# Build and run
docker build -t shell-tester -f local_testing/bash_shell.Dockerfile .
docker run --rm -it -v $(pwd):/app shell-tester make test_bash
