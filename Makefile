.PHONY: release build

current_version_number := $(shell git tag --list "v*" | sort -V | tail -n 1 | cut -c 2-)
next_version_number := $(shell echo $$(($(current_version_number)+1)))

docs:
	(sleep 0.5 && open http://localhost:6060/pkg/github.com/codecrafters-io/shell-tester/internal/)
	godoc -http=:6060

release:
	git tag v$(next_version_number)
	git push origin main v$(next_version_number)

build:
	go build -o dist/main.out ./cmd/tester

test:
	go test -count=1 -p 1 -v ./internal/...

# ToDo: Update stage slugs
test_bash: build
	CODECRAFTERS_SUBMISSION_DIR=./internal/test_helpers/bash \
	CODECRAFTERS_TEST_CASES_JSON="[ \
		{\"slug\":\"init\",\"tester_log_prefix\":\"tester::#DX1\",\"title\":\"Stage #1: Init\"}, \
		{\"slug\":\"missing-command\",\"tester_log_prefix\":\"tester::#AXY\",\"title\":\"Stage #2: Missing Command\"}, \
		{\"slug\":\"repl\",\"tester_log_prefix\":\"tester::#CX3\",\"title\":\"Stage #3: REPL\"}, \
		{\"slug\":\"echo\",\"tester_log_prefix\":\"tester::#FG3\",\"title\":\"Stage #5: Echo\"}, \
		{\"slug\":\"type1\",\"tester_log_prefix\":\"tester::#BX7\",\"title\":\"Stage #6: Type built-in\"}, \
		{\"slug\":\"type2\",\"tester_log_prefix\":\"tester::#DI8\",\"title\":\"Stage #7: Type for executables\"}, \
		{\"slug\":\"run\",\"tester_log_prefix\":\"tester::#P0D\",\"title\":\"Stage #8: Run a program\"}, \
		{\"slug\":\"pwd\",\"tester_log_prefix\":\"tester::#EXT1\",\"title\":\"Stage #9: PWD\"}, \
		{\"slug\":\"cd1\",\"tester_log_prefix\":\"tester::#EXT2\",\"title\":\"Stage #10: CD-1\"}, \
		{\"slug\":\"cd2\",\"tester_log_prefix\":\"tester::#EXT3\",\"title\":\"Stage #10: CD-2\"}, \
		{\"slug\":\"cd3\",\"tester_log_prefix\":\"tester::#EXT4\",\"title\":\"Stage #10: CD-3\"} \
	]" \
	dist/main.out

test_dash: build
	CODECRAFTERS_SUBMISSION_DIR=./internal/test_helpers/dash \
	CODECRAFTERS_TEST_CASES_JSON="[ \
		{\"slug\":\"init\",\"tester_log_prefix\":\"tester::#DX1\",\"title\":\"Stage #1: Init\"}, \
		{\"slug\":\"missing-command\",\"tester_log_prefix\":\"tester::#AXY\",\"title\":\"Stage #2: Missing Command\"}, \
		{\"slug\":\"repl\",\"tester_log_prefix\":\"tester::#CX3\",\"title\":\"Stage #3: REPL\"}, \
		{\"slug\":\"echo\",\"tester_log_prefix\":\"tester::#FG3\",\"title\":\"Stage #5: Echo\"}, \
		{\"slug\":\"type1\",\"tester_log_prefix\":\"tester::#BX7\",\"title\":\"Stage #6: Type built-in\"}, \
		{\"slug\":\"type2\",\"tester_log_prefix\":\"tester::#DI8\",\"title\":\"Stage #7: Type for executables\"}, \
		{\"slug\":\"run\",\"tester_log_prefix\":\"tester::#P0D\",\"title\":\"Stage #8: Run a program\"}, \
		{\"slug\":\"pwd\",\"tester_log_prefix\":\"tester::#EXT1\",\"title\":\"Stage #9: PWD\"}, \
		{\"slug\":\"cd1\",\"tester_log_prefix\":\"tester::#EXT2\",\"title\":\"Stage #10: CD-1\"}, \
		{\"slug\":\"cd2\",\"tester_log_prefix\":\"tester::#EXT3\",\"title\":\"Stage #10: CD-2\"}, \
		{\"slug\":\"cd3\",\"tester_log_prefix\":\"tester::#EXT4\",\"title\":\"Stage #10: CD-3\"} \
	]" \
	dist/main.out

test_ryan: build
	CODECRAFTERS_SUBMISSION_DIR=./internal/test_helpers/ryan_shell \
	CODECRAFTERS_TEST_CASES_JSON="[ \
		{\"slug\":\"init\",\"tester_log_prefix\":\"tester::#DX1\",\"title\":\"Stage #1: Init\"}, \
		{\"slug\":\"missing-command\",\"tester_log_prefix\":\"tester::#AXY\",\"title\":\"Stage #2: Missing Command\"}, \
		{\"slug\":\"repl\",\"tester_log_prefix\":\"tester::#CX3\",\"title\":\"Stage #3: REPL\"}, \
		{\"slug\":\"echo\",\"tester_log_prefix\":\"tester::#FG3\",\"title\":\"Stage #5: Echo\"}, \
		{\"slug\":\"type1\",\"tester_log_prefix\":\"tester::#BX7\",\"title\":\"Stage #6: Type built-in\"}, \
		{\"slug\":\"type2\",\"tester_log_prefix\":\"tester::#DI8\",\"title\":\"Stage #7: Type for executables\"}, \
		{\"slug\":\"run\",\"tester_log_prefix\":\"tester::#P0D\",\"title\":\"Stage #8: Run a program\"} \
	]" \
	dist/main.out

test_all_success: test_bash test_dash

test_failure: build
	CODECRAFTERS_SUBMISSION_DIR=./internal/test_helpers/failure \
	CODECRAFTERS_TEST_CASES_JSON="[ \
		{\"slug\":\"init\",\"tester_log_prefix\":\"tester::#DX1\",\"title\":\"Stage #1: Init\"}, \
		{\"slug\":\"missing-command\",\"tester_log_prefix\":\"tester::#AX2\",\"title\":\"Stage #2: Missing Command\"}, \
		{\"slug\":\"repl\",\"tester_log_prefix\":\"tester::#CX3\",\"title\":\"Stage #3: REPL\"}, \
		{\"slug\":\"exit\",\"tester_log_prefix\":\"tester::#FG3\",\"title\":\"Stage #4: Exit\"} \
	]" \
	dist/main.out

# Removes ALL zsh related config files across the system
test_zsh_dangerously: build
	CODECRAFTERS_SUBMISSION_DIR=./internal/test_helpers/zsh \
	CODECRAFTERS_TEST_CASES_JSON="[ \
		{\"slug\":\"init\",\"tester_log_prefix\":\"tester::#DX1\",\"title\":\"Stage #1: Init\"}, \
		{\"slug\":\"missing-command\",\"tester_log_prefix\":\"tester::#AXY\",\"title\":\"Stage #2: Missing Command\"}, \
		{\"slug\":\"repl\",\"tester_log_prefix\":\"tester::#CX3\",\"title\":\"Stage #3: REPL\"}, \
		{\"slug\":\"exit\",\"tester_log_prefix\":\"tester::#FG3\",\"title\":\"Stage #4: Exit\"} \
	]" \
	dist/main.out


record_fixtures:
	CODECRAFTERS_RECORD_FIXTURES=true make test

update_tester_utils:
	go get -u github.com/codecrafters-io/tester-utils
