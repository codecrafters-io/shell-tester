#!/bin/sh
exec ash -c '
while true; do
    printf "$ "
    read -r cmd
    if [ -z "$cmd" ]; then
        continue
    fi
    if [ "$cmd" = "exit" ] || [ "$cmd" = "exit 0" ]; then
        exit 0
    fi
    # Handle line continuation in single quotes by escaping backslashes
    if echo "$cmd" | grep -q "echo.*'\''.*\\\\\\\\n.*'\''"; then
        # Double the backslashes to prevent line continuation
        modified_cmd=$(echo "$cmd" | sed "s/\\\\\\\\/\\\\\\\\\\\\\\\\/g")
        eval "$modified_cmd"
    else
        eval "$cmd"
    fi
done
'
