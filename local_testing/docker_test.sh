#!/usr/bin/env bash
set -euo pipefail

usage() {
    echo "Usage: $0 [bash|ash|zsh|all] [test|record_fixtures]"
    exit 1
}

if [[ $# -ne 2 ]]; then
    usage
fi

SHELL_TYPE="$1"
MODE="$2"

if [[ "$MODE" != "test" && "$MODE" != "record_fixtures" ]]; then
    usage
fi

# Ensure we're in repo root
cd "$(dirname "$0")/.."

run_one() {
    local shell="$1"
    
    case "$shell" in
        bash)
            dockerfile="local_testing/bash_shell.Dockerfile"
            image="shell-tester-bash"
            make_target="test_bash"
            ;;
        ash)
            dockerfile="local_testing/ash_shell.Dockerfile"
            image="shell-tester-ash"
            make_target="test_ash"
            ;;
        zsh)
            dockerfile="local_testing/zsh_shell.Dockerfile"
            image="shell-tester-zsh"
            make_target="test_zsh"
            ;;
        *)
            echo "Unknown shell: $shell"
            exit 1
            ;;
    esac
    
    echo "==> Building $shell image..."
    docker build -t "$image" -f "$dockerfile" .
    
    echo "==> Running $shell tests ($MODE)..."
    
    local env_flags=()
    if [[ "$MODE" == "record_fixtures" ]]; then
        env_flags=(-e CODECRAFTERS_RECORD_FIXTURES=true)
    fi
    
    docker run --rm -it \
        ${env_flags[@]+"${env_flags[@]}"} \
        -v "$(pwd)":/home/runner/work/shell-tester/shell-tester \
        "$image" \
        make "$make_target"
}

if [[ "$SHELL_TYPE" == "all" ]]; then
    run_one bash
    run_one ash
    run_one zsh
else
    run_one "$SHELL_TYPE"
fi