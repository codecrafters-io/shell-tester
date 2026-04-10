package internal

import (
	"strings"

	"github.com/codecrafters-io/shell-tester/internal/utils"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

const ASH_PATH = "internal/test_helpers/ash/your_shell.sh"
const DASH_PATH = "internal/test_helpers/dash/your_shell.sh"

func isTestingTesterUsingBusyboxOnAlpine(stageHarness *test_case_harness.TestCaseHarness) bool {
	path := stageHarness.Executable.Path
	isTestingTesterUsingBusybox := strings.HasSuffix(path, ASH_PATH) || strings.HasSuffix(path, DASH_PATH)

	return isTestingTesterUsingBusybox && utils.IsOnAlpine()
}
