#!/bin/bash

for i in $(seq 1 100)
do
    echo "Running iteration $i"
    make test_debug > /tmp/test
    if [ $? -ne 0 ]; then
        echo "make test_debug failed on iteration $i"
        exit 1
    fi
done