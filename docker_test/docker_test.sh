#!/usr/bin/env bash
set -euo pipefail

usage() {
    echo "Usage: $0 test [zsh|bash|ash|all]"
    echo "       $0 record_fixtures"
    echo "       $0 tests_excluding_ash"
    echo "  With no shell type, 'test' uses ash image and runs 'make test' (same as record_fixtures flow)."
    exit 1
}

if [[ $# -lt 1 ]]; then
    usage
fi

MODE="$1"

if [[ "$MODE" != "test" && "$MODE" != "record_fixtures" && "$MODE" != "tests_excluding_ash" ]]; then
    usage
fi

if [[ "$MODE" == "record_fixtures" || "$MODE" == "tests_excluding_ash" ]]; then
    if [[ $# -ne 1 ]]; then
        usage
    fi
else
    if [[ $# -eq 2 ]]; then
        SHELL_TYPE="$2"
        if [[ "$SHELL_TYPE" != "zsh" && "$SHELL_TYPE" != "bash" && "$SHELL_TYPE" != "ash" && "$SHELL_TYPE" != "all" ]]; then
            usage
        fi
    elif [[ $# -ne 1 ]]; then
        usage
    else
        SHELL_TYPE=""
    fi
fi

# Script dir and repo root (dirname "$0" is the dir containing this script)
DOCKER_TEST_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$DOCKER_TEST_DIR/.." && pwd)"
cd "$REPO_ROOT"

DOCKER_RUN_OPTS="--rm"

if [[ "$MODE" == "record_fixtures" ]]; then
    echo "==> Building ash image..."
    docker build -t shell-tester-ash -f "$DOCKER_TEST_DIR/ash_shell.Dockerfile" .
    echo "==> Running record_fixtures..."
    docker run $DOCKER_RUN_OPTS -v "$(pwd)":/home/runner/work/shell-tester/shell-tester shell-tester-ash make record_fixtures
    exit 0
fi

if [[ "$MODE" == "tests_excluding_ash" ]]; then
    echo "==> Building ash image..."
    docker build -t shell-tester-ash -f "$(dirname "$0")"/ash_shell.Dockerfile .
    echo "==> Running tests_excluding_ash..."
    docker run $DOCKER_RUN_OPTS -v "$(pwd)":/home/runner/work/shell-tester/shell-tester shell-tester-ash make tests_excluding_ash
    exit 0
fi

if [[ -z "${SHELL_TYPE:-}" ]]; then
    echo "==> Building ash image..."
    docker build -t shell-tester-ash -f "$(dirname "$0")"/ash_shell.Dockerfile .
    echo "==> Running test..."
    docker run $DOCKER_RUN_OPTS -v "$(pwd)":/home/runner/work/shell-tester/shell-tester shell-tester-ash make test
    exit 0
fi

run_one() {
    local shell="$1"
    local dockerfile image make_target
    case "$shell" in
        bash)
            dockerfile="$DOCKER_TEST_DIR/bash_shell.Dockerfile"
            image="shell-tester-bash"
            make_target="test_bash"
            ;;
        ash)
            dockerfile="$DOCKER_TEST_DIR/ash_shell.Dockerfile"
            image="shell-tester-ash"
            make_target="test_ash"
            ;;
        zsh)
            dockerfile="$DOCKER_TEST_DIR/zsh_shell.Dockerfile"
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