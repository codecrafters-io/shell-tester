package internal

import (
	"os"
	"strings"

	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

const ASH_PATH = "internal/test_helpers/ash/your_shell.sh"
const DASH_PATH = "internal/test_helpers/dash/your_shell.sh"

func isTestingTesterUsingAshOnAlpine(stageHarness *test_case_harness.TestCaseHarness) bool {
	path := stageHarness.Executable.Path
	isTestingTesterUsingBusyboxAsh := strings.HasSuffix(path, ASH_PATH)

	_, err := os.Stat("/etc/alpine-release")
	isOnAlpine := err == nil

	return isTestingTesterUsingBusyboxAsh && isOnAlpine
}

func isTestingTesterUsingDashOnAlpine(stageHarness *test_case_harness.TestCaseHarness) bool {
	path := stageHarness.Executable.Path
	isTestingTesterUsingBusyboxDash := strings.HasSuffix(path, DASH_PATH)

	_, err := os.Stat("/etc/alpine-release")
	isOnAlpine := err == nil

	return isTestingTesterUsingBusyboxDash && isOnAlpine
}
