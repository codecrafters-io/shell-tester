#!/usr/bin/env bash
set -euo pipefail

usage() {
    echo "Usage: $0 record_fixtures"
    echo "       $0 test [bash|zsh|ash|all]"
    exit 1
}

if [[ $# -lt 1 ]]; then
    usage
fi

CMD="$1"
SHELL_TYPE="${2:-}"

# Ensure we're in repo root
cd "$(dirname "$0")/.."

case "$CMD" in
    record_fixtures)
        if [[ $# -ne 1 ]]; then
            usage
        fi
        make record_fixtures
        ;;
    test)
        if [[ $# -eq 1 ]]; then
            make test
            exit 0
        fi
        if [[ $# -ne 2 ]]; then
            usage
        fi
        # test + shell: run in docker (current behavior)
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

            echo "==> Running $shell tests..."
            docker run --rm -it \
                -v "$(pwd)":/app \
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
        ;;
    *)
        usage
        ;;
esac
