#!/bin/sh
# A script that allocates memory in a loop until killed.
# Used to test memory limiting functionality.

data=""
while true; do
    # Allocate ~100MB chunks by appending to a string
    data="$data$(head -c $((100 * 1024 * 1024)) /dev/zero | tr '\0' 'x')"
    echo "Length: $((${#data} / 1024 / 1024)) MB"
done 