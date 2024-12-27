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
    elif [ "$command" = "cd" ]; then
        # Handle home directory expansion for cd
        args=$(echo "$*" | sed "s|^~/|$HOME/|;s|^~$|$HOME|")
        cd "$args" 2>&1
    else
        # Try to execute command and handle errors
        if ! command -v "$command" >/dev/null 2>&1; then
            echo "$command: command not found"
        else
            "$command" "$@" 2>&1
        fi
    fi
done
'
