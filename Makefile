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

test_ls_against_bsd_ls:
	TESTER_DIR=$(shell pwd) go test -count=1 -p 1 -v ./internal/custom_executable/ls/... -system

test_cat_against_bsd_cat:
	TESTER_DIR=$(shell pwd) go test -count=1 -p 1 -v ./internal/custom_executable/cat/... -system

record_fixtures:
	CODECRAFTERS_RECORD_FIXTURES=true make test

update_tester_utils:
	go get -u github.com/codecrafters-io/tester-utils

copy_course_file:
	gh api repos/codecrafters-io/build-your-own-shell/contents/course-definition.yml \
		| jq -r .content \
		| base64 -d \
		> internal/test_helpers/course_definition.yml

TEST_TARGET ?= test_bash
RUNS ?= 100
test_flakiness:
	@$(foreach i,$(shell seq 1 $(RUNS)), \
		echo "Running iteration $(i)/$(RUNS) of $(TEST_TARGET)" ; \
		make $(TEST_TARGET) > /tmp/test ; \
		if [ "$$?" -ne 0 ]; then \
			echo "Test failed on iteration $(i)" ; \
			cat /tmp/test ; \
			exit 1 ; \
		fi ;\
	)

build_executables:
	oses="darwin linux" ; \
	arches="arm64 amd64" ; \
	for os in $$oses; do \
		for arch in $$arches; do \
		GOOS="$$os" GOARCH="$$arch" go build -o built_executables/ls_$${os}_$${arch} ./internal/custom_executable/ls/ls.go; \
		GOOS="$$os" GOARCH="$$arch" go build -o built_executables/cat_$${os}_$${arch} ./internal/custom_executable/cat/cat.go; \
		done; \
	done

