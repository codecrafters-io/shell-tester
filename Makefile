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
	TESTER_DIR=$(shell pwd) go test -count=1 -p 1 -v ./internal/...

test_bash: build
	CODECRAFTERS_REPOSITORY_DIR=./internal/test_helpers/bash \
	CODECRAFTERS_TEST_CASES_JSON="[ \
		{\"slug\":\"oo8\",\"tester_log_prefix\":\"tester::#oo8\",\"title\":\"Stage #1: Init\"}, \
		{\"slug\":\"cz2\",\"tester_log_prefix\":\"tester::#cz2\",\"title\":\"Stage #2: Missing Command\"}, \
		{\"slug\":\"ff0\",\"tester_log_prefix\":\"tester::#ff0\",\"title\":\"Stage #3: REPL\"}, \
		{\"slug\":\"pn5\",\"tester_log_prefix\":\"tester::#pn5\",\"title\":\"Stage #4: Exit\"}, \
		{\"slug\":\"iz3\",\"tester_log_prefix\":\"tester::#iz3\",\"title\":\"Stage #5: Echo\"}, \
		{\"slug\":\"ez5\",\"tester_log_prefix\":\"tester::#ez5\",\"title\":\"Stage #6: Type built-in\"}, \
		{\"slug\":\"mg5\",\"tester_log_prefix\":\"tester::#mg5\",\"title\":\"Stage #7: Type for executables\"}, \
		{\"slug\":\"ip1\",\"tester_log_prefix\":\"tester::#ip1\",\"title\":\"Stage #8: Run a program\"}, \
		{\"slug\":\"ei0\",\"tester_log_prefix\":\"tester::#ei0\",\"title\":\"Stage #9: PWD\"}, \
		{\"slug\":\"ra6\",\"tester_log_prefix\":\"tester::#ra6\",\"title\":\"Stage #10: CD-1\"}, \
		{\"slug\":\"gq9\",\"tester_log_prefix\":\"tester::#gq9\",\"title\":\"Stage #11: CD-2\"}, \
		{\"slug\":\"gp4\",\"tester_log_prefix\":\"tester::#gp4\",\"title\":\"Stage #12: CD-3\"} \
	]" \
	TESTER_DIR=$(shell pwd) dist/main.out

test_dash: build
	CODECRAFTERS_REPOSITORY_DIR=./internal/test_helpers/dash \
	CODECRAFTERS_TEST_CASES_JSON="[ \
		{\"slug\":\"oo8\",\"tester_log_prefix\":\"tester::#oo8\",\"title\":\"Stage #1: Init\"}, \
		{\"slug\":\"cz2\",\"tester_log_prefix\":\"tester::#cz2\",\"title\":\"Stage #2: Missing Command\"}, \
		{\"slug\":\"ff0\",\"tester_log_prefix\":\"tester::#ff0\",\"title\":\"Stage #3: REPL\"}, \
		{\"slug\":\"pn5\",\"tester_log_prefix\":\"tester::#pn5\",\"title\":\"Stage #4: Exit\"}, \
		{\"slug\":\"iz3\",\"tester_log_prefix\":\"tester::#iz3\",\"title\":\"Stage #5: Echo\"}, \
		{\"slug\":\"ez5\",\"tester_log_prefix\":\"tester::#ez5\",\"title\":\"Stage #6: Type built-in\"}, \
		{\"slug\":\"mg5\",\"tester_log_prefix\":\"tester::#mg5\",\"title\":\"Stage #7: Type for executables\"}, \
		{\"slug\":\"ip1\",\"tester_log_prefix\":\"tester::#ip1\",\"title\":\"Stage #8: Run a program\"}, \
		{\"slug\":\"ei0\",\"tester_log_prefix\":\"tester::#ei0\",\"title\":\"Stage #9: PWD\"}, \
		{\"slug\":\"ra6\",\"tester_log_prefix\":\"tester::#ra6\",\"title\":\"Stage #10: CD-1\"}, \
		{\"slug\":\"gq9\",\"tester_log_prefix\":\"tester::#gq9\",\"title\":\"Stage #11: CD-2\"}, \
		{\"slug\":\"gp4\",\"tester_log_prefix\":\"tester::#gp4\",\"title\":\"Stage #12: CD-3\"} \
	]" \
	TESTER_DIR=$(shell pwd) dist/main.out

