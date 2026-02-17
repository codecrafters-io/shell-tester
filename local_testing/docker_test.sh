#!/usr/bin/env bash
set -euo pipefail

usage() {
    echo "Usage: $0 test [zsh|bash|ash|all]"
    echo "       $0 record_fixtures"
    exit 1
}

if [[ $# -lt 1 ]]; then
    usage
fi

MODE="$1"

if [[ "$MODE" != "test" && "$MODE" != "record_fixtures" ]]; then
    usage
fi

if [[ "$MODE" == "record_fixtures" ]]; then
    if [[ $# -ne 1 ]]; then
        usage
    fi
else
    if [[ $# -ne 2 ]]; then
        usage
    fi
    SHELL_TYPE="$2"
    if [[ "$SHELL_TYPE" != "zsh" && "$SHELL_TYPE" != "bash" && "$SHELL_TYPE" != "ash" && "$SHELL_TYPE" != "all" ]]; then
        usage
    fi
fi

# Ensure we're in repo root
cd "$(dirname "$0")/.."

# Use -it only when stdin is a TTY (e.g. local); CI has no TTY
DOCKER_RUN_OPTS="--rm $([ -t 0 ] && echo -it || echo -i)"

if [[ "$MODE" == "record_fixtures" ]]; then
    echo "==> Building ash image..."
    docker build -t shell-tester-ash -f local_testing/ash_shell.Dockerfile .
    echo "==> Running record_fixtures..."
    docker run $DOCKER_RUN_OPTS -v "$(pwd)":/home/runner/work/shell-tester/shell-tester shell-tester-ash make record_fixtures
    exit 0
fi

run_one() {
    local shell="$1"
    local dockerfile image make_target
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
    echo "==> Running $make_target..."
    docker run $DOCKER_RUN_OPTS -v "$(pwd)":/home/runner/work/shell-tester/shell-tester "$image" make "$make_target"
}

if [[ "$SHELL_TYPE" == "all" ]]; then
    run_one bash
    run_one ash
    run_one zsh
else
    run_one "$SHELL_TYPE"
fi