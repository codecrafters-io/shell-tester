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
    # Split command and arguments while preserving quotes
    set -- $cmd
    command="$1"
    shift
    # Handle line continuation in single quotes for echo
    if [ "$command" = "echo" ] && echo "$*" | grep -q "'\''.*\\\\\\\\n.*'\''"; then
        # Execute echo with preserved quotes
        echo "$*"
    else
        # Execute command with arguments
        command "$@" 2>&1 || echo "$command: command not found"
    fi
done
'