test_ryan: build
	CODECRAFTERS_REPOSITORY_DIR=./internal/test_helpers/ryan_shell \
	CODECRAFTERS_TEST_CASES_JSON="[ \
		{\"slug\":\"oo8\",\"tester_log_prefix\":\"tester::#oo8\",\"title\":\"Stage #1: Init\"}, \
		{\"slug\":\"cz2\",\"tester_log_prefix\":\"tester::#cz2\",\"title\":\"Stage #2: Missing Command\"}, \
		{\"slug\":\"ff0\",\"tester_log_prefix\":\"tester::#ff0\",\"title\":\"Stage #3: REPL\"}, \
		{\"slug\":\"iz3\",\"tester_log_prefix\":\"tester::#iz3\",\"title\":\"Stage #5: Echo\"}, \
		{\"slug\":\"ez5\",\"tester_log_prefix\":\"tester::#ez5\",\"title\":\"Stage #6: Type built-in\"}, \
		{\"slug\":\"mg5\",\"tester_log_prefix\":\"tester::#mg5\",\"title\":\"Stage #7: Type for executables\"}, \
		{\"slug\":\"ip1\",\"tester_log_prefix\":\"tester::#ip1\",\"title\":\"Stage #8: Run a program\"} \
	]" \
	TESTER_DIR=$(shell pwd) dist/main.out

test_all_success: test_bash test_dash

test_failure: build
	CODECRAFTERS_REPOSITORY_DIR=./internal/test_helpers/failure \
	CODECRAFTERS_TEST_CASES_JSON="[ \
		{\"slug\":\"oo8\",\"tester_log_prefix\":\"tester::#oo8\",\"title\":\"Stage #1: Init\"}, \
		{\"slug\":\"cz2\",\"tester_log_prefix\":\"tester::#cz2\",\"title\":\"Stage #2: Missing Command\"}, \
		{\"slug\":\"ff0\",\"tester_log_prefix\":\"tester::#ff0\",\"title\":\"Stage #3: REPL\"}, \
		{\"slug\":\"pn5\",\"tester_log_prefix\":\"tester::#pn5\",\"title\":\"Stage #4: Exit\"} \
	]" \
	dist/main.out

# Removes ALL zsh related config files across the system
test_zsh_dangerously: build
	CODECRAFTERS_REPOSITORY_DIR=./internal/test_helpers/zsh \
	CODECRAFTERS_TEST_CASES_JSON="[ \
		{\"slug\":\"oo8\",\"tester_log_prefix\":\"tester::#oo8\",\"title\":\"Stage #1: Init\"}, \
		{\"slug\":\"cz2\",\"tester_log_prefix\":\"tester::#cz2\",\"title\":\"Stage #2: Missing Command\"}, \
		{\"slug\":\"ff0\",\"tester_log_prefix\":\"tester::#ff0\",\"title\":\"Stage #3: REPL\"}, \
		{\"slug\":\"pn5\",\"tester_log_prefix\":\"tester::#pn5\",\"title\":\"Stage #4: Exit\"} \
	]" \
	dist/main.out


record_fixtures:
	CODECRAFTERS_RECORD_FIXTURES=true make test

update_tester_utils:
	go get -u github.com/codecrafters-io/tester-utils

copy_course_file:
	gh api repos/codecrafters-io/build-your-own-shell/contents/course-definition.yml \
		| jq -r .content \
		| base64 -d \
		> internal/test_helpers/course_definition.yml
