#!/bin/bash

# Build the Docker image
docker build -t shell-tester .

# Run the tests in the container
docker run --rm -it shell-tester make test_zsh