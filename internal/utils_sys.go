package internal

import (
	"os"
	"strings"

	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

const ASH_PATH = "internal/test_helpers/ash/your_shell.sh"

func isTestingTesterUsingBusyboxAshOnAlpine(stageHarness *test_case_harness.TestCaseHarness) bool {
	path := stageHarness.Executable.Path
	isTestingTesterUsingBusyboxAsh := strings.HasSuffix(path, ASH_PATH)

	_, err := os.Stat("/etc/alpine-release")
	isOnAlpine := err == nil

	return isTestingTesterUsingBusyboxAsh && isOnAlpine
}
