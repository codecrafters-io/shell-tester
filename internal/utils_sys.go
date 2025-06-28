package internal

import (
	"os"
	"strings"

	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func isTestingTesterUsingBusyboxOnAlpine(stageHarness *test_case_harness.TestCaseHarness) bool {
	path := stageHarness.Executable.Path
	isTestingTesterUsingBusybox := strings.HasSuffix(path, "internal/test_helpers/ash/your_shell.sh")

	_, err := os.Stat("/etc/alpine-release")
	isOnAlpine := err == nil

	return isTestingTesterUsingBusybox && isOnAlpine
}
