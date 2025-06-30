package internal

import (
	"os"
	"strings"

	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

var BUSYBOX_SHELL_PATHS = []string{
	"internal/test_helpers/ash/your_shell.sh",
	"internal/test_helpers/dash/your_shell.sh",
}

func isTestingTesterUsingBusyboxOnAlpine(stageHarness *test_case_harness.TestCaseHarness) bool {
	path := stageHarness.Executable.Path
	isTestingTesterUsingBusybox := false
	for _, busyboxShellPath := range BUSYBOX_SHELL_PATHS {
		if strings.HasSuffix(path, busyboxShellPath) {
			isTestingTesterUsingBusybox = true
			break
		}
	}

	_, err := os.Stat("/etc/alpine-release")
	isOnAlpine := err == nil

	return isTestingTesterUsingBusybox && isOnAlpine
}