BASE_STAGES = [ \
	{\"slug\":\"oo8\",\"tester_log_prefix\":\"tester::\#oo8\",\"title\":\"Stage\#1: Init\"}, \
	{\"slug\":\"cz2\",\"tester_log_prefix\":\"tester::\#cz2\",\"title\":\"Stage\#2: Invalid Command\"}, \
	{\"slug\":\"ff0\",\"tester_log_prefix\":\"tester::\#ff0\",\"title\":\"Stage\#3: REPL\"}, \
	{\"slug\":\"pn5\",\"tester_log_prefix\":\"tester::\#pn5\",\"title\":\"Stage\#4: Exit\"}, \
	{\"slug\":\"iz3\",\"tester_log_prefix\":\"tester::\#iz3\",\"title\":\"Stage\#5: Echo\"}, \
	{\"slug\":\"ez5\",\"tester_log_prefix\":\"tester::\#ez5\",\"title\":\"Stage\#6: Type built-in\"}, \
	{\"slug\":\"mg5\",\"tester_log_prefix\":\"tester::\#mg5\",\"title\":\"Stage\#7: Type for executables\"}, \
	{\"slug\":\"ip1\",\"tester_log_prefix\":\"tester::\#ip1\",\"title\":\"Stage\#8: Run a program\"} \
]

NAVIGATION_STAGES = [ \
	{\"slug\":\"ei0\",\"tester_log_prefix\":\"tester::\#ei0\",\"title\":\"Stage\#9: PWD\"}, \
	{\"slug\":\"ra6\",\"tester_log_prefix\":\"tester::\#ra6\",\"title\":\"Stage\#10: CD-1\"}, \
	{\"slug\":\"gq9\",\"tester_log_prefix\":\"tester::\#gq9\",\"title\":\"Stage\#11: CD-2\"}, \
	{\"slug\":\"gp4\",\"tester_log_prefix\":\"tester::\#gp4\",\"title\":\"Stage\#12: CD-3\"} \
]

QUOTING_STAGES = [ \
	{\"slug\":\"ni6\",\"tester_log_prefix\":\"tester::\#ni6\",\"title\":\"Stage\#13: Quoting with single quotes\"}, \
	{\"slug\":\"tg6\",\"tester_log_prefix\":\"tester::\#tg6\",\"title\":\"Stage\#14: Quoting with double quotes\"}, \
	{\"slug\":\"yt5\",\"tester_log_prefix\":\"tester::\#yt5\",\"title\":\"Stage\#15: Quoting with backslashes\"}, \
	{\"slug\":\"le5\",\"tester_log_prefix\":\"tester::\#le5\",\"title\":\"Stage\#16: Quoting with single and double quotes\"}, \
	{\"slug\":\"gu3\",\"tester_log_prefix\":\"tester::\#gu3\",\"title\":\"Stage\#17: Quoting with mixed quotes\"}, \
	{\"slug\":\"qj0\",\"tester_log_prefix\":\"tester::\#qj0\",\"title\":\"Stage\#18: Quoting program names\"} \
]

REDIRECTIONS_STAGES = [ \
	{\"slug\":\"jv1\",\"tester_log_prefix\":\"tester::\#jv1\",\"title\":\"Stage\#19: Redirect stdout\"}, \
	{\"slug\":\"vz4\",\"tester_log_prefix\":\"tester::\#vz4\",\"title\":\"Stage\#20: Redirect stderr\"}, \
	{\"slug\":\"el9\",\"tester_log_prefix\":\"tester::\#el9\",\"title\":\"Stage\#21: Append stdout\"}, \
	{\"slug\":\"un3\",\"tester_log_prefix\":\"tester::\#un3\",\"title\":\"Stage\#22: Append stderr\"} \
]

define COMMON_COMPLETION_STAGES
[ \
  {"slug":"qp2","tester_log_prefix":"tester::#qp2","title":"Stage#1: builtins completion"}, \
  {"slug":"gm9","tester_log_prefix":"tester::#gm9","title":"Stage#2: completion with args"}, \
  {"slug":"qm8","tester_log_prefix":"tester::#qm8","title":"Stage#3: completion with invalid command"}, \
  {"slug":"gy5","tester_log_prefix":"tester::#gy5","title":"Stage#4: valid command"}, \
  {"slug":"wt6","tester_log_prefix":"tester::#wt6","title":"Stage#6: partial completions"} \
]
endef

define COMPLETIONS_STAGE5
{"slug":"wh6","tester_log_prefix":"tester::#wh6","title":"Stage#5: completion with multiple executables"}
endef

# Use eval to properly escape the strings
COMPLETIONS_STAGES = $(shell echo '$(COMMON_COMPLETION_STAGES)' | sed 's/]$$/, $(COMPLETIONS_STAGE5)]/' | sed 's/"/\\"/g')
COMPLETIONS_STAGES_FOR_ZSH = $(shell echo '$(COMMON_COMPLETION_STAGES)' | sed 's/"/\\"/g')
REDIRECTIONS_STAGES = $(shell echo '$(_REDIRECTIONS_STAGES)' | sed 's/"/\\"/g')

define run_test
	CODECRAFTERS_REPOSITORY_DIR=./internal/test_helpers/$(2) \
	CODECRAFTERS_TEST_CASES_JSON="$(1)" \
	dist/main.out
endef

test_base_w_ash: build
	$(call run_test,$(BASE_STAGES),ash)

test_nav_w_ash: build
	$(call run_test,$(NAVIGATION_STAGES),ash)

test_quoting_w_ash: build
	$(call run_test,$(QUOTING_STAGES),ash)

test_redirections_w_ash: build
	$(call run_test,$(REDIRECTIONS_STAGES),ash)

test_completions_w_ash: build
	$(call run_test,$(COMPLETIONS_STAGES),ash)

test_base_w_bash: build
	$(call run_test,$(BASE_STAGES),bash)

test_nav_w_bash: build
	$(call run_test,$(NAVIGATION_STAGES),bash)

test_quoting_w_bash: build
	$(call run_test,$(QUOTING_STAGES),bash)

test_redirections_w_bash: build
	$(call run_test,$(REDIRECTIONS_STAGES),bash)

test_completions_w_bash: build
	$(call run_test,$(COMPLETIONS_STAGES),bash)

test_base_w_dash: build
	$(call run_test,$(BASE_STAGES),dash)

test_nav_w_dash: build
	$(call run_test,$(NAVIGATION_STAGES),dash)

test_quoting_w_dash: build
	$(call run_test,$(QUOTING_STAGES),dash)

test_redirections_w_dash: build
	$(call run_test,$(REDIRECTIONS_STAGES),dash)

test_base_w_zsh: build
	$(call run_test,$(BASE_STAGES),zsh)

test_nav_w_zsh: build
	$(call run_test,$(NAVIGATION_STAGES),zsh)

test_quoting_w_zsh: build
	$(call run_test,$(QUOTING_STAGES),zsh)

test_redirections_w_zsh: build
	$(call run_test,$(REDIRECTIONS_STAGES),zsh)

test_completions_w_zsh: build
	$(call run_test,$(COMPLETIONS_STAGES_FOR_ZSH),zsh)

test_ash:
	make test_base_w_ash
	make test_nav_w_ash
	make test_quoting_w_ash
	make test_redirections_w_ash
	make test_completions_w_ash

test_bash:
	make test_base_w_bash
	make test_nav_w_bash
	make test_quoting_w_bash
	make test_redirections_w_bash
	make test_completions_w_bash

# Doesn't support completions
test_dash:
	make test_base_w_dash
	make test_nav_w_dash
	make test_quoting_w_dash
	make test_redirections_w_dash

# Removes ALL zsh config files across the system
test_zsh:
	make test_base_w_zsh
	make test_nav_w_zsh
	make test_quoting_w_zsh
	make test_redirections_w_zsh
	make test_completions_w_zsh