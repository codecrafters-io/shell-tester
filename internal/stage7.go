package internal

import (
	"fmt"
	"os/exec"
	"regexp"

	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/shell-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func getPath(executable string) string {
	path, err := exec.LookPath(executable)
	if err != nil {
		return fmt.Sprintf(`%s[:]? not found`, executable)
	} else {
		return path
	}
}

func testType2(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	shell := shell_executable.NewShellExecutable(stageHarness)

	executables := []string{"cat", "cp", "mkdir", "awk", "whoami", "go", "python3", "nonexistent"}

	if err := shell.Start(); err != nil {
		return err
	}

	for _, executable := range executables {
		command := fmt.Sprintf("type %s", executable)
		expectedPattern := getPath(executable)

		testCase := test_cases.RegexTestCase{
			Command:                    command,
			ExpectedPattern:            regexp.MustCompile(expectedPattern + "\r\n"),
			ExpectedPatternExplanation: fmt.Sprintf("match %q", expectedPattern),
			SuccessMessage:             "Received expected response",
		}
		if err := testCase.Run(shell, logger); err != nil {
			return err
		}
	}

	// ToDo: Add check for shell still running

	return nil
}
